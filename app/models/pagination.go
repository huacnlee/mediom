package models

import (
	"github.com/jinzhu/gorm"
	"math"
)

type Pagination struct {
	Query        *gorm.DB
	TotalEntites int
	PerPage      int
	Path         string
	Page         int
	TotalPages   int
}

func (p *Pagination) Paginate(page int) *gorm.DB {
	p.Page = page
	p.Query.Count(&p.TotalEntites)
	if p.TotalEntites == 0 {
		return p.Query
	}

	p.TotalPages = int(math.Ceil(float64(p.TotalEntites) / float64(p.PerPage)))

	if !(p.Page > 0 && p.Page <= p.TotalPages) {
		p.Page = 1
	}

	query := p.Query.Offset((p.Page - 1) * p.PerPage).Limit(p.PerPage)

	return query

}
