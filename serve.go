package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"ssh_manage/config"
	_ "ssh_manage/config"
	"ssh_manage/controller"
	"ssh_manage/controller/middleware"
	_ "ssh_manage/database" //初始化Mysql/Redis连接池
)

var run_mode = config.Config.Web.Model
var web_port = config.Config.Web.Port

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	gin.SetMode(run_mode)
	gin.DisableConsoleColor()
	router := gin.Default()
	router.Use(gin.Recovery())
	router.LoadHTMLGlob("view/*")
	router.Static("/static", "./static")
	router.GET("/", func(c *gin.Context) {
		c.Redirect(302, "/login")
	})
	router.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})
	router.GET("/console", func(c *gin.Context) {
		c.HTML(http.StatusOK, "console.html", nil)
	})
	router.GET("/servers", func(c *gin.Context) {
		c.HTML(http.StatusOK, "s_list.html", nil)
	})
	router.GET("/add", func(c *gin.Context) {
		c.HTML(http.StatusOK, "add.html", nil)
	})
	router.GET("/setpass", func(c *gin.Context) {
		c.HTML(http.StatusOK, "reset.html", nil)
	})
	router.GET("/openterm", func(c *gin.Context) {
		c.HTML(http.StatusOK, "open_term.html", nil)
	})
	router.GET("/term", func(c *gin.Context) {
		c.HTML(http.StatusOK, "term.html", nil)
	})

	api := router.Group("/v1")
	{
		api.POST("/login", controller.Login)
		api.POST("/send", controller.Send)
		api.GET("/term/:sid", controller.WsSsh)
		api.GET("/sftp/:sid", controller.Sftp_ssh)
		api.Use(middleware.Auth()).GET("/userinfo", controller.Info)
		api.Use(middleware.Auth()).POST("/nickname", controller.UpdataNick)
		api.Use(middleware.Auth()).POST("/addser", controller.Addser)
		api.Use(middleware.Auth()).POST("/repass", controller.Resetpass)
		api.Use(middleware.Auth()).POST("/delete", controller.Del)
		api.Use(middleware.Auth()).POST("/getterm", controller.GetTerm)
	}
	if err := router.Run(web_port); err != nil {
		log.Panicf("Web Serve Start Err : %s", err.Error())
	}
}
