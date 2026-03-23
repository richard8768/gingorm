package dto

type UserCreateRequest struct {
	Username string `json:"username" binding:"required,min=4,max=20"`
	Password string `json:"password" binding:"required,min=8"`
	Email    string `json:"email"    binding:"required,email"`
	Mobile   string `json:"mobile"   binding:"required,mobile"`
}

type UserResponse struct {
	ID         uint   `json:"id"`
	MemberName string `json:"username"`
	Email      string `json:"email"`
}

type UserLoginRequest struct {
	Username string `json:"username" binding:"required,min=4,max=20"`
	Password string `json:"password" binding:"required,min=8"`
}

type UserLoginResponse struct {
	Token string `json:"token"`
}
