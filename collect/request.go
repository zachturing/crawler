package collect

// Request 根据URL获取内容并解析
type Request struct {
	URL       string
	ParseFunc func([]byte) ParseResult
}

// ParseResult 解析结果
type ParseResult struct {
	Requests []*Request
	Items    []interface{}
}
