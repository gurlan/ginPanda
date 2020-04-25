package model

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"os"
	"sync"
)

type BaseModel struct {
}

var db *gorm.DB
var onceDb sync.Once
var err error

func (BaseModel) GetDbInstance() *gorm.DB {
	password := os.Getenv("DbPassword")
	if len(password) < 1 {
		password = "123456"
	}
	onceDb.Do(func() {
		db, err = gorm.Open("mysql", "root:"+password+"@tcp(localhost:3306)/panda?charset=utf8&parseTime=True")
		if err != nil {
			panic(err)
		}
	})
	return db
}
