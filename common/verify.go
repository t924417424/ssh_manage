package common

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"log"
	"regexp"
	"ssh_manage/database"
	"strings"
)

type verifyImpl interface {
	Verify() (key, code string)
}

func Verify(v verifyImpl) (is_verify bool) {
	phone, code := v.Verify()
	cache := database.Cache.Get()
	defer cache.Close()
	s_code, err := redis.String(cache.Do("GET", phone))
	if err != nil {
		log.Println("Verify Err:", err.Error())
		return
	}
	if code != s_code {
		log.Println(fmt.Sprintf("手机号：%s -- 验证码：%s 校验失败", phone, code))
		return
	}
	return true
}

func CheckIp(ip string) bool {
	addr := strings.Trim(ip, " ")
	regStr := `^(([1-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.)(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){2}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$`
	if match, _ := regexp.MatchString(regStr, addr); match {
		return true
	}
	return false
}