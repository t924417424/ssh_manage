package controller

import (
	"github.com/gin-gonic/gin"
	"ssh_manage/common"
	"ssh_manage/database"
	"ssh_manage/errcode"
	"ssh_manage/model"
	"ssh_manage/model/Apiform"
)

func Login(c *gin.Context) {
	//common.Sendsms()
	//log.Print(db.DB.Exec("select * from products"))
	//token := c.MustGet("token").(string)
	//c.JSON(200, gin.H{"token": token})
	var resp Apiform.Resp
	resp.Code = errcode.C_from_err
	resp.Msg = "手机号和验证码不能为空！"
	var user Apiform.Login
	if c.ShouldBind(&user) == nil {
		if common.Verify(&user) {
			var userinfo model.User
			db := database.Get()
			defer db.Close()
			db.DB.Where(model.User{Phone: user.Phone}).FirstOrCreate(&userinfo)
			newToken, err := common.ReleaseToken(userinfo.ID)
			if err == nil && userinfo.ID > 0 {
				resp.Code = errcode.C_nil_err
				resp.Msg = "登陆成功"
				resp.Data = userinfo
				resp.Token = newToken
			} else {
				resp.Code = errcode.S_auth_err
				resp.Msg = "Token创建失败"
			}
		} else {
			resp.Code = errcode.S_Verify_err
			resp.Msg = "验证码校验失败"
		}
	}
	//log.Printf(c.ClientIP())
	c.JSON(200, resp)
}
