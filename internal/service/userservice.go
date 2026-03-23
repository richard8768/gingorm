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
	UserReg(context *gin.Context, d *dto.UserCreateRequest) (*dto.UserResponse, error)
	UserLogin(context *gin.Context, req *dto.UserLoginRequest) (*dto.UserLoginResponse, error)
	GetUserInfo(context *gin.Context) (*dto.UserResponse, error)
	GetUserId(context *gin.Context) (uint, error)
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

	token, err := util.GenToken(user, "user")
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
