package data

import "time"

// Issue represents a GitHub issue.
type Issue struct {
	Number        int
	Title         string
	State         string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Author        string
	Labels        []Label
	Assignees     []string
	Milestone     string
	CommentCount  int
	ReactionCount int
	Body          string
}

// Label represents a GitHub label.
type Label struct {
	Name  string
	Color string // hex color without '#'
}

// PageInfo holds cursor-based pagination state from GraphQL.
type PageInfo struct {
	HasNextPage bool
	EndCursor   string
}

// IssueListResult is the result of listing issues.
type IssueListResult struct {
	Issues   []Issue
	PageInfo PageInfo
}

// IssueListOptions configures an issue list query.
type IssueListOptions struct {
	States  []string // "OPEN", "CLOSED"
	Labels  []string
	OrderBy IssueOrder
	First   int
	After   string
}

// IssueOrder specifies how to sort issues.
type IssueOrder struct {
	Field     string // "CREATED_AT", "UPDATED_AT", "COMMENTS"
	Direction string // "ASC", "DESC"
}
