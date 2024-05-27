package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"reflect"
	"regexp"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/antchfx/htmlquery"
	"github.com/chromedp/chromedp"
	"github.com/zachturing/crawler/collect"
	"github.com/zachturing/crawler/parse/doubangroup"
)

func main() {
	var workList []*collect.Request
	for i := 25; i <= 100; i += 25 {
		str := fmt.Sprintf("https://www.douban.com/group/szsh/discussion?start=%d", i)
		workList = append(workList, &collect.Request{
			URL:       str,
			ParseFunc: doubangroup.ParseURL,
		})
	}

	var f collect.Fetcher = &collect.BrowserFetch{
		Timeout: 3000 * time.Millisecond,
	}

	for len(workList) > 0 {
		items := workList
		workList = nil
		for _, item := range items {
			fmt.Printf("current url:%s\n", item.URL)
			body, err := f.Get(item.URL)
			if err != nil {
				log.Printf("fetch %s error: %v", item.URL, err)
				continue
			}

			res := item.ParseFunc(body)
			for _, item := range res.Items {
				fmt.Printf("%s\n", item.(string))
			}

			workList = append(workList, res.Requests...)
		}
	}
}

func createQuery(q interface{}) string {
	fmt.Println(reflect.ValueOf(q).Kind())
	if reflect.ValueOf(q).Kind() != reflect.Struct {
		return "xxx"
	}

	t := reflect.TypeOf(q).Name() // 获取结构体名 Student
	query := fmt.Sprintf("insert into %s values(", t)
	v := reflect.ValueOf(q)

	for i := 0; i < v.NumField(); i++ {
		switch v.Field(i).Kind() {
		case reflect.Int:
			if i == 0 {
				query = fmt.Sprintf("%s%d", query, v.Field(i).Int())
			} else {
				query = fmt.Sprintf("%s,%d", query, v.Field(i).Int())
			}
		case reflect.String:
			if i == 0 {
				query = fmt.Sprintf("%s\"%s\"", query, v.Field(i).String())
			} else {
				query = fmt.Sprintf("%s,\\'%s\\'", query, v.Field(i).String())
			}
		}
	}
	query = fmt.Sprintf("%s)", query)
	fmt.Println(query)
	return query
}

func fetch() {
	url := "https://book.douban.com/subject/1007305/"
	var f collect.Fetcher = &collect.BrowserFetch{}
	body, err := f.Get(url)
	if err != nil {
		fmt.Printf("read content failed, err:%v\n", err)
		return
	}
	fmt.Println(string(body))
	goQueryParse(body)
	// reParse(body)
	// xPathParse(body)
}

func chrome() {
	// 1、创建谷歌浏览器实例
	ctx, cancel := chromedp.NewContext(
		context.Background())
	defer cancel()

	// 2、设置context超时时间
	ctx, cancel = context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	// 3、爬取页面，等待某一个元素出现,接着模拟鼠标点击，最后获取数据
	var example string
	err := chromedp.Run(ctx,
		chromedp.Navigate(`https://pkg.go.dev/time`),           // 爬取的网站
		chromedp.WaitVisible(`body > footer`),                  // 等待当前标签可见
		chromedp.Click(`#example-After`, chromedp.NodeVisible), // 模拟对某一个标签的点击事件
		chromedp.Value(`#example-After textarea`, &example),
	)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Go's time.After example:\\n%s", example)
}

func goQueryParse(body []byte) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		fmt.Printf("read content failed:%v\n", err)
	}
	doc.Find(`div.ant-col.ant-col-6 h2`).Each(func(i int, selection *goquery.Selection) {
		title := selection.Text()
		fmt.Printf("review:%d %s\n", i, title)
	})
}

func xPathParse(body []byte) {
	doc, err := htmlquery.Parse(bytes.NewReader(body))
	if err != nil {
		fmt.Printf("htmlquery parse failed:%v\n", err)
	}

	nodes := htmlquery.Find(doc, `//div[@class="ant-col ant-col-6"]//h2`)
	for i, node := range nodes {
		fmt.Printf("fetch card:%v, %v\n", i, htmlquery.InnerText(node))
	}
}

var titleRe = regexp.MustCompile(`<div class="small_[\s\S]*?<h2>([\s\S]*?)</h2>`)

func reParse(body []byte) {
	matches := titleRe.FindAllSubmatch(body, -1)
	for _, m := range matches {
		fmt.Printf("fetch title:%v\n", string(m[1]))
	}
	fmt.Println("done.")
}
