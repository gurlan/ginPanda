package model

import (
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"io"
	"log"
	"net"
	"net/http"
	"reflect"
	"regexp"
	"strings"
	"sync"
	"time"
)

type ItPandaModel struct {
	Book
	Mutex sync.Mutex
}

func (it *ItPandaModel) GetCollector() {

	collector = colly.NewCollector(
		//colly.Debugger(&debug.LogDebugger{}),
		colly.UserAgent("Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/78.0.3904.87 Mobile Safari/537.36"),
		colly.MaxDepth(3),
		//colly.Async(true),
		colly.URLFilters(

			regexp.MustCompile(`www.itpanda.net`),


		),
	)

	collector.WithTransport(&http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   120 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       100 * time.Second,
		TLSHandshakeTimeout:   60 * time.Second,
		ExpectContinueTimeout: 100 * time.Second,
	})
	// Limit the number of threads started by colly to two
	// when visiting links which domains' matches "*httpbin.*" glob
	collector.Limit(&colly.LimitRule{
		DomainGlob:  "*www.itpanda.net*",
		Parallelism: 5,
		//	RandomDelay: 2 * time.Second,
	})

	// Set error handler
	collector.OnError(func(r *colly.Response, err error) {
		log.Println(r.Request.URL.String())
		collector.Visit(r.Request.AbsoluteURL(r.Request.URL.String()))
		fmt.Println("Request URL Fail:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

}

func (it *ItPandaModel) List() {
	it.GetCollector()
	var book Book
	// On every a element which has href attribute call callback
	collector.OnHTML(".nav.flex-column", func(e *colly.HTMLElement) {
		dom := e.DOM
		dom.Find("a").Each(func(index int, ele *goquery.Selection) {
			link, _ := ele.Attr("href")
			fmt.Printf("Link found: %d ->%s\n", index, link)

			collector.Visit(e.Request.AbsoluteURL(link))

		})
	})

	collector.OnHTML(".list-unstyled", func(e *colly.HTMLElement) {
		dom := e.DOM
		dom.Find("a").Each(func(index int, ele *goquery.Selection) {
			link, _ := ele.Attr("href")
			fmt.Printf("Link found: %d ->%s\n", index, link)

			collector.Visit(e.Request.AbsoluteURL(link))

		})
	})
	collector.OnHTML(".justify-content-center", func(e *colly.HTMLElement) {
		dom := e.DOM
		title := dom.Find(".my-3").Eq(0).Text()                                         //标题
		author := dom.Find(".media-body").Find("p").Eq(0).Contents().Not("span").Text() //标题
		image, _ := dom.Find(".media img").Attr("src")                                  //封面
		introduce := dom.Find(".text-danger").Eq(0).Next().Text()                       //介绍
		OriginalUrl := e.Request.URL.String()

		reg := regexp.MustCompile(`https://www.itpanda.net/book/\d+$`)
		result := reg.FindAllStringSubmatch(OriginalUrl, -1)
		book = Book{}
		if len(result) > 0 {
			book = Book{
				Title:       title,
				Author:      author,
				Image:       image,
				Introduce:   introduce,
				OriginalUrl: OriginalUrl,
				Catid:       1,
				Userid:      1,
				Username:    "admin",
				AddTime:     time.Now().Unix(),
				Createtime:  time.Now().Unix(),
				Status:      1,
			}

			downloadUrl, _ := dom.Find(".mr-2").Last().Attr("href") //站内下载链接

			log.Println("downloadUrl:" + downloadUrl)
			it.Mutex.Lock()
			it.insert(book, e.Request.AbsoluteURL(downloadUrl))
		}

	})

	// Before making a request print "Visiting ..."
	collector.OnRequest(func(r *colly.Request) {
		log.Println("Visiting:" + r.URL.String())
	})

	collector.Visit("https://www.itpanda.net/")

}
func (it *ItPandaModel) insert(book Book, downLoadUrl string) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
			it.insert(book,downLoadUrl)
		}
	}()
	html := it.Get(downLoadUrl)
//	log.Println(html)
	dom, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		log.Println(err.Error())
	}

	baiduUrl, _ := dom.Find(".alert-success").Find("a").Eq(1).Attr("href")
	baiduPasswordStr := dom.Find(".alert-success").Find("p").Contents().Not("a").Text()
	baiduPassword := ""
	if len(baiduPasswordStr) > 4 {
		baiduPassword = baiduPasswordStr[len(baiduPasswordStr)-4:]
	}
	book.BaiduUrl = baiduUrl
	book.BaiduPassword = baiduPassword
	SmartPrint(book)
	db.NewRecord(book)
	db.Create(&book)

	it.Mutex.Unlock()
}

// 发送GET请求
// url：         请求地址
// response：    请求返回的内容
func (it *ItPandaModel) Get(url string) string {


	// 超时时间：5秒
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	var buffer [512]byte
	result := bytes.NewBuffer(nil)
	for {
		n, err := resp.Body.Read(buffer[0:])
		result.Write(buffer[0:n])
		if err != nil && err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
	}

	return result.String()
}
func SmartPrint(i interface{}) {
	var kv = make(map[string]interface{})
	vValue := reflect.ValueOf(i)
	vType := reflect.TypeOf(i)
	for i := 0; i < vValue.NumField(); i++ {
		kv[vType.Field(i).Name] = vValue.Field(i)
	}
	fmt.Println("获取到数据:")
	for k, v := range kv {
		fmt.Print(k)
		fmt.Print(":")
		fmt.Print(v)
		fmt.Println()
	}
}
