package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	netUrl "net/url"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
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
		reflect.ValueOf(isLocal),
	}
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
						result_location,
					}

					elem_json, err := json.Marshal(temp_elem_json)
					if err != nil {
						panic(err.Error())
					}

					results = append(results, Result{
						runtime.Name,
						result_title,
						result_url,
						elem_json,
					})
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

		url := "https://api.lever.co/v0/postings/mitte?&mode=json"

		type Jobs []struct {
			AdditionalPlain string `json:"additionalPlain"`
			Additional      string `json:"additional"`
			Categories      struct {
				Commitment string `json:"commitment"`
				Department string `json:"department"`
				Location   string `json:"location"`
				Team       string `json:"team"`
			} `json:"categories"`
			CreatedAt        int64  `json:"createdAt"`
			DescriptionPlain string `json:"descriptionPlain"`
			Description      string `json:"description"`
			ID               string `json:"id"`
			Lists            []struct {
				Text    string `json:"text"`
				Content string `json:"content"`
			} `json:"lists"`
			Text      string `json:"text"`
			HostedURL string `json:"hostedUrl"`
			ApplyURL  string `json:"applyUrl"`
		}

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

		var jsonJobs Jobs
		err := json.Unmarshal(body, &jsonJobs)
		if err != nil {
			panic(err.Error())
		}

		for _, elem := range jsonJobs {
			result_title := elem.Text
			result_url := elem.HostedURL

			elem_json, err := json.Marshal(elem)
			if err != nil {
				panic(err.Error())
			}

			results = append(results, Result{
				runtime.Name,
				result_title,
				result_url,
				elem_json,
			})
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
						result_location,
					}

					elem_json, err := json.Marshal(temp_elem_json)
					if err != nil {
						panic(err.Error())
					}

					results = append(results, Result{
						runtime.Name,
						result_title,
						result_url,
						elem_json,
					})
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
						result_url,
					}
					elem_json, err := json.Marshal(temp_elem_json)
					if err != nil {
						panic(err.Error())
					}

					results = append(results, Result{
						runtime.Name,
						result_title,
						result_url,
						elem_json,
					})
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

		type Jobs struct {
			Data []struct {
				JobCategories []string  `json:"job_categories"`
				UpdatedAt     time.Time `json:"updated_at"`
				Offices       []string  `json:"offices"`
				ID            int       `json:"id"`
				Title         string    `json:"title"`
				Entity        string    `json:"entity"`
			} `json:"data"`
			Facets struct {
				Offices       []string `json:"offices"`
				JobCategories []string `json:"job_categories"`
				ContractTypes []string `json:"contract_types"`
				EntryLevels   []string `json:"entry_levels"`
				Entity        []string `json:"entity"`
			} `json:"facets"`
			Total int    `json:"total"`
			First string `json:"first"`
			Last  string `json:"last"`
			Next  string `json:"next"`
		}

		var jsonJobs_1 Jobs

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
				var tempJsonJobs_2 Jobs
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
			result_url := z_base_result_url + strconv.Itoa(elem.ID)

			elem_json, err := json.Marshal(elem)
			if err != nil {
				panic(err.Error())
			}

			results = append(results, Result{
				runtime.Name,
				result_title,
				result_url,
				elem_json,
			})
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

		type JsonJobs struct {
			Count string `json:"count"`
			Jobs  []struct {
				CompanyID      string    `json:"company_id"`
				CompanyName    string    `json:"company_name"`
				Description    string    `json:"description"`
				JobID          string    `json:"job_id"`
				JobTitle       string    `json:"job_title"`
				Locations      []string  `json:"locations"`
				LocationsCount string    `json:"locations_count"`
				PublishDate    time.Time `json:"publish_date"`
				Summary        string    `json:"summary"`
			} `json:"jobs"`
			NextPage string `json:"next_page"`
			PageSize string `json:"page_size"`
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

			result_title := elem.JobTitle
			result_url := g_base_result_url + strings.Split(elem.JobID, "/")[1]

			elem_json, err := json.Marshal(elem)
			if err != nil {
				panic(err.Error())
			}

			results = append(results, Result{
				runtime.Name,
				result_title,
				result_url,
				elem_json,
			})
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
				result_department := e.ChildText(tag_department)

				e.ForEach("div", func(_ int, el *colly.HTMLElement) {
					result_title := el.ChildText(tag_title)
					result_url := el.ChildAttr(tag_url, "href")
					result_location := el.ChildText(tag_location)

					_, err := netUrl.ParseRequestURI(result_url)
					if err == nil {

						temp_elem_json := Job{
							result_title,
							result_url,
							result_department,
							result_location,
						}

						elem_json, err := json.Marshal(temp_elem_json)
						if err != nil {
							panic(err.Error())
						}

						results = append(results, Result{
							runtime.Name,
							result_title,
							result_url,
							elem_json,
						})
					}
				})
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

		type JsonJobs struct {
			Status    int `json:"status"`
			Hits      int `json:"hits"`
			TotalHits int `json:"totalHits"`
			Data      struct {
				Jobs []struct {
					Country            string      `json:"country"`
					SubCategory        string      `json:"subCategory"`
					Industry           interface{} `json:"industry"`
					Title              string      `json:"title"`
					MultiLocation      []string    `json:"multi_location"`
					Type               interface{} `json:"type"`
					OrgFunction        interface{} `json:"orgFunction"`
					Experience         string      `json:"experience"`
					Locale             string      `json:"locale"`
					MultiLocationArray []struct {
						Location string `json:"location"`
					} `json:"multi_location_array"`
					JobSeqNo             string      `json:"jobSeqNo"`
					PostedDate           string      `json:"postedDate"`
					SearchresultsDisplay interface{} `json:"searchresults_display"`
					DescriptionTeaser    string      `json:"descriptionTeaser"`
					DateCreated          string      `json:"dateCreated"`
					State                string      `json:"state"`
					TargetLevel          string      `json:"targetLevel"`
					JdDisplay            interface{} `json:"jd_display"`
					ReqID                interface{} `json:"reqId"`
					Badge                string      `json:"badge"`
					JobID                string      `json:"jobId"`
					IsMultiLocation      bool        `json:"isMultiLocation"`
					JobVisibility        []string    `json:"jobVisibility"`
					Mostpopular          float64     `json:"mostpopular"`
					Location             string      `json:"location"`
					Category             string      `json:"category"`
					LocationLatlong      interface{} `json:"locationLatlong"`
				}
			} `json:"data"`
			Eid string `json:"eid"`
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
			fmt.Println("Visiting", m_base_url)
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

			results_per_page := 10 // len(jsonJobs_1.Data.Jobs)

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
			result_url := m_base_result_url + elem.JobID

			elem_json, err := json.Marshal(elem)
			if err != nil {
				panic(err.Error())
			}

			results = append(results, Result{
				runtime.Name,
				result_title,
				result_url,
				elem_json,
			})
		}
	}
	return
}

