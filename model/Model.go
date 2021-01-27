package model

import (
	"github.com/Gre-Z/common/jtime"
)

type Model struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt jtime.JsonTime
	UpdatedAt jtime.JsonTime
	DeletedAt jtime.JsonTime `sql:"index"`
}