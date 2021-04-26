package dto

type ListData struct {
	Total int64       `json:"total"`
	Data  interface{} `json:"data"`
}
