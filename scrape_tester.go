package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
	// "strconv"
	"context"
	"strings"

	// "github.com/gocolly/colly/v2"
	"github.com/chromedp/chromedp"
)

func mainTTT() {
	/**
		t := &http.Transport{}
		t.RegisterProtocol("file", http.NewFileTransport(http.Dir("/")))
		dir, err := os.Getwd()
		if err != nil {
			panic(err.Error())
		}

		c := colly.NewCollector()

		c.OnHTML(".tab-content", func(e *colly.HTMLElement) {
			e.ForEach("p", func(_ int, el *colly.HTMLElement) {
				result_title := el.ChildText("a")
				result_url := el.ChildAttr("a", "href")
				fmt.Println(result_title, result_url)
			})
	    })
	*/

	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	var res string
	var initialResponse string
	if err := chromedp.Run(ctx,
		chromedp.Navigate("https://job.bytedance.com/en/position?limit=10"),
		chromedp.Sleep(5*time.Second),
		chromedp.EvaluateAsDevTools(`document.cookie`, &res),
		chromedp.OuterHTML("html", &initialResponse),
	); err != nil {
		panic(err)
	}
	SaveResponseToFileWithFileName(initialResponse, "bytedance.html")

	token := strings.Split(res, "atsx-csrf-token=")[1]
	fmt.Println(res)
	fmt.Println(token)

	client := &http.Client{}
	data := strings.NewReader(`{"limit":30}`)
	req, err := http.NewRequest("POST", "https://job.bytedance.com/api/v1/search/job/posts", data)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("x-csrf-token", strings.ReplaceAll(token, "%3D", "="))
	req.Header.Set("Cookie", "channel=overseas; atsx-csrf-token="+token)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", bodyText)

	SaveResponseToFileWithFileName(string(bodyText), "bytedancef.json")

	/**
	    start_url := "https://coachhub-jobs.personio.de/"
	    tag_main := ".panel-container"
	    tag_job_section := ".recent-job-list"
	    tag_title := "h6"
	    tag_location := "p"

	    _ = start_url

	    type Job struct {
	        Url        string
	        Title      string
	        Location   string
	    }

	    c := colly.NewCollector()

	    c.OnHTML(tag_main, func(e *colly.HTMLElement) {
	        e.ForEach(tag_job_section, func(_ int, el *colly.HTMLElement) {
	            result_url := el.ChildAttr("a", "href")
	            result_title := el.ChildText(tag_title)
	            result_location := strings.Split(el.ChildText(tag_location), "·")[1]
	            fmt.Println(result_url)
	            fmt.Println("\t", result_title)
	            fmt.Println("\t\t", result_location)

	            _ = result_title
	            _ = result_location
	        })
	    })

		c.WithTransport(t)
	    c.Visit("file:" + dir + "/response.html")
	*/
}

func SaveResponseToFileWithFileName(response string, filename string) {
	dir, err := os.Getwd()
	if err != nil {
		panic(err.Error())
	}
	f, err := os.Create(dir + "/" + filename)
	if err != nil {
		panic(err.Error())
	}
	defer f.Close()
	f.WriteString(response)
}
