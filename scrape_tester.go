package main

import (
	"fmt"
	"net/http"
	"os"
	// "strconv"
	// "strings"
	// "context"

	"github.com/gocolly/colly/v2"
	// "github.com/chromedp/chromedp"
)

func mainTTTT() {
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

	/**
	    ctx, cancel := chromedp.NewContext(context.Background())
		defer cancel()

		var initialResponse string
		if err := chromedp.Run(ctx,
			chromedp.Navigate("https://coachhub-jobs.personio.de/"),
			chromedp.WaitReady("body", chromedp.ByQuery),
			chromedp.OuterHTML("html", &initialResponse),
		); err != nil {
			panic(err)
		}
	    SaveResponseToFileWithFileName(initialResponse, "coachhub.html")

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
	*/

	c.WithTransport(t)
	c.Visit("file:" + dir + "/response.html")
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
