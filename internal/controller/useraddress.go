package controller

import (
	"gin_demo/internal/dto"
	"gin_demo/internal/service"
	"gin_demo/internal/util"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type UserAddressHandler struct {
	IUserAddressService service.IUserAddressService
}

// get user address list
// @Summary AddressList
// @Schemes http https
// @Description AddressList
// @Tags UserAddress
// @Accept json
// @Produce json
// @Param page query int false "page"
// @Param page_size query int false "page_size"
// @Param keyword query string false "keyword"
// @Success 200 {object} dto.UserAddressListResponse
// @Router /useraddress/index [get]
func (h *UserAddressHandler) AddressList(context *gin.Context) {
	var req dto.UserAddressSearchRequest
	if err := util.CheckReqBindQuery(context, &req); err != nil {
		util.HttpResponse(context, 500, err, nil)
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
// @Schemes http https
// @Description AddressInfo
// @Tags UserAddress
// @Accept json
// @Produce json
// @Param id query int true "id"
// @Success 200 {object} dto.UserAddressResponse
// @Router /useraddress/info [get]
func (h *UserAddressHandler) AddressInfo(context *gin.Context) {
	var req dto.UserAddressGetRequest
	if err := util.CheckReqBindQuery(context, &req); err != nil {
		util.HttpResponse(context, 500, err, nil)
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
// @Schemes http https
// @Description AddAddress
// @Tags UserAddress
// @Accept json
// @Produce json
// @Param body body dto.UserAddressCreateRequest true "请求body"
// @Success 200 {object} dto.UserAddressResponse
// @Router /useraddress/add [post]
func (h *UserAddressHandler) AddAddress(context *gin.Context) {
	var req dto.UserAddressCreateRequest
	if err := util.CheckReqBindJson(context, &req); err != nil {
		util.HttpResponse(context, 500, err, nil)
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
// @Schemes http https
// @Description UpdateAddress
// @Tags UserAddress
// @Accept json
// @Produce json
// @Param body body dto.UserAddressUpdateRequest true "请求body"
// @Success 200 {object} dto.UserAddressResponse
// @Router /useraddress/edit [post]
func (h *UserAddressHandler) UpdateAddress(context *gin.Context) {
	var req dto.UserAddressUpdateRequest
	if err := util.CheckReqBindJson(context, &req); err != nil {
		util.HttpResponse(context, 500, err, nil)
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
// @Schemes http https
// @Description DeleteAddress
// @Tags UserAddress
// @Accept json
// @Produce json
// @Param body body dto.UserAddressDeleteRequest true "请求body"
// @Success 200 {object} dto.UserAddressDeleteResponse
// @Router /useraddress/del [post]
func (h *UserAddressHandler) DeleteAddress(context *gin.Context) {
	var req dto.UserAddressDeleteRequest
	if err := util.CheckReqBindJson(context, &req); err != nil {
		util.HttpResponse(context, 500, err, nil)
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
// @Schemes http https
// @Description SetDefaultAddress
// @Tags UserAddress
// @Accept json
// @Produce json
// @Param body body dto.UserAddressRequest true "请求body"
// @Success 200 {object} dto.UserAddressResponse
// @Router /useraddress/setdefault [post]
func (h *UserAddressHandler) SetDefaultAddress(context *gin.Context) {
	var req dto.UserAddressRequest
	if err := util.CheckReqBindJson(context, &req); err != nil {
		util.HttpResponse(context, 500, err, nil)
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

// upload address excel file
// @Summary upload address excel file
// @Schemes http https
// @Description upload address excel file
// @Tags UserAddress
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "文件"
// @Success 200 {string} ok
// @Router /useraddress/upload [post]
func (h *UserAddressHandler) Upload(context *gin.Context) {
	validate, ok := binding.Validator.Engine().(*validator.Validate)
	if ok {
		validate.RegisterStructValidation(util.FileUploadValidation, dto.UserAddressUploadRequest{})
	}

	var req dto.UserNormalFileUploadRequest
	if err := context.ShouldBind(&req); err != nil {
		util.HttpResponse(context, 500, err.Error(), nil)
		return
	}

	//file, err := context.FormFile("file")
	//if err != nil {
	//	util.HttpResponse(context, 500, err.Error(), nil)
	//	return
	//}

	form, err := context.MultipartForm()
	if err != nil {
		util.HttpResponse(context, 500, err.Error(), nil)
		return
	}
	files := form.File["file"]
	if files == nil {
		util.HttpResponse(context, 500, "file is empty", nil)
		return
	}

	file := files[0]
	_, err = h.IUserAddressService.Upload(context, file)
	if err != nil {
		util.HttpResponse(context, 500, err.Error(), nil)
		return
	}

	util.HttpResponse(context, 200, "ok", "ok")
	return
}

// download address excel file
// @Summary download address excel file
// @Schemes http https
// @Description download address excel file
// @Tags UserAddress
// @Success 200 {object} gin.Context
// @Router /useraddress/download [get]
func (h *UserAddressHandler) Download(context *gin.Context) {
	userAddressDownloadListResponse, titleList, err := h.IUserAddressService.Download(context)
	if err != nil {
		util.HttpResponse(context, 500, err.Error(), nil)
		return
	}

	fileName := util.GenFileName() + "_export.xlsx"
	util.ExportToExcel(context, titleList, userAddressDownloadListResponse, fileName)
}
