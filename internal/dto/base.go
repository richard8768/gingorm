package dto

type HttpResponse struct {
	Code    int `json:"code"`
	Message any `json:"message"`
	Data    any `json:"data"`
}
