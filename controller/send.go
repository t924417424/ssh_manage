package controller

import (
	"github.com/gin-gonic/gin"
	"ssh_manage/common"
	"ssh_manage/errcode"
	"ssh_manage/model/Apiform"
)

func Send(c *gin.Context) {
	var resp Apiform.Resp
	var send Apiform.Send
	resp.Code = errcode.C_phone_err
	resp.Msg = "手机号未提交！"
	if c.ShouldBind(&send) == nil {
		if common.VerifyMobileFormat(send.Phone) {
			if err := send.SendCaptcha(c.ClientIP()); err != nil {
				resp.Code = errcode.S_send_err
				resp.Msg = err.Error()
			} else {
				resp.Code = errcode.C_nil_err
				resp.Msg = "发送成功！"
			}
		} else {
			resp.Code = errcode.C_phone_err
			resp.Msg = "手机号验证失败！"
		}
	}
	c.JSON(200, resp)
}
