package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func load_func(nOffset int) {
	db_page_load("1.db")
	nTick := 0
	var builder strings.Builder
	builder.WriteString("扮演数学专家角色，规则从36个数据中选择7个号码，我按照\"[期数] 号码\"的格式,提供100组历史数据，如下：\n")

	for _, item := range g_PageMainData {
		nOffset--
		if nOffset >= 0 {
			continue
		}

		nTick++
		builder.WriteString(fmt.Sprintf("[%s] %s\n", item.DrawNumber, item.Numbers))
		fmt.Printf("[%s] %s\n", item.DrawNumber, item.Numbers)
		if nTick >= 100 {
			break
		}
	}

	builder.WriteString("请根据这些数据历史记录，输出新选择的5组数据，不要告诉我过程")
	ioutil.WriteFile("1.txt", []byte(builder.String()), 0777)
	// os.Exit(0)
}

func extractCombinationNumbers(input string) []string {
	re := regexp.MustCompile(`组合 \d+: ([\d\s]+)`)
	match := re.FindAllStringSubmatch(input, -1)

	numbers := make([]string, len(match))
	for i, m := range match {
		numbers[i] = m[1]
	}

	return numbers
}

func extractNumbers(line string) string {
	parts := strings.Split(line, "] ")
	if len(parts) != 2 {
		return ""
	}

	numbers := strings.TrimSpace(parts[1])
	return numbers
}

func extractNumbersSlice(data string) []string {
	lines := strings.Split(data, "\n")
	numbersSlice := make([]string, 0, len(lines))
	for _, line := range lines {
		numbers := extractNumbers(line)
		if numbers != "" {
			numbersSlice = append(numbersSlice, numbers)
		}
	}

	return numbersSlice
}

func convertToStringSlice(intSlice []int) string {
	strSlice := make([]string, len(intSlice))

	for i, num := range intSlice {
		strSlice[i] = fmt.Sprintf("%02d", num)
	}

	return strings.Join(strSlice, " ")
}
func appendToFile(filepath, content string) error {
	// 以追加模式打开文件，如果文件不存在则创建
	file, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// 写入内容到文件
	_, err = file.WriteString(content)
	if err != nil {
		return err
	}

	return nil
}
func check_func(nOffset int) {

	// 加载选择的数据
	bin, _ := ioutil.ReadFile("2.txt")
	strContent := string(bin)
	strChoose := extractNumbersSlice(strContent)
	var csa [][]int
	for _, item := range strChoose {
		var cs []int
		item = strings.ReplaceAll(item, "\n", "")
		for _, str := range strings.Split(item, " ") {
			n, _ := strconv.Atoi(str)
			cs = append(cs, n)
		}

		if len(cs) != 7 {
			continue
		}

		csa = append(csa, cs)
	}
	// fmt.Println(strChoose)
	fmt.Println(csa)
	// os.Exit(0)

	// 加载历史数据
	db_page_load("1.db")
	nTick := 0
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("测试 %d\n", nOffset))

	for _, cs := range csa {
		for _, item := range g_PageMainData {
			nOffset--
			if nOffset >= 0 {
				continue
			}

			nTick++
			if nTick <= 100 {
				continue
			}

			// item.Check()
			_ = item
			nSame := item.Check(cs)
			if nSame >= 4 {
				strContent = fmt.Sprintf("期数:%s 相同数量:%d 选:%s 开奖:%s\n", item.DrawDate, nSame, convertToStringSlice(cs), convertToStringSlice(item.GetNumber()))
				builder.WriteString(strContent)
				fmt.Printf(strContent)
			}
			break
		}
	}
	// ioutil.WriteFile("2.result.txt", []byte(builder.String()), 0777)
	appendToFile("2.result.txt", builder.String())
	// os.Exit(0)
}

func readFileToString(filepath string) string {
	data, _ := ioutil.ReadFile(filepath)
	return string(data)
}

func main() {
	nTick := 0
	for {
		// 加载最近100期
		load_func(nTick)
		// chatpgt推荐号码5组
		chatgpt_call(readFileToString("1.txt"))
		// 计算下一期的匹配结果
		check_func(100 + nTick)
		nTick++
	}

	os.Exit(0)

	// 启动一个 HTTP 服务器来提供 pprof 数据
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	go start_http_server()
	time.Sleep(time.Second / 10)

	msg_loop_init()
	g_mq.Register(MSG_PAGE_MAIN_BASE, page_main_data_process_msg)
	g_mq.Register(MSG_PAGE_DETAIL_BASE, page_main_detail_process_msg)
	g_mq.Add(MSG_PAGE_MAIN_GET, "http://www.fjtc.com.cn/36x7xq/index.html", nil)
	// var data PageMainData
	// g_mq.Add(MSG_PAGE_DETAIL_GET, "http://www.fjtc.com.cn//36x7xq/442909.html", &data)
	msg_loop_run()

	// writeCSV(data, "金山小.csv")
}
