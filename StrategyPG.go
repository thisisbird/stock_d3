package main
import (
	. "StrategyMain"
)

func StrategyA() {
	kk := G_kLine_array
	if len(kk) > 10 {
		if kk[K-3].C < kk[K-2].C && kk[K-2].C < kk[K-1].C && kk[K-1].C < kk[K].C { //連漲四k
			Buy("buy_test").Lots(1).Price(1).NextBar()
		}

		if kk[K-3].C > kk[K-2].C && kk[K-2].C > kk[K-1].C && kk[K-1].C > kk[K].C { //連跌四k
			Sell("sell_test").Lots(2).Price(1).NextBar()
		}
	}

}
func main(){
	Main(StrategyA)
}