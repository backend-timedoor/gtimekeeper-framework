package paginate

import (
	"math"

	"github.com/backend-timedoor/gtimekeeper/app"
	"gorm.io/gorm"
)

type PaginationRequest struct {
	Page  int `json:"page" query:"page"`
	Limit int `json:"Limit" query:"limit"`
}

type PaginationResponse struct {
	Data       any        `json:"data"`
	Pagination Pagination `json:"pagination"`
}

type Pagination struct {
	Page       int `json:"page"`
	TotalPages int `json:"total_pages"`
	TotalData  int `json:"total_data"`
	Limit      int `json:"limit"`
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

func Paginate(model interface{}, p *Pagination, req *PaginationRequest) func(db *gorm.DB) *gorm.DB {
	req.GetPaginationRequest()

	offset := (req.Page - 1) * req.Limit
	var totalData int64
	app.DB.Model(model).Count(&totalData)

	totalPages := math.Ceil(float64(totalData) / float64(req.Limit))
	p.TotalPages = int(totalPages)
	p.TotalData = int(totalData)
	p.Page = req.Page
	p.Limit = req.Limit

	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(offset).Limit(req.Limit)
	}
}
