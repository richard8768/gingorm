package service

import (
	"errors"
	"fmt"
	"gin_demo/internal/dto"
	"gin_demo/internal/model"
	"gin_demo/internal/util"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/spf13/cast"
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
		page = 0
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
	err = copier.Copy(&userAddressDtoList, &userAddressList)
	if err != nil {
		return nil, err
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
	fmt.Println(userAddressModel)
	result := s.db.Create(userAddressModel)
	if result.Error != nil {
		return nil, result.Error
	}
	fmt.Println(userAddressModel)

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
