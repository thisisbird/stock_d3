package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	// "strings"
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

var total = 0
var path = "public/data/"
var k_file = "0845_300min"

func main() {
	readCSV(path + k_file + ".csv")
}

func readCSV(fileName string) {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	finalName := path + k_file + "_action.csv"
	scanner := bufio.NewScanner(file)
	os.Remove(finalName)
	
	options := os.O_WRONLY | os.O_CREATE //開啟檔案的選項
	finalFile, err := os.OpenFile(finalName, options, os.FileMode(0600))
	check(err)
	i := -1       //天數
	lots := 1     //總口數
	buyPrice := 0 //進場價格(結算用)
	kLine_array := []kLine{}
	myStock_array := []myStock{}
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
		// data.percentChange, _ = strconv.Atoi(sli[10])

		kLine_array = append(kLine_array, data)
		// if i < 4883 {//跳過前面的
		// 	i++
		// 	continue
		// }
		lots, buyPrice = Strategy(finalFile, kLine_array, myStock_array, i, lots, buyPrice) //策略判斷
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

func Strategy(finalFile *os.File, kk []kLine, myStocka []myStock, i int, lots int, buyPrice int) (int, int) {
	myStock := myStock{}
	buy := false
	sell := false
	action := false
	if len(kk) > 10 {
		if kk[i-3].c < kk[i-2].c && kk[i-2].c < kk[i-1].c && kk[i-1].c < kk[i].c { //連漲四天
			buy = true
		}

		if kk[i-3].c > kk[i-2].c && kk[i-2].c > kk[i-1].c && kk[i-1].c > kk[i].c { //連跌四天
			sell = true
		}
	}

	if buy && lots > 0 {
		myStock.action = "buy"
		myStock.price = kk[i].c
		myStock.lots = 1
		myStock.balance = 0
		buyPrice = kk[i].c
		action = true
		lots--
	}

	if sell && lots == 0 {
		myStock.action = "sell"
		myStock.price = kk[i].c
		myStock.lots = 1
		myStock.balance = kk[i].c - buyPrice
		lots++
		total += myStock.balance
		action = true
	}

	if action { //執行動作
		myStock.action_name = "action_name"
		myStock.datetime = kk[i].date + " " + kk[i].time
		// dataToCSV := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(myStock)), ","), "{}")
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
	return lots, buyPrice
}
