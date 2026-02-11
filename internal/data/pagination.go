package data

// Paginator manages cursor-based pagination state.
type Paginator struct {
	cursor      string
	pageSize    int
	totalLoaded int
	hasNextPage bool
}

// NewPaginator creates a new paginator with the given page size.
func NewPaginator(pageSize int) *Paginator {
	if pageSize <= 0 {
		pageSize = 50
	}
	return &Paginator{
		pageSize:    pageSize,
		hasNextPage: true,
	}
}

// PageRequest holds the parameters for a page request.
type PageRequest struct {
	First int
	After string
}

// NextPageRequest returns the parameters for the next page, or nil if exhausted.
func (p *Paginator) NextPageRequest() *PageRequest {
	if !p.hasNextPage {
		return nil
	}
	req := &PageRequest{First: p.pageSize}
	if p.cursor != "" {
		req.After = p.cursor
	}
	return req
}

// Update records the result of a page fetch.
func (p *Paginator) Update(pageInfo PageInfo, count int) {
	p.hasNextPage = pageInfo.HasNextPage
	p.cursor = pageInfo.EndCursor
	p.totalLoaded += count
}

// HasNextPage returns whether more pages are available.
func (p *Paginator) HasNextPage() bool {
	return p.hasNextPage
}

// TotalLoaded returns the total number of items loaded so far.
func (p *Paginator) TotalLoaded() int {
	return p.totalLoaded
}

// Reset resets the paginator to its initial state.
func (p *Paginator) Reset() {
	p.cursor = ""
	p.totalLoaded = 0
	p.hasNextPage = true
}
