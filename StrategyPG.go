package main

import (
	. "StrategyMain"
	"fmt"
)

func main() {
	Max_lots = 1                         //做多最大口數
	Min_lots = 0                         //做空最大口數
	K_file = "0845_300min"               //來源檔案
	FinalName = K_file + "_action.csv"   //依時間排序
	FinalName2 = K_file + "_action2.csv" //進出場對應排序
	fmt.Println(macd)
	Main(StrategyA, setReady)
}

var ems20 []float32
var ems12 []float32
var ems26 []float32
var dif []float32
var macd []float32
var d_m []float32

func setReady() {
	SET(&ems20, EMA(C, 20))
	SET(&ems12, EMA(C, 12))
	SET(&ems26, EMA(C, 26))
	SET(&dif, ems12[0]-ems26[0])
	SET(&macd, EMA(dif, 26))
	SET(&d_m, dif[0]-macd[0])
}

func StrategyA() {
	// cond2 := d_m[0] > 0
	cond3 := d_m[0] > d_m[1]
	// cond4 := dif[0] > 0 &&  macd[0] > 0
	cond7 := ems20[0] > ems20[1] && ems20[1] > ems20[2]

	if cond3 && cond7{
		Buy("MacdLE").Lots(1).NextBar().At("market")
	}
	
	// if C[3] < C[2] && C[2] < C[1] && C[1] < C[0] { //連漲四k
	// 	Buy("buy_test").Lots(1).NextBar().At("market")
	// }

	if C[3] > C[2] && C[2] > C[1] && C[1] > C[0] { //連跌四k
		// Sell("sell_test").Lots(1).NextBar().At("limit", EntryPrice+50)
		Sell("sell_test").Lots(1).NextBar().At("market")
	}

}
