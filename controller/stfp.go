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

type sftpReq struct {
	Type  string `json:"type"`
	Token string `json:"token"`
}

type sftpResp struct {
	Code int    `json:"code"`
	Type string `json:"type"`
	Msg  string `json:"msg"`
	Data string `json:"data"`
}

func SftpSsh(c *gin.Context) {
	wsConn, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if core.HandleError(c, err) {
		return
	}
	defer wsConn.Close()

	var auth Apiform.WsAuth

	if c.ShouldBindUri(&auth) != nil {
		_ = wsConn.WriteMessage(websocket.TextMessage, []byte("参数错误\r\n"))
		_ = wsConn.Close()
		return
	}

	for {
		_, wsData, err := wsConn.ReadMessage()
		if err != nil {
			log.Println(err.Error())
			_ = wsConn.Close()
			//logrus.WithError(err).Error("reading webSocket message failed")
			return
		}
		//unmashal bytes into struct
		msgObj := sftpReq{}
		if err := json.Unmarshal(wsData, &msgObj); err != nil {
			log.Println("Auth : unmarshal websocket message failed:", string(wsData))
			continue
		}
		respMsg := sftpResp{}
		token := msgObj.Token
		claims, err := common.ParseToken(token)
		valid := claims.Valid()
		if valid != nil || err != nil {
			respMsg.Code = errcode.S_auth_fmt_err
			respMsg.Msg = "身份令牌校验不通过"
			respMsg.Data = err.Error()
			msg, _ := json.Marshal(respMsg)
			if err := wsConn.WriteMessage(websocket.TextMessage, msg); err != nil {
				log.Println("sftp token fmt err:", err)
			}
			_ = wsConn.Close()
			return
		}
		if claims.Userid != sftp_clients.Client.C[auth.Sid].Uid { //身份与缓存不符合
			respMsg.Code = errcode.S_auth_fmt_err
			respMsg.Msg = "用户权限不通过"
			respMsg.Data = err.Error()
			msg, _ := json.Marshal(respMsg)
			if err := wsConn.WriteMessage(websocket.TextMessage, msg); err != nil {
				log.Println("sftp server_user err:", err)
			}
			_ = wsConn.Close()
			return
		}

		path, err := sftp_clients.Client.C[auth.Sid].Sftp.Getwd()
		if err != nil {
			respMsg.Code = errcode.S_send_err
			respMsg.Type = "connect"
			respMsg.Msg = "SFTP连接失败"
			msg, _ := json.Marshal(respMsg)
			if err := wsConn.WriteMessage(websocket.TextMessage, msg); err != nil {
				log.Println("sftp connect err:", err)
			}
			return
		}

		respMsg.Code = 200
		respMsg.Type = "connect"
		respMsg.Msg = "连接成功"
		respMsg.Data = path
		msg, _ := json.Marshal(respMsg)
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
