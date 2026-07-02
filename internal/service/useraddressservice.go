package service

import (
	"errors"
	"fmt"
	"gin_demo/internal/config"
	"gin_demo/internal/dto"
	"gin_demo/internal/model"
	"gin_demo/internal/util"
	"mime/multipart"
	"os"
	"slices"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/spf13/cast"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

type IUserAddressService interface {
	Create(db *gorm.DB) *UserAddressService
	AddressList(context *gin.Context, req *dto.UserAddressSearchRequest) (*dto.UserAddressListResponse, error)
	AddressInfo(context *gin.Context, req *dto.UserAddressGetRequest) (*dto.UserAddressResponse, error)
	AddAddress(context *gin.Context, req *dto.UserAddressCreateRequest) (*dto.UserAddressResponse, error)
	UpdateAddress(context *gin.Context, req *dto.UserAddressUpdateRequest) (*dto.UserAddressResponse, error)
	DeleteAddress(context *gin.Context, req *dto.UserAddressDeleteRequest) (*dto.UserAddressDeleteResponse, error)
	SetDefaultAddress(context *gin.Context, req *dto.UserAddressRequest) (*dto.UserAddressResponse, error)
	Upload(context *gin.Context, file *multipart.FileHeader) (string, error)
	Download(context *gin.Context) ([][]string, []string, error)
}

type UserAddressService struct {
	db *gorm.DB
}

func (s *UserAddressService) Create(db *gorm.DB) *UserAddressService {
	return &UserAddressService{db: db}
}
func (s *UserAddressService) AddressList(context *gin.Context, req *dto.UserAddressSearchRequest) (*dto.UserAddressListResponse, error) {
	userId, err := util.GetUserId(context)
	if err != nil {
		return nil, errors.New("意外的错误")
	}
	var page = req.Page
	var pageSize = req.PageSize
	var keyword = req.Keyword
	if page <= 0 {
		page = 1
	}
	page--
	if pageSize <= 0 {
		pageSize = 20
	}

	var userAddress model.MemberAddress
	var userAddressList []model.MemberAddress
	var query = s.db.Model(&userAddress).Where("member_id = ?", userId)
	if keyword != "" {
		query.Where(
			s.db.Where("address like ?", ""+keyword+"%").Or("consignee_name like ?", ""+keyword+"%"),
		)
	}

	result := query.Limit(pageSize).Offset(page).Find(&userAddressList)
	if result.Error != nil {
		return nil, result.Error
	}

	total := result.RowsAffected

	var userAddressDtoList []dto.UserAddressResponse
	if len(userAddressList) > 0 {
		for _, singleUserAddress := range userAddressList {
			var typeSingleUserAddress dto.UserAddressResponse
			_ = copier.Copy(&typeSingleUserAddress, singleUserAddress)
			userAddressDtoList = append(userAddressDtoList, typeSingleUserAddress)
		}
	}

	response := &dto.UserAddressListResponse{
		Total:           total,
		Page:            page + 1,
		PageSize:        pageSize,
		UserAddressList: userAddressDtoList,
	}

	return response, nil
}

func (s *UserAddressService) AddressInfo(context *gin.Context, req *dto.UserAddressGetRequest) (*dto.UserAddressResponse, error) {
	userId, err := util.GetUserId(context)
	if err != nil {
		return nil, errors.New("意外的错误")
	}

	userAddressId := req.ID
	var userAddressModel model.MemberAddress
	result := s.db.Where("id = ?", userAddressId).Where("member_id = ?", userId).First(&userAddressModel)
	if result.Error != nil {
		return nil, result.Error
	}

	response := &dto.UserAddressResponse{}
	err = copier.Copy(response, userAddressModel)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (s *UserAddressService) AddAddress(context *gin.Context, req *dto.UserAddressCreateRequest) (*dto.UserAddressResponse, error) {
	userId, err := util.GetUserId(context)
	if err != nil {
		return nil, errors.New("意外的错误")
	}

	userAddressModel := &model.MemberAddress{MemberID: uint64(userId), Address: req.Address, Tel: req.Tel,
		ConsigneeName: req.ConsigneeName, Post: req.Post, ProvinceID: cast.ToInt64(req.ProvinceID), CityID: cast.ToInt64(req.CityID), AreaID: cast.ToInt64(req.AreaID)}
	result := s.db.Create(userAddressModel)
	if result.Error != nil {
		return nil, result.Error
	}

	response := &dto.UserAddressResponse{}
	err = copier.Copy(response, userAddressModel)
	if err != nil {
		return nil, err
	}

	return response, nil

}

func (s *UserAddressService) UpdateAddress(context *gin.Context, req *dto.UserAddressUpdateRequest) (*dto.UserAddressResponse, error) {
	userId, err := util.GetUserId(context)
	if err != nil {
		return nil, errors.New("意外的错误")
	}
	userAddressId := req.ID
	var userAddressModel model.MemberAddress
	result := s.db.Where("id = ?", userAddressId).Where("member_id = ?", userId).First(&userAddressModel)
	if result.Error != nil {
		return nil, result.Error
	}

	result = s.db.Model(&userAddressModel).
		Where("id = ?", userAddressId).Where("member_id = ?", userId).Updates(req)
	if result.Error != nil {
		return nil, result.Error
	}

	response := &dto.UserAddressResponse{}
	err = copier.Copy(response, userAddressModel)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (s *UserAddressService) DeleteAddress(context *gin.Context, req *dto.UserAddressDeleteRequest) (*dto.UserAddressDeleteResponse, error) {
	userId, err := util.GetUserId(context)
	if err != nil {
		return nil, errors.New("意外的错误")
	}

	userAddressId := req.ID
	var userAddressModel model.MemberAddress
	result := s.db.Where("id = ?", userAddressId).Where("member_id = ?", userId).First(&userAddressModel)
	if result.Error != nil {
		return nil, result.Error
	}

	result = s.db.Delete(&userAddressModel)
	if result.Error != nil {
		return nil, result.Error
	}

	////if u want restore the data use the follow code
	//result = s.db.Unscoped().Model(&model.MemberAddress{}).Where("id = ?", userAddressId).Where("member_id = ?", userId).Update("deleted_at", 0)
	//if result.Error != nil {
	//	return nil, result.Error
	//}

	////if u want force delete the data use the follow code
	//result = s.db.Unscoped().Where("id = ?", userAddressId).Where("member_id = ?", userId).Delete(&model.MemberAddress{})
	//if result.Error != nil {
	//	return nil, result.Error
	//}

	response := &dto.UserAddressDeleteResponse{
		ID:           uint(userAddressModel.ID),
		DeleteStatus: result.RowsAffected,
	}
	return response, nil
}
func (s *UserAddressService) SetDefaultAddress(context *gin.Context, req *dto.UserAddressRequest) (*dto.UserAddressResponse, error) {
	userId, err := util.GetUserId(context)
	if err != nil {
		return nil, errors.New("意外的错误")
	}

	userAddressId := req.ID
	var userAddressModel model.MemberAddress
	result := s.db.Where("id = ?", userAddressId).Where("member_id = ?", userId).First(&userAddressModel)
	if result.Error != nil {
		return nil, result.Error
	}
	tx := s.db.Begin()
	result = s.db.Model(&model.MemberAddress{}).Where("member_id = ?", userId).Update("is_default", 0)
	if result.Error != nil {
		tx.Rollback()
		return nil, result.Error
	}
	result = s.db.Model(&userAddressModel).
		Where("id = ?", userAddressId).Where("member_id = ?", userId).Updates(model.MemberAddress{IsDefault: 1})
	if result.Error != nil {
		tx.Rollback()
		return nil, result.Error
	}
	tx.Commit()

	response := &dto.UserAddressResponse{}
	err = copier.Copy(response, userAddressModel)
	if err != nil {
		return nil, err
	}

	return response, nil
}
func (s *UserAddressService) Upload(context *gin.Context, file *multipart.FileHeader) (string, error) {
	userId, err := util.GetUserId(context)
	if err != nil {
		return "", errors.New("意外的错误")
	}

	titleList := []string{"地址", "电话", "联系人", "邮编"}
	columnList := []string{"address", "tel", "consignee_name", "post"}
	var insertColumnList []string

	fileName := util.GenFileName()
	fileType := util.GetImageFileType(file.Header["Content-Type"][0])
	finalFileName := fileName + fileType
	saveFilePath, _ := config.GetLocalUploadPath()
	saveFilePath = saveFilePath + finalFileName
	context.SaveUploadedFile(file, saveFilePath)

	f, err := excelize.OpenFile(saveFilePath)
	if err != nil {
		return "", err
	}
	defer f.Close()
	index := -1
	//读取表头和数据
	rows, err := f.GetRows("Sheet1")
	if err != nil {
		return "", err
	}
	var userAddressList []model.MemberAddress
	var userAddress model.MemberAddress
	for rowIndex, row := range rows {
		userAddress = model.MemberAddress{}
		for colIdx, colCell := range row {
			if rowIndex == 0 {
				if !slices.Contains(titleList, colCell) {
					return "", errors.New("表头不匹配")
				}
				index = slices.Index(titleList, colCell)
				if index != -1 {
					insertColumnList = append(insertColumnList, columnList[index])
				}
			} else {
				if insertColumnList[colIdx] == "address" {
					userAddress.Address = colCell
				} else if insertColumnList[colIdx] == "tel" {
					userAddress.Tel = colCell
				} else if insertColumnList[colIdx] == "consignee_name" {
					userAddress.ConsigneeName = colCell
				} else if insertColumnList[colIdx] == "post" {
					userAddress.Post = colCell
				}
				userAddress.ProvinceID = 123456
				userAddress.CityID = 147852
				userAddress.AreaID = 159753
				userAddress.MemberID = uint64(userId)
			}
		}
		if rowIndex != 0 {
			userAddressList = append(userAddressList, userAddress)
		}
	}
	//写入数据至db
	result := s.db.CreateInBatches(userAddressList, 1000)
	if result.Error != nil {
		return "", result.Error
	}

	_ = os.Remove(saveFilePath)

	return "ok", nil
}

func (s *UserAddressService) Download(context *gin.Context) ([][]string, []string, error) {
	userId, err := util.GetUserId(context)
	if err != nil {
		return nil, nil, errors.New("意外的错误")
	}
	var userAddress model.MemberAddress
	var userAddressList []model.MemberAddress
	var query = s.db.Model(&userAddress).Select("id", "address", "tel", "consignee_name", "post").Where("member_id = ?", userId)

	result := query.Find(&userAddressList)
	if result.Error != nil {
		return nil, nil, result.Error
	}

	var userAddressDownloadDtoList [][]string
	if len(userAddressList) > 0 {
		for _, singleUserAddress := range userAddressList {
			var typeSingleUserAddress []string
			typeSingleUserAddress = append(typeSingleUserAddress, fmt.Sprintf("%d", singleUserAddress.ID))
			typeSingleUserAddress = append(typeSingleUserAddress, singleUserAddress.Address, singleUserAddress.Tel, singleUserAddress.ConsigneeName, singleUserAddress.Post)
			userAddressDownloadDtoList = append(userAddressDownloadDtoList, typeSingleUserAddress)
		}
	}

	titleList := []string{"ID", "地址", "电话", "联系人", "邮编"}

	return userAddressDownloadDtoList, titleList, nil
}
