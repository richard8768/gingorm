package controller

import (
	"errors"
	"gin_demo/internal/dto"
	"gin_demo/internal/service"
	"gin_demo/internal/util"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type UserHandler struct {
	IUserService           service.IUserService
	IUserSingleFileService service.IUserSingleFileService
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

// user logout
// @Summary UserLogout
// @Schemes
// @Description UserLogout
// @Tags UserLogout
// @Accept json
// @Produce json
// @Success 200 {object} dto.UserLoginResponse
// @Router /user/logout [get]
func (h *UserHandler) UserLogout(context *gin.Context) {
	userLogoutResponse, err := h.IUserService.UserLogout(context)
	if err != nil {
		util.HttpResponse(context, 500, err.Error(), nil)
		return
	}
	util.HttpResponse(context, 200, "ok", userLogoutResponse)
	return
}

func (h *UserHandler) UserBindLoginMobile(context *gin.Context) {
	var req dto.UserBindLoginMobileRequest
	if err := util.CheckReqBind(context, &req); err != nil {
		util.HttpResponse(context, 500, err, nil)
		return
	}
	_, err := h.IUserService.UserBindLoginMobile(context, &req)
	if err != nil {
		util.HttpResponse(context, 500, err.Error(), nil)
		return
	}
	util.HttpResponse(context, 200, "ok", nil)
}

func (h *UserHandler) UserBindLoginEmail(context *gin.Context) {
	var req dto.UserBindLoginEmailRequest
	if err := util.CheckReqBind(context, &req); err != nil {
		util.HttpResponse(context, 500, err, nil)
		return
	}
	_, err := h.IUserService.UserBindLoginEmail(context, &req)
	if err != nil {
		util.HttpResponse(context, 500, err.Error(), nil)
		return
	}
	util.HttpResponse(context, 200, "ok", nil)
}

func (h *UserHandler) UserCheckBindMobileEmail(context *gin.Context) {
	var req dto.UserCheckBindMobileEmailRequest
	if err := util.CheckReqBind(context, &req); err != nil {
		util.HttpResponse(context, 500, err, nil)
		return
	}
	UserCheckBindMobileEmailResponse, err := h.IUserService.UserCheckBindMobileEmail(context, &req)
	if err != nil {
		util.HttpResponse(context, 500, err.Error(), nil)
		return
	}
	util.HttpResponse(context, 200, "ok", UserCheckBindMobileEmailResponse)
}

func (h *UserHandler) UserChangePwd(context *gin.Context) {
	var req dto.UserChangePwdRequest
	if err := util.CheckReqBind(context, &req); err != nil {
		util.HttpResponse(context, 500, err, nil)
		return
	}
	_, err := h.IUserService.UserChangePwd(context, &req)
	if err != nil {
		util.HttpResponse(context, 500, err.Error(), nil)
		return
	}
	util.HttpResponse(context, 200, "ok", nil)
}

func (h *UserHandler) UserUpdateProfile(context *gin.Context) {
	var req dto.UserUpdateProfileRequest
	if err := util.CheckReqBind(context, &req); err != nil {
		util.HttpResponse(context, 500, err, nil)
		return
	}
	_, err := h.IUserService.UserUpdateProfile(context, &req)
	if err != nil {
		util.HttpResponse(context, 500, err.Error(), nil)
		return
	}
	util.HttpResponse(context, 200, "ok", nil)
}

func (h *UserHandler) UserResetPwd(context *gin.Context) {
	var req dto.UserResetPwdRequest
	if err := util.CheckReqBind(context, &req); err != nil {
		util.HttpResponse(context, 500, err, nil)
		return
	}
	_, err := h.IUserService.UserResetPwd(context, &req)
	if err != nil {
		util.HttpResponse(context, 500, err.Error(), nil)
		return
	}
	util.HttpResponse(context, 200, "ok", nil)
}

// upload single file
// @Summary upload single file
// @Schemes
// @Description upload single file
// @Tags UploadHandler
// @Accept json
// @Produce json
// @Param body body dto.UserSingleFileUploadRequest true "请求body"
// @Success 200 {object} dto.UserSingleFileUploadResponse
// @Router /user/upload [post]
func (h *UserHandler) UserUpload(context *gin.Context) {
	file, err := handleFileUpload(context, "file")
	if err != nil {
		util.HttpResponse(context, 500, err.Error(), nil)
		return
	}
	userUploadResponse, err := h.IUserSingleFileService.Upload(context, file, true)
	if err != nil {
		util.HttpResponse(context, 500, err.Error(), nil)
		return
	}
	util.HttpResponse(context, 200, "ok", userUploadResponse)
	return
}

func (h *UserHandler) UserUploadAvatar(context *gin.Context) {
	file, err := handleFileUpload(context, "image")
	if err != nil {
		util.HttpResponse(context, 500, err.Error(), nil)
		return
	}
	userUploadAvatarResponse, err := h.IUserSingleFileService.Upload(context, file, false)
	if err != nil {
		util.HttpResponse(context, 500, err.Error(), nil)
		return
	}
	util.HttpResponse(context, 200, "ok", userUploadAvatarResponse)
	return
}

func handleFileUpload(context *gin.Context, fileType string) (*multipart.FileHeader, error) {
	if fileType != "file" && fileType != "image" {
		return nil, errors.New("invalid file type")
	}
	validate, ok := binding.Validator.Engine().(*validator.Validate)
	if fileType == "file" {
		if ok {
			validate.RegisterStructValidation(util.FileUploadValidation, dto.UserSingleFileUploadRequest{})
		}
	} else {
		if ok {
			validate.RegisterStructValidation(util.FileUploadValidation, dto.UserAvatarUploadRequest{})
		}
	}
	var req dto.UserSingleFileUploadRequest
	if err := context.ShouldBind(&req); err != nil {
		util.HttpResponse(context, 500, err.Error(), nil)
		return nil, err
	}

	//file, err := context.FormFile("file")
	//if err != nil {
	//	util.HttpResponse(context, 500, err.Error(), nil)
	//	return
	//}

	form, err := context.MultipartForm()
	if err != nil {
		util.HttpResponse(context, 500, err.Error(), nil)
		return nil, err
	}
	files := form.File["file"]
	if files == nil {
		util.HttpResponse(context, 500, "file is empty", nil)
		return nil, err
	}

	file := files[0]
	return file, nil
}

// download single file
// @Summary download single file
// @Schemes
// @Description download single file
// @Tags UploadHandler
// @Accept json
// @Produce json
// @Param body body dto.UserSingleFileDownloadRequest true "下载文件ID"
// @Success 200 {object} dto.UserSingleFileDownloadRequest
// @Router /user/download [get]
func (h *UserHandler) UserDownload(context *gin.Context) {
	var req dto.UserSingleFileDownloadRequest
	if err := context.ShouldBindQuery(&req); err != nil {
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			util.HttpResponse(context, 500, err.Error(), nil)
			return
		}
		util.HttpResponse(context, 500, util.RemoveTopStruct(errs.Translate(util.Trans)), nil)
		return
	}

	filename, filepath, err := h.IUserSingleFileService.Download(context, &req)
	if err != nil {
		util.HttpResponse(context, 500, err.Error(), nil)
		return
	}

	f, err := os.Open(filepath)
	if err != nil {
		http.Error(context.Writer, "File not found", http.StatusNotFound)
		return

	}
	defer f.Close()

	stat, _ := f.Stat()
	context.Writer.Header().Set("Content-Type", "application/octet-stream")
	context.Writer.Header().Set("Content-Disposition", `attachment; filename=`+filename+``)
	context.Writer.Header().Set("Content-Length", strconv.FormatInt(stat.Size(), 10))
	io.Copy(context.Writer, f)
}
func (h *UserHandler) UserChunkUpload(context *gin.Context) {
	util.HttpResponse(context, 200, "ok", nil)
}

func (h *UserHandler) UserChunkMerge(context *gin.Context) {
	util.HttpResponse(context, 200, "ok", nil)
}

func (h *UserHandler) UserChunkDownload(context *gin.Context) {
	util.HttpResponse(context, 200, "ok", nil)
}
