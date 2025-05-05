package common

import "regexp"

const (
	DefaultPageSize    = 50
	DefaultCurrentPage = 1
)

// Paging represents sorter, paging
type Paging struct {
	Total       int64 `form:"total" json:"total"`
	PageSize    int   `form:"page_size" json:"page_size"`
	CurrentPage int   `form:"current_page" json:"current_page"`
}

func (s *Paging) Norm() {
	if s.PageSize == 0 {
		s.PageSize = 10
	}

	if s.CurrentPage == 0 {
		s.CurrentPage = 1
	}
}

// GetPageSize ...
func (s *Paging) GetPageSize() int {
	if s.PageSize == 0 {
		return 10
	}
	return s.PageSize
}

// GetCurrentPage ...
func (s *Paging) GetCurrentPage() int {
	if s.CurrentPage == 0 {
		return 1
	}
	return s.CurrentPage
}

// GetOffset ...
func (s *Paging) GetOffset() int {
	return s.GetPageSize() * (s.GetCurrentPage() - 1)
}

// TableListParams represents sorter, paging
type TableListParams struct {
	Sorter string `form:"sorter" json:"sorter"`
	Paging
}

func (t TableListParams) CreateCopy() TableListParams {
	return TableListParams{
		Sorter: t.Sorter,
		Paging: Paging{
			CurrentPage: t.Paging.CurrentPage,
			PageSize:    t.Paging.PageSize,
			Total:       t.Paging.Total,
		},
	}
}

func (t TableListParams) SetPaging() TableListParams {
	paging := Paging{
		Total:       t.Paging.Total,
		CurrentPage: DefaultCurrentPage,
		PageSize:    DefaultPageSize,
	}
	if t.CurrentPage > 0 {
		paging.CurrentPage = t.Paging.CurrentPage
	}

	if t.PageSize > 0 {
		paging.PageSize = t.Paging.PageSize
	}

	return TableListParams{
		Sorter: t.Sorter,
		Paging: paging,
	}
}

// SQLTagRegex const regex for sql tag in db struct
var SQLTagRegex = regexp.MustCompile("sql:\"([a-z].*?)\"")

// JSONTagRegex const regex for json tag in db struct
var JSONTagRegex = regexp.MustCompile("json:\"(.*?)\"")