func (runtime Runtime) Twitter(
	version int, isLocal bool) (response Response, results []Result) {
	if version == 1 {

		t_base_url := "https://careers.twitter.com/content/careers-twitter/en/jobs.careers.search.json?limit=100&offset="

		results_per_page := 100

		type Jobs struct {
			Results []struct {
				Title       string `json:"title"`
				Description string `json:"description"`
				Modified    int64  `json:"modified"`
				InternalJob bool   `json:"internalJob"`
				URL         string `json:"url"`
				Locations   []struct {
					ID    string `json:"id"`
					Title string `json:"title"`
				} `json:"locations"`
				Team struct {
					ID    string `json:"id"`
					Title string `json:"title"`
				} `json:"team"`
			} `json:"results"`
			PageCount  int `json:"pageCount"`
			TotalCount int `json:"totalCount"`
		}

		var jsonJobs_1 Jobs

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
			fmt.Println("Visiting ", t_base_url+"0")
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
				var tempJsonJobs_2 Jobs
				err = json.Unmarshal(temp_body, &tempJsonJobs_2)
				if err != nil {
					panic(err.Error())
				}
				jsonJobs_1.Results = append(
					jsonJobs_1.Results, tempJsonJobs_2.Results...)
				time.Sleep(SecondsSleep * time.Second)
			}
		}

		response_json, err := json.Marshal(jsonJobs_1)
		if err != nil {
			panic(err.Error())
		}
		response = Response{[]byte(response_json)}

		for _, elem := range jsonJobs_1.Results {

			result_title := elem.Title
			result_url := elem.URL

			elem_json, err := json.Marshal(elem)
			if err != nil {
				panic(err.Error())
			}

			results = append(results, Result{
				runtime.Name,
				result_title,
				result_url,
				elem_json,
			})
		}
	}
	return
}

