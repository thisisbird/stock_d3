package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	fileName := "TXF1-分鐘-成交價.txt"
	fileName = "o_data/kevin/" + fileName
	readCSV(fileName)
}

func readCSV(fileName string) {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(file)

	options := os.O_WRONLY | os.O_CREATE //開啟檔案的選項
	file2, err := os.OpenFile("day_kline.csv", options, os.FileMode(0600))
	check(err)
	date := ""
	o := 0
	h := 0
	l := 0
	c := 0
	v := 0

	str := "Date,Open,High,Low,Close,TotalVolume"
	_, err = fmt.Fprintln(file2, str)
	check(err)

	for scanner.Scan() {
		sli := strings.Split(scanner.Text(), ",")

		if len(sli) <= 1 {
			continue
		}
		vv, _ := strconv.Atoi(sli[6])

		if date != sli[0] { //新的一天寫入資料
			if c != 0 {
				str := date + "," + strconv.Itoa(o) + "," + strconv.Itoa(h) + "," + strconv.Itoa(l) + "," + strconv.Itoa(c) + "," + strconv.Itoa(v)
				_, err = fmt.Fprintln(file2, str)
				check(err)
			}

			date = sli[0] //會直接執行下方條件
			oo, _ := strconv.ParseFloat(sli[2], 64)
			cc, _ := strconv.ParseFloat(sli[5], 64)
			o = int(oo)
			h = 0
			l = 0
			c = int(cc)
			v = 0
		}

		if date == sli[0] { //壓k棒的 高 低 量
			hh, _ := strconv.ParseFloat(sli[3], 64)
			ll, _ := strconv.ParseFloat(sli[4], 64)
			h = max(h, int(hh))
			l = min(l, int(ll))
			v += vv
		}
		if "13:30:00" == sli[1] || "13:45:00" == sli[1] {
			cc, _ := strconv.ParseFloat(sli[5], 64)
			c = int(cc) //取最後一根的收
		}
	}
	str = date + "," + strconv.Itoa(o) + "," + strconv.Itoa(h) + "," + strconv.Itoa(l) + "," + strconv.Itoa(c) + "," + strconv.Itoa(v)
	_, err = fmt.Fprintln(file2, str) //最後一筆資料寫入
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
