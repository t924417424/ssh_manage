# ssh_manage
go版本多用户webssh管理工具

项目仅用于学习交流，未经允许禁止任何其他用途

ssh2ws部分代码修改自https://github.com/hequan2017/go-webssh

服务端不保存用户明文密码，且不保存解密秘钥，如需对其他用户开放，请不要修改此部分代码，以免造成不必要的损失！


## 开发框架
Gin + gorm

##开发计划
✔ ssh功能

× 服务器文件管理

## 在线演示
[点击进入SSH云管理平台](https://www.do18.cn)
 
##环境
> Mysql
> Redis
> Go

##配置文件
>  修改config.toml的相关参数，短信接口使用阿里云短信
```toml
#配置文件
[Web]
model = "release"       #debug  release  test
port = "0.0.0.0:8082"   #服务要运行的端口

[Database]
host = "127.0.0.1"
port = 3306
username = "root"       #数据库账号
password = "root"       #数据库密码
dbname = "ssh"          #数据库名
poolsize = 10           #Mysql连接池大小

[Redis]
host = "127.0.0.1"
port = 6379
password = ""           #没有则不填
poolsize = 10           #Redis连接池大小

[Alisms]
accessid = "—"
accesskey = "-"
signname = "-"  #短信签名
template = "-"  #模板代码

```
##运行
###（Mysql会在首次使用时自动初始化）
```shell script
go build & ./ssh_manage
go run server.go 
```

## 前端
> Lauyi + Xterm.js


##补充说明
如需要使用Nginx等进行反代，请确保可以正常代理websocket

##免责声明  
本软件按“原样”提供，不提供任何形式的明示或暗示担保，包括但不限于对适销性，特定目的的适用性和非侵权性的担保。无论是由于软件，使用或其他方式产生的，与之有关或与之有关的合同，侵权或其他形式的任何索赔，损害或其他责任，作者或版权所有者概不负责。
