package controller

import (
	"gin_demo/internal/dto"
	"gin_demo/internal/service"
	"gin_demo/internal/util"

	"github.com/gin-gonic/gin"
)

type MessageHandler struct {
	IMessageService service.IMessageService
}

// send email code
// @Summary send email code
// @Schemes http https
// @Description send email code
// @Accept json
// @Produce json
// @Param body body dto.SendEmailRequest true "请求body"
// @Success 200 {string} Email sent successfully
// @Router /sendEmail [post]
func (h *MessageHandler) SendEmailHandler(context *gin.Context) {
	var req dto.SendEmailRequest
	if err := util.CheckReqBindJson(context, &req); err != nil {
		util.HttpResponse(context, 500, err, nil)
		return
	}

	_, err := h.IMessageService.SendEmail(context, &req)
	if err != nil {
		util.HttpResponse(context, 500, err.Error(), nil)
		return
	}

	util.HttpResponse(context, 200, "ok", "Email sent successfully")
	return
}

// send sms code
// @Summary send sms code
// @Schemes http https
// @Description send sms code
// @Accept json
// @Produce json
// @Param body body dto.SendSmsRequest true "请求body"
// @Success 200 {string} SMS sent successfully
// @Router /sendSms [post]
func (h *MessageHandler) SendSmsHandler(context *gin.Context) {
	var req dto.SendSmsRequest
	if err := util.CheckReqBindJson(context, &req); err != nil {
		util.HttpResponse(context, 500, err, nil)
		return
	}

	_, err := h.IMessageService.SendSms(context, &req)
	if err != nil {
		util.HttpResponse(context, 500, err.Error(), nil)
		return
	}

	util.HttpResponse(context, 200, "ok", "SMS sent successfully")
	return
}
