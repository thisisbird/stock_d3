package main
import (
	. "StrategyMain"
)

func main(){
	Max_lots = 2 //做多最大口數
	Min_lots = 0 //做空最大口數
	K_file = "0845_300min" //來源檔案
	FinalName = K_file + "_action.csv"   //依時間排序
	FinalName2 = K_file + "_action2.csv" //進出場對應排序

	a := SquareSum(2)
	b := a(3)
	c := b(4)
	panic(c)

	Main(StrategyA)
}


func StrategyA() {
		if C(3)< C(2) && C(2) < C(1) && C(1) < C(0){ //連漲四k
			Buy("buy_test").Lots(1).Price(1).NextBar()
		}

		if C(3) > C(2) && C(2) > C(1) && C(1) > C(0){ //連跌四k
			Sell("sell_test").Lots(2).Price(1).NextBar()
		}

}
