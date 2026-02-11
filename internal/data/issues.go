package data

import "time"

// Querier abstracts a GraphQL client for testability.
type Querier interface {
	Do(query string, variables map[string]interface{}, response interface{}) error
}

// IssueClient fetches issue data via GraphQL.
type IssueClient struct {
	querier Querier
	owner   string
	repo    string
}

// NewIssueClient creates an IssueClient for the given repository.
func NewIssueClient(q Querier, owner, repo string) *IssueClient {
	return &IssueClient{querier: q, owner: owner, repo: repo}
}

// List fetches a page of issues matching the given options.
func (c *IssueClient) List(opts IssueListOptions) (IssueListResult, error) {
	if opts.First == 0 {
		opts.First = 50
	}
	if opts.OrderBy.Field == "" {
		opts.OrderBy.Field = "CREATED_AT"
		opts.OrderBy.Direction = "DESC"
	}

	vars := map[string]interface{}{
		"owner":   c.owner,
		"name":    c.repo,
		"first":   opts.First,
		"orderBy": map[string]interface{}{"field": opts.OrderBy.Field, "direction": opts.OrderBy.Direction},
	}
	if opts.After != "" {
		vars["after"] = opts.After
	}
	if len(opts.States) > 0 {
		vars["states"] = opts.States
	}
	if len(opts.Labels) > 0 {
		vars["labels"] = opts.Labels
	}

	var resp listIssuesResponse
	if err := c.querier.Do(listIssuesQuery, vars, &resp); err != nil {
		return IssueListResult{}, err
	}

	return resp.toResult(), nil
}

// Get fetches a single issue by number, including its body.
func (c *IssueClient) Get(number int) (Issue, error) {
	vars := map[string]interface{}{
		"owner":  c.owner,
		"name":   c.repo,
		"number": number,
	}

	var resp getIssueResponse
	if err := c.querier.Do(getIssueQuery, vars, &resp); err != nil {
		return Issue{}, err
	}

	node := resp.Repository.Issue
	return node.toIssue(), nil
}

// GraphQL queries

const listIssuesQuery = `query ListIssues($owner: String!, $name: String!, $first: Int!, $after: String, $states: [IssueState!], $labels: [String!], $orderBy: IssueOrder!) {
  repository(owner: $owner, name: $name) {
    issues(first: $first, after: $after, states: $states, labels: $labels, orderBy: $orderBy) {
      pageInfo { hasNextPage endCursor }
      nodes {
        number
        title
        state
        createdAt
        updatedAt
        author { login }
        labels(first: 10) { nodes { name color } }
        assignees(first: 5) { nodes { login } }
        milestone { title }
        comments { totalCount }
        reactions { totalCount }
      }
    }
  }
}`

const getIssueQuery = `query GetIssue($owner: String!, $name: String!, $number: Int!) {
  repository(owner: $owner, name: $name) {
    issue(number: $number) {
      number
      title
      state
      createdAt
      updatedAt
      author { login }
      labels(first: 10) { nodes { name color } }
      assignees(first: 5) { nodes { login } }
      milestone { title }
      comments { totalCount }
      reactions { totalCount }
      body
    }
  }
}`

// Internal response structs mirroring GraphQL JSON shape.

type listIssuesResponse struct {
	Repository struct {
		Issues struct {
			PageInfo graphqlPageInfo `json:"pageInfo"`
			Nodes    []issueNode    `json:"nodes"`
		} `json:"issues"`
	} `json:"repository"`
}

func (r *listIssuesResponse) toResult() IssueListResult {
	issues := make([]Issue, len(r.Repository.Issues.Nodes))
	for i, n := range r.Repository.Issues.Nodes {
		issues[i] = n.toIssue()
	}
	pi := r.Repository.Issues.PageInfo
	return IssueListResult{
		Issues: issues,
		PageInfo: PageInfo{
			HasNextPage: pi.HasNextPage,
			EndCursor:   pi.EndCursor,
		},
	}
}

type getIssueResponse struct {
	Repository struct {
		Issue issueNode `json:"issue"`
	} `json:"repository"`
}

type issueNode struct {
	Number    int       `json:"number"`
	Title     string    `json:"title"`
	State     string    `json:"state"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Author    struct {
		Login string `json:"login"`
	} `json:"author"`
	Labels struct {
		Nodes []struct {
			Name  string `json:"name"`
			Color string `json:"color"`
		} `json:"nodes"`
	} `json:"labels"`
	Assignees struct {
		Nodes []struct {
			Login string `json:"login"`
		} `json:"nodes"`
	} `json:"assignees"`
	Milestone *struct {
		Title string `json:"title"`
	} `json:"milestone"`
	Comments struct {
		TotalCount int `json:"totalCount"`
	} `json:"comments"`
	Reactions struct {
		TotalCount int `json:"totalCount"`
	} `json:"reactions"`
	Body string `json:"body"`
}

func (n *issueNode) toIssue() Issue {
	labels := make([]Label, len(n.Labels.Nodes))
	for i, l := range n.Labels.Nodes {
		labels[i] = Label{Name: l.Name, Color: l.Color}
	}

	assignees := make([]string, len(n.Assignees.Nodes))
	for i, a := range n.Assignees.Nodes {
		assignees[i] = a.Login
	}

	milestone := ""
	if n.Milestone != nil {
		milestone = n.Milestone.Title
	}

	return Issue{
		Number:        n.Number,
		Title:         n.Title,
		State:         n.State,
		CreatedAt:     n.CreatedAt,
		UpdatedAt:     n.UpdatedAt,
		Author:        n.Author.Login,
		Labels:        labels,
		Assignees:     assignees,
		Milestone:     milestone,
		CommentCount:  n.Comments.TotalCount,
		ReactionCount: n.Reactions.TotalCount,
		Body:          n.Body,
	}
}

type graphqlPageInfo struct {
	HasNextPage bool   `json:"hasNextPage"`
	EndCursor   string `json:"endCursor"`
}