func (runtime Runtime) Shopify(
	version int, isLocal bool) (
	response Response, results []Result) {
	if version == 1 {

		url := "https://api.lever.co/v0/postings/shopify?mode=json"

		type Jobs []struct {
			AdditionalPlain string `json:"additionalPlain"`
			Additional      string `json:"additional"`
			Categories      struct {
				Commitment string `json:"commitment"`
				Department string `json:"department"`
				Location   string `json:"location"`
				Team       string `json:"team"`
			} `json:"categories"`
			CreatedAt        int64  `json:"createdAt"`
			DescriptionPlain string `json:"descriptionPlain"`
			Description      string `json:"description"`
			ID               string `json:"id"`
			Lists            []struct {
				Text    string `json:"text"`
				Content string `json:"content"`
			} `json:"lists"`
			Text      string `json:"text"`
			HostedURL string `json:"hostedUrl"`
			ApplyURL  string `json:"applyUrl"`
		}

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

		var jsonJobs Jobs
		err := json.Unmarshal(body, &jsonJobs)
		if err != nil {
			panic(err.Error())
		}

		response = Response{body}

		for _, elem := range jsonJobs {
			result_title := elem.Text
			result_url := elem.HostedURL

			elem_json, err := json.Marshal(elem)
			if err != nil {
				panic(err.Error())
			}

			results = append(results, Result{
				runtime.Name,
				result_title,
				result_url,
				elem_json,
			})
		}
	}
	return
}

