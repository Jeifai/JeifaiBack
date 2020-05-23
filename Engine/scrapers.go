package main

import (
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"
)

type Runtime struct {
	Name string
}

type Response struct {
	Html []byte
}

type Result struct {
	CompanyName string
	ScrapingUrl string
	Title       string
	ResultUrl   string
}

func Scrape(
	scraper_name string, scraper_version int, isLocal bool) (
	response Response, result []Result) {
	fmt.Println("Starting Scrape...")
	runtime := Runtime{scraper_name}
	strucReflected := reflect.ValueOf(runtime)
	method := strucReflected.MethodByName(scraper_name)
	params := []reflect.Value{
		reflect.ValueOf(scraper_version),
		reflect.ValueOf(isLocal)}
	function_output := method.Call(params)
	response = function_output[0].Interface().(Response)
	results = function_output[1].Interface().([]Result)
	results = Unique(result)
	fmt.Println("Number of results scraped: " + strconv.Itoa(len(results)))
	return
}

func (runtime Runtime) Kununu(
	version int, isLocal bool) (
	response Response, results []Result) {
	c := colly.NewCollector()
	if isLocal {
		t := &http.Transport{}
		t.RegisterProtocol("file", http.NewFileTransport(http.Dir("/")))
		c.WithTransport(t)
	}
	if version == 1 {
		url := "https://www.kununu.com/at/kununu/jobs"
		main_tag := "div"
		main_tag_attr := "class"
		main_tag_value := "company-profile-job-item"
		tag_title := "a"
		tag_url := "a"
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
		c.OnResponse(func(r *colly.Response) {
			response = Response{r.Body}
		})
		c.OnRequest(func(r *colly.Request) {
			fmt.Println("Visiting", r.URL.String())
		})
		c.OnError(func(r *colly.Response, err error) {
			fmt.Println(
				"Request URL:", r.Request.URL,
				"failed with response:", r,
				"\nError:", err)
		})
		if isLocal {
			dir, err := os.Getwd()
			if err != nil {
				fmt.Println(err)
			}
			c.Visit("file:" + dir + "/response.html")
		} else {
			c.Visit(url)
		}
	}
	return
}

func (runtime Runtime) Mitte(
	version int, isLocal bool) (
	response Response, results []Result) {
	var body []byte
	url := "https://api.lever.co/v0/postings/mitte?group=team&mode=json"
	if isLocal {
		dir, err := os.Getwd()
		if err != nil {
			fmt.Println(err)
		}
		temp_body, err := ioutil.ReadFile(dir + "/response.html")
		fmt.Println("Visiting", dir+"/response.html")
		if err != nil {
			panic(err.Error())
		}
		body = temp_body
	} else {
		res, err := http.Get(url)
		fmt.Println("Visiting", url)
		if err != nil {
			panic(err.Error())
		}
		temp_body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			panic(err.Error())
		}
		body = temp_body
	}
	if version == 1 {
		response = Response{body}
		type Postings struct {
			Title string `json:"text"`
			Url   string `json:"hostedUrl"`
		}
		type JsonJob struct {
			Positions []Postings `json:"postings"`
		}
		type JsonJobs []JsonJob
		var jsonJobs JsonJobs
		err := json.Unmarshal(body, &jsonJobs)
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

func (runtime Runtime) IMusician(
	version int, isLocal bool) (
	response Response, results []Result) {
	c := colly.NewCollector()
	if isLocal {
		t := &http.Transport{}
		t.RegisterProtocol("file", http.NewFileTransport(http.Dir("/")))
		c.WithTransport(t)
	}
	if version == 1 {
		url := "https://imusician-digital-jobs.personio.de/"
		main_tag := "a"
		main_tag_attr := "class"
		main_tag_value := "job-box-link"
		tag_title := ".jb-title"

		c.OnHTML(main_tag, func(e *colly.HTMLElement) {
			if strings.Contains(e.Attr(main_tag_attr), main_tag_value) {
				result_title := e.ChildText(tag_title)
				result_url := e.Attr("href")
				results = append(results, Result{
					runtime.Name,
					url,
					result_title,
					result_url})
			}
		})
		c.OnResponse(func(r *colly.Response) {
			response = Response{r.Body}
		})
		c.OnRequest(func(r *colly.Request) {
			fmt.Println("Visiting", r.URL.String())
		})
		c.OnError(func(r *colly.Response, err error) {
			fmt.Println(
				"Request URL:", r.Request.URL,
				"failed with response:", r,
				"\nError:", err)
		})
		if isLocal {
			dir, err := os.Getwd()
			if err != nil {
				fmt.Println(err)
			}
			c.Visit("file:" + dir + "/response.html")
		} else {
			c.Visit(url)
		}
	}
	return
}

func (runtime Runtime) Babelforce(
	version int, isLocal bool) (
	response Response, results []Result) {
	c := colly.NewCollector()
	if isLocal {
		t := &http.Transport{}
		t.RegisterProtocol("file", http.NewFileTransport(http.Dir("/")))
		c.WithTransport(t)
	}
	if version == 1 {
		url := "https://www.babelforce.com/jobs/"
		main_tag := "div"
		main_tag_attr := "class"
		main_tag_value := "qodef-portfolio"
		tag_title := "h5"
		tag_url := "a"

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
		c.OnResponse(func(r *colly.Response) {
			response = Response{r.Body}
		})
		c.OnRequest(func(r *colly.Request) {
			fmt.Println("Visiting", r.URL.String())
		})
		c.OnError(func(r *colly.Response, err error) {
			fmt.Println(
				"Request URL:", r.Request.URL,
				"failed with response:", r,
				"\nError:", err)
		})
		if isLocal {
			dir, err := os.Getwd()
			if err != nil {
				fmt.Println(err)
			}
			c.Visit("file:" + dir + "/response.html")
		} else {
			c.Visit(url)
		}
	}
	return
}
