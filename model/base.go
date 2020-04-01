package model

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type BaseModel struct {
	db *gorm.DB
}

var err error

func (b *BaseModel) GetDbInstance() *gorm.DB {
	if b.db != nil {
		return b.db
	}
	b.db, err = gorm.Open("mysql", "root:123456@/panda?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}

	return b.db
}
