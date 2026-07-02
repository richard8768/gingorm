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
type UserBindLoginMobileRequest struct {
	Mobile      string `json:"mobile" binding:"required,mobile"`
	CaptchaCode string `json:"captcha_code" binding:"required"`
}

type UserBindLoginEmailRequest struct {
	Email       string `json:"email" binding:"required,email"`
	CaptchaCode string `json:"captcha_code" binding:"required"`
}

type UserCheckBindMobileEmailRequest struct {
	Mobile      string `json:"mobile" binding:"omitempty,mobile,required_without=Email"`
	Email       string `json:"email" binding:"omitempty,email,required_without=Mobile"`
	CaptchaCode string `json:"captcha_code" binding:"required"`
}

type UserCheckBindMobileEmailResponse struct {
	AccountType string `json:"account_type"`
	Account     string `json:"account"`
}

type UserChangePwdRequest struct {
	OldPassword     string `json:"old_password" binding:"required,min=8"`
	NewPassword     string `json:"new_password" binding:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" binding:"required,min=8,eqfield=NewPassword"`
}

type UserResetPwdRequest struct {
	Mobile          string `json:"mobile" binding:"omitempty,mobile,required_without=Email"`
	Email           string `json:"email" binding:"omitempty,email,required_without=Mobile"`
	NewPassword     string `json:"new_password" binding:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" binding:"required,min=8,eqfield=NewPassword"`
}

type UserUpdateProfileRequest struct {
	NickName   string `json:"nick_name" binding:"required,min=2,max=100"`
	TrueName   string `json:"true_name" binding:"required,min=2,max=100"`
	Sex        string `json:"sex" binding:"required,oneof=male female other"`
	Mobile     string `json:"mobile" binding:"required,min=11,max=11,mobile"`
	ProvinceID int    `json:"province_id" binding:"required,number"`
	CityID     int    `json:"city_id" binding:"required,number"`
	AreaID     int    `json:"area_id" binding:"required,number"`
	Address    string `json:"address" binding:"required,min=10,max=200"`
	HeadImg    string `json:"head_img" binding:"required"`
	Age        int    `json:"age" binding:"required,number,min=18,max=130"`
}
