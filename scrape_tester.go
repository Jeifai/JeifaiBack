package main

import (
	"fmt"
	// "io/ioutil"
	// "log"
	"net/http"
	"os"
	// "time"
	// "strconv"
	// "context"
	// "strings"

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

	// start_url := "https://jobs.sap.com/search/?q=&sortColumn=referencedate&sortDirection=desc&startrow=%d"
	// base_url := "https://jobs.sap.com/%s"
	// number_results_per_page := 25
	// counter := 0

	base_job_url := "https://improbable.io%s"
	c.OnHTML("html", func(e *colly.HTMLElement) {
		e.ForEach(".vacancy-card-module--root--b8XxX", func(_ int, el *colly.HTMLElement) {
			result_title := el.ChildText("h2 > a")
			result_url := fmt.Sprintf(base_job_url, el.ChildAttr("h2 > a", "href"))
			result_location := el.ChildTexts("li")[0]
			result_department := el.ChildTexts("li")[1]
			result_description := el.ChildText(".vacancy-card-module--body--2CJvD")
			fmt.Println(result_title)
			fmt.Println(result_url)
			fmt.Println(result_location)
			fmt.Println(result_department)
			fmt.Println(result_description)
			fmt.Printf("\n")
		})
	})
	c.WithTransport(t)
	c.Visit("file:" + dir + "/sonomotors.html")
	// c.Visit(fmt.Sprintf(start_url, 0))
	///////////////////////////////////////////////////////////////////////////////// CHROMEDP BLOCK

	/**
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// var res string
	var initialResponse string
	if err := chromedp.Run(ctx,
		chromedp.Navigate("https://improbable.io/careers/opportunities/"),
		chromedp.Sleep(5*time.Second),
		// chromedp.EvaluateAsDevTools(`document.cookie`, &res),
		chromedp.OuterHTML("html", &initialResponse),
	); err != nil {
		panic(err)
	}
	SaveResponseToFileWithFileName(initialResponse, "sonomotors.html")
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
