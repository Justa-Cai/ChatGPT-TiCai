package main

import (
	"strconv"
	"strings"
)

type PageMainData struct {
	// 开奖日期
	DrawDate string `selector:"span" gorm:"column:draw_date" gorm:"primaryKey"`
	// 开奖期数
	DrawNumber string `selector:"a" attr:"title" gorm:"column:draw_number"`
	// 链接
	Link string `selector:"a" attr:"href" gorm:"column:link"`
	//
	Title string `selector:"a" attr:"title" gorm:"column:title"`

	// 开奖号码
	Numbers string `gorm:"numbers"`
	Number  []int  `selector:"-" gorm:"-"`
}

func (that *PageMainData) Process() {
	that.DrawDate = strings.ReplaceAll(that.DrawDate, "(", "")
	that.DrawDate = strings.ReplaceAll(that.DrawDate, ")", "")
	that.DrawNumber = strings.ReplaceAll(that.DrawNumber, "中国体育彩票福建22选5第", "")
	that.DrawNumber = strings.ReplaceAll(that.DrawNumber, "中国体育彩票福建省22选5第", "")
	that.DrawNumber = strings.ReplaceAll(that.DrawNumber, "期开奖公告", "")
	that.Link = "http://www.fjtc.com.cn/" + that.Link
}

func (that *PageMainData) GetDrawNumber() int {
	n, _ := strconv.Atoi(that.DrawNumber)
	return n
}

func (that *PageMainData) GetNumber() []int {
	if len(that.Number) != 0 {
		return that.Number
	}

	for _, item := range strings.Split(that.Numbers, " ") {
		n, _ := strconv.Atoi(item)
		that.Number = append(that.Number, n)
	}
	return that.Number
}

func countEqualValues(arr1, arr2 []int) int {
	count := 0

	// 遍历第二个数组
	for _, num2 := range arr2 {
		// 在第一个数组中查找相等的数值
		for _, num1 := range arr1 {
			if num1 == num2 {
				count++
				break
			}
		}
	}

	return count
}

func (that *PageMainData) Check(cs []int) int {
	return countEqualValues(that.GetNumber(), cs)
}
