package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	allFileName := []string{
		"1998_fut.csv", "1999_fut.csv", "2000_fut.csv", "2001_fut.csv",
		"2002_fut.csv", "2003_fut.csv", "2004_fut.csv", "2005_fut.csv",
		"2006_fut.csv", "2007_fut.csv", "2008_fut.csv", "2009_fut.csv",
		"2010_fut.csv", "2011_fut.csv", "2012_fut.csv", "2013_fut.csv",
		"2014_fut.csv", "2015_fut.csv", "2016_fut.csv", "2017_fut.csv",
		"2018_fut.csv", "2019_fut.csv", "2020_fut.csv",
	}

	for _, fileName := range allFileName {
		fileName = "o_data/"+fileName
		readCSV(fileName)
	}

}

func readCSV(fileName string) {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(file)
	options := os.O_WRONLY | os.O_APPEND | os.O_CREATE //開啟檔案的選項
	file2, err := os.OpenFile("new.csv", options, os.FileMode(0600))
	check(err)
	// i := 0
	for scanner.Scan() {
		// if i < 0 {
		// 	_, err = fmt.Fprintln(file2, scanner.Text())
		// 	check(err)
		// 	i++
		// 	continue
		// }
		sli := strings.Split(scanner.Text(), ",")
		if len(sli) > 1 {
			ok:=false
			if(len(sli) < 18){
				ok= true
			}else{
				if(sli[17] == "一般"){
					ok=true
				}
			}
			if sli[1] == "MTX" && ok {
				if strings.Trim(sli[2], " ") == finalDay(sli[0]) { // 第三週的禮拜三結算，隔天要抓下個月的值
					_, err = fmt.Fprintln(file2, scanner.Text())
					check(err)
				}
			}
		} else {
			fmt.Println(scanner.Text())
		}

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

func finalDay(sli string) string {
	a := strings.Split(sli, "/")
	year, _ := strconv.Atoi(a[0])
	month, _ := strconv.Atoi(a[1])
	day, _ := strconv.Atoi(a[2])
	b := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local)
	c := b.Weekday() //每個月一號是禮拜幾
	final := 0       //結算日
	if c == 0 {
		final = 18
	}
	if c == 1 {
		final = 17
	}
	if c == 2 {
		final = 16
	}
	if c == 3 {
		final = 15
	}
	if c == 4 {
		final = 21
	}
	if c == 5 {
		final = 20
	}
	if c == 6 {
		final = 19
	}
	gg := ""
	if day <= final {
		gg = a[0] + a[1]
	} else {
		if month+1 < 10 {
			gg = a[0] + "0" + strconv.Itoa(month+1)
		} else if month+1 > 12 {
			gg = strconv.Itoa(year+1) + "01"
		} else {
			gg = a[0] + strconv.Itoa(month+1)

		}
	}

	return gg
}
