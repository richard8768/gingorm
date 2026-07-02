package service

import (
	"errors"
	"gin_demo/internal/dto"
	"gin_demo/internal/model"
	"gin_demo/internal/util"

	"github.com/gin-gonic/gin"
	"github.com/golang-module/carbon"
	"github.com/jinzhu/copier"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type IUserService interface {
	Create(db *gorm.DB) *UserService
	UserReg(context *gin.Context, req *dto.UserCreateRequest) (*dto.UserResponse, error)
	UserLogin(context *gin.Context, req *dto.UserLoginRequest) (*dto.UserLoginResponse, error)
	GetUserInfo(context *gin.Context) (*dto.UserResponse, error)
	GetUserId(context *gin.Context) (uint, error)
	UserLogout(context *gin.Context) (*dto.UserLoginResponse, error)
	UserBindLoginMobile(context *gin.Context, req *dto.UserBindLoginMobileRequest) (string, error)
	UserBindLoginEmail(context *gin.Context, req *dto.UserBindLoginEmailRequest) (string, error)
	UserCheckBindMobileEmail(context *gin.Context, req *dto.UserCheckBindMobileEmailRequest) (*dto.UserCheckBindMobileEmailResponse, error)
	UserChangePwd(context *gin.Context, req *dto.UserChangePwdRequest) (string, error)
	UserResetPwd(context *gin.Context, req *dto.UserResetPwdRequest) (string, error)
	UserUpdateProfile(context *gin.Context, req *dto.UserUpdateProfileRequest) (string, error)
}

type UserService struct {
	db *gorm.DB
}

var userLoginField = "member_name = ? or email = ?"

func (s *UserService) Create(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

func (s *UserService) UserReg(context *gin.Context, req *dto.UserCreateRequest) (*dto.UserResponse, error) {
	var existUser model.Member
	if err := s.db.Where(userLoginField, req.Username, req.Email).First(&existUser).Error; err == nil {
		return nil, errors.New("用户名或邮箱已存在")
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	userModel := &model.Member{
		MemberName: req.Username,
		Email:      req.Email,
		MemberPass: string(hashedPassword),
		Mobile:     req.Mobile,
	}

	tx := s.db.Begin()
	result := tx.Create(userModel)
	if result.Error != nil {
		tx.Rollback()
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		tx.Rollback()
		return nil, errors.New("用户创建失败")
	}
	userProfileModel := &model.MemberProfile{
		MemberID: userModel.ID,
	}
	result = tx.Create(userProfileModel)
	if result.Error != nil {
		tx.Rollback()
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		tx.Rollback()
		return nil, errors.New("用户创建失败.")
	}
	userAccountModel := &model.MemberAccount{
		MemberID: int64(userModel.ID),
	}
	result = tx.Create(userAccountModel)
	if result.Error != nil {
		tx.Rollback()
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		tx.Rollback()
		return nil, errors.New("用户创建失败..")
	}
	tx.Commit()

	response := &dto.UserResponse{}
	err := copier.Copy(response, userModel)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (s *UserService) UserLogin(context *gin.Context, req *dto.UserLoginRequest) (*dto.UserLoginResponse, error) {
	var user model.Member
	if err := s.db.Where(userLoginField, req.Username, req.Username).First(&user).Error; err != nil {
		return nil, errors.New("错误的用户信息")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.MemberPass), []byte(req.Password)); err != nil {
		return nil, errors.New("用户名或密码错误")
	}

	remoteIP := context.RemoteIP()
	s.db.Model(&user).Where(userLoginField, req.Username, req.Username).
		Updates(model.Member{LoginTime: carbon.Now().Timestamp(), LoginIP: remoteIP}).
		UpdateColumn("LoginTimes", gorm.Expr("LoginTimes + ?", 1))

	token, err := util.GenToken(user, "user", "login")
	if err != nil {
		return nil, err
	}
	response := &dto.UserLoginResponse{
		Token: token,
	}

	return response, nil
}

func (s *UserService) GetUserId(context *gin.Context) (uint, error) {
	return util.GetUserId(context)
}

func (s *UserService) UserLogout(context *gin.Context) (*dto.UserLoginResponse, error) {
	userId, err := s.GetUserId(context)
	if err != nil {
		return nil, errors.New("意外的错误")
	}
	var user model.Member
	if err := s.db.Where("id = ?", userId).First(&user).Error; err != nil {
		return nil, err
	}
	token, err := util.GenToken(user, "user", "logout")
	if err != nil {
		return nil, err
	}
	response := &dto.UserLoginResponse{
		Token: token,
	}

	return response, nil
}

func (s *UserService) GetUserInfo(context *gin.Context) (*dto.UserResponse, error) {
	userId, err := s.GetUserId(context)
	if err != nil {
		return nil, errors.New("意外的错误")
	}
	var user model.Member
	if err := s.db.Where("id = ?", userId).First(&user).Error; err != nil {
		return nil, err
	}

	response := &dto.UserResponse{}
	err = copier.Copy(response, user)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (s *UserService) UserBindLoginMobile(context *gin.Context, req *dto.UserBindLoginMobileRequest) (string, error) {
	userId, err := s.GetUserId(context)
	if err != nil {
		return "", errors.New("意外的错误")
	}

	var captchaCodeModel model.CaptchaCode
	s.UpdateCaptchaCodeExpired(&captchaCodeModel)

	captchaCheckField := "captcha_type = 1 and captcha_account = ? and captcha_code = ?"
	if result := s.db.Where(captchaCheckField, req.Mobile, req.CaptchaCode).First(&captchaCodeModel); result.Error != nil {
		return "", errors.New("验证码错误")
	}
	if captchaCodeModel.IsExpired == 1 {
		return "", errors.New("验证码已过期")
	}
	if captchaCodeModel.IsUsed == 1 {
		return "", errors.New("验证码已使用")
	}
	s.db.Model(&captchaCodeModel).Where(captchaCheckField, req.Mobile, req.CaptchaCode).
		Updates(model.CaptchaCode{IsUsed: 1})

	var user model.Member
	s.db.Model(&user).Where("id = ?", userId).
		Updates(model.Member{Mobile: req.Mobile})

	return "操作成功", nil
}

func (s *UserService) UserBindLoginEmail(context *gin.Context, req *dto.UserBindLoginEmailRequest) (string, error) {
	userId, err := s.GetUserId(context)
	if err != nil {
		return "", errors.New("意外的错误")
	}

	var captchaCodeModel model.CaptchaCode
	s.UpdateCaptchaCodeExpired(&captchaCodeModel)

	captchaCheckField := "captcha_type = 2 and captcha_account = ? and captcha_code = ?"
	if result := s.db.Where(captchaCheckField, req.Email, req.CaptchaCode).First(&captchaCodeModel); result.Error != nil {
		return "", errors.New("验证码错误")
	}
	if captchaCodeModel.IsExpired == 1 {
		return "", errors.New("验证码已过期")
	}
	if captchaCodeModel.IsUsed == 1 {
		return "", errors.New("验证码已使用")
	}
	s.db.Model(&captchaCodeModel).Where(captchaCheckField, req.Email, req.CaptchaCode).
		Updates(model.CaptchaCode{IsUsed: 1})

	var user model.Member
	s.db.Model(&user).Where("id = ?", userId).
		Updates(model.Member{Email: req.Email})

	return "操作成功", nil
}

func (s *UserService) UserCheckBindMobileEmail(context *gin.Context, req *dto.UserCheckBindMobileEmailRequest) (*dto.UserCheckBindMobileEmailResponse, error) {
	captchaType := 0
	captchaAccount := ""
	accountType := ""
	if req.Mobile != "" {
		captchaType = 1
		captchaAccount = req.Mobile
		accountType = "mobile"
	} else if req.Email != "" {
		captchaType = 2
		captchaAccount = req.Email
		accountType = "email"
	}
	if captchaType == 0 {
		return &dto.UserCheckBindMobileEmailResponse{}, errors.New("意外的错误")
	}

	var captchaCodeModel model.CaptchaCode
	s.UpdateCaptchaCodeExpired(&captchaCodeModel)

	captchaCheckField := "captcha_type = ? and captcha_account = ? and captcha_code = ?"
	if result := s.db.Where(captchaCheckField, captchaType, captchaAccount, req.CaptchaCode).First(&captchaCodeModel); result.Error != nil {
		return &dto.UserCheckBindMobileEmailResponse{}, errors.New("验证码错误")
	}
	if captchaCodeModel.IsExpired == 1 {
		return &dto.UserCheckBindMobileEmailResponse{}, errors.New("验证码已过期")
	}
	if captchaCodeModel.IsUsed == 1 {
		return &dto.UserCheckBindMobileEmailResponse{}, errors.New("验证码已使用")
	}
	s.db.Model(&captchaCodeModel).Where(captchaCheckField, captchaType, captchaAccount, req.CaptchaCode).
		Updates(model.CaptchaCode{IsUsed: 1})

	return &dto.UserCheckBindMobileEmailResponse{
		Account:     captchaAccount,
		AccountType: accountType,
	}, nil
}

func (s *UserService) UserChangePwd(context *gin.Context, req *dto.UserChangePwdRequest) (string, error) {
	userId, err := s.GetUserId(context)
	if err != nil {
		return "", errors.New("意外的错误")
	}
	var user model.Member
	if err = s.db.Where("id=?", userId).First(&user).Error; err != nil {
		return "", errors.New("错误的用户信息")
	}
	if err = bcrypt.CompareHashAndPassword([]byte(user.MemberPass), []byte(req.OldPassword)); err != nil {
		return "", errors.New("用户名或密码错误")
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	s.db.Model(&user).Where("id = ?", userId).
		Updates(model.Member{MemberPass: string(hashedPassword)})

	return "操作成功", nil
}

func (s *UserService) UserResetPwd(context *gin.Context, req *dto.UserResetPwdRequest) (string, error) {
	loginAccountField := ""
	loginAccount := ""
	if req.Email != "" {
		loginAccountField = "email = ?"
		loginAccount = req.Email
	} else if req.Mobile != "" {
		loginAccountField = "mobile = ?"
		loginAccount = req.Mobile
	}
	if loginAccountField == "" {
		return "", errors.New("错误的登录账号")
	}

	var user model.Member
	if err := s.db.Where(loginAccountField, loginAccount).First(&user).Error; err != nil {
		return "", errors.New("错误的用户信息")
	}
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	s.db.Model(&user).Where(loginAccountField, loginAccount).
		Updates(model.Member{MemberPass: string(hashedPassword)})

	return "操作成功", nil
}

func (s *UserService) UserUpdateProfile(context *gin.Context, req *dto.UserUpdateProfileRequest) (string, error) {
	userId, err := s.GetUserId(context)
	if err != nil {
		return "", errors.New("意外的错误")
	}

	var userProfile model.MemberProfile
	sex := 0
	if req.Sex == "male" {
		sex = 1
	} else if req.Sex == "female" {
		sex = 2
	}
	s.db.Model(&userProfile).Where("member_id = ?", userId).
		Updates(model.MemberProfile{
			NickName:   req.NickName,
			TrueName:   req.TrueName,
			Sex:        int64(sex),
			Mobile:     req.Mobile,
			ProvinceID: uint64(req.ProvinceID),
			CityID:     uint64(req.CityID),
			AreaID:     uint64(req.AreaID),
			Address:    req.Address,
			HeadImg:    req.HeadImg,
			Age:        int64(req.Age),
		})

	return "操作成功", nil
}

func (s *UserService) UpdateCaptchaCodeExpired(captchaCodeModel *model.CaptchaCode) {
	s.db.Model(&captchaCodeModel).Where("is_used=0 and is_expired=0 and expired_time < ?", carbon.Now().Timestamp()).
		Updates(model.CaptchaCode{IsExpired: 1})
}
