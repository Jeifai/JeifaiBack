package main

import (
	"fmt"
	// "io/ioutil"
	// "log"
	// "net/http"
	"os"
	"time"
	// "strconv"
	// "context"
	// "strings"

	"github.com/gocolly/colly/v2"
	// "github.com/chromedp/chromedp"
)

func mainTest() {
	/**
	  t := &http.Transport{}
	  t.RegisterProtocol("file", http.NewFileTransport(http.Dir("/")))
	  dir, err := os.Getwd()
	  if err != nil {
	      panic(err.Error())
	  }
	*/

	c := colly.NewCollector()

	start_url := "https://jobs.sap.com/search/?q=&sortColumn=referencedate&sortDirection=desc&startrow=%d"
	base_url := "https://jobs.sap.com/%s"
	number_results_per_page := 25
	counter := 0

	c.OnHTML(".html5", func(e *colly.HTMLElement) {
		e.ForEach(".data-row", func(_ int, el *colly.HTMLElement) {
			result_url := fmt.Sprintf(base_url, el.ChildAttr("a", "href"))
			fmt.Println(result_url)
		})
		// temp_pages := strings.Split(e.ChildText(".srHelp"), " ")
		// s_temp_pages := temp_pages[len(temp_pages)-1]
		// total_pages, err := strconv.Atoi(s_temp_pages)
		// if err != nil {
		//     panic(err.Error())
		// }

		if counter >= 3 {
			return
		} else {
			counter++
			time.Sleep(2 * time.Second)
			// fmt.Println(fmt.Sprintf(start_url, counter * number_results_per_page))
			c.Visit(fmt.Sprintf(start_url, counter*number_results_per_page))
		}
	})
	// c.WithTransport(t)
	// c.Visit("file:" + dir + "/response.html")
	c.Visit(fmt.Sprintf(start_url, 0))

	///////////////////////////////////////////////////////////////////////////////// CHROMEDP BLOCK
	/**
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
