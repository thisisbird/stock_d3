package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

var timeMap = map[string]string{}
var startTime = "08:45"
var minK = 300

func main() {
	timeMap = timeMapping(startTime, minK)
	finaFilelName := startTime + "_" + strconv.Itoa(minK) + "min.csv"
	fileName := "TXF1-分鐘-成交價.txt"
	fileName = "o_data/kevin/" + fileName
	readCSV(fileName, finaFilelName)
}

/**
* start 開始時間
* count 幾分k
 */
func readCSV(fileName string, finaFilelName string) {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(file)

	options := os.O_WRONLY | os.O_CREATE //開啟檔案的選項
	file2, err := os.OpenFile(finaFilelName, options, os.FileMode(0600))
	check(err)
	date := ""
	time := ""
	finalTime := ""
	o := 0
	h := 0
	l := 0
	c := 0
	v := 0

	str := "Date,Time,Open,High,Low,Close,TotalVolume"
	_, err = fmt.Fprintln(file2, str)
	check(err)

	for scanner.Scan() {
		sli := strings.Split(scanner.Text(), ",")

		if len(sli) <= 1 || timeMap[sli[1]] == "" {
			continue
		}
		vv, _ := strconv.Atoi(sli[6])
		if date != sli[0] || finalTime != timeMap[sli[1]] { //新的一k寫入資料
			if c != 0 {
				str := date + "," + finalTime + "," + strconv.Itoa(o) + "," + strconv.Itoa(h) + "," + strconv.Itoa(l) + "," + strconv.Itoa(c) + "," + strconv.Itoa(v)
				_, err = fmt.Fprintln(file2, str)
				check(err)
			}

			date = sli[0]             //會直接執行下方條件
			time = sli[1]             //會直接執行下方條件
			finalTime = timeMap[time] //會直接執行下方條件

			oo, _ := strconv.ParseFloat(sli[2], 64)
			cc, _ := strconv.ParseFloat(sli[5], 64)
			o = int(oo)
			h = 0
			l = 0
			c = int(cc)
			v = 0
		}

		if date == sli[0] && timeMap[sli[1]] == finalTime { //壓k棒的 高 低 量
			hh, _ := strconv.ParseFloat(sli[3], 64)
			ll, _ := strconv.ParseFloat(sli[4], 64)
			h = max(h, int(hh))
			l = min(l, int(ll))
			v += vv
		}
		if finalTime == sli[1] {
			cc, _ := strconv.ParseFloat(sli[5], 64)
			c = int(cc) //取最後一根的收
		}
	}
	str = date + "," + finalTime + "," + strconv.Itoa(o) + "," + strconv.Itoa(h) + "," + strconv.Itoa(l) + "," + strconv.Itoa(c) + "," + strconv.Itoa(v)
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

func timeMapping(start string, count int) map[string]string {
	data := map[string]string{}
	finalTime := start
	for i := 1; i <= 999; i++ {
		if i%count == 1 {
			finalTime = timePlus(finalTime, count)
		}
		time := timePlus(start, i)
		data[time] = finalTime

		if time == "13:45:00" {
			break
		}
	}

	return data
}

func timePlus(time string, plus int) string {
	sli := strings.Split(time, ":")
	hour, _ := strconv.Atoi(sli[0])
	min, _ := strconv.Atoi(sli[1])
	totalMin := min + plus
	min = totalMin % 60
	hourPlus := totalMin / 60
	hour += hourPlus
	return intTOString(hour) + ":" + intTOString(min) + ":00"
}

func intTOString(x int) string {
	if x < 10 {
		return "0" + strconv.Itoa(x)
	}
	return strconv.Itoa(x)
}
