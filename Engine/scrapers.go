package main

import (
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly"
	"io/ioutil"
	"net/http"
	netUrl "net/url"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type Runtime struct {
	Name string
}

type Response struct {
	Html []byte
}

type Result struct {
	CompanyName string
	Title       string
	ResultUrl   string
	Data        json.RawMessage
}

const SecondsSleep = 2 // Seconds between pagination

func Scrape(
	scraper_name string, scraper_version int, isLocal bool) (
	response Response, results []Result) {
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
	results = Unique(results)
	fmt.Println("Number of results scraped: " + strconv.Itoa(len(results)))
	return
}

func (runtime Runtime) Kununu(
	version int, isLocal bool) (
	response Response, results []Result) {

	if version == 1 {

		c := colly.NewCollector()

		url := "https://www.kununu.com/at/kununu/jobs"
		main_tag := "div"
		main_tag_attr := "class"
		main_tag_value := "company-profile-job-item"
		tag_title := "a"
		tag_url := "a"
		tag_location := ".item-location"

		type Job struct {
			Title    string
			Url      string
			Location string
		}

		c.OnHTML(main_tag, func(e *colly.HTMLElement) {
			if strings.Contains(e.Attr(main_tag_attr), main_tag_value) {
				result_title := e.ChildText(tag_title)
				result_url := e.ChildAttr(tag_url, "href")
				result_location := e.ChildText(tag_location)

				_, err := netUrl.ParseRequestURI(result_url)
				if err == nil {

					temp_elem_json := Job{
						result_title,
						result_url,
						result_location}

					elem_json, err := json.Marshal(temp_elem_json)
					if err != nil {
						panic(err.Error())
					}

					results = append(results, Result{
						runtime.Name,
						result_title,
						result_url,
						elem_json})
				}
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
			t := &http.Transport{}
			t.RegisterProtocol("file", http.NewFileTransport(http.Dir("/")))
			c.WithTransport(t)
			dir, err := os.Getwd()
			if err != nil {
				panic(err.Error())
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

	if version == 1 {

		url := "https://api.lever.co/v0/postings/mitte?group=team&mode=json"

		type Category struct {
			Commitment string `json:"commitment"`
			Location   string `json:"location"`
			Team       string `json:"team"`
		}

		type List struct {
			Text    string `json:"text"`
			Content string `json:"content"`
		}

		type Job struct {
			Title            string   `json:"text"`
			Url              string   `json:"hostedUrl"`
			AdditionalPlain  string   `json:"additionalPlain"`
			Additional       string   `json:"additional"`
			Categories       Category `json:"categories"`
			CreatedAt        int      `json:"createdAt"`
			DescriptionPlain string   `json:"descriptionPlain"`
			Description      string   `json:"description"`
			Id               string   `json:"id"`
			Lists            []List   `json:"lists"`
		}

		type Postings struct {
			Jobs []Job `json:"postings"`
		}

		type JsonJobs []Postings

		var body []byte
		if isLocal {
			dir, err := os.Getwd()
			if err != nil {
				panic(err.Error())
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

		response = Response{body}

		var jsonJobs JsonJobs
		err := json.Unmarshal(body, &jsonJobs)
		if err != nil {
			panic(err.Error())
		}

		for _, elem := range jsonJobs {
			for _, subElem := range elem.Jobs {

				result_title := subElem.Title
				result_url := subElem.Url

				_, err := netUrl.ParseRequestURI(result_url)
				if err == nil {

					elem_json, err := json.Marshal(subElem)
					if err != nil {
						panic(err.Error())
					}

					results = append(results, Result{
						runtime.Name,
						result_title,
						result_url,
						elem_json})
				}
			}
		}
	}
	return
}

func (runtime Runtime) IMusician(
	version int, isLocal bool) (
	response Response, results []Result) {

	if version == 1 {

		c := colly.NewCollector()

		url := "https://imusician-digital-jobs.personio.de/"
		main_tag := "a"
		main_tag_attr := "class"
		main_tag_value := "job-box-link"
		tag_title := ".jb-title"
		tag_description := "span"
		tag_location := "span"

		type Job struct {
			Title       string
			Url         string
			Description string
			Location    string
		}

		c.OnHTML(main_tag, func(e *colly.HTMLElement) {
			if strings.Contains(e.Attr(main_tag_attr), main_tag_value) {
				result_title := e.ChildText(tag_title)
				result_url := e.Attr("href")
				result_description := e.ChildTexts(tag_description)[0]
				result_location := e.ChildTexts(tag_location)[2]

				_, err := netUrl.ParseRequestURI(result_url)
				if err == nil {

					temp_elem_json := Job{
						result_title,
						result_url,
						result_description,
						result_location}

					elem_json, err := json.Marshal(temp_elem_json)
					if err != nil {
						panic(err.Error())
					}

					results = append(results, Result{
						runtime.Name,
						result_title,
						result_url,
						elem_json})
				}
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
			t := &http.Transport{}
			t.RegisterProtocol("file", http.NewFileTransport(http.Dir("/")))
			c.WithTransport(t)
			dir, err := os.Getwd()
			if err != nil {
				panic(err.Error())
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

	if version == 1 {

		c := colly.NewCollector()

		url := "https://www.babelforce.com/jobs/"
		main_tag := "div"
		main_tag_attr := "class"
		main_tag_value := "qodef-portfolio"
		tag_title := "h5"
		tag_url := "a"

		type Job struct {
			Title string
			Url   string
		}

		c.OnHTML(main_tag, func(e *colly.HTMLElement) {
			if strings.Contains(e.Attr(main_tag_attr), main_tag_value) {

				result_title := e.ChildText(tag_title)
				result_url := e.ChildAttr(tag_url, "href")

				_, err := netUrl.ParseRequestURI(result_url)
				if err == nil {

					temp_elem_json := Job{
						result_title,
						result_url}
					elem_json, err := json.Marshal(temp_elem_json)
					if err != nil {
						panic(err.Error())
					}

					results = append(results, Result{
						runtime.Name,
						result_title,
						result_url,
						elem_json})
				}
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
			t := &http.Transport{}
			t.RegisterProtocol("file", http.NewFileTransport(http.Dir("/")))
			c.WithTransport(t)
			dir, err := os.Getwd()
			if err != nil {
				panic(err.Error())
			}
			c.Visit("file:" + dir + "/response.html")
		} else {
			c.Visit(url)
		}
	}
	return
}

func (runtime Runtime) Zalando(
	version int, isLocal bool) (response Response, results []Result) {

	if version == 1 {

		z_base_url := "https://jobs.zalando.com/api/jobs/?limit=100&offset="

		type Job struct {
			Id            int      `json:"id"`
			Title         string   `json:"title"`
			JobCategories []string `json:"job_categories"`
			UpdatedAt     string   `json:"updated_at"`
			Officies      []string `json:"offices"`
			Entity        string   `json:"entity"`
		}

		type JsonJobs struct {
			Data []Job  `json:"data"`
			Last string `json:last`
		}

		var jsonJobs_1 JsonJobs

		if isLocal {
			dir, err := os.Getwd()
			if err != nil {
				panic(err.Error())
			}
			body, err := ioutil.ReadFile(dir + "/response.html")
			fmt.Println("Visiting", dir+"/response.html")
			if err != nil {
				panic(err.Error())
			}
			err = json.Unmarshal(body, &jsonJobs_1)
			if err != nil {
				panic(err.Error())
			}
		} else {
			var first_body []byte
			res, err := http.Get(z_base_url + "0")
			fmt.Println("Visiting ", z_base_url+"0")
			if err != nil {
				panic(err.Error())
			}
			temp_body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				panic(err.Error())
			}
			first_body = temp_body
			err = json.Unmarshal(first_body, &jsonJobs_1)
			if err != nil {
				panic(err.Error())
			}

			offset, err := strconv.Atoi(
				strings.Split(jsonJobs_1.Last, "offset=")[1])
			if err != nil {
				panic(err.Error())
			}

			for i := 1; i < (offset/100)+1; i++ {
				temp_z_url := z_base_url + strconv.Itoa(i*100)
				res, err := http.Get(temp_z_url)
				fmt.Println("Visiting ", temp_z_url)
				if err != nil {
					panic(err.Error())
				}
				temp_body, err := ioutil.ReadAll(res.Body)
				if err != nil {
					panic(err.Error())
				}
				var tempJsonJobs_2 JsonJobs
				err = json.Unmarshal(temp_body, &tempJsonJobs_2)
				if err != nil {
					panic(err.Error())
				}
				jsonJobs_1.Data = append(
					jsonJobs_1.Data, tempJsonJobs_2.Data...)
				time.Sleep(SecondsSleep * time.Second)
			}

			response_json, err := json.Marshal(jsonJobs_1)
			if err != nil {
				panic(err.Error())
			}
			response = Response{[]byte(response_json)}
		}

		for _, elem := range jsonJobs_1.Data {

			result_title := elem.Title
			z_base_result_url := "https://jobs.zalando.com/de/jobs/"
			result_url := z_base_result_url + strconv.Itoa(elem.Id)

			_, err := netUrl.ParseRequestURI(result_url)
			if err == nil {

				elem_json, err := json.Marshal(elem)
				if err != nil {
					panic(err.Error())
				}

				results = append(results, Result{
					runtime.Name,
					result_title,
					result_url,
					elem_json})
			}
		}
	}
	return
}

func (runtime Runtime) Google(
	version int, isLocal bool) (response Response, results []Result) {

	if version == 1 {

		g_base_url := "https://careers.google.com/api/jobs/jobs-v1/search/?page_size=100&page="

		g_base_result_url := "https://careers.google.com/jobs/results/"

		results_per_page := 100

		type Job struct {
			Id            string   `json:"job_id"`
			Title         string   `json:"job_title"`
			Summary       string   `json:"summary"`
			PublishDate   string   `json:"publish_date"`
			LocationCount string   `json:"locations_count"`
			Locations     []string `json:"locations"`
			Description   string   `json:"description"`
			CompanyName   string   `json:"company_name"`
			CompanyId     string   `json:"company_name"`
		}

		type JsonJobs struct {
			Jobs  []Job  `json:"jobs"`
			Count string `json:"count"`
			Next  string `json:"next_page"`
		}

		var jsonJobs_1 JsonJobs

		if isLocal {
			dir, err := os.Getwd()
			if err != nil {
				panic(err.Error())
			}
			body, err := ioutil.ReadFile(dir + "/response.html")
			fmt.Println("Visiting", dir+"/response.html")
			if err != nil {
				panic(err.Error())
			}
			err = json.Unmarshal(body, &jsonJobs_1)
			if err != nil {
				panic(err.Error())
			}
		} else {
			var first_body []byte
			res, err := http.Get(g_base_url + "1")
			if err != nil {
				panic(err.Error())
			}
			fmt.Println("Visiting ", g_base_url+"1")
			temp_body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				panic(err.Error())
			}
			first_body = temp_body
			err = json.Unmarshal(first_body, &jsonJobs_1)
			if err != nil {
				panic(err.Error())
			}

			in_count_results, err := strconv.Atoi(jsonJobs_1.Count)
			for i := 2; i < in_count_results/results_per_page+2; i++ {
				temp_g_url := g_base_url + strconv.Itoa(i)
				res, err := http.Get(temp_g_url)
				fmt.Println("Visiting", temp_g_url)
				if err != nil {
					panic(err.Error())
				}
				temp_body, err := ioutil.ReadAll(res.Body)
				if err != nil {
					panic(err.Error())
				}
				var tempJsonJobs_2 JsonJobs
				err = json.Unmarshal(temp_body, &tempJsonJobs_2)
				if err != nil {
					panic(err.Error())
				}
				jsonJobs_1.Jobs = append(jsonJobs_1.Jobs, tempJsonJobs_2.Jobs...)
				time.Sleep(SecondsSleep * time.Second)
			}
		}

		response_json, err := json.Marshal(jsonJobs_1)
		if err != nil {
			panic(err.Error())
		}
		response = Response{[]byte(response_json)}

		for _, elem := range jsonJobs_1.Jobs {

			result_title := elem.Title
			result_url := g_base_result_url + strings.Split(elem.Id, "/")[1]

			_, err := netUrl.ParseRequestURI(result_url)
			if err == nil {

				elem_json, err := json.Marshal(elem)
				if err != nil {
					panic(err.Error())
				}

				results = append(results, Result{
					runtime.Name,
					result_title,
					result_url,
					elem_json})
			}
		}
	}
	return
}

func (runtime Runtime) Soundcloud(
	version int, isLocal bool) (
	response Response, results []Result) {

	if version == 1 {

		c := colly.NewCollector()

		url := "https://boards.greenhouse.io/embed/job_board?for=soundcloud71"
		main_tag := "section"
		main_tag_attr := "class"
		main_tag_value := "level-0"
		tag_title := "a"
		tag_url := "a"
		tag_department := "h3"
		tag_location := "span"

		type Job struct {
			Title      string
			Url        string
			Department string
			Location   string
		}

		c.OnHTML(main_tag, func(e *colly.HTMLElement) {
			if strings.Contains(e.Attr(main_tag_attr), main_tag_value) {
				result_title := e.ChildText(tag_title)
				result_url := e.ChildAttr(tag_url, "href")
				result_department := e.ChildText(tag_department)
				result_location := e.ChildText(tag_location)

				_, err := netUrl.ParseRequestURI(result_url)
				if err == nil {

					temp_elem_json := Job{
						result_title,
						result_url,
						result_department,
						result_location}

					elem_json, err := json.Marshal(temp_elem_json)
					if err != nil {
						panic(err.Error())
					}

					results = append(results, Result{
						runtime.Name,
						result_title,
						result_url,
						elem_json})
				}
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
			t := &http.Transport{}
			t.RegisterProtocol("file", http.NewFileTransport(http.Dir("/")))
			c.WithTransport(t)
			dir, err := os.Getwd()
			if err != nil {
				panic(err.Error())
			}
			c.Visit("file:" + dir + "/response.html")
		} else {
			c.Visit(url)
		}
	}
	return
}

func (runtime Runtime) Microsoft(
	version int, isLocal bool) (
	response Response, results []Result) {

	if version == 1 {

		m_base_url := "https://careers.microsoft.com/us/en/search-results?s=1&from="

		m_base_result_url := "https://careers.microsoft.com/us/en/job/"

		first_part_json := `"eagerLoadRefineSearch":`
		second_part_json := `}; phApp.sessionParams`

		results_per_page := 50

		type Job struct {
			Country           string `json:"country"`
			SubCategory       string `json:"subCategory"`
			Title             string `json:"title"`
			Experience        string `json:"experience"`
			JobSeqNo          string `json:"jobSeqNo"`
			PostedDate        string `json:"postedDate"`
			DescriptionTeaser string `json:"descriptionTeaser"`
			DateCreated       string `json:"dateCreated"`
			State             string `json:"state"`
			JobId             string `json:"jobId"`
			Location          string `json:"location"`
			Category          string `json:"category"`
		}

		type Data struct {
			Jobs []Job `json:"jobs"`
		}

		type JsonJobs struct {
			Data      Data `json:"data"`
			TotalHits int  `json:"totalHits"`
		}

		var jsonJobs_1 JsonJobs

		if isLocal {
			dir, err := os.Getwd()
			if err != nil {
				panic(err.Error())
			}
			body, err := ioutil.ReadFile(dir + "/response.html")
			fmt.Println("Visiting", dir+"/response.html")
			if err != nil {
				panic(err.Error())
			}
			err = json.Unmarshal(body, &jsonJobs_1)
			if err != nil {
				panic(err.Error())
			}
		} else {
			res, err := http.Get(m_base_url)
			if err != nil {
				panic(err.Error())
			}

			temp_body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				panic(err.Error())
			}

			temp_resultsJson := strings.Split(string(temp_body), first_part_json)[1]
			resultsJson := strings.Split(temp_resultsJson, second_part_json)[0]

			err = json.Unmarshal([]byte(resultsJson), &jsonJobs_1)
			if err != nil {
				panic(err.Error())
			}

			number_pages := jsonJobs_1.TotalHits / results_per_page

			for i := 1; i <= number_pages; i++ {
				temp_m_url := m_base_url + strconv.Itoa(i*results_per_page)
				res, err := http.Get(temp_m_url)
				if err != nil {
					panic(err.Error())
				}
				fmt.Println("Visiting", temp_m_url)
				temp_body, err := ioutil.ReadAll(res.Body)
				if err != nil {
					panic(err.Error())
				}
				temp_resultsJson := strings.Split(string(temp_body), first_part_json)[1]
				resultsJson := strings.Split(temp_resultsJson, second_part_json)[0]
				var jsonJobs_2 JsonJobs
				err = json.Unmarshal([]byte(resultsJson), &jsonJobs_2)
				if err != nil {
					panic(err.Error())
				}
				jsonJobs_1.Data.Jobs = append(jsonJobs_1.Data.Jobs, jsonJobs_2.Data.Jobs...)
				time.Sleep(SecondsSleep * time.Second)
			}
		}

		response_json, err := json.Marshal(jsonJobs_1)
		if err != nil {
			panic(err.Error())
		}
		response = Response{[]byte(response_json)}

		for _, elem := range jsonJobs_1.Data.Jobs {

			result_title := elem.Title
			result_url := m_base_result_url + elem.JobId

			_, err := netUrl.ParseRequestURI(result_url)
			if err == nil {

				elem_json, err := json.Marshal(elem)
				if err != nil {
					panic(err.Error())
				}

				results = append(results, Result{
					runtime.Name,
					result_title,
					result_url,
					elem_json})
			}
		}
	}
	return
}

func (runtime Runtime) Twitter(
	version int, isLocal bool) (response Response, results []Result) {

	if version == 1 {

		t_base_url := "https://careers.twitter.com/content/careers-twitter/en/jobs.careers.search.json?limit=100&offset="

        results_per_page := 100
        
		type Location struct {
			Id         string   `json:"id"`
			Title       string   `json:"title"`
		}

        type Team struct {
			Id         string   `json:"id"`
			Title       string   `json:"title"`
		}

		type Job struct {
			Title         string   `json:"title"`
			Description       string   `json:"description"`
			Modified   int   `json:"modified"`
			InternalJob bool   `json:"internalJob"`
			Url     string `json:"url"`
            Team   Team   `json:"team"`
            Locations []Location   `json:"locations"`
		}

		type JsonJobs struct {
			Jobs  []Job  `json:"results"`
			TotalCount int `json:"totalCount"`
			PageCount  int `json:"pageCount"`
		}

		var jsonJobs_1 JsonJobs

		if isLocal {
			dir, err := os.Getwd()
			if err != nil {
				panic(err.Error())
			}
			body, err := ioutil.ReadFile(dir + "/response.html")
			fmt.Println("Visiting", dir+"/response.html")
			if err != nil {
				panic(err.Error())
			}
			err = json.Unmarshal(body, &jsonJobs_1)
			if err != nil {
				panic(err.Error())
			}
		} else {
			var first_body []byte
			res, err := http.Get(t_base_url + "0")
			if err != nil {
				panic(err.Error())
			}
			fmt.Println("Visiting ", t_base_url+"1")
			temp_body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				panic(err.Error())
			}
			first_body = temp_body
			err = json.Unmarshal(first_body, &jsonJobs_1)
			if err != nil {
				panic(err.Error())
			}

			for i := 1; i < jsonJobs_1.TotalCount/results_per_page+1; i++ {
                temp_t_url := t_base_url + strconv.Itoa(i*100)
				res, err := http.Get(temp_t_url)
				fmt.Println("Visiting", temp_t_url)
				if err != nil {
					panic(err.Error())
				}
				temp_body, err := ioutil.ReadAll(res.Body)
				if err != nil {
					panic(err.Error())
				}
				var tempJsonJobs_2 JsonJobs
				err = json.Unmarshal(temp_body, &tempJsonJobs_2)
				if err != nil {
					panic(err.Error())
				}
				jsonJobs_1.Jobs = append(jsonJobs_1.Jobs, tempJsonJobs_2.Jobs...)
                time.Sleep(SecondsSleep * time.Second)
			}
		}

		response_json, err := json.Marshal(jsonJobs_1)
		if err != nil {
			panic(err.Error())
		}
		response = Response{[]byte(response_json)}

		for _, elem := range jsonJobs_1.Jobs {

			result_title := elem.Title
			result_url := elem.Url

			_, err := netUrl.ParseRequestURI(result_url)
			if err == nil {

				elem_json, err := json.Marshal(elem)
				if err != nil {
					panic(err.Error())
				}

				results = append(results, Result{
					runtime.Name,
					result_title,
					result_url,
					elem_json})
			}
        }
	}
	return
}

func (runtime Runtime) Shopify(
	version int, isLocal bool) (
	response Response, results []Result) {

	if version == 1 {

		url := "https://api.lever.co/v0/postings/shopify?group=team&mode=json"

		type Category struct {
			Commitment string `json:"commitment"`
			Location   string `json:"location"`
            Team       string `json:"team"`
            Department  string  `json:"department"`
		}

		type List struct {
			Text    string `json:"text"`
			Content string `json:"content"`
		}

		type Job struct {
			Title            string   `json:"text"`
			Url              string   `json:"hostedUrl"`
			AdditionalPlain  string   `json:"additionalPlain"`
			Additional       string   `json:"additional"`
			Categories       Category `json:"categories"`
			CreatedAt        int      `json:"createdAt"`
			DescriptionPlain string   `json:"descriptionPlain"`
			Description      string   `json:"description"`
			Id               string   `json:"id"`
			Lists            []List   `json:"lists"`
		}

		type Postings struct {
			Jobs []Job `json:"postings"`
		}

		type JsonJobs []Postings

		var body []byte
		if isLocal {
			dir, err := os.Getwd()
			if err != nil {
				panic(err.Error())
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

		response = Response{body}

		var jsonJobs JsonJobs
		err := json.Unmarshal(body, &jsonJobs)
		if err != nil {
			panic(err.Error())
		}

		for _, elem := range jsonJobs {
			for _, subElem := range elem.Jobs {

				result_title := subElem.Title
				result_url := subElem.Url

				_, err := netUrl.ParseRequestURI(result_url)
				if err == nil {

					elem_json, err := json.Marshal(subElem)
					if err != nil {
						panic(err.Error())
					}

					results = append(results, Result{
						runtime.Name,
						result_title,
						result_url,
						elem_json})
				}
			}
		}
	}
	return
}