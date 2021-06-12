package StrategyMain

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type fn func()
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
	price_max       int //紀錄最高值
	price_min       int //紀錄最低值
}

//交易編號,委託單編號,類型,訊號,成交時間,價格,數量,獲利,獲利(%),累積獲利,累積獲利(%),最大可能獲利,最大可能獲利(%),最大可能虧損,最大可能虧損(%)
//Date,Time,Open,High,Low,Close,TotalVolume
type kLine struct {
	Date string
	Time string
	O    int
	H    int
	L    int
	C    int
	V    int
}

type Buy string
type Sell string

var C = []float32{}
var g_myStock_array = []myStock{}
var g_kLine_array = []kLine{}

var total = 0
var lots = 0 //初始口數
var K = -1   //天數

//以下可調整參數
var Max_lots = 2 //做多最大口數
var Min_lots = 0 //做空最大口數
var path = "public/data/"
var K_file = "0845_300min"

//以上可調整參數

var FinalName = K_file + "_action.csv"   //依時間排序
var FinalName2 = K_file + "_action2.csv" //進出場對應排序

var finalFile *os.File
var finalFile2 *os.File

// func Strategy() { //在其他檔案撰寫策略
// 	kk := g_kLine_array

// 	if len(kk) > 10 {
// 		if kk[K-3].C < kk[K-2].C && kk[K-2].C < kk[K-1].C && kk[K-1].C < kk[K].C { //連漲四k
// 			Buy("buy_test").Lots(1).Price(1).NextBar()
// 		}

// 		if kk[K-3].C > kk[K-2].C && kk[K-2].C > kk[K-1].C && kk[K-1].C > kk[K].C { //連跌四k
// 			Sell("sell_test").Lots(2).Price(1).NextBar()
// 		}
// 	}

// }

func Main(myfunction fn, setReady fn) {
	os.Remove(path + FinalName)
	os.Remove(path + FinalName2)
	finalFile, _ = os.OpenFile(path+FinalName, os.O_WRONLY|os.O_CREATE, os.FileMode(0600))
	finalFile2, _ = os.OpenFile(path+FinalName2, os.O_WRONLY|os.O_CREATE, os.FileMode(0600))
	Min_lots *= -1
	readCSV(path+K_file+".csv", myfunction, setReady)
}

