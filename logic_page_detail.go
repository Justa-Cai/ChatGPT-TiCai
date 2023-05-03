package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"strings"

	"github.com/gocolly/colly/v2"
)

func ticai_page_detail_process_num(strNumber string) []string {
	// 　　本期开奖号码： 01 02 04 05 07 18 26　　　　　 23
	strNumber = strings.ReplaceAll(strNumber, "本期开奖号码：", "")
	strNumber = strings.ReplaceAll(strNumber, "本期出球顺序： ", "")
	strNumber = strings.ReplaceAll(strNumber, "　　　　　", " ")
	strNumber = strings.ReplaceAll(strNumber, "  ", " ")
	strNumber = strings.ReplaceAll(strNumber, "\t", " ")
	re := regexp.MustCompile(`[\x{00A0}]`)
	strNumber = re.ReplaceAllString(strNumber, " ")
	strNumber = strings.TrimSpace(strNumber)

	// StringToHex(strNumber)
	ret := strings.SplitN(strNumber, " ", -1)
	var newRet []string
	for _, item := range ret {
		// StringToHex(item)
		if len(item) == 0 {
			continue
		}
		newRet = append(newRet, strings.TrimSpace(item))
	}

	return newRet
}

func ticai_page_detail(strUrl string) string {
	var ret []string
	// 创建一个新的Collector
	c := colly.NewCollector()

	// 设置User-Agent标头
	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.93 Safari/537.36"

	// 在请求之前执行的操作
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", randomUserAgent())
		// fmt.Println("Visiting", r.URL.String())
	})

	// 在访问响应之后执行的操作
	c.OnResponse(func(r *colly.Response) {
		// fmt.Println("Visited", r.Request.URL.String())
		ioutil.WriteFile("1.html", r.Body, 0777)
	})

	// 查找匹配的链接并访问它们
	// body > main > article > section.section3 > p:nth-child(4) > span
	bFound := false
	bBegin := false
	strNumber := ""
	_ = bBegin
	_ = strNumber

	c.OnHTML("span", func(e *colly.HTMLElement) {
		// fmt.Println(e.Text)
		if strings.Index(e.Text, "本期开奖号码") != -1 {
			bBegin = true
			strNumber = e.Text
			ret = ticai_page_detail_process_num(e.Text)
			// log.Println(ret)
			if len(ret) == 5 {
				log.Println(ret)
				bBegin = false
				bFound = true
			}
		}
	})

	if bFound == false {
		c.OnHTML("body > main > article > section.section3 > div:nth-child(3)", func(e *colly.HTMLElement) {
			// fmt.Println(e.Text)
			if strings.Index(e.Text, "本期开奖号码") != -1 {
				bBegin = true
				strNumber = e.Text
				ret = ticai_page_detail_process_num(e.Text)
				// log.Println(ret)
				if len(ret) == 5 {
					log.Println(ret)
					bBegin = false
					bFound = true
				}
			}
		})
	}

	if bFound == false {
		c.OnHTML("body > main > article > section.section3 > div:nth-child(4)", func(e *colly.HTMLElement) {
			// fmt.Println(e.Text)
			if strings.Index(e.Text, "本期开奖号码") != -1 {
				bBegin = true
				strNumber = e.Text
				ret = ticai_page_detail_process_num(e.Text)
				// log.Println(ret)
				if len(ret) == 5 {
					log.Println(ret)
					bBegin = false
					bFound = true
				}
			}
		})
	}

	if bFound == false {
		c.OnHTML("body > main > article > section.section3 > p:nth-child(4)", func(e *colly.HTMLElement) {
			// fmt.Println(e.Text)
			if strings.Index(e.Text, "本期开奖号码") != -1 {
				bBegin = true
				strNumber = e.Text
				ret = ticai_page_detail_process_num(e.Text)
				// log.Println(ret)
				if len(ret) == 5 {
					log.Println(ret)
					bBegin = false
					bFound = true
				}
			}
		})
	}

	if bFound == false {
		c.OnHTML("body > main > article > section.section3 > p", func(e *colly.HTMLElement) {
			// fmt.Println(e.Text)
			// fmt.Println(strings.Split(e.Text, "\n"))
			for _, item := range strings.Split(e.Text, "\n") {
				if strings.Index(e.Text, "本期开奖号码") == -1 {
					continue
				}
				ret = ticai_page_detail_process_num(item)
				if len(ret) == 5 {
					log.Println(ret)
					bBegin = false
					bFound = true
				}

			}
		})
	}

	// 开始访问起始URL
	err := c.Visit(strUrl)
	if err != nil {
		log.Println(err.Error())
		return ""
	}

	if len(ret) != 5 {
		log.Printf("url: %s parse failed len:%d\n", strUrl, len(ret))
	}

	return strings.Join(ret, " ")
}

func StringToHex(str string) {
	for _, c := range str {
		fmt.Printf("0x%X ", c)
	}
	fmt.Println()
}

func page_main_detail_process_msg(msg *Msg) {
	switch msg.Id {
	case MSG_PAGE_DETAIL_GET:
		ret := ticai_page_detail(msg.Param)
		data, _ := msg.Data.(*PageMainData)
		data.Numbers = ret
	}
}
