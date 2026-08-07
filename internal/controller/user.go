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
// @Schemes http https
// @Description UserReg
// @Tags User
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
// @Schemes http https
// @Description UserLogin
// @Tags User
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
// @Schemes http https
// @Description UserIndex
// @Tags User
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
// @Schemes http https
// @Description UserLogout
// @Tags User
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

// user bind login mobile
// @Summary user bind login mobile
// @Schemes http https
// @Description user bind login mobile
// @Tags User
// @Accept json
// @Produce json
// @Param body body dto.UserBindLoginMobileRequest true "请求body"
// @Success 200 {string} ok
// @Router /user/bindLoginMobile [post]
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
	util.HttpResponse(context, 200, "ok", "ok")
}

// user bind login email
// @Summary user bind login email
// @Schemes http https
// @Description user bind login email
// @Tags User
// @Accept json
// @Produce json
// @Param body body dto.UserBindLoginEmailRequest true "请求body"
// @Success 200 {string} ok
// @Router /user/bindLoginEmail [post]
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
	util.HttpResponse(context, 200, "ok", "ok")
}

// user check bind login email/mobile
// @Summary user check bind login email/mobile
// @Schemes http https
// @Description user check bind login email/mobile
// @Tags User
// @Accept json
// @Produce json
// @Param body body dto.UserCheckBindMobileEmailRequest true "请求body"
// @Success 200 {object} dto.UserCheckBindMobileEmailResponse
// @Router /user/checkBindMobile [post]
func (h *UserHandler) UserCheckBindMobile(context *gin.Context) {
	h.handleUserCheckBindMobileEmail(context)
}

// user check bind login email/mobile
// @Summary user check bind login email/mobile
// @Schemes http https
// @Description user check bind login email/mobile
// @Tags User
// @Accept json
// @Produce json
// @Param body body dto.UserCheckBindMobileEmailRequest true "请求body"
// @Success 200 {object} dto.UserCheckBindMobileEmailResponse
// @Router /user/checkBindEmail [post]
func (h *UserHandler) UserCheckBindEmail(context *gin.Context) {
	h.handleUserCheckBindMobileEmail(context)
}

// user check bind login email/mobile
// @Summary user check bind login email/mobile
// @Schemes http https
// @Description user check bind login email/mobile
// @Tags User
// @Accept json
// @Produce json
// @Param body body dto.UserCheckBindMobileEmailRequest true "请求body"
// @Success 200 {object} dto.UserCheckBindMobileEmailResponse
// @Router /user/checkBindMobileEmail [post]
func (h *UserHandler) UserCheckBindMobileEmail(context *gin.Context) {
	h.handleUserCheckBindMobileEmail(context)
}

func (h *UserHandler) handleUserCheckBindMobileEmail(context *gin.Context) {
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

// user change password
// @Summary user change password
// @Schemes http https
// @Description user change password
// @Tags User
// @Accept json
// @Produce json
// @Param body body dto.UserChangePwdRequest true "请求body"
// @Success 200 {string} ok
// @Router /user/changePwd [post]
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
	util.HttpResponse(context, 200, "ok", "ok")
}

// user update profile
// @Summary user update profile
// @Schemes http https
// @Description user update profile
// @Tags User
// @Accept json
// @Produce json
// @Param body body dto.UserUpdateProfileRequest true "请求body"
// @Success 200 {string} ok
// @Router /user/updateProfile [post]
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
	util.HttpResponse(context, 200, "ok", "ok")
}

// user reset password
// @Summary user reset password
// @Schemes http https
// @Description user reset password
// @Tags User
// @Accept json
// @Produce json
// @Param body body dto.UserResetPwdRequest true "请求body"
// @Success 200 {string} ok
// @Router /user/resetPwd [post]
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
	util.HttpResponse(context, 200, "ok", "ok")
}

// upload  file
// @Summary upload  file
// @Schemes http https
// @Description upload  file
// @Tags User
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "文件"
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

// upload  avatar
// @Summary upload  avatar
// @Schemes http https
// @Description upload  avatar
// @Tags User
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "文件"
// @Success 200 {object} dto.UserNormalFileUploadResponse
// @Router /user/uploadAvatar [post]
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

// download  file
// @Summary download  file
// @Schemes http https
// @Description download  file
// @Tags User
// @Param id query uint true "下载文件ID"
// @Success 200 {object} gin.Context
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

// user large file chunk upload init
// @Summary user large file chunk upload init
// @Schemes http https
// @Description user large file chunk upload init
// @Tags User
// @Accept json
// @Produce json
// @Param body body dto.UserLargeFileUploadInitRequest true "请求body"
// @Success 200 {object} dto.UserLargeFileUploadInitResponse
// @Router /user/chunkInit [post]
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

// user large file chunk upload
// @Summary user large file chunk upload
// @Schemes http https
// @Description user large file chunk upload
// @Tags User
// @Accept application/x-www-form-urlencoded
// @Produce json
// @Param upload_id header string true "upload_id"
// @Param chunk_index header int true "chunk_index"
// @Param chunk_md5 header string true "chunk_md5"
// @Param total_chunks header int true "total_chunks"
// @Param body body byte true "文件分块二进制内容"
// @Success 200 {object} dto.UserLargeFileUploadResponse
// @Router /user/chunkUpload [post]
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

// user large file chunk upload query
// @Summary user large file chunk upload query
// @Schemes http https
// @Description user large file chunk upload query
// @Tags User
// @Produce json
// @Param upload_id query string true "upload_id"
// @Success 200 {object} dto.UserLargeFileUploadResponse
// @Router /user/chunkUploadQuery [get]
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

// user large file chunk merge
// @Summary user large file chunk merge
// @Schemes http https
// @Description user large file chunk merge
// @Tags User
// @Accept json
// @Param body body dto.UserChunkUploadIdRequest true "请求body"
// @Produce json
// @Success 200 {string} ok
// @Router /user/chunkMerge [post]
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
	util.HttpResponse(context, 200, "ok", "ok")
}

// user large file download
// @Summary user large file download
// @Schemes http https
// @Description user large file download
// @Tags User
// @Param id query uint true "下载文件ID"
// @Success 200 {object} gin.Context
// @Router /user/chunkDownload [get]
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
