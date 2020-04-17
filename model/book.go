package model

import "github.com/gocolly/colly"

var book Book
var nowPage = 1
var collector *colly.Collector
type Book struct {
	ID            int64 `gorm:"primary_key"`
	Title         string
	Author        string  `gorm:default:"-"`
	PressTime     string  `gorm:default:"-"`
	Press         string  `gorm:default:"-"`
	Score         float32 `gorm:default:0.0`
	Introduce     string  `gorm:default:"-"`
	Image         string  `gorm:default:""`
	BaiduUrl      string  `gorm:default:"-"`
	BaiduPassword string  `gorm:default:"-"`
	OriginalUrl   string  `gorm:default:"-"`
	TagId         string  `gorm:default:0`
	TagName       string  `gorm:default:"-"`
	Catid         int64   `gorm:default:"1"`
	Userid        int64   `gorm:default:"1"`
	Username      string  `gorm:default:"admin"`
	AddTime       int64
	Createtime    int64
	Status        int64 `gorm:default:"1"`
}

func (Book) TableName() string {
	return "book"
}
