package data

import "time"

// Comment represents a GitHub issue comment.
type Comment struct {
	Author    string
	Body      string
	CreatedAt time.Time
	UpdatedAt time.Time
	Reactions int
}

// CommentListResult is the result of listing comments.
type CommentListResult struct {
	Comments []Comment
	PageInfo PageInfo
}

// CommentClient fetches comment data via GraphQL.
type CommentClient struct {
	querier Querier
	owner   string
	repo    string
}

// NewCommentClient creates a CommentClient for the given repository.
func NewCommentClient(q Querier, owner, repo string) *CommentClient {
	return &CommentClient{querier: q, owner: owner, repo: repo}
}

// List fetches comments for an issue.
func (c *CommentClient) List(issueNumber, first int, after string) (CommentListResult, error) {
	if first == 0 {
		first = 25
	}

	vars := map[string]interface{}{
		"owner":  c.owner,
		"name":   c.repo,
		"number": issueNumber,
		"first":  first,
	}
	if after != "" {
		vars["after"] = after
	}

	var resp listCommentsResponse
	if err := c.querier.Do(listCommentsQuery, vars, &resp); err != nil {
		return CommentListResult{}, err
	}

	return resp.toResult(), nil
}

const listCommentsQuery = `query ListComments($owner: String!, $name: String!, $number: Int!, $first: Int!, $after: String) {
  repository(owner: $owner, name: $name) {
    issue(number: $number) {
      comments(first: $first, after: $after) {
        pageInfo { hasNextPage endCursor }
        nodes {
          author { login }
          body
          createdAt
          updatedAt
          reactions { totalCount }
        }
      }
    }
  }
}`

type listCommentsResponse struct {
	Repository struct {
		Issue struct {
			Comments struct {
				PageInfo graphqlPageInfo `json:"pageInfo"`
				Nodes    []commentNode   `json:"nodes"`
			} `json:"comments"`
		} `json:"issue"`
	} `json:"repository"`
}

func (r *listCommentsResponse) toResult() CommentListResult {
	comments := make([]Comment, len(r.Repository.Issue.Comments.Nodes))
	for i, n := range r.Repository.Issue.Comments.Nodes {
		author := n.Author.Login
		if author == "" {
			author = "[deleted]"
		}
		comments[i] = Comment{
			Author:    author,
			Body:      n.Body,
			CreatedAt: n.CreatedAt,
			UpdatedAt: n.UpdatedAt,
			Reactions: n.Reactions.TotalCount,
		}
	}
	pi := r.Repository.Issue.Comments.PageInfo
	return CommentListResult{
		Comments: comments,
		PageInfo: PageInfo(pi),
	}
}

type commentNode struct {
	Author struct {
		Login string `json:"login"`
	} `json:"author"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Reactions struct {
		TotalCount int `json:"totalCount"`
	} `json:"reactions"`
}
