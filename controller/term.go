package controller

import (
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"ssh_manage/common"
	"ssh_manage/common/core"
	"ssh_manage/common/sftp_clients"
	"ssh_manage/database"
	"ssh_manage/model/Apiform"
	"strconv"
)

var upGrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024 * 1024 * 10,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type AuthMsg struct {
	Type  string `json:"type"`
	Token string `json:"token"`
}

// handle webSocket connection.
// first,we establish a ssh connection to ssh server when a webSocket comes;
// then we deliver ssh data via ssh connection between browser and ssh server.
// That is, read webSocket data from browser (e.g. 'ls' command) and send data to ssh server via ssh connection;
// the other hand, read returned ssh data from ssh server and write back to browser via webSocket API.
func WsSsh(c *gin.Context) {
	wsConn, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if core.HandleError(c, err) {
		return
	}
	defer wsConn.Close()

	cols, err := strconv.Atoi(c.DefaultQuery("cols", "120"))
	if core.WshandleError(wsConn, err) {
		return
	}
	rows, err := strconv.Atoi(c.DefaultQuery("rows", "32"))
	if core.WshandleError(wsConn, err) {
		return
	}

	var serInfo Apiform.SerInfo //接收反序列化数据
	var auth Apiform.WsAuth

	if c.ShouldBindUri(&auth) != nil {
		_ = wsConn.WriteMessage(websocket.BinaryMessage, []byte("参数错误\r\n"))
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
		msgObj := AuthMsg{}
		if err := json.Unmarshal(wsData, &msgObj); err != nil {
			log.Println("Auth : unmarshal websocket message failed:", string(wsData))
			continue
		}
		token := msgObj.Token
		claims, err := common.ParseToken(token)
		valid := claims.Valid()
		if valid != nil || err != nil {
			_ = wsConn.WriteMessage(websocket.BinaryMessage, []byte("身份验证失败\r\n"))
			_ = wsConn.Close()
			return
		}
		cache := database.Cache.Get()
		defer cache.Close()
		//log.Println(auth)
		sInfo, err := redis.Bytes(cache.Do("GET", auth.Sid))
		//log.Println(string(s_info))
		if err != nil || len(sInfo) == 0 {
			_ = wsConn.WriteMessage(websocket.BinaryMessage, []byte("连接超时，请重试！\r\n"))
			_ = wsConn.Close()
			return
		}
		if json.Unmarshal(sInfo, &serInfo) != nil {
			_ = wsConn.WriteMessage(websocket.BinaryMessage, []byte("服务器信息获取失败，请重试！\r\n"))
			_ = wsConn.Close()
			return
		}
		//log.Println(ser_info)
		if claims.Userid != serInfo.BindUser { //验证权限
			_ = wsConn.WriteMessage(websocket.BinaryMessage, []byte("权限验证失败，请重试！\r\n"))
			_ = wsConn.Close()
			return
		}
		break
		//break
	}
	client, err := core.NewSshClient(core.Server{serInfo.Ip, serInfo.Port, serInfo.Username, serInfo.Password})
	if core.WshandleError(wsConn, err) {
		return
	}
	defer client.Close()
	//startTime := time.Now()
	ssConn, err := core.NewSshConn(cols, rows, client) //加入sftp客户端
	if core.WshandleError(wsConn, err) {
		return
	}
	sftp_clients.Client.Lock()
	sftp_clients.Client.C[auth.Sid] = &sftp_clients.MyClient{serInfo.BindUser, ssConn.SftpClient}
	sftp_clients.Client.Unlock()
	defer func() {
		sftp_clients.Client.Lock()
		delete(sftp_clients.Client.C, auth.Sid) //释放SFTP客户端
		sftp_clients.Client.Unlock()
	}()
	defer ssConn.Close()
	quitChan := make(chan bool, 3)

	// most messages are ssh output, not webSocket input
	go ssConn.ReceiveWsMsg(wsConn, quitChan)
	go ssConn.SendComboOutput(wsConn, quitChan)
	go ssConn.SessionWait(quitChan)

	<-quitChan //任意协程退出则结束
	fmt.Println("Exit")
	log.Println("websocket finished")
}
