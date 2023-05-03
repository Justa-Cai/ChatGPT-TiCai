package main

import (
	"log"
	"os"
	"strings"

	"github.com/gocolly/colly/v2"
)

func ticai_page_main(strUrl string) []*PageMainData {
	var data []*PageMainData
	// 创建一个新的Collector
	c := colly.NewCollector()

	// 设置User-Agent标头
	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.93 Safari/537.36"

	// 在请求之前执行的操作
	c.OnRequest(func(r *colly.Request) {
		// r.Headers.Set("User-Agent", randomUserAgent())
		// fmt.Println("Visiting", r.URL.String())
	})

	// 在访问响应之后执行的操作
	c.OnResponse(func(r *colly.Response) {
		// fmt.Println("Visited", r.Request.URL.String())
		// ioutil.WriteFile("1.html", r.Body, 0777)
	})

	// 查找匹配的链接并访问它们
	c.OnHTML("div.con", func(e *colly.HTMLElement) {
		e.ForEach("li.clearfix", func(i int, e *colly.HTMLElement) {
			var main *PageMainData = new(PageMainData)
			err := e.Unmarshal(main)
			if err != nil {
				log.Println(err.Error())
				return
			}
			main.Process()
			if main.GetDrawNumber() <= 18030 {
				return
			}

			data = append(data, main)
			// log.Println(main)
		})

	})

	// 如果有下一页，需要继续访问
	// bGetNext := false
	c.OnHTML("body > main > article > div.fy > div > a", func(e *colly.HTMLElement) {
		if strings.Index(e.Text, "下一页") != -1 {
			g_mq.Add(MSG_PAGE_MAIN_GET, "http://www.fjtc.com.cn/36x7xq/"+e.Attr("href"), nil)
			// bGetNext = true
			return
		}
	})

	// if bGetNext == false {
	// 	g_mq.Add(MSG_PAGE_MAIN_SAVE, "", nil)
	// }

	c.OnError(func(r *colly.Response, err error) {
		// 错误处理逻辑
		log.Printf("请求发生错误: %s\n", err)
	})

	// 开始访问起始URL
	err := c.Visit(strUrl)
	if err != nil {
		log.Println(err.Error())
		return nil
	}

	return data
}

var g_PageMainData []*PageMainData
var g_nIdelTick = 0

func page_main_data_process_msg(msg *Msg) {
	switch msg.Id {
	case MSG_PAGE_MAIN_GET:
		ret := ticai_page_main(msg.Param)
		for _, item := range ret {
			g_mq.Add(MSG_PAGE_DETAIL_GET, item.Link, item)
		}
		g_PageMainData = append(g_PageMainData, ret...)

		// g_mq.Add(MSG_PAGE_MAIN_SAVE, "", nil)
		g_nIdelTick = 0

	case MSG_PAGE_MAIN_SAVE:
		// 等待所有线程退出
		// log.Printf("wait finish...")
		// g_mq.wg.Wait()
		// 保存到文件
		log.Printf("save db...")
		db_page_save("1.db")
		os.Exit(0)

	case MSG_PAGE_MAIN_IDLE:
		g_nIdelTick++
		if g_nIdelTick == 5 {
			g_mq.Add(MSG_PAGE_MAIN_SAVE, "", nil)
		}
	}

}
