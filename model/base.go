package model

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"sync"
)

type BaseModel struct {
}

var db *gorm.DB
var onceDb sync.Once
var err error

func (BaseModel) GetDbInstance() *gorm.DB {
	onceDb.Do(func() {
		db, err = gorm.Open("mysql", "root:123456@/panda?charset=utf8&parseTime=True&loc=Local")
		if err != nil {
			panic(err)
		}
	})
	return db
}
