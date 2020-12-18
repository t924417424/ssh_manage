package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"ssh_manage/common"
	"ssh_manage/common/core"
	"ssh_manage/common/sftp_clients"
	"ssh_manage/errcode"
	"ssh_manage/model/Apiform"
)

type sftp_req struct {
	Type  string `json:"type"`
	Token string `json:"token"`
}

type sftp_resp struct {
	Code int    `json:"code"`
	Type string `json:"type"`
	Msg  string `json:"msg"`
	Data string `json:"data"`
}

func Sftp_ssh(c *gin.Context) {
	wsConn, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if core.HandleError(c, err) {
		return
	}
	defer wsConn.Close()

	var auth Apiform.WsAuth

	if c.ShouldBindUri(&auth) != nil {
		wsConn.WriteMessage(websocket.TextMessage, []byte("参数错误\r\n"))
		wsConn.Close()
		return
	}

	for {
		_, wsData, err := wsConn.ReadMessage()
		if err != nil {
			log.Println(err.Error())
			wsConn.Close()
			//logrus.WithError(err).Error("reading webSocket message failed")
			return
		}
		//unmashal bytes into struct
		msgObj := sftp_req{}
		if err := json.Unmarshal(wsData, &msgObj); err != nil {
			log.Println("Auth : unmarshal websocket message failed:", string(wsData))
			continue
		}
		resp_msg := sftp_resp{}
		token := msgObj.Token
		claims, err := common.ParseToken(token)
		valid := claims.Valid()
		if valid != nil || err != nil {
			resp_msg.Code = errcode.S_auth_fmt_err
			resp_msg.Msg = "身份令牌校验不通过"
			resp_msg.Data = err.Error()
			msg, _ := json.Marshal(resp_msg)
			if err := wsConn.WriteMessage(websocket.TextMessage, msg); err != nil {
				log.Println("sftp token fmt err:", err)
			}
			wsConn.Close()
			return
		}
		if claims.Userid != sftp_clients.Client.C[auth.Sid].Uid { //身份与缓存不符合
			resp_msg.Code = errcode.S_auth_fmt_err
			resp_msg.Msg = "用户权限不通过"
			resp_msg.Data = err.Error()
			msg, _ := json.Marshal(resp_msg)
			if err := wsConn.WriteMessage(websocket.TextMessage, msg); err != nil {
				log.Println("sftp server_user err:", err)
			}
			wsConn.Close()
			return
		}

		path, err := sftp_clients.Client.C[auth.Sid].Sftp.Getwd()
		if err != nil {
			resp_msg.Code = errcode.S_send_err
			resp_msg.Type = "connect"
			resp_msg.Msg = "SFTP连接失败"
			msg, _ := json.Marshal(resp_msg)
			if err := wsConn.WriteMessage(websocket.TextMessage, msg); err != nil {
				log.Println("sftp connect err:", err)
			}
			return
		}

		resp_msg.Code = 200
		resp_msg.Type = "connect"
		resp_msg.Msg = "连接成功"
		resp_msg.Data = path
		msg, _ := json.Marshal(resp_msg)
		if err := wsConn.WriteMessage(websocket.TextMessage, msg); err != nil {
			log.Println("sftp return err:", err)
			return
		}

		break
		//break
	}
	quitChan := make(chan bool, 2)
	go sftp_clients.Client.C[auth.Sid].ReceiveWsMsg(wsConn, quitChan)
	<-quitChan //任意协程退出则结束
	fmt.Println("Sftp Exit")
	log.Println("sftp websocket finished")
}
