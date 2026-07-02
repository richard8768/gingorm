package dto

type SendEmailRequest struct {
	Email string `json:"email"  binding:"required,email"`
}

type SendSmsRequest struct {
	Phone string `json:"phone" binding:"required,mobile"`
}
