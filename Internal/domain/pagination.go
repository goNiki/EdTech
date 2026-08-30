package domain

type Pagination struct {
	Page     int64 `json:"page"`
	PageSize int64 `json:"page_size"`
}

func (p *Pagination) Sanitize() {
	if p.Page < 1 {
		p.Page = 1
	}
	if p.PageSize < 1 {
		p.PageSize = 10
	} else if p.PageSize > 100 {
		p.PageSize = 100
	}
}

func (p Pagination) Offset() int64 {
	return (p.Page - 1) * p.PageSize
}
