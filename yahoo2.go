//可用
package main

import (
	"crypto/tls"
	"fmt"

	"net/http"
	"net/http/cookiejar"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/queue"
)

func main() {
	jar, _ := cookiejar.New(nil)

	// Instantiate default collector
	c := colly.NewCollector(colly.AllowedDomains("tw.stock.yahoo.com"))

	c.WithTransport(&http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	})

	//setup our client based on the cookies data
	c.SetCookieJar(jar)

	q, _ := queue.New(
		1, // Number of consumer threads
		&queue.InMemoryQueueStorage{MaxSize: 100000}, // Use default queue storage
	)
	q.AddURL("https://tw.stock.yahoo.com/q/q?s=2330")
	q.AddURL("https://tw.stock.yahoo.com/q/q?s=2317")
	q.AddURL("https://tw.stock.yahoo.com/q/q?s=3008")

	c.OnHTML("body", func(e *colly.HTMLElement) {

		date := e.DOM.Find("table").Eq(0).Find("td").Eq(1).Find("font").Eq(0).Text()
		time := e.DOM.Find("table").Eq(2).Find("td").Eq(1).Text()
		name := e.DOM.Find("table").Eq(2).Find("td").Eq(0).Find("a").Eq(0).Text()
		o := e.DOM.Find("table").Eq(2).Find("td").Eq(8).Text()
		h := e.DOM.Find("table").Eq(2).Find("td").Eq(9).Text()
		l := e.DOM.Find("table").Eq(2).Find("td").Eq(10).Text()
		c := e.DOM.Find("table").Eq(2).Find("td").Eq(7).Text()
		q := e.DOM.Find("table").Eq(2).Find("td").Eq(6).Text() //量
		fmt.Println(date, time, name, o, h, l, c, q)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	q.Run(c)
}
