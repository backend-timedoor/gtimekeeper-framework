package paginate

import (
	"math"

	"github.com/backend-timedoor/gtimekeeper-framework/app"
	"gorm.io/gorm"
)

type PaginationRequest struct {
	Page  int `json:"page" query:"page"`
	Limit int `json:"Limit" query:"limit"`
}

type PaginationResponse struct {
	Data       any        `json:"data"`
	Pagination Pagination `json:"meta"`
}

type Pagination struct {
	CurrentPage int `json:"current_page"`
	LastPage    int `json:"last_page"`
	Total       int `json:"total"`
	PerPage     int `json:"per_page"`
}

func (r *PaginationRequest) GetPaginationRequest() error {
	page, limit := 1, 10

	if r.Page <= 0 {
		r.Page = page
	}

	if r.Limit <= 0 || r.Limit >= 50 {
		r.Limit = limit
	}

	return nil
}

func Paginate(model interface{}, p *Pagination, req *PaginationRequest, scope ...func(*gorm.DB) *gorm.DB) func(db *gorm.DB) *gorm.DB {
	req.GetPaginationRequest()

	offset := (req.Page - 1) * req.Limit
	var totalData int64
	app.DB.Model(model).Scopes(scope...).Count(&totalData)

	totalPages := math.Ceil(float64(totalData) / float64(req.Limit))
	p.LastPage = int(totalPages)
	p.Total = int(totalData)
	p.CurrentPage = req.Page
	p.PerPage = req.Limit

	return func(db *gorm.DB) *gorm.DB {
		return db.Scopes(scope...).Offset(offset).Limit(req.Limit)
	}
}
