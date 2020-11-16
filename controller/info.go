package controller

import (
	"github.com/Gre-Z/common/jtime"
	"github.com/gin-gonic/gin"
	"ssh_manage/database"
	"ssh_manage/errcode"
	"ssh_manage/model"
	"ssh_manage/model/Apiform"
	"time"
)

func Info(c *gin.Context) {
	var resp Apiform.Resp
	new_token := c.MustGet("token").(string)
	if new_token != "" { //更新Token逻辑
		resp.Token = new_token
	}
	uid := c.MustGet("uid").(uint)
	var limit Apiform.Slist
	if c.MustGet("uid").(uint) > 0 {
		if c.ShouldBind(&limit) == nil {
			//var user model.User
			var list Apiform.List_resp
			var server model.Server
			server.BindUser = uid
			db := database.Get()
			defer db.Close()
			db.DB.Model(&model.Server{}).Where(&server).Count(&list.Count).Offset((limit.Page - 1) * limit.Limit).Limit(limit.Limit).Find(&list.List)
			//db.DB.Model(&user).Related(&servers,"Servers").Count(&list.Count).Offset((limit.Page - 1) * limit.Limit).Limit(limit.Limit).Find(&list.List)
			resp.Code = 200
			resp.Data = list
			resp.Msg = "查询成功"
		} else {
			resp.Code = errcode.C_from_err
			resp.Msg = "数据格式错误"
		}
	} else {
		resp.Code = errcode.S_Verify_err
		resp.Msg = "Token信息错误"
	}
	c.JSON(200, resp)
}

func UpdataNick(c *gin.Context) {
	var resp Apiform.Resp
	var edit Apiform.Edit
	new_token := c.MustGet("token").(string)
	if new_token != "" { //更新Token逻辑
		resp.Token = new_token
	}
	uid := c.MustGet("uid").(uint)
	//nickname, name_exist := c.GetPostForm("nickname")
	//sidstr, sid_exist := c.GetPostForm("id")
	//sid, err := strconv.Atoi(sidstr)
	//log.Println(c.ShouldBind(&edit))
	if c.ShouldBind(&edit) == nil {
		//server.Nickname = nickname
		var server model.Server
		server.ID = edit.ID
		server.BindUser = uid
		db := database.Get()
		defer db.Close()
		result := db.DB.Model(&model.Server{}).Where(&server).Update(model.Server{Nickname: edit.Nickname, Ip: edit.Ip, Port: edit.Port, Username: edit.Username})
		if result.RowsAffected == 1 && result.Error == nil {
			resp.Code = errcode.C_nil_err
			resp.Msg = "保存成功"
		} else {
			resp.Code = errcode.S_Db_err
			resp.Msg = "修改失败"
		}
	} else {
		resp.Code = errcode.C_from_err
		resp.Msg = "提交字段错误"
	}
	c.JSON(200, resp)
}

func Resetpass(c *gin.Context) {
	var resp Apiform.Resp
	var edit Apiform.EditPass
	new_token := c.MustGet("token").(string)
	if new_token != "" { //更新Token逻辑
		resp.Token = new_token
	}
	uid := c.MustGet("uid").(uint)
	//nickname, name_exist := c.GetPostForm("nickname")
	//sidstr, sid_exist := c.GetPostForm("id")
	//sid, err := strconv.Atoi(sidstr)
	//log.Println(c.ShouldBind(&edit))
	if c.ShouldBind(&edit) == nil {
		//server.Nickname = nickname
		var server model.Server
		server.ID = edit.ID
		server.BindUser = uid
		db := database.Get()
		defer db.Close()
		result := db.DB.Model(&model.Server{}).Where(&server).Update(model.Server{Password: edit.Password})
		if result.RowsAffected == 1 && result.Error == nil {
			resp.Code = errcode.C_nil_err
			resp.Msg = "保存成功"
		} else {
			resp.Code = errcode.S_Db_err
			resp.Msg = "修改失败"
		}
	} else {
		resp.Code = errcode.C_from_err
		resp.Msg = "提交字段错误"
	}
	c.JSON(200, resp)
}

func Del(c *gin.Context) {
	var resp Apiform.Resp
	var del Apiform.Edit
	new_token := c.MustGet("token").(string)
	if new_token != "" { //更新Token逻辑
		resp.Token = new_token
	}
	uid := c.MustGet("uid").(uint)
	if c.ShouldBind(&del) == nil {
		//server.Nickname = nickname
		var server model.Server
		server.ID = del.ID
		server.BindUser = uid
		db := database.Get()
		defer db.Close()
		result := db.DB.Where(&server).Delete(&model.Server{})
		if result.RowsAffected == 1 && result.Error == nil {
			resp.Code = errcode.C_nil_err
			resp.Msg = "删除成功"
		} else {
			resp.Code = errcode.S_Db_err
			resp.Msg = "操作失败"
		}
	} else {
		resp.Code = errcode.C_from_err
		resp.Msg = "提交字段错误"
	}
	c.JSON(200, resp)
}

func GetTerm(c *gin.Context) {
	var resp Apiform.Resp
	var term Apiform.GetTerm
	new_token := c.MustGet("token").(string)
	if new_token != "" { //更新Token逻辑
		resp.Token = new_token
	}
	uid := c.MustGet("uid").(uint)
	resp.Code = errcode.C_from_err
	resp.Msg = "表单错误"
	if c.ShouldBind(&term) == nil {
		var server model.Server
		server.ID = term.ID
		server.BindUser = uid
		db := database.Get()
		defer db.Close()
		result := db.DB.Model(&model.Server{}).First(&server)
		if result.RowsAffected == 1 && result.Error == nil {
			db.DB.Model(&model.Server{}).Where(&server).Update(model.Server{BeforeTime: jtime.JsonTime{time.Now()}})
			sid, err := term.Decode(server)
			//log.Println(sid)
			if err == nil {
				resp.Code = errcode.C_nil_err
				resp.Data = sid
				resp.Msg = "OK"
			} else {
				resp.Code = errcode.S_Verify_err
				resp.Msg = "秘钥解密失败"
			}
		} else {
			resp.Code = errcode.S_Db_err
			resp.Msg = "服务器信息检索失败"
		}
	}
	c.JSON(200,resp)
}