func (runtime Runtime) Urbansport(
	version int, isLocal bool) (
	response Response, results []Result) {
	if version == 1 {

		c := colly.NewCollector()

		url := "https://boards.greenhouse.io/urbansportsclub"
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
				result_department := e.ChildText(tag_department)

				e.ForEach("div", func(_ int, el *colly.HTMLElement) {
					result_title := el.ChildText(tag_title)
					result_url := el.ChildAttr(tag_url, "href")
					result_location := el.ChildText(tag_location)

					_, err := netUrl.ParseRequestURI(result_url)
					if err == nil {

						temp_elem_json := Job{
							result_title,
							result_url,
							result_department,
							result_location,
						}

						elem_json, err := json.Marshal(temp_elem_json)
						if err != nil {
							panic(err.Error())
						}

						results = append(results, Result{
							runtime.Name,
							result_title,
							result_url,
							elem_json,
						})
					}
				})
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

func (runtime Runtime) N26(version int, isLocal bool) (response Response, results []Result) {
	if version == 1 {

		c := colly.NewCollector()
		l := c.Clone()

		url := "https://n26.com/en/careers"
		n_base_url := "https://www.n26.com/"

		main_tag := "a"
		main_attr := "href"
		string_location_url := "locations"
		string_result_url := "positions"

		tag_title := "div"
		tag_details := "dd"

		if isLocal {

			type JsonJob struct {
				CompanyName string          `json:"CompanyName"`
				Title       string          `json:"Title"`
				Url         string          `json:"ResultUrl"`
				Data        json.RawMessage `json:"Data"`
			}

			dir, err := os.Getwd()
			if err != nil {
				panic(err.Error())
			}
			body, err := ioutil.ReadFile(dir + "/response.html")
			fmt.Println("Visiting", dir+"/response.html")
			if err != nil {
				panic(err.Error())
			}

			jobs := make([]JsonJob, 0)
			json.Unmarshal(body, &jobs)

			for _, elem := range jobs {
				results = append(results, Result{
					runtime.Name,
					elem.Title,
					elem.Url,
					elem.Data,
				})
			}
		} else {

			type Job struct {
				Title    string
				Url      string
				Location string
				Contract string
			}

			c.OnHTML(main_tag, func(e *colly.HTMLElement) {
				if strings.Contains(e.Attr(main_attr), string_location_url) {
					temp_location_url := e.Attr(main_attr)
					location_url := n_base_url + temp_location_url
					fmt.Println("Visiting", location_url)
					l.Visit(location_url)
				}
			})

			l.OnHTML(main_tag, func(e *colly.HTMLElement) {
				if strings.Contains(e.Attr(main_attr), string_result_url) {
					temp_result_url := e.Attr(main_attr)
					result_url := n_base_url + temp_result_url

					goquerySelection := e.DOM

					var titles []string
					goquerySelection.Find(tag_title).Each(func(i int, s *goquery.Selection) {
						titles = append(titles, s.Nodes[0].FirstChild.Data)
					})
					result_title := titles[0]

					details_nodes := goquerySelection.Find(tag_details).Nodes
					location := details_nodes[0].FirstChild.Data
					contract := ""
					if len(details_nodes) > 1 {
						contract = details_nodes[1].FirstChild.Data
					}

					temp_elem_json := Job{result_title, result_url, location, contract}

					elem_json, err := json.Marshal(temp_elem_json)
					if err != nil {
						panic(err.Error())
					}

					results = append(results, Result{
						runtime.Name,
						result_title,
						result_url,
						elem_json,
					})
				}
			})

			c.OnScraped(func(r *colly.Response) {
				response_json, err := json.Marshal(results)
				if err != nil {
					panic(err.Error())
				}
				response = Response{[]byte(response_json)}
			})

			c.OnError(func(r *colly.Response, err error) {
				fmt.Println(
					"Request URL:", r.Request.URL,
					"failed with response:", r,
					"\nError:", err)
			})

			l.OnError(func(r *colly.Response, err error) {
				fmt.Println(
					"Request URL:", r.Request.URL,
					"failed with response:", r,
					"\nError:", err)
			})

			c.Visit(url)
		}
	}
	return
}

func (runtime Runtime) Blinkist(
	version int, isLocal bool) (
	response Response, results []Result) {
	if version == 1 {

		url := "https://api.lever.co/v0/postings/blinkist?mode=json"

		type Jobs []struct {
			AdditionalPlain string `json:"additionalPlain"`
			Additional      string `json:"additional"`
			Categories      struct {
				Commitment string `json:"commitment"`
				Department string `json:"department"`
				Location   string `json:"location"`
				Team       string `json:"team"`
			} `json:"categories"`
			CreatedAt        int64  `json:"createdAt"`
			DescriptionPlain string `json:"descriptionPlain"`
			Description      string `json:"description"`
			ID               string `json:"id"`
			Lists            []struct {
				Text    string `json:"text"`
				Content string `json:"content"`
			} `json:"lists"`
			Text      string `json:"text"`
			HostedURL string `json:"hostedUrl"`
			ApplyURL  string `json:"applyUrl"`
		}

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

		var jsonJobs Jobs
		err := json.Unmarshal(body, &jsonJobs)
		if err != nil {
			panic(err.Error())
		}

		for _, elem := range jsonJobs {
			result_title := elem.Text
			result_url := elem.HostedURL

			elem_json, err := json.Marshal(elem)
			if err != nil {
				panic(err.Error())
			}

			results = append(results, Result{
				runtime.Name,
				result_title,
				result_url,
				elem_json,
			})
		}
	}
	return
}

func (runtime Runtime) Deutschebahn(
	version int, isLocal bool) (
	response Response, results []Result) {
	if version == 1 {

		if isLocal {
			fmt.Println("I AM LOCAL")
		} else {

			c := colly.NewCollector()

			start_url := "https://karriere.deutschebahn.com/service/search/karriere-de/2653760?pageNum="
			d_base_url := "https://karriere.deutschebahn.com/"

			main_section_tag := "ul"
			main_section_attr := "class"
			main_section_value := "result-items"

			type Job struct {
				Url         string
				Title       string
				Location    string
				Entity      string
				Publication string
				Description string
			}

			c.OnHTML(main_section_tag, func(e *colly.HTMLElement) {
				if strings.Contains(e.Attr(main_section_attr), main_section_value) {
					e.ForEach("li", func(_ int, el *colly.HTMLElement) {
						temp_job_url, exists := el.DOM.Find("div[class=info]").Find("a").Attr("href")
						_ = exists
						job_title := el.DOM.Find("span[class=title]").Text()
						job_location := strings.TrimSpace(el.DOM.Find("span[class=location]").Text())
						job_entity := strings.TrimSpace(el.DOM.Find("span[class=entity]").Text())
						job_publication := strings.TrimSpace(el.DOM.Find("span[class=publication]").Text())
						job_description := strings.TrimSpace(el.DOM.Find("p[class=responsibilities-text]").Text())

						temp_job_url = d_base_url + temp_job_url
						u, err := netUrl.Parse(temp_job_url)
						if err != nil {
							panic(err.Error())
						}
						u.RawQuery = ""
						job_url := u.String()

						temp_elem_json := Job{
							job_url,
							job_title,
							job_location,
							job_entity,
							job_publication,
							job_description,
						}

						elem_json, err := json.Marshal(temp_elem_json)
						if err != nil {
							panic(err.Error())
						}

						results = append(results, Result{
							runtime.Name,
							job_title,
							job_url,
							elem_json,
						})
					})
				}
			})

			// Find and visit next page links
			c.OnHTML("a[class=active]", func(e *colly.HTMLElement) {
				next_page_url := start_url + e.Text
				fmt.Println("Visiting", next_page_url)
				e.Request.Visit(next_page_url)
			})
			fmt.Println("Visiting", start_url+"0")
			c.Visit(start_url + "0")
		}

		response_json, err := json.Marshal(results)
		if err != nil {
			panic(err.Error())
		}
		response = Response{[]byte(response_json)}

	}
	return
}
