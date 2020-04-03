package model

import (
	"github.com/astaxie/goredis"
)

var redisClient goredis.Client

const BookUrls = "d4j_book_urls"
const HasSaveUrls = "d4j_has_save_urls"


func init() {
	redisClient.Addr = "127.0.0.1:6379"

}
func PutInQueen(url string) {
	err := redisClient.Lpush(BookUrls, []byte(url))
	if err != nil {
		panic(err)
	}
}
func RPopFromQueen() string {
	keys := make([]string, 1)
	keys = append(keys, BookUrls)
	_, val, err := redisClient.Brpop(keys, 30)
	if err != nil {
		return ""
	}
	return string(val)
}
func PutInSet(url string) {
	redisClient.Sadd(HasSaveUrls, []byte(url))
}
func HasInSet(url string) bool {
	flag, _ := redisClient.Sismember(HasSaveUrls, []byte(url))
	if flag {
		return true
	}
	return false
}
