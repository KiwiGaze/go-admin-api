package dto

import "testing"

func TestPaginationDefaults(t *testing.T) {
	t.Run("page index defaults to one", func(t *testing.T) {
		if got := (&Pagination{}).GetPageIndex(); got != 1 {
			t.Fatalf("GetPageIndex() = %d, want 1", got)
		}
	})

	t.Run("page size defaults to ten", func(t *testing.T) {
		if got := (&Pagination{}).GetPageSize(); got != 10 {
			t.Fatalf("GetPageSize() = %d, want 10", got)
		}
	})

	t.Run("positive values are preserved", func(t *testing.T) {
		p := &Pagination{PageIndex: 3, PageSize: 25}
		if got := p.GetPageIndex(); got != 3 {
			t.Fatalf("GetPageIndex() = %d, want 3", got)
		}
		if got := p.GetPageSize(); got != 25 {
			t.Fatalf("GetPageSize() = %d, want 25", got)
		}
	})
}
