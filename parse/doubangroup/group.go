package doubangroup

import (
	"github.com/zachturing/crawler/collect"
	"regexp"
)

const urlListRe = `(https://www.douban.com/group/topic/[0-9a-z]+/)"[^>]*>([^<]+)</a>`
const ContentRe = `<div class="topic-content">[\s\S]*?阳台[\s\S]*?<div class="aside">`

func ParseURL(contents []byte) collect.ParseResult {
	re := regexp.MustCompile(urlListRe)

	// 匹配所有帖子链接
	matches := re.FindAllSubmatch(contents, -1)
	result := collect.ParseResult{}

	for _, m := range matches {
		u := string(m[1])
		result.Requests = append(result.Requests, &collect.Request{
			URL: u,
			ParseFunc: func(c []byte) collect.ParseResult {
				return GetContent(c, u)
			},
		})
	}

	return result
}

func GetContent(contents []byte, url string) collect.ParseResult {
	// 匹配所有带阳台的帖子
	re := regexp.MustCompile(ContentRe)
	ok := re.Match(contents)
	if !ok {
		return collect.ParseResult{Items: []interface{}{}}
	}
	result := collect.ParseResult{Items: []interface{}{url}}
	return result
}
