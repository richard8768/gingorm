package service

import (
	"errors"
	"gin_demo/internal/config"
	"gin_demo/internal/dto"
	"gin_demo/internal/model"
	"gin_demo/internal/util"
	"mime/multipart"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type IUserSingleFileService interface {
	Create(db *gorm.DB) *UserSingleFileService
	Upload(context *gin.Context, file *multipart.FileHeader, isSaveFileToDb bool) (dto.UserSingleFileUploadResponse, error)
	Download(context *gin.Context, req *dto.UserSingleFileDownloadRequest) (string, string, error)
}

type UserSingleFileService struct {
	db *gorm.DB
}

func (s *UserSingleFileService) Create(db *gorm.DB) *UserSingleFileService {
	return &UserSingleFileService{db: db}
}

func (s *UserSingleFileService) Upload(context *gin.Context, file *multipart.FileHeader, isSaveFileToDb bool) (dto.UserSingleFileUploadResponse, error) {
	var rsp dto.UserSingleFileUploadResponse
	var userId uint
	var err error
	userId, err = util.GetUserId(context)
	if err != nil {
		return dto.UserSingleFileUploadResponse{FileName: "", FilePath: ""}, errors.New("意外的错误")
	}

	fileName := util.GenFileName()
	fileType := util.GetImageFileType(file.Header["Content-Type"][0])
	finalFileName := fileName + fileType
	saveFilePath, _ := config.GetLocalUploadPath()
	saveFilePath = saveFilePath + finalFileName
	context.SaveUploadedFile(file, saveFilePath)

	userIdUint64 := uint64(userId)
	if isSaveFileToDb == true {
		userUploadModel := &model.MemberUpload{MemberID: userIdUint64, FileName: finalFileName, SaveFilePath: saveFilePath}
		result := s.db.Create(userUploadModel)
		if result.Error != nil {
			return dto.UserSingleFileUploadResponse{FileName: "", FilePath: ""}, result.Error
		}
	}

	rsp = dto.UserSingleFileUploadResponse{FileName: finalFileName, FilePath: saveFilePath}
	return rsp, nil
}

func (s *UserSingleFileService) Download(context *gin.Context, req *dto.UserSingleFileDownloadRequest) (string, string, error) {
	userId, err := util.GetUserId(context)
	if err != nil {
		return "", "", errors.New("意外的错误")
	}

	fileId := req.ID

	var userUploadModel model.MemberUpload
	memberFileField := " member_id=? and id=?"
	result := s.db.Model(&userUploadModel).Where(memberFileField, uint64(userId), uint64(fileId)).Select("file_name", "save_file_path").Find(&userUploadModel)
	if result.Error != nil {
		return "", "", result.Error
	}

	return userUploadModel.FileName, userUploadModel.SaveFilePath, nil
}
