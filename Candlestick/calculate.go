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
	MA5           float32
	MA10          float32
	MA20          float32
	MA60          float32
}

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
	options := os.O_WRONLY | os.O_APPEND | os.O_CREATE //開啟檔案的選項
	file2, err := os.OpenFile("ma.csv", options, os.FileMode(0600))
	check(err)
	i := -1
	array := []myStock{}

	for scanner.Scan() {
		if i == -1 { //跳過第一行標題
			i++
			_, err = fmt.Fprintln(file2, "交易日期,開盤價,最高價,最低價,收盤價,成交量,漲跌價,漲跌%,5MA,10MA,20MA,60MA")
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

		data.MA5 = MA(array, 5, data.close)
		data.MA10 = MA(array, 10, data.close)
		data.MA20 = MA(array, 20, data.close)
		data.MA60 = MA(array, 60, data.close)

		array = append(array, data)
		dataToCSV := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(array[i])), ", "), "{}")

		_, err = fmt.Fprintln(file2, dataToCSV)
		i++
	}
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

func MA(array []myStock, count int, close int) float32 {
	count--
	total := len(array)
	first := total - count
	sum := 0

	if first >= 0 {
		array = array[first:total]
		for i := 0; i < count; i++ {
			sum += array[i].close
		}
		return float32(sum+close) / float32(count+1)
	}

	for i := 0; i < total; i++ {
		sum += array[i].close
	}
	return float32(sum+close) / float32(total+1)

}
