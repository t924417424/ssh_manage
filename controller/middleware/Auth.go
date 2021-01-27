package middleware

import (
	"github.com/gin-gonic/gin"
	"ssh_manage/common"
	"ssh_manage/database"
	"ssh_manage/errcode"
	"ssh_manage/model"
	"ssh_manage/model/Apiform"
	"strings"
	"time"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		var resp Apiform.Resp
		jwtToken := c.GetHeader("Authorization")
		//log.Println(jwt_token)
		//log.Println(strings.HasPrefix(jwt_token, "Bearer "))
		if jwtToken == "" || !strings.HasPrefix(jwtToken, "Bearer ") {
			resp.Code = errcode.S_auth_fmt_err
			resp.Msg = "Token不正确"
			c.JSON(200, resp)
			c.Abort()
			return
		}
		jwtToken = jwtToken[7:]
		claims, err := common.ParseToken(jwtToken)
		if err != nil {
			resp.Code = errcode.S_auth_err
			resp.Msg = "Token错误，请重新登录"
			c.JSON(200, resp)
			c.Abort()
			return
		}
		valid := claims.Valid()
		if valid != nil {
			resp.Code = errcode.S_auth_err
			resp.Msg = "用户登录超时，请重新登录"
			c.JSON(200, resp)
			c.Abort()
			return
		}
		var userInfo model.User
		db := database.Get()
		defer db.Close()
		userInfo.ID = claims.Userid
		db.DB.Where(userInfo).First(&userInfo)
		if userInfo.Phone == 0 {
			resp.Code = errcode.S_auth_err
			resp.Msg = "用户不存在，请重新登录"
			c.JSON(200, resp)
			c.Abort()
			return
		}
		c.Set("uid", claims.Userid)
		c.Set("token", "")
		newToken, err := common.ReleaseToken(claims.Userid)
		if time.Now().Add(24*time.Hour).Unix() > claims.ExpiresAt { //如果过期时间小于一天，则更新客户端token
			c.Set("token", newToken)
		}
		c.Next()
	}
}
