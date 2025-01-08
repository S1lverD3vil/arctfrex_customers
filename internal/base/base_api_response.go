package base

type ApiResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
	Time    string `json:"time"`
}
