package database

import (
	"github.com/jinzhu/gorm"
	"sync"
)

//type pool interface {
//	newpool(newdb func()*gorm.DB,size int) *sqlpool
//	get() (db *gorm.DB)
//	put(db *gorm.DB)
//}

type sqlpool struct {
	new func() *gorm.DB
	db  []*gorm.DB
	sync.Mutex
}

func newpool(newdb func() *gorm.DB, size int) *sqlpool {
	return &sqlpool{newdb, make([]*gorm.DB, 0, size), sync.Mutex{}}
}

func (s *sqlpool) get() (db *gorm.DB) {
	s.Lock()
	defer s.Unlock()
	//log.Printf("before len:%d", len(s.db))
	if len(s.db) > 0 {
		db = s.db[len(s.db)-1]
		s.db = s.db[:len(s.db)-1]
	} else {
		db = s.new()
	}
	//log.Printf("after len:%d", len(s.db))
	return db
}

func (s *sqlpool) put(db *gorm.DB) {
	s.Lock()
	defer s.Unlock()
	if len(s.db) < cap(s.db) {
		s.db = append(s.db, db)
	} else {
		db.Close()
	}
}
