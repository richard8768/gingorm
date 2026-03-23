package controller

import (
	"gin_demo/internal/dto"
	"gin_demo/internal/service"
	"gin_demo/internal/util"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type UserHandler struct {
	IUserService service.IUserService
}

// user reg
// @Summary UserReg
// @Schemes
// @Description UserReg
// @Tags UserReg
// @Accept json
// @Produce json
// @Param body body dto.UserCreateRequest true "请求body"
// @Success 200 {object} dto.UserResponse
// @Router /user/reg [post]
func (h *UserHandler) UserReg(context *gin.Context) {
	var req dto.UserCreateRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			util.HttpResponse(context, 500, err.Error(), nil)
			return
		}
		util.HttpResponse(context, 500, util.RemoveTopStruct(errs.Translate(util.Trans)), nil)
		return
	}

	userRegResponse, err := h.IUserService.UserReg(context, &req)
	if err != nil {
		util.HttpResponse(context, 500, err.Error(), nil)
		return
	}

	util.HttpResponse(context, 200, "ok", userRegResponse)
	return
}

// user login
// @Summary UserLogin
// @Schemes
// @Description UserLogin
// @Tags UserLogin
// @Accept json
// @Produce json
// @Param body body dto.UserLoginRequest true "请求body"
// @Success 200 {object} dto.UserLoginResponse
// @Router /user/login [post]
func (h *UserHandler) UserLogin(context *gin.Context) {
	var req dto.UserLoginRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			util.HttpResponse(context, 500, err.Error(), nil)
			return
		}
		util.HttpResponse(context, 500, util.RemoveTopStruct(errs.Translate(util.Trans)), nil)
		return
	}

	userLoginResponse, err := h.IUserService.UserLogin(context, &req)
	if err != nil {
		util.HttpResponse(context, 500, err.Error(), nil)
		return
	}
	util.HttpResponse(context, 200, "ok", userLoginResponse)
	return
}

// user index
// @Summary UserIndex
// @Schemes
// @Description UserIndex
// @Tags UserIndex
// @Accept json
// @Produce json
// @Success 200 {object} dto.UserResponse
// @Router /user/index [get]
func (h *UserHandler) UserIndex(context *gin.Context) {
	userInfo, err := h.IUserService.GetUserInfo(context)
	if err != nil {
		util.HttpResponse(context, 500, err.Error(), nil)
		return
	}
	util.HttpResponse(context, 200, "ok", userInfo)
	return
}
