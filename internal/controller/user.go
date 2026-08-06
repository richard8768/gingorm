package controller

import (
	"gin_demo/internal/dto"
	"gin_demo/internal/service"
	"gin_demo/internal/util"
	"io"
	"net/http"
	"os"
	"strconv"
	"sync/atomic"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	IUserService           service.IUserService
	IUserNormalFileService service.IUserNormalFileService
	IUserLargeFileService  service.IUserLargeFileService
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
	if err := util.CheckReqBindJson(context, &req); err != nil {
		util.HttpResponse(context, 500, err, nil)
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
	if err := util.CheckReqBindJson(context, &req); err != nil {
		util.HttpResponse(context, 500, err, nil)
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
	if err := util.CheckReqBindJson(context, &req); err != nil {
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
	if err := util.CheckReqBindJson(context, &req); err != nil {
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
	if err := util.CheckReqBindJson(context, &req); err != nil {
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
	if err := util.CheckReqBindJson(context, &req); err != nil {
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
	if err := util.CheckReqBindJson(context, &req); err != nil {
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
	if err := util.CheckReqBindJson(context, &req); err != nil {
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
// @Param body body dto.UserNormalFileUploadRequest true "请求body"
// @Success 200 {object} dto.UserNormalFileUploadResponse
// @Router /user/upload [post]
func (h *UserHandler) UserUpload(context *gin.Context) {
	file, err := util.HandleFileUpload(context, "file")
	if err != nil {
		util.HttpResponse(context, 500, err.Error(), nil)
		return
	}
	userUploadResponse, err := h.IUserNormalFileService.Upload(context, file, true)
	if err != nil {
		util.HttpResponse(context, 500, err.Error(), nil)
		return
	}
	util.HttpResponse(context, 200, "ok", userUploadResponse)
	return
}

func (h *UserHandler) UserUploadAvatar(context *gin.Context) {
	file, err := util.HandleFileUpload(context, "image")
	if err != nil {
		util.HttpResponse(context, 500, err.Error(), nil)
		return
	}
	userUploadAvatarResponse, err := h.IUserNormalFileService.Upload(context, file, false)
	if err != nil {
		util.HttpResponse(context, 500, err.Error(), nil)
		return
	}
	util.HttpResponse(context, 200, "ok", userUploadAvatarResponse)
	return
}

// download single file
// @Summary download single file
// @Schemes
// @Description download single file
// @Tags UploadHandler
// @Accept json
// @Produce json
// @Param body body dto.UserNormalFileDownloadRequest true "下载文件ID"
// @Success 200 {object} dto.UserNormalFileDownloadRequest
// @Router /user/download [get]
func (h *UserHandler) UserDownload(context *gin.Context) {
	var req dto.UserNormalFileDownloadRequest
	if err := util.CheckReqBindQuery(context, &req); err != nil {
		util.HttpResponse(context, 500, err, nil)
		return
	}

	filename, filepath, err := h.IUserNormalFileService.Download(context, &req)
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
	context.Writer.Header().Set("Content-Disposition", `attachment; filename="`+filename+`"`)
	context.Writer.Header().Set("Content-Length", strconv.FormatInt(stat.Size(), 10))
	io.Copy(context.Writer, f)
}

func (h *UserHandler) UserChunkInit(context *gin.Context) {
	var req dto.UserLargeFileUploadInitRequest
	if err := util.CheckReqBind(context, &req); err != nil {
		util.HttpResponse(context, 500, err, nil)
		return
	}
	response, err := h.IUserLargeFileService.ChunkInit(context, &req)
	if err != nil {
		util.HttpResponse(context, 500, err.Error(), nil)
		return
	}
	util.HttpResponse(context, 200, "ok", response)
}

var counter int64

func (h *UserHandler) UserChunkUploadList(context *gin.Context) {
	util.InitUpload(context)
	count := atomic.LoadInt64(&counter)
	if int(count) > 4 {
		context.AbortWithStatusJSON(500, gin.H{
			"code":    500,
			"message": "too many connections",
			"data":    nil,
		})
		return
	}
	atomic.AddInt64(&counter, 1)
	var req dto.UserLargeFileUploadRequest
	if err := util.CheckReqBindHeader(context, &req); err != nil {
		atomic.AddInt64(&counter, -1)
		util.HttpResponse(context, 500, err, nil)
		return
	}

	response, err := h.IUserLargeFileService.ChunkUpload(context, &req)
	_ = context.Request.Body.Close()
	atomic.AddInt64(&counter, -1)
	if err != nil {
		util.HttpResponse(context, 500, err.Error(), nil)
		return
	}
	util.HttpResponse(context, 200, "ok", response)
}

func (h *UserHandler) UserChunkUploadQuery(context *gin.Context) {
	var req dto.UserChunkUploadIdRequest
	if err := util.CheckReqBindQuery(context, &req); err != nil {
		util.HttpResponse(context, 500, err, nil)
		return
	}

	response, err := h.IUserLargeFileService.UserChunkUploadQuery(context, &req)
	if err != nil {
		util.HttpResponse(context, 500, err.Error(), nil)
		return
	}
	util.HttpResponse(context, 200, "ok", response)
}

func (h *UserHandler) UserChunkMerge(context *gin.Context) {
	var req dto.UserChunkUploadIdRequest
	if err := util.CheckReqBind(context, &req); err != nil {
		util.HttpResponse(context, 500, err, nil)
		return
	}
	_, err := h.IUserLargeFileService.ChunkMerge(context, &req, true)
	if err != nil {
		util.HttpResponse(context, 500, err.Error(), nil)
		return
	}
	util.HttpResponse(context, 200, "ok", "请求成功,文件正在合并中...")
}

func (h *UserHandler) UserChunkDownload(context *gin.Context) {
	var req dto.UserNormalFileDownloadRequest
	if err := util.CheckReqBindQuery(context, &req); err != nil {
		util.HttpResponse(context, 500, err, nil)
		return
	}

	filename, filepath, err := h.IUserNormalFileService.Download(context, &req)
	if err != nil {
		util.HttpResponse(context, 500, err.Error(), nil)
		return
	}

	file, err := os.Open(filepath)
	if err != nil {
		util.HttpResponse(context, 500, err.Error(), nil)
		return
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		util.HttpResponse(context, 500, err.Error(), nil)
		return
	}

	fileSize := strconv.FormatInt(fileInfo.Size(), 10)

	context.Writer.Header().Set("Content-Type", "application/octet-stream")
	context.Writer.Header().Set("Content-Disposition", `attachment; filename="`+filename+`"`)
	context.Writer.Header().Set("Content-Length", fileSize)

	http.ServeContent(context.Writer, context.Request, filename, fileInfo.ModTime(), file)

}
