package service

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"sync"

	"github.com/axgle/mahonia"
	"github.com/opesun/goquery"
	"github.com/solenovex/web-tutor/exception"
	"github.com/solenovex/web-tutor/model"
	"github.com/solenovex/web-tutor/result"
)

var (
	domain = "https://www.qq.com/"
	reg    = `[a-zA-z]+://new.qq.com/omn/[0-9]*/[^\s]*.html`
	//放到全局防止panic
	reFirst = regexp.MustCompile(reg)
)

func GetUrl() [][]string {

	res, err := http.Get(domain)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer res.Body.Close()

	bytes, _ := ioutil.ReadAll(res.Body)
	html := string(bytes)
	//fmt.Println(html)

	data := reFirst.FindAllStringSubmatch(html, -1)
	fmt.Println(len(data))

	return data
}

func GetNews(data [][]string) error {
	var limit = make(chan int, len(data))
	lock := sync.Mutex{}

	for i := 0; i < len(data); i++ {
		limit <- 1
		go func(i int) {
			lock.Lock()
			date := string([]rune(data[i][0])[23:31])
			gid := string([]rune(data[i][0])[32:48])
			result, err := GetNewText(data[i][0], date, gid)
			if err != nil {
				log.Fatalln(err.Error())
			}
			fmt.Println(result)
			//校验数据库是否有这条新闻
			judge := CheckGid(gid)
			if judge != 1 && len(result.Title) != 0 && len(result.Text) != 0 {
				InsertNews(result)

			}
			judge = CheckDate(result.Date)
			if judge != 1 {
				//扔入布隆过滤器
				MyBloom.Add(result.Date)
			}
			lock.Unlock()
			<-limit
		}(i)
	}

	fmt.Println("结束存储！")
	return nil
}
func GetNewText(url string, date string, gid string) (model.News, error) {
	//获取信息
	p, err := goquery.ParseUrl(url)
	if err != nil {
		return model.News{}, exception.CreateNewException(result.ParaseUrlError, result.GetStautsString(result.ParaseUrlError))
	}
	ptitle := p.Find("h1").Text()
	//对title进行转码
	title := ConvertToString(ptitle, "gbk", "utf-8")
	pcontext := p.Find("p").Text()
	context := ConvertToString(pcontext, "gbk", "utf-8")
	data := model.News{Nid: 0, Gid: gid, Text: context, Title: title, Date: date}
	return data, nil

}

//转码函数，将腾讯新闻转码
func ConvertToString(src string, srcCode string, tagCode string) string {

	srcCoder := mahonia.NewDecoder(srcCode)

	srcResult := srcCoder.ConvertString(src)

	tagCoder := mahonia.NewDecoder(tagCode)

	_, cdata, _ := tagCoder.Translate([]byte(srcResult), true)

	result := string(cdata)

	return result

}
