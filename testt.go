package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gocolly/colly/v2"
)

func mainesttT() {
	c := colly.NewCollector()

	c.OnResponse(func(r *colly.Response) {
		results_sections := strings.Split(string(r.Body), "job_list_row")
		for i := 1; i < len(results_sections); i++ {
			elem := results_sections[i]
			result_url := strings.Split(strings.Split(elem, `<p><a href="`)[1], `"`)[0]
			result_title := strings.Split(strings.Split(elem, `class="job_link font_bold">`)[1], `</a>`)[0]
			result_location := strings.Split(strings.Split(elem, `<span class="location">`)[1], `</span>`)[0]
			result_category := strings.Split(strings.Split(elem, `<span class="category">`)[1], `</span>`)[0]
			result_description := strings.Split(strings.Split(elem, `<p class="jlr_description">`)[1], `</p>`)[0]
			fmt.Println(result_url)
			fmt.Println("\t", result_title)
			fmt.Println("\t\t", result_location)
			fmt.Println("\t\t\t", result_category)
			fmt.Println("\t\t\t\t", result_description)
		}
	})

	t := &http.Transport{}
	t.RegisterProtocol("file", http.NewFileTransport(http.Dir("/")))
	c.WithTransport(t)
	dir, err := os.Getwd()
	if err != nil {
		panic(err.Error())
	}
	c.Visit("file:" + dir + "/bodyText.html")
}
