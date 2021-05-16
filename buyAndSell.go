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
	date          string
	open          int
	high          int
	low           int
	close         int
	volume        int
	change        int
	percentChange float32
	action        string //buy sell
	price         int    //買賣價格
	lots          int    //口數
	balance       int
}

var total = 0

func main() {
	readCSV("new.csv")
}

//[0交易日期 1契約 2到期月份(週別) 3開盤價 4最高價 5最低價 6收盤價 7漲跌價 8漲跌% 9成交量 10結算價 11未沖銷契約數 12最後最佳買價 13最後最佳賣價 14歷史最高價 15歷史最低價 16是否因訊息面暫停交易 17交易時段 18價差對單式委託成交量]

func readCSV(fileName string) {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(file)
	options := os.O_WRONLY | os.O_CREATE //開啟檔案的選項
	file2, err := os.OpenFile("buyAndSell.csv", options, os.FileMode(0600))
	check(err)
	i := -1       //天數
	lots := 1     //總口數
	buyPrice := 0 //進場價格(結算用)
	array := []myStock{}
	for scanner.Scan() {
		if i == -1 { //跳過第一行標題
			i++
			_, err = fmt.Fprintln(file2, "交易日期,開盤價,最高價,最低價,收盤價,成交量,漲跌價,漲跌%,買賣,價格,口數,收益")
			continue
		}

		sli := strings.Split(scanner.Text(), ",")
		data := myStock{}
		data.date = sli[0]
		data.open, _ = strconv.Atoi(sli[3])
		data.high, _ = strconv.Atoi(sli[4])
		data.low, _ = strconv.Atoi(sli[5])
		data.close, _ = strconv.Atoi(sli[6])
		data.volume, _ = strconv.Atoi(sli[9])
		data.change, _ = strconv.Atoi(sli[7])
		// data.percentChange, _ = strconv.Atoi(sli[10])

		array = append(array, data)
		// if i < 4883 {//跳過前面的
		// 	i++
		// 	continue
		// }
		lots, buyPrice = Strategy(file2, array, i, lots, buyPrice) //策略判斷
		i++
	}

	_, err = fmt.Fprintln(file2, total)
	check(err)

	err = file.Close()
	check(err)
	err = file2.Close()
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

func Strategy(file2 *os.File, array []myStock, i int, lots int, buyPrice int) (int, int) {
	buy := false
	sell := false
	action := false
	if len(array) > 10 {
		if array[i-3].close < array[i-2].close && array[i-2].close < array[i-1].close && array[i-1].close < array[i].close { //連漲四天
			buy = true
		}

		if array[i-3].close > array[i-2].close && array[i-2].close > array[i-1].close && array[i-1].close > array[i].close { //連跌四天
			sell = true
		}
	}

	if buy && lots > 0 {
		array[i].action = "buy"
		array[i].price = array[i].close
		array[i].lots = 1
		array[i].balance = 0
		buyPrice = array[i].close
		action = true
		lots--
	}

	if sell && lots == 0 {
		array[i].action = "sell"
		array[i].price = array[i].close
		array[i].lots = 1
		array[i].balance = array[i].close - buyPrice
		lots++
		total += array[i].balance
		action = true
	}

	if action { //執行動作
		dataToCSV := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(array[i])), ", "), "{}")
		_, err := fmt.Fprintln(file2, dataToCSV)
		check(err)
	}
	return lots, buyPrice
}