func readCSV(kLine_fileName string, StrategyFunction fn, setReady fn) {
	file, err := os.Open(kLine_fileName)
	check(err)

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {

		if K == -1 { //跳過第一行標題
			K++
			fmt.Fprintln(finalFile, "交易編號,委託單編號,類型,訊號,成交時間,價格,數量,獲利,獲利(%),累積獲利,累積獲利(%),最大可能獲利,最大可能獲利(%),最大可能虧損,最大可能虧損(%)")
			fmt.Fprintln(finalFile2, "交易編號,委託單編號,類型,訊號,成交時間,價格,數量,獲利,獲利(%),累積獲利,累積獲利(%),最大可能獲利,最大可能獲利(%),最大可能虧損,最大可能虧損(%)")
			continue
		}

		sli := strings.Split(scanner.Text(), ",")
		data := kLine{}
		data.Date = sli[0]
		data.Time = sli[1]
		data.O, _ = strconv.Atoi(sli[2])
		data.H, _ = strconv.Atoi(sli[3])
		data.L, _ = strconv.Atoi(sli[4])
		data.C, _ = strconv.Atoi(sli[5])
		data.V, _ = strconv.Atoi(sli[6])
		g_kLine_array = prependKline(g_kLine_array, data)
		recordHighAndLow(data)
		if temp_ready { //隔天執行用
			action()
			temp_ready = false
		}

		SET(&C, float32(data.C))
		setReady() //設置變數

		// Strategy() //策略判斷
		if K > 26 {

			StrategyFunction()
		}
		K++
	}
	err = file.Close()
	check(err)
	err = finalFile.Close()
	check(err)
	err = finalFile2.Close()
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

var temp_action string
var temp_action_name string
var temp_balance_final = 0 //累積獲利價格
var temp_price int         //成交價格
var temp_lots int          //成交口數
var temp_ready = false     //明天是否執行

func (s Sell) Price(price int) Sell {
	temp_price = price
	return s
}
func (s Sell) Lots(lots int) Sell {
	temp_lots = lots
	return s
}
func (s Sell) NextBar() {
	temp_action_name = string(s)
	temp_action = "sell"
	temp_ready = true
	// action()
}

func (b Buy) Price(price int) Buy {
	temp_price = price
	return b
}
func (b Buy) Lots(lots int) Buy {
	temp_lots = lots
	return b
}
func (b Buy) NextBar() {
	temp_action_name = string(b)
	temp_action = "buy"
	temp_ready = true
}

func action() { //執行策略(紀錄)
	is_action := false
	myStock := myStock{}

	if 0 <= lots && (lots+temp_lots) <= Max_lots && temp_action == "buy" {
		myStock.price = g_kLine_array[0].O//當天/隔天開盤價
		// myStock.balance = 0
		lots += temp_lots
		is_action = true
	}
	if 0 <= (lots-temp_lots) && lots <= Max_lots && temp_action == "sell" {
		myStock.price = g_kLine_array[0].C//當天or隔天收盤價

		lots -= temp_lots

		if lots == 0 {
			total += myStock.balance
		}
		is_action = true
	}

	if is_action { //執行動作
		myStock.datetime = g_kLine_array[0].Date + " " + g_kLine_array[0].Time
		myStock.action = temp_action
		myStock.action_name = temp_action_name
		myStock.lots = temp_lots
		action1(myStock)
		action2(myStock)
	}

}
func action1(myStock myStock) { //依時間排序
	saveRow(finalFile, myStock)
}
func action2(myStock myStock) { //進出場對應排序
	count := len(g_myStock_array)
	if count == 0 {
		g_myStock_array = append(g_myStock_array, myStock)
		return
	}
	if count > 0 && g_myStock_array[0].action == myStock.action { //買賣別一樣就加入陣列
		g_myStock_array = append(g_myStock_array, myStock)
		return
	}
	if count > 0 && g_myStock_array[0].action != myStock.action { //買賣別不一樣

		if g_myStock_array[0].lots == myStock.lots { //口數比對

			myStock.balance = balance(g_myStock_array[0], myStock)
			myStock.balance_final = temp_balance_final
			myStock.balance_max = balanceMax(g_myStock_array[0], myStock)
			myStock.balance_min = balanceMin(g_myStock_array[0], myStock)

			saveRow(finalFile2, g_myStock_array[0])
			saveRow(finalFile2, myStock)
			g_myStock_array = g_myStock_array[1:]
			return
		}
		if g_myStock_array[0].lots > myStock.lots { //口數比對
			new_lots := g_myStock_array[0].lots - myStock.lots
			g_myStock_array[0].lots = myStock.lots

			myStock.balance = balance(g_myStock_array[0], myStock)
			myStock.balance_final = temp_balance_final
			myStock.balance_max = balanceMax(g_myStock_array[0], myStock)
			myStock.balance_min = balanceMin(g_myStock_array[0], myStock)
			saveRow(finalFile2, g_myStock_array[0])
			saveRow(finalFile2, myStock)

			g_myStock_array[0].lots = new_lots
			return
		}
		if g_myStock_array[0].lots < myStock.lots { //口數比對
			x := len(g_myStock_array)
			for x > 0 {
				new_lots := myStock.lots - g_myStock_array[0].lots
				myStock.lots = g_myStock_array[0].lots

				myStock.balance = balance(g_myStock_array[0], myStock)
				myStock.balance_final = temp_balance_final
				myStock.balance_max = balanceMax(g_myStock_array[0], myStock)
				myStock.balance_min = balanceMin(g_myStock_array[0], myStock)

				saveRow(finalFile2, g_myStock_array[0])
				saveRow(finalFile2, myStock)
				myStock.lots = new_lots
				g_myStock_array = g_myStock_array[1:]

				x = len(g_myStock_array)
			}
			return
		}

	}
}
func saveRow(file *os.File, myStock myStock) {
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

	_, err := fmt.Fprintln(file, dataToCSV)
	check(err)
}

func recordHighAndLow(k kLine) {
	if len(g_myStock_array) == 0 {
		return
	}
	for i := 0; i < len(g_myStock_array); i++ {
		g_myStock_array[i].price_max = max(g_myStock_array[i].price_max, k.H)
		g_myStock_array[i].price_min = min(g_myStock_array[i].price_min, k.L)
	}
}

func balance(entry myStock, exit myStock) int {
	balance := exit.lots*exit.price - entry.lots*entry.price
	temp_balance_final += balance
	return balance
}
func balanceMax(entry myStock, exit myStock) int {
	balance := entry.lots*entry.price_max - entry.lots*entry.price
	return balance
}
func balanceMin(entry myStock, exit myStock) int {
	balance := entry.lots*entry.price_min - entry.lots*entry.price
	return balance
}

func max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a int, b int) int {
	if a == 0 {
		return b
	}
	if a < b {
		return a
	}
	return b
}

func prependKline(x []kLine, y kLine) []kLine {
	x = append(x, y)
	copy(x[1:], x)
	x[0] = y
	return x
}

func SET(x *[]float32, y float32) {
	*x = append(*x, 0.0)
	copy((*x)[1:], *x)
	(*x)[0] = y
}

func EMA(array []float32, x int) float32 {
	if len(array) < x {
		x = len(array)
	}
	total := (x + 1) * x / 2
	var ema float32
	for i := 0; i < x; i++ {
		ema += array[i] * float32(x-i)
	}
	return ema / float32(total)
}
