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

func (h *MessageHandler) SendEmailHandler(context *gin.Context) {
	var req dto.SendEmailRequest
	if err := util.CheckReqBind(context, &req); err != nil {
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

func (h *MessageHandler) SendSmsHandler(context *gin.Context) {
	var req dto.SendSmsRequest
	if err := util.CheckReqBind(context, &req); err != nil {
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
