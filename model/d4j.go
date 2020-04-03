package model

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/queue"
	"log"
	"net"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

type D4jModel struct {
	Book
}

var book Book
var nowPage = 1
var collector *colly.Collector

func (d *D4jModel) GetCollector() {

	collector = colly.NewCollector(
		//colly.Debugger(&debug.LogDebugger{}),
		colly.UserAgent("Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/78.0.3904.87 Mobile Safari/537.36"),
		colly.MaxDepth(3),
		//colly.Async(true),
		colly.URLFilters(
			regexp.MustCompile(`d4j.cn/wp-login.php$`),
			regexp.MustCompile(`d4j.cn$`),
			regexp.MustCompile(`d4j.cn/page/\d*`),
			regexp.MustCompile(`d4j.cn/\d*.html$`),
			regexp.MustCompile(`d4j.cn/download.php`),

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
		DomainGlob:  "*www.d4j.cn*",
		Parallelism: 5,
		//RandomDelay: 2 * time.Second,
	})

	// authenticate
	err := collector.Post("https://www.d4j.cn/wp-login.php", map[string]string{"log": "136911578@qq.com", "pwd": "denet789", "rememberme": "forever", "wp-submit": "登录", "redirect_to": "http://www.d4j.cn", "testcookie": "1"})
	if err != nil {
		log.Fatal(err)
	}

	// Set error handler
	collector.OnError(func(r *colly.Response, err error) {
		log.Println(r.Request.URL.String())
		collector.Visit(r.Request.AbsoluteURL(r.Request.URL.String()))
		fmt.Println("Request URL Fail:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

}

func (d *D4jModel) List() {
	d.GetCollector()

	// On every a element which has href attribute call callback
	collector.OnHTML("#kratos-blog-post", func(e *colly.HTMLElement) {
		dom := e.DOM
		dom.Find(".kratos-entry-title-new a").Each(func(index int, ele *goquery.Selection) {
			link, _ := ele.Attr("href")
			fmt.Printf("Link found: %d ->%s\n", index, link)
			reg := regexp.MustCompile(`https://www.d4j.cn/\d*.html$`)
			result := reg.FindAllStringSubmatch(link, -1)
			if len(result) > 0 {
				fmt.Println("Inserting", result[0][0])
				PutInQueen(result[0][0])
			} else {
				fmt.Println("Passing", link)
			}
			if index >= 15 {
				nowPage += 1
				collector.Visit("https://www.d4j.cn/page/" + strconv.Itoa(nowPage))
			}
		})

	})

	// Before making a request print "Visiting ..."
	collector.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	collector.Visit("https://www.d4j.cn/page/" + strconv.Itoa(nowPage))
}

func (d *D4jModel) Detail() {

	d.GetCollector()
	// Before making a request print "Visiting ..."
	collector.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	collector.OnHTML("#container", func(e *colly.HTMLElement) {

		dom := e.DOM
		title := dom.Find(".kratos-entry-title").Text() //标题
		//fmt.Println("title", title)
		image, _ := dom.Find(".kratos-post-content img").Attr("src") //封面
		//fmt.Println("image", image)
		introduce := dom.Find(".title-h2").Next().Text() //介绍
		//fmt.Println("introduce", introduce)
		downloadUrl, _ := dom.Find(".downbtn").Attr("href") //站内下载链接

		book = Book{
			Title:       title,
			Introduce:   introduce,
			Image:       image,
			OriginalUrl: e.Request.URL.String(),
			Catid:       1,
			Userid:      1,
			Username:    "admin",
			AddTime:     time.Now().Unix(),
			Createtime:  time.Now().Unix(),
			Status:      1,
		}
		d.download(downloadUrl, book)

	})

	q, _ := queue.New(
		2,                                           // Number of consumer threads
		&queue.InMemoryQueueStorage{MaxSize: 10000}, // Use default queue storage
	)

	for {
		// Add URLs to the queue
		url := RPopFromQueen()
		if url == "" {
			break
		}
		q.AddURL(url)
	}
	// Consume URLs
	q.Run(collector)
}

func (d *D4jModel) download(url string, book Book) {
	d.GetCollector()
	collector.Visit(url)
	collector.OnHTML(".wrap", func(e *colly.HTMLElement) {

		dom := e.DOM
		baiduUrl, _ := dom.Find(".downfile").Eq(3).Find("a").Attr("href")
		baiduPassword := dom.Find(".plus_l").Find("li").Eq(3).Children().Text()

		authorStr := dom.Find(".plus_l").Find("li").Eq(2).Text()
		reg := regexp.MustCompile(`作者信息 ：【(.*?)】`)
		result := reg.FindAllStringSubmatch(authorStr, -1)
		author := "-"
		if len(result) > 0 {
			author = result[0][1]
		}

		book.Author = author
		book.BaiduUrl = baiduUrl
		book.BaiduPassword = baiduPassword

		db.NewRecord(book)
		db.Create(&book)
		book = Book{}
	})
}
