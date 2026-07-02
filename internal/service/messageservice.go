package service

import (
	"gin_demo/internal/dto"
	"gin_demo/internal/model"
	"gin_demo/internal/util"
	MRand "math/rand"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-module/carbon"
	"gorm.io/gorm"
)

type IMessageService interface {
	Create(db *gorm.DB) *MessageService
	SendEmail(context *gin.Context, req *dto.SendEmailRequest) (string, error)
	SendSms(context *gin.Context, req *dto.SendSmsRequest) (string, error)
}

type MessageService struct {
	db *gorm.DB
}

func (s *MessageService) Create(db *gorm.DB) *MessageService {
	return &MessageService{db: db}
}

func (s *MessageService) SendEmail(context *gin.Context, req *dto.SendEmailRequest) (string, error) {
	r := MRand.New(MRand.NewSource(time.Now().UnixNano()))
	code := util.GenRandStrings(r, 6, "number")
	captchaCodeModel := &model.CaptchaCode{
		CaptchaType:    2,
		CaptchaAccount: req.Email,
		CaptchaCode:    code,
		IsExpired:      0,
		IsUsed:         0,
		ExpiredTime:    carbon.Now().AddSeconds(60).Timestamp(),
	}
	result := s.db.Create(captchaCodeModel)
	if result.Error != nil {
		return "", result.Error
	}

	return "Email sent successfully", nil
}
func (s *MessageService) SendSms(context *gin.Context, req *dto.SendSmsRequest) (string, error) {
	r := MRand.New(MRand.NewSource(time.Now().UnixNano()))
	code := util.GenRandStrings(r, 6, "number")
	captchaCodeModel := &model.CaptchaCode{
		CaptchaType:    1,
		CaptchaAccount: req.Phone,
		CaptchaCode:    code,
		IsExpired:      0,
		IsUsed:         0,
		ExpiredTime:    carbon.Now().AddSeconds(60).Timestamp(),
	}
	result := s.db.Create(captchaCodeModel)
	if result.Error != nil {
		return "", result.Error
	}
	return "SMS sent successfully", nil
}
