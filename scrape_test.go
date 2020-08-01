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

    /**
    ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	var initialResponse string
	if err := chromedp.Run(ctx,
		chromedp.Navigate("https://ninox.com/en/jobs"),
		chromedp.WaitReady("body", chromedp.ByQuery),
		chromedp.OuterHTML("html", &initialResponse),
	); err != nil {
		panic(err)
	}
    SaveResponseToFileWithFileName(initialResponse, "ninox.html")
    */

    start_url := "https://ninox.com/en/jobs"
    job_base_url := "https://ninox.com/%s"
    tag_job_section := ".job-new"
    tag_title := "h4"
    tag_location := ".jobs-j-openinglugar"

    _ = start_url


    type Job struct {
        Url        string
        Title      string
        Location   string
    }

    c := colly.NewCollector()

    c.OnHTML(tag_job_section, func(e *colly.HTMLElement) {
        result_url := fmt.Sprintf(job_base_url, e.ChildAttr("a", "href"))
        result_title := e.ChildText(tag_title)
        result_location := e.ChildText(tag_location)
        fmt.Println(result_url)
        fmt.Println("\t", result_title)
        fmt.Println("\t\t", result_location)
    })

    c.WithTransport(t)
    c.Visit("file:" + dir + "/ninox.html")
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
