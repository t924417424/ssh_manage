package core

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"time"
)

func JsonError(c *gin.Context, msg interface{}) {
	c.AbortWithStatusJSON(200, gin.H{"ok": false, "msg": msg})
}

func HandleError(c *gin.Context, err error) bool {
	if err != nil {
		//logrus.WithError(err).Error("gin context http handler error")
		JsonError(c, err.Error())
		return true
	}
	return false
}

func WshandleError(ws *websocket.Conn, err error) bool {
	if err != nil {
		log.Println("handler ws ERROR:",err.Error())
		dt := time.Now().Add(time.Second)
		if err := ws.WriteControl(websocket.CloseMessage, []byte(err.Error()), dt); err != nil {
			log.Println("websocket writes control message failed:",err.Error())
		}
		return true
	}
	return false
}
