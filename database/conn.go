package database

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
	"ssh_manage/config"
	"ssh_manage/model"
)

var dbConf = config.Config.Database

var pool *sqlpool

type MyDb struct {
	DB *gorm.DB
}

func init() {
	pool = newpool(newDb, dbConf.Poolsize)
}

func newDb() *gorm.DB {
	db, err := gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", dbConf.Username, dbConf.Password, dbConf.Host, dbConf.Port, dbConf.Dbname))
	if err != nil {
		log.Panicf("db open err :%s", err.Error())
	}
	if !db.HasTable(&model.User{}){
		log.Println("init table")
		db.CreateTable(&model.Server{},&model.User{})
	}
	return db
}

func (s *MyDb) Close() {
	pool.put(s.DB)
}

func Get() *MyDb {
	return &MyDb{pool.get()}
}
