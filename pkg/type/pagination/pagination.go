package pagination

type Pagination struct {
	limit uint32
	page  uint32
}

func New(limit uint32, page uint32) *Pagination {

	if limit == 0 {
		limit = 10
	}

	if page == 0 {
		limit = 1
	}

	return &Pagination{limit: limit, page: page}
}

func (p *Pagination) Limit() uint32 {
	return p.limit
}

func (p *Pagination) Page() uint32 {
	return p.page
}

func (p *Pagination) Offset() uint32 {
	return (p.page - 1) * p.limit
}
