package main

import (
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"time"
)

type Runtime struct {
	Name string
}
type Result struct {
	CompanyName string
	CompanyUrl  string
	Title       string
	ResultUrl   string
}

func runner(scraper_name string, scraper_version int) (result []Result) {
	runtime := Runtime{scraper_name}
	strucReflected := reflect.ValueOf(runtime)
	method := strucReflected.MethodByName(scraper_name)
	params := []reflect.Value{reflect.ValueOf(scraper_version)}
	temp_result := method.Call(params)
	result = temp_result[0].Interface().([]Result)
	return
}

func getJson(url string, target interface{}) error {
	var client = &http.Client{Timeout: 10 * time.Second}
	r, err := client.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(target)
}

func (runtime Runtime) Kununu(version int) (results []Result) {
	if version == 1 {
		url := "https://www.kununu.com/at/kununu/jobs"
		main_tag := "div"
		main_tag_attr := "class"
		main_tag_value := "company-profile-job-item"
		tag_title := "a"
		tag_url := "a"

		c := colly.NewCollector()
		c.OnHTML(main_tag, func(e *colly.HTMLElement) {
			if strings.Contains(e.Attr(main_tag_attr), main_tag_value) {
				result_title := e.ChildText(tag_title)
				result_url := e.ChildAttr(tag_url, "href")
				results = append(results, Result{
					runtime.Name,
					url,
					result_title,
					result_url})
			}
		})
		c.Visit(url)
	}
	return
}

func (runtime Runtime) Mitte(version int) (results []Result) {
	if version == 1 {
		url := "https://api.lever.co/v0/postings/mitte?group=team&mode=json"
		res, err := http.Get(url)
		if err != nil {
			panic(err.Error())
		}
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			panic(err.Error())
		}

		type Postings struct {
			Title string `json:"text"`
			Url   string `json:"hostedUrl"`
		}
		type JsonJob struct {
			Positions []Postings `json:"postings"`
		}
		type JsonJobs []JsonJob
		var jsonJobs JsonJobs
		err = json.Unmarshal(body, &jsonJobs)
		if err != nil {
			fmt.Println(err)
		}
		for _, elem := range jsonJobs {
			result_title := elem.Positions[0].Title
			result_url := elem.Positions[0].Url
			results = append(results, Result{
				runtime.Name,
				url,
				result_title,
				result_url})
		}
	}
	return
}
