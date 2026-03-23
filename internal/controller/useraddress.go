package controller

import (
	"gin_demo/internal/dto"
	"gin_demo/internal/service"
	"gin_demo/internal/util"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type UserAddressHandler struct {
	IUserAddressService service.IUserAddressService
}

// get user address list
// @Summary AddressList
// @Schemes
// @Description AddressList
// @Tags AddressList
// @Accept json
// @Produce json
// @Param page query int false "page"
// @Param page_size query int false "page_size"
// @Param keyword query string false "keyword"
// @Success 200 {object} dto.UserAddressListResponse
// @Router /useraddress/index [get]
func (h *UserAddressHandler) AddressList(context *gin.Context) {
	var req dto.UserAddressSearchRequest
	if err := context.ShouldBindQuery(&req); err != nil {
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			util.HttpResponse(context, 500, err.Error(), nil)
			return
		}
		util.HttpResponse(context, 500, util.RemoveTopStruct(errs.Translate(util.Trans)), nil)
		return
	}

	addressListResponse, err := h.IUserAddressService.AddressList(context, &req)
	if err != nil {
		util.HttpResponse(context, 500, err.Error(), nil)
		return
	}

	util.HttpResponse(context, 200, "ok", addressListResponse)
	return
}

// get user address info
// @Summary AddressInfo
// @Schemes
// @Description AddressInfo
// @Tags AddressInfo
// @Accept json
// @Produce json
// @Param id query int true "id"
// @Success 200 {object} dto.UserAddressResponse
// @Router /useraddress/info [get]
func (h *UserAddressHandler) AddressInfo(context *gin.Context) {
	var req dto.UserAddressGetRequest
	if err := context.ShouldBindQuery(&req); err != nil {
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			util.HttpResponse(context, 500, err.Error(), nil)
			return
		}
		util.HttpResponse(context, 500, util.RemoveTopStruct(errs.Translate(util.Trans)), nil)
		return
	}

	addressInfoResponse, err := h.IUserAddressService.AddressInfo(context, &req)
	if err != nil {
		util.HttpResponse(context, 500, err.Error(), nil)
		return
	}

	util.HttpResponse(context, 200, "ok", addressInfoResponse)
}

// add user address
// @Summary AddAddress
// @Schemes
// @Description AddAddress
// @Tags AddAddress
// @Accept json
// @Produce json
// @Param body body dto.UserAddressCreateRequest true "请求body"
// @Success 200 {object} dto.UserAddressResponse
// @Router /useraddress/add [post]
func (h *UserAddressHandler) AddAddress(context *gin.Context) {
	var req dto.UserAddressCreateRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			util.HttpResponse(context, 500, err.Error(), nil)
			return
		}
		util.HttpResponse(context, 500, util.RemoveTopStruct(errs.Translate(util.Trans)), nil)
		return
	}

	addAddressResponse, err := h.IUserAddressService.AddAddress(context, &req)
	if err != nil {
		util.HttpResponse(context, 500, err.Error(), nil)
		return
	}

	util.HttpResponse(context, 200, "ok", addAddressResponse)
	return
}

// edit user address
// @Summary UpdateAddress
// @Schemes
// @Description UpdateAddress
// @Tags UpdateAddress
// @Accept json
// @Produce json
// @Param body body dto.UserAddressUpdateRequest true "请求body"
// @Success 200 {object} dto.UserAddressResponse
// @Router /useraddress/edit [post]
func (h *UserAddressHandler) UpdateAddress(context *gin.Context) {
	var req dto.UserAddressUpdateRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			util.HttpResponse(context, 500, err.Error(), nil)
			return
		}
		util.HttpResponse(context, 500, util.RemoveTopStruct(errs.Translate(util.Trans)), nil)
		return
	}

	updateAddressResponse, err := h.IUserAddressService.UpdateAddress(context, &req)
	if err != nil {
		util.HttpResponse(context, 500, err.Error(), nil)
		return
	}

	util.HttpResponse(context, 200, "ok", updateAddressResponse)
	return
}

// delete user address
// @Summary DeleteAddress
// @Schemes
// @Description DeleteAddress
// @Tags DeleteAddress
// @Accept json
// @Produce json
// @Param body body dto.UserAddressDeleteRequest true "请求body"
// @Success 200 {object} dto.UserAddressDeleteResponse
// @Router /useraddress/del [post]
func (h *UserAddressHandler) DeleteAddress(context *gin.Context) {
	var req dto.UserAddressDeleteRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			util.HttpResponse(context, 500, err.Error(), nil)
			return
		}
		util.HttpResponse(context, 500, util.RemoveTopStruct(errs.Translate(util.Trans)), nil)
		return
	}

	deleteAddressResponse, err := h.IUserAddressService.DeleteAddress(context, &req)
	if err != nil {
		util.HttpResponse(context, 500, err.Error(), nil)
		return
	}

	util.HttpResponse(context, 200, "ok", deleteAddressResponse)
	return
}

// set user default address
// @Summary SetDefaultAddress
// @Schemes
// @Description SetDefaultAddress
// @Tags SetDefaultAddress
// @Accept json
// @Produce json
// @Param body body dto.UserAddressRequest true "请求body"
// @Success 200 {object} dto.UserAddressResponse
// @Router /useraddress/setdefault [post]
func (h *UserAddressHandler) SetDefaultAddress(context *gin.Context) {
	var req dto.UserAddressRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			util.HttpResponse(context, 500, err.Error(), nil)
			return
		}
		util.HttpResponse(context, 500, util.RemoveTopStruct(errs.Translate(util.Trans)), nil)
		return
	}

	setDefaultAddressResponse, err := h.IUserAddressService.SetDefaultAddress(context, &req)
	if err != nil {
		util.HttpResponse(context, 500, err.Error(), nil)
		return
	}

	util.HttpResponse(context, 200, "ok", setDefaultAddressResponse)
	return
}
