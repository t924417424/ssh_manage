package model

import (
	"github.com/Gre-Z/common/jtime"
)

type Server struct {
	Model
	Nickname   string
	Ip         string `grom:"size:15"`
	Port       int
	Username   string `grom:"size:255"`
	Password   string `gorm:"type:longtext" json:"-"` //存放加密后的登录密码，解密密码存放到用户浏览器
	BindUser   uint    `json:"-"`
	BeforeTime jtime.JsonTime
}
