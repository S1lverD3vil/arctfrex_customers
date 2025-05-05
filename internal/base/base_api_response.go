package base

import "arctfrex-customers/internal/common"

type ApiResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
	Time    string `json:"time"`
}

type ApiPaginatedResponse struct {
	Message    string                 `json:"message"`
	Data       interface{}            `json:"data"`
	Pagination common.TableListParams `json:"pagination"`
	Time       string                 `json:"time"`
}
