package dto

type Pagination struct {
	PageIndex int `form:"pageIndex"`
	PageSize  int `form:"pageSize"`
}

func (p *Pagination) GetPageIndex() int {
	if p.PageIndex <= 0 {
		return 1
	}
	return p.PageIndex
}

func (p *Pagination) GetPageSize() int {
	if p.PageSize <= 0 {
		return 10
	}
	return p.PageSize
}