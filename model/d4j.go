package model

import (
	"fmt"
	"github.com/gocolly/colly"
	"log"
	"regexp"
	"time"
)

type D4jModel struct {
	Book
	collector *colly.Collector
}

var book *Book

func (d *D4jModel) getInstance() *colly.Collector {
	if d.collector != nil {
		return d.collector
	}
	book = new(Book)
	d.collector = colly.NewCollector(
		//colly.Debugger(&debug.LogDebugger{}),
		colly.UserAgent("Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/78.0.3904.87 Mobile Safari/537.36"),
		colly.MaxDepth(3),
		//colly.Async(true),
		colly.URLFilters(
			regexp.MustCompile("d4j\\.cn"),
		),
	)
	// Limit the number of threads started by colly to two
	// when visiting links which domains' matches "*httpbin.*" glob
	d.collector.Limit(&colly.LimitRule{
		DomainGlob:  "*www.d4j.cn.*",
		Parallelism: 2,
		RandomDelay: 5 * time.Second,
	})

	// authenticate
	err := d.collector.Post("https://www.d4j.cn/wp-login.php", map[string]string{"log": "136911578@qq.com", "pwd": "denet789", "rememberme": "forever", "wp-submit": "登录", "redirect_to": "http://www.d4j.cn", "testcookie": "1"})
	if err != nil {
		log.Fatal(err)
	}

	// Set error handler
	d.collector.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})
	return d.collector
}

func (d *D4jModel) List() {

	// Before making a request print "Visiting ..."
	d.getInstance().OnRequest(func(r *colly.Request) {
		reg := regexp.MustCompile(`d4j\.cn/\d*\.html`)
		result := reg.FindAllStringSubmatch(r.URL.String(), -1)
		if len(result) > 0 {
			//çresultStr := string(result[0][1])
			fmt.Println("Visiting", result[0])
		}

	})
	// start scraping
	d.getInstance().Visit("https://www.d4j.cn")
	//d.getInstance().Wait()
}

func (d *D4jModel) Detail() {
	url := "https://www.d4j.cn/16107.html"

	d.getInstance().OnHTML("#container", func(e *colly.HTMLElement) {
		dom := e.DOM
		title := dom.Find(".kratos-entry-title").Text() //标题
		//fmt.Println("title", title)
		image, _ := dom.Find(".kratos-post-content img").Attr("src") //封面
		//fmt.Println("image", image)
		introduce := dom.Find(".title-h2").Next().Text() //介绍
		//fmt.Println("introduce", introduce)
		downloadUrl, _ := dom.Find(".downbtn").Attr("href") //站内下载链接

		book = &Book{
			Title:       title,
			Introduce:   introduce,
			Image:       image,
			OriginalUrl: url,
			Catid:       1,
			Userid:      1,
			Username:    "admin",
			AddTime:     time.Now().Unix(),
			Createtime:  time.Now().Unix(),
			Status:      1,
		}
		d.getInstance().Visit(e.Request.AbsoluteURL(downloadUrl))
	})

	d.getInstance().OnHTML(".wrap", func(e *colly.HTMLElement) {
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

		d.GetDbInstance().NewRecord(book)
		d.GetDbInstance().Create(&book)
	})
	// Before making a request print "Visiting ..."
	d.getInstance().OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})
	d.getInstance().Visit(url)

}
