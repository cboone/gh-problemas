package data

import "testing"

func TestPaginator_InitialState(t *testing.T) {
	p := NewPaginator(50)
	if !p.HasNextPage() {
		t.Error("expected HasNextPage true initially")
	}
	if p.TotalLoaded() != 0 {
		t.Errorf("expected TotalLoaded 0, got %d", p.TotalLoaded())
	}
	req := p.NextPageRequest()
	if req == nil {
		t.Fatal("expected non-nil request")
	}
	if req.First != 50 {
		t.Errorf("expected First 50, got %d", req.First)
	}
	if req.After != "" {
		t.Errorf("expected empty After, got %q", req.After)
	}
}

func TestPaginator_SequentialPages(t *testing.T) {
	p := NewPaginator(10)

	// First page
	p.Update(PageInfo{HasNextPage: true, EndCursor: "cursor1"}, 10)
	if p.TotalLoaded() != 10 {
		t.Errorf("expected TotalLoaded 10, got %d", p.TotalLoaded())
	}

	req := p.NextPageRequest()
	if req == nil {
		t.Fatal("expected non-nil request after first page")
	}
	if req.After != "cursor1" {
		t.Errorf("expected After cursor1, got %q", req.After)
	}

	// Second page
	p.Update(PageInfo{HasNextPage: true, EndCursor: "cursor2"}, 10)
	if p.TotalLoaded() != 20 {
		t.Errorf("expected TotalLoaded 20, got %d", p.TotalLoaded())
	}

	req = p.NextPageRequest()
	if req.After != "cursor2" {
		t.Errorf("expected After cursor2, got %q", req.After)
	}
}

func TestPaginator_Exhausted(t *testing.T) {
	p := NewPaginator(10)
	p.Update(PageInfo{HasNextPage: false, EndCursor: ""}, 5)

	if p.HasNextPage() {
		t.Error("expected HasNextPage false after exhausted")
	}
	req := p.NextPageRequest()
	if req != nil {
		t.Error("expected nil request when exhausted")
	}
}

func TestPaginator_Reset(t *testing.T) {
	p := NewPaginator(10)
	p.Update(PageInfo{HasNextPage: false, EndCursor: "cursor1"}, 10)

	p.Reset()
	if !p.HasNextPage() {
		t.Error("expected HasNextPage true after reset")
	}
	if p.TotalLoaded() != 0 {
		t.Errorf("expected TotalLoaded 0 after reset, got %d", p.TotalLoaded())
	}
	req := p.NextPageRequest()
	if req == nil {
		t.Fatal("expected non-nil request after reset")
	}
	if req.After != "" {
		t.Errorf("expected empty After after reset, got %q", req.After)
	}
}

func TestPaginator_DefaultPageSize(t *testing.T) {
	p := NewPaginator(0)
	req := p.NextPageRequest()
	if req.First != 50 {
		t.Errorf("expected default page size 50, got %d", req.First)
	}
}
