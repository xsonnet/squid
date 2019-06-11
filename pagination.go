package squid

import (
	"math"
	"net/http"
	"net/url"
	"strconv"
)

type Pagination struct {
	Request 	*http.Request
	Total   	int
	Size  		int
}

type PaginationLinks struct {
	First, Prev, Next, Last 	string
	Pages [						]string
}

func InitPagination(req *http.Request, total, size int) *Pagination {
	return &Pagination{
		Request: req,
		Total:   total,
		Size:  size,
	}
}

func (p *Pagination) Make() interface{} {
	queryParams := p.Request.URL.Query()
	//获取page
	page := queryParams.Get("page")
	if page == "" {
		page = "1"
	}
	number, _ := strconv.Atoi(page)
	if number == 0 {
		return nil
	}
	//计算总页数
	var totalPageNum = int(math.Ceil(float64(p.Total) / float64(p.Size)))
	links := PaginationLinks{}

	//首页和上一页链接
	if number > 1 {
		links.First = p.pageURL("1")
		links.Prev = p.pageURL(strconv.Itoa(number -1))
	} else {
		links.First = "#"
		links.Prev = "#"
	}

	//末页和下一页
	if number < totalPageNum {
		links.Last = p.pageURL(strconv.Itoa(totalPageNum))
		links.Next = p.pageURL(strconv.Itoa(number +1))
	} else {
		links.Last = "#"
		links.Next = "#"
	}

	//生成中间页码链接
	links.Pages = make([]string, 0, 10)
	startPos := number - 3
	endPos := number + 3
	if startPos < 1 {
		endPos = endPos + int(math.Abs(float64(startPos))) + 1
		startPos = 1
	}
	if endPos > totalPageNum {
		endPos = totalPageNum
	}
	for i := startPos; i <= endPos; i++ {
		links.Pages = append(links.Pages, p.pageURL(strconv.Itoa(i)))
	}

	return links
}

func (p *Pagination) pageURL(page string) string {
	u, _ := url.Parse(p.Request.URL.String())
	q := u.Query()
	q.Set("page", page)
	u.RawQuery = q.Encode()
	return u.String()
}