package common

import (
	"errors"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	"log"
	"math/rand"
	"regexp"
	"ssh_manage/config"
	"time"
)

var aliconfig = config.Config.Alisms

func VerifyMobileFormat(mobileNum string) bool {
	regular := "^((13[0-9])|(14[5,7])|(15[0-3,5-9])|(17[0,3,5-8])|(18[0-9])|166|198|199|(147))\\d{8}$"
	reg := regexp.MustCompile(regular)
	return reg.MatchString(mobileNum)
}

func Sendsms(phone string) (captcha string,err error) {
	captcha = createCaptcha()
	client, err := dysmsapi.NewClientWithAccessKey("cn-hangzhou", aliconfig.Accessid, aliconfig.Accesskey)
	request := dysmsapi.CreateSendBatchSmsRequest()
	request.Scheme = "https"
	request.PhoneNumberJson = fmt.Sprintf("[\"%s\"]", phone)
	request.SignNameJson = fmt.Sprintf("[\"%s\"]", aliconfig.Signname)
	request.TemplateCode = aliconfig.Template
	request.TemplateParamJson = fmt.Sprintf("[{\"code\":\"%s\"}]", captcha)
	response, err := client.SendBatchSms(request)
	//if err != nil {
	//	log.Println(err.Error())
	//}
	if response.Code != "OK" {
		err = errors.New("短信服务器错误")
		log.Println(response)
	}
	return
}

func createCaptcha() string {
	return fmt.Sprintf("%08v", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(100000000))
}
