package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type myStock struct {
	index           int
	index_action    int
	action          string //buy sell
	action_name     string //訊號
	datetime        string
	price           int //價格
	lots            int //口數
	balance         int
	balance_p       float32
	balance_final   int
	balance_final_p float32
	balance_max     int
	balance_max_p   float32
	balance_min     int
	balance_min_p   float32
}

//交易編號,委託單編號,類型,訊號,成交時間,價格,數量,獲利,獲利(%),累積獲利,累積獲利(%),最大可能獲利,最大可能獲利(%),最大可能虧損,最大可能虧損(%)
//Date,Time,Open,High,Low,Close,TotalVolume
type kLine struct {
	date string
	time string
	o    int
	h    int
	l    int
	c    int
	v    int
}

type Buy string
type Sell string

var g_myStock_array = []myStock{}
var g_kLine_array = []kLine{}
var total = 0
var lots = 0 //初始口數
var i = -1   //天數

//以下可調整參數
var max_lots = 1 //做多最大口數
var min_lots = 0 //做空最大口數
var buyPrice = 0 //進場價格(結算用)
var path = "public/data/"
var k_file = "0845_300min"

//以上可調整參數

var finalName = k_file + "_action.csv"
var _ = os.Remove(path + finalName)
var finalFile, err = os.OpenFile(path+finalName, os.O_WRONLY|os.O_CREATE, os.FileMode(0600))

func main() {
	min_lots *= -1
	readCSV(path + k_file + ".csv")
}

func readCSV(kLine_fileName string) {
	file, err := os.Open(kLine_fileName)
	check(err)

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {

		if i == -1 { //跳過第一行標題
			i++
			_, err = fmt.Fprintln(finalFile, "交易編號,委託單編號,類型,訊號,成交時間,價格,數量,獲利,獲利(%),累積獲利,累積獲利(%),最大可能獲利,最大可能獲利(%),最大可能虧損,最大可能虧損(%)")
			continue
		}

		sli := strings.Split(scanner.Text(), ",")
		data := kLine{}
		data.date = sli[0]
		data.time = sli[1]
		data.o, _ = strconv.Atoi(sli[2])
		data.h, _ = strconv.Atoi(sli[3])
		data.l, _ = strconv.Atoi(sli[4])
		data.c, _ = strconv.Atoi(sli[5])
		data.v, _ = strconv.Atoi(sli[6])

		g_kLine_array = append(g_kLine_array, data)

		Strategy(i, lots, buyPrice) //策略判斷
		i++
	}

	_, err = fmt.Fprintln(finalFile, total)
	check(err)
	err = file.Close()
	check(err)
	err = finalFile.Close()
	check(err)

	if scanner.Err() != nil {
		log.Fatal(scanner.Err())
	}

}
func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func Strategy(i int, lots int, buyPrice int) {
	kk := g_kLine_array

	if len(kk) > 10 {
		if kk[i-3].c < kk[i-2].c && kk[i-2].c < kk[i-1].c && kk[i-1].c < kk[i].c { //連漲四天
			// buy = true
			Buy("buy_test").Lots(2).Price(1).nextBar()
		}

		if kk[i-3].c > kk[i-2].c && kk[i-2].c > kk[i-1].c && kk[i-1].c > kk[i].c { //連跌四天
			// sell = true
			Sell("sell_test").Lots(2).Price(1).nextBar()
		}
	}

}

func (s Sell) Price(price int) Sell {
	// fmt.Println("price", s, price)
	return s
}
func (s Sell) Lots(lots int) Sell {
	// fmt.Println("lots", s, lots)
	return s
}
func (s Sell) nextBar() {
	// fmt.Println("next", s)
	action("sell", string(s))
}
func (b Buy) Price(price int) Buy {
	// fmt.Println("price", s, price)
	return b
}
func (b Buy) Lots(lots int) Buy {
	// fmt.Println("lots", s, lots)
	return b
}
func (b Buy) nextBar() {
	// fmt.Println("next", s)
	action("buy", string(b))
}

func action(action string, action_name string) {
	is_action := false
	myStock := myStock{}
	myStock.action = action
	myStock.action_name = action_name

	if 0 <= lots && lots < max_lots && action == "buy" {
		myStock.price = g_kLine_array[i].c
		myStock.lots = 1
		myStock.balance = 0
		buyPrice = g_kLine_array[i].c
		lots++
		is_action = true
	}
	fmt.Println(lots, max_lots, action)
	if 0 < lots && lots <= max_lots && action == "sell" {
		myStock.action = "sell"
		myStock.price = g_kLine_array[i].c
		myStock.lots = 1
		myStock.balance = g_kLine_array[i].c - buyPrice
		lots--
		total += myStock.balance
		is_action = true
	}
	if is_action { //執行動作
		myStock.datetime = g_kLine_array[i].date + " " + g_kLine_array[i].time
		dataToCSV := fmt.Sprint(myStock.index) + ","
		dataToCSV += fmt.Sprint(myStock.index_action) + ","
		dataToCSV += myStock.action + ","
		dataToCSV += myStock.action_name + ","
		dataToCSV += myStock.datetime + ","
		dataToCSV += fmt.Sprint(myStock.price) + ","
		dataToCSV += fmt.Sprint(myStock.lots) + ","
		dataToCSV += fmt.Sprint(myStock.balance) + ","
		dataToCSV += fmt.Sprint(myStock.balance_p) + ","
		dataToCSV += fmt.Sprint(myStock.balance_final) + ","
		dataToCSV += fmt.Sprint(myStock.balance_final_p) + ","
		dataToCSV += fmt.Sprint(myStock.balance_max) + ","
		dataToCSV += fmt.Sprint(myStock.balance_max_p) + ","
		dataToCSV += fmt.Sprint(myStock.balance_min) + ","
		dataToCSV += fmt.Sprint(myStock.balance_min_p)

		_, err := fmt.Fprintln(finalFile, dataToCSV)
		check(err)
	}
}
