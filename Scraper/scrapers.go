package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	netUrl "net/url"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
	. "github.com/logrusorgru/aurora"
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
	fmt.Println(Gray(8-1, "Starting Scrape..."))
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
	return
}

func (runtime Runtime) Dreamingjobs(
	version int, isLocal bool) (
	response Response, results []Result) {
	switch version {
	case 1:

		c := colly.NewCollector()

		url := "https://robimalco.github.io/dreamingjobs.github.io/"
		main_tag := "ul"
		main_tag_attr := "class"
		main_tag_value := "position"
		tag_title := "h2"
		tag_url := "a"
		tag_department := "li[class=department]"
		tag_type := "li[class=type]"
		tag_location := "li[class=location]"

		type Job struct {
			Title      string
			Url        string
			Department string
			Type       string
			Location   string
		}

		c.OnHTML(main_tag, func(e *colly.HTMLElement) {
			if strings.Contains(e.Attr(main_tag_attr), main_tag_value) {
				e.ForEach("li", func(_ int, el *colly.HTMLElement) {
					result_title := el.ChildText(tag_title)
					result_url := url + el.ChildAttr(tag_url, "href")
					result_department := el.ChildText(tag_department)
					result_type := el.ChildText(tag_type)
					result_location := el.ChildText(tag_location)

					_, err := netUrl.ParseRequestURI(result_url)
					if err == nil {

						temp_elem_json := Job{
							result_title,
							result_url,
							result_department,
							result_type,
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
			fmt.Println(Gray(8-1, "Visiting"), Gray(8-1, r.URL.String()))
		})

		c.OnError(func(r *colly.Response, err error) {
			fmt.Println(
				Red("Request URL:"), Red(r.Request.URL),
				Red("failed with response:"), Red(r),
				Red("\nError:"), Red(err))
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

func (runtime Runtime) Kununu(
	version int, isLocal bool) (
	response Response, results []Result) {
	switch version {
	case 1:

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
			fmt.Println(Gray(8-1, "Visiting"), Gray(8-1, r.URL.String()))
		})

		c.OnError(func(r *colly.Response, err error) {
			fmt.Println(
				Red("Request URL:"), Red(r.Request.URL),
				Red("failed with response:"), Red(r),
				Red("\nError:"), Red(err))
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
	switch version {
	case 1:

		c := colly.NewCollector()

		s_start_url := "https://api.lever.co/v0/postings/mitte?&mode=json"

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

		var jsonJobs Jobs

		c.OnResponse(func(r *colly.Response) {
			var tempJson Jobs
			err := json.Unmarshal(r.Body, &tempJson)
			if err != nil {
				panic(err.Error())
			}

			for _, elem := range tempJson {

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

			jsonJobs = append(jsonJobs, tempJson...)
		})

		c.OnRequest(func(r *colly.Request) {
			fmt.Println(Gray(8-1, "Visiting"), Gray(8-1, r.URL.String()))
		})

		c.OnScraped(func(r *colly.Response) {
			jsonJobs_marshal, err := json.Marshal(jsonJobs)
			if err != nil {
				panic(err.Error())
			}
			response = Response{[]byte(jsonJobs_marshal)}
		})

		c.OnError(func(r *colly.Response, err error) {
			fmt.Println(
				Red("Request URL:"), Red(r.Request.URL),
				Red("failed with response:"), Red(r),
				Red("\nError:"), Red(err))
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
			c.Visit(s_start_url)
		}
	}
	return
}

func (runtime Runtime) IMusician(
	version int, isLocal bool) (
	response Response, results []Result) {
	switch version {
	case 1:

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
			fmt.Println(Gray(8-1, "Visiting"), Gray(8-1, r.URL.String()))
		})

		c.OnError(func(r *colly.Response, err error) {
			fmt.Println(
				Red("Request URL:"), Red(r.Request.URL),
				Red("failed with response:"), Red(r),
				Red("\nError:"), Red(err))
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
	switch version {
	case 1:

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
			fmt.Println(Gray(8-1, "Visiting"), Gray(8-1, r.URL.String()))
		})

		c.OnError(func(r *colly.Response, err error) {
			fmt.Println(
				Red("Request URL:"), Red(r.Request.URL),
				Red("failed with response:"), Red(r),
				Red("\nError:"), Red(err))
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
	switch version {
	case 1:

		c := colly.NewCollector()

		z_start_url := "https://jobs.zalando.com/api/jobs/?limit=100&offset=0"
		z_base_url := "https://jobs.zalando.com"
		z_base_result_url := "https://jobs.zalando.com/de/jobs/"

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

		var jsonJobs Jobs

		c.OnResponse(func(r *colly.Response) {
			var tempJsonJobs Jobs
			err := json.Unmarshal(r.Body, &tempJsonJobs)
			if err != nil {
				panic(err.Error())
			}

			for _, elem := range tempJsonJobs.Data {

				result_title := elem.Title
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

			jsonJobs.Data = append(jsonJobs.Data, tempJsonJobs.Data...)

			if tempJsonJobs.Next != "" {
				time.Sleep(SecondsSleep * time.Second)
				c.Visit(z_base_url + tempJsonJobs.Next)
			}
		})

		c.OnRequest(func(r *colly.Request) {
			fmt.Println(Gray(8-1, "Visiting"), Gray(8-1, r.URL.String()))
		})

		c.OnScraped(func(r *colly.Response) {
			jsonJobs_marshal, err := json.Marshal(jsonJobs)
			if err != nil {
				panic(err.Error())
			}
			response = Response{[]byte(jsonJobs_marshal)}
		})

		c.OnError(func(r *colly.Response, err error) {
			fmt.Println(
				Red("Request URL:"), Red(r.Request.URL),
				Red("failed with response:"), Red(r),
				Red("\nError:"), Red(err))
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
			c.Visit(z_start_url)
		}
	}
	return
}

func (runtime Runtime) Google(
	version int, isLocal bool) (response Response, results []Result) {
	switch version {
	case 1:

		c := colly.NewCollector()

		g_start_url := "https://careers.google.com/api/v2/jobs/search/?page_size=100&page=1"
		g_base_url := "https://careers.google.com/api/v2/jobs/search/?page_size=100&page="
		g_base_result_url := "https://careers.google.com/jobs/results/"
		number_results_per_page := 100

		type JsonJobs struct {
			Count    int `json:"count"`
			NextPage int `json:"next_page"`
			Jobs     []struct {
				Description   string    `json:"description"`
				CompanyID     string    `json:"company_id"`
				Locations     []string  `json:"locations"`
				Summary       string    `json:"summary"`
				LocationCount int       `json:"location_count"`
				PublishDate   time.Time `json:"publish_date"`
				CompanyName   string    `json:"company_name"`
				JobTitle      string    `json:"job_title"`
				JobID         string    `json:"job_id"`
			} `json:"jobs"`
			PageSize int `json:"page_size"`
		}

		var jsonJobs JsonJobs

		c.OnResponse(func(r *colly.Response) {
			var tempJsonJobs JsonJobs
			err := json.Unmarshal(r.Body, &tempJsonJobs)
			if err != nil {
				panic(err.Error())
			}

			for _, elem := range tempJsonJobs.Jobs {

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

			jsonJobs.Jobs = append(jsonJobs.Jobs, tempJsonJobs.Jobs...)

			/**
						total_count, err := strconv.Atoi(tempJsonJobs.Count)
						if err != nil {
							total_count = 0
						}

						next_page, err := strconv.Atoi(tempJsonJobs.NextPage)
						if err != nil {
							next_page = 0
						}

			            total_pages := total_count/number_results_per_page + 2
						if total_pages <= next_page {
							return
						}

						if next_page != 0 {
							time.Sleep(SecondsSleep * time.Second)
							c.Visit(g_base_url + tempJsonJobs.NextPage)
			            }
			*/

			total_pages := tempJsonJobs.Count/number_results_per_page + 2

			if total_pages <= tempJsonJobs.NextPage {
				return
			}

			if tempJsonJobs.NextPage != 0 {
				time.Sleep(SecondsSleep * time.Second)
				c.Visit(g_base_url + strconv.Itoa(tempJsonJobs.NextPage))
			}
		})

		c.OnRequest(func(r *colly.Request) {
			fmt.Println(Gray(8-1, "Visiting"), Gray(8-1, r.URL.String()))
		})

		c.OnScraped(func(r *colly.Response) {
			jsonJobs_marshal, err := json.Marshal(jsonJobs)
			if err != nil {
				panic(err.Error())
			}
			response = Response{[]byte(jsonJobs_marshal)}
		})

		c.OnError(func(r *colly.Response, err error) {
			fmt.Println(
				Red("Request URL:"), Red(r.Request.URL),
				Red("failed with response:"), Red(r),
				Red("\nError:"), Red(err))
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
			c.Visit(g_start_url)
		}
	}
	return
}

func (runtime Runtime) Soundcloud(
	version int, isLocal bool) (
	response Response, results []Result) {
	switch version {
	case 1:

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
			fmt.Println(Gray(8-1, "Visiting"), Gray(8-1, r.URL.String()))
		})

		c.OnError(func(r *colly.Response, err error) {
			fmt.Println(
				Red("Request URL:"), Red(r.Request.URL),
				Red("failed with response:"), Red(r),
				Red("\nError:"), Red(err))
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
	switch version {
	case 1:

		c := colly.NewCollector()

		m_start_url := "https://careers.microsoft.com/us/en/search-results?s=1&from=1"
		m_base_url := "https://careers.microsoft.com/us/en/search-results?s=1&from="
		m_base_result_url := "https://careers.microsoft.com/us/en/job/"
		first_part_json := `"eagerLoadRefineSearch":`
		second_part_json := `}; phApp.sessionParams`
		counter := 0
		number_results_per_page := 10 // len(jsonJobs.Data.Jobs)

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

		var jsonJobs JsonJobs

		c.OnResponse(func(r *colly.Response) {
			var resultsJson []byte
			if isLocal {
				resultsJson = r.Body
			} else {
				temp_resultsJson := strings.Split(string(r.Body), first_part_json)[1]
				s_resultsJson := strings.Split(temp_resultsJson, second_part_json)[0]
				resultsJson = []byte(s_resultsJson)
			}

			var tempJsonJobs JsonJobs
			err := json.Unmarshal(resultsJson, &tempJsonJobs)
			if err != nil {
				panic(err.Error())
			}

			for _, elem := range tempJsonJobs.Data.Jobs {

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

			jsonJobs.Data.Jobs = append(jsonJobs.Data.Jobs, tempJsonJobs.Data.Jobs...)

			total_pages := tempJsonJobs.TotalHits/number_results_per_page + 2

			if isLocal {
				return
			} else {
				if counter >= total_pages {
					return
				} else {
					counter = counter + 1
					time.Sleep(SecondsSleep * time.Second)
					temp_m_url := m_base_url + strconv.Itoa(counter*number_results_per_page)
					c.Visit(temp_m_url)
				}
			}
		})

		c.OnRequest(func(r *colly.Request) {
			fmt.Println(Gray(8-1, "Visiting"), Gray(8-1, r.URL.String()))
		})

		c.OnScraped(func(r *colly.Response) {
			jsonJobs_marshal, err := json.Marshal(jsonJobs)
			if err != nil {
				panic(err.Error())
			}
			response = Response{[]byte(jsonJobs_marshal)}
		})

		c.OnError(func(r *colly.Response, err error) {
			fmt.Println(
				Red("Request URL:"), Red(r.Request.URL),
				Red("failed with response:"), Red(r),
				Red("\nError:"), Red(err))
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
			c.Visit(m_start_url)
		}
	}
	return
}

func (runtime Runtime) Twitter(
	version int, isLocal bool) (response Response, results []Result) {
	switch version {
	case 1:

		c := colly.NewCollector()

		t_start_url := "https://careers.twitter.com/content/careers-twitter/en/jobs.careers.search.json?limit=100&offset=0"
		t_base_url := "https://careers.twitter.com/content/careers-twitter/en/jobs.careers.search.json?limit=100&offset="
		counter := 0
		number_results_per_page := 100

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

		var jsonJobs Jobs

		c.OnResponse(func(r *colly.Response) {
			var tempJsonJobs Jobs
			err := json.Unmarshal(r.Body, &tempJsonJobs)
			if err != nil {
				panic(err.Error())
			}

			for _, elem := range tempJsonJobs.Results {

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

			jsonJobs.Results = append(jsonJobs.Results, tempJsonJobs.Results...)

			total_pages := tempJsonJobs.TotalCount/number_results_per_page + 1

			if isLocal {
				return
			} else {
				if counter >= total_pages {
					return
				} else {
					counter = counter + 1
					time.Sleep(SecondsSleep * time.Second)
					temp_t_url := t_base_url + strconv.Itoa(counter*100)
					c.Visit(temp_t_url)
				}
			}
		})

		c.OnRequest(func(r *colly.Request) {
			fmt.Println(Gray(8-1, "Visiting"), Gray(8-1, r.URL.String()))
		})

		c.OnScraped(func(r *colly.Response) {
			jsonJobs_marshal, err := json.Marshal(jsonJobs)
			if err != nil {
				panic(err.Error())
			}
			response = Response{[]byte(jsonJobs_marshal)}
		})

		c.OnError(func(r *colly.Response, err error) {
			fmt.Println(
				Red("Request URL:"), Red(r.Request.URL),
				Red("failed with response:"), Red(r),
				Red("\nError:"), Red(err))
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
			c.Visit(t_start_url)
		}
	}
	return
}

func (runtime Runtime) Shopify(
	version int, isLocal bool) (
	response Response, results []Result) {
	switch version {
	case 1:

		c := colly.NewCollector()

		s_start_url := "https://api.lever.co/v0/postings/shopify?mode=json"

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

		var jsonJobs Jobs

		c.OnResponse(func(r *colly.Response) {
			var tempJson Jobs
			err := json.Unmarshal(r.Body, &tempJson)
			if err != nil {
				panic(err.Error())
			}

			for _, elem := range tempJson {

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

			jsonJobs = append(jsonJobs, tempJson...)
		})

		c.OnRequest(func(r *colly.Request) {
			fmt.Println(Gray(8-1, "Visiting"), Gray(8-1, r.URL.String()))
		})

		c.OnScraped(func(r *colly.Response) {
			jsonJobs_marshal, err := json.Marshal(jsonJobs)
			if err != nil {
				panic(err.Error())
			}
			response = Response{[]byte(jsonJobs_marshal)}
		})

		c.OnError(func(r *colly.Response, err error) {
			fmt.Println(
				Red("Request URL:"), Red(r.Request.URL),
				Red("failed with response:"), Red(r),
				Red("\nError:"), Red(err))
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
			c.Visit(s_start_url)
		}
	}
	return
}

func (runtime Runtime) Urbansport(
	version int, isLocal bool) (
	response Response, results []Result) {
	switch version {
	case 1:

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
			fmt.Println(Gray(8-1, "Visiting"), Gray(8-1, r.URL.String()))
		})

		c.OnError(func(r *colly.Response, err error) {
			fmt.Println(
				Red("Request URL:"), Red(r.Request.URL),
				Red("failed with response:"), Red(r),
				Red("\nError:"), Red(err))
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

func (runtime Runtime) N26(
	version int, isLocal bool) (
	response Response, results []Result) {
	switch version {
	case 1:

		c := colly.NewCollector()
		l := c.Clone()

		n_start_url := "https://n26.com/en/careers"
		n_base_url := "https://www.n26.com/"
		main_tag := "a"
		main_attr := "href"
		string_location_url := "locations"
		string_result_url := "positions"
		tag_title := "div"
		tag_details := "dd"

		c.OnHTML(main_tag, func(e *colly.HTMLElement) {
			if strings.Contains(e.Attr(main_attr), string_location_url) {
				temp_location_url := e.Attr(main_attr)
				location_url := n_base_url + temp_location_url
				l.Visit(location_url)
			}
		})

		c.OnResponse(func(r *colly.Response) {
			if isLocal {
				type JsonJob struct {
					CompanyName string          `json:"CompanyName"`
					Title       string          `json:"Title"`
					Url         string          `json:"ResultUrl"`
					Data        json.RawMessage `json:"Data"`
				}

				jobs := make([]JsonJob, 0)
				json.Unmarshal(r.Body, &jobs)

				for _, elem := range jobs {
					results = append(results, Result{
						runtime.Name,
						elem.Title,
						elem.Url,
						elem.Data,
					})
				}
			}
		})

		c.OnRequest(func(r *colly.Request) {
			fmt.Println(Gray(8-1, "Visiting"), Gray(8-1, r.URL.String()))
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
				Red("Request URL:"), Red(r.Request.URL),
				Red("failed with response:"), Red(r),
				Red("\nError:"), Red(err))
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

				type Job struct {
					Title    string
					Url      string
					Location string
					Contract string
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

		l.OnRequest(func(r *colly.Request) {
			fmt.Println(Gray(8-1, "Visiting"), Gray(8-1, r.URL.String()))
		})

		l.OnError(func(r *colly.Response, err error) {
			fmt.Println(
				Red("Request URL:"), Red(r.Request.URL),
				Red("failed with response:"), Red(r),
				Red("\nError:"), Red(err))
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
			c.Visit(n_start_url)
		}
	}
	return
}

func (runtime Runtime) Blinkist(
	version int, isLocal bool) (
	response Response, results []Result) {
	switch version {
	case 1:

		c := colly.NewCollector()

		s_start_url := "https://api.lever.co/v0/postings/blinkist?&mode=json"

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

		var jsonJobs Jobs

		c.OnResponse(func(r *colly.Response) {
			var tempJson Jobs
			err := json.Unmarshal(r.Body, &tempJson)
			if err != nil {
				panic(err.Error())
			}

			for _, elem := range tempJson {

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

			jsonJobs = append(jsonJobs, tempJson...)
		})

		c.OnRequest(func(r *colly.Request) {
			fmt.Println(Gray(8-1, "Visiting"), Gray(8-1, r.URL.String()))
		})

		c.OnScraped(func(r *colly.Response) {
			jsonJobs_marshal, err := json.Marshal(jsonJobs)
			if err != nil {
				panic(err.Error())
			}
			response = Response{[]byte(jsonJobs_marshal)}
		})

		c.OnError(func(r *colly.Response, err error) {
			fmt.Println(
				Red("Request URL:"), Red(r.Request.URL),
				Red("failed with response:"), Red(r),
				Red("\nError:"), Red(err))
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
			c.Visit(s_start_url)
		}
	}
	return
}

func (runtime Runtime) Deutschebahn(
	version int, isLocal bool) (
	response Response, results []Result) {
	switch version {
	case 1:

		c := colly.NewCollector()

		start := 301
		end := 400
		counter := start
		d_start_url := "https://karriere.deutschebahn.com/service/search/karriere-de/2653760?pageNum=" + strconv.Itoa(start)
		d_base_url := "https://karriere.deutschebahn.com/service/search/karriere-de/2653760?pageNum="
		d_job_url := "https://karriere.deutschebahn.com/"
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
					job_publication := strings.TrimSpace(
						el.DOM.Find("span[class=publication]").Text(),
					)
					job_description := strings.TrimSpace(
						el.DOM.Find("p[class=responsibilities-text]").Text(),
					)

					temp_job_url = d_job_url + temp_job_url
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
			next_page_url := d_base_url + e.Text

			counter = counter + 1

			if counter > end {
				return
			} else {
				e.Request.Visit(next_page_url)
			}
		})

		c.OnResponse(func(r *colly.Response) {
			if isLocal {

				type JsonJob struct {
					CompanyName string
					Title       string
					Url         string
					Data        json.RawMessage
				}

				jobs := make([]JsonJob, 0)
				json.Unmarshal(r.Body, &jobs)

				for _, elem := range jobs {
					results = append(results, Result{
						runtime.Name,
						elem.Title,
						elem.Url,
						elem.Data,
					})
				}
			}
		})

		c.OnRequest(func(r *colly.Request) {
			fmt.Println(Gray(8-1, "Visiting"), Gray(8-1, r.URL.String()))
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
				Red("Request URL:"), Red(r.Request.URL),
				Red("failed with response:"), Red(r),
				Red("\nError:"), Red(err))
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
			c.Visit(d_start_url)
		}
	}
	return
}

func (runtime Runtime) Celo(
	version int, isLocal bool) (
	response Response, results []Result) {
	switch version {
	case 1:

		c := colly.NewCollector()

		c_start_url := "https://api.lever.co/v0/postings/celo?mode=json"

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

		var jsonJobs Jobs

		c.OnResponse(func(r *colly.Response) {
			var tempJson Jobs
			err := json.Unmarshal(r.Body, &tempJson)
			if err != nil {
				panic(err.Error())
			}

			for _, elem := range tempJson {

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

			jsonJobs = append(jsonJobs, tempJson...)
		})

		c.OnRequest(func(r *colly.Request) {
			fmt.Println(Gray(8-1, "Visiting"), Gray(8-1, r.URL.String()))
		})

		c.OnScraped(func(r *colly.Response) {
			jsonJobs_marshal, err := json.Marshal(jsonJobs)
			if err != nil {
				panic(err.Error())
			}
			response = Response{[]byte(jsonJobs_marshal)}
		})

		c.OnError(func(r *colly.Response, err error) {
			fmt.Println(
				Red("Request URL:"), Red(r.Request.URL),
				Red("failed with response:"), Red(r),
				Red("\nError:"), Red(err))
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
			c.Visit(c_start_url)
		}
	}
	return
}

func (runtime Runtime) Penta(
	version int, isLocal bool) (
	response Response, results []Result) {
	switch version {
	case 1:

		c := colly.NewCollector()

		p_start_url := "https://penta.recruitee.com/api/offers"

		type Jobs struct {
			Offers []struct {
				ID                 int           `json:"id"`
				Slug               string        `json:"slug"`
				Position           int           `json:"position"`
				Status             string        `json:"status"`
				OptionsPhone       string        `json:"options_phone"`
				OptionsPhoto       string        `json:"options_photo"`
				OptionsCoverLetter string        `json:"options_cover_letter"`
				OptionsCv          string        `json:"options_cv"`
				Remote             interface{}   `json:"remote"`
				CountryCode        string        `json:"country_code"`
				StateCode          string        `json:"state_code"`
				PostalCode         string        `json:"postal_code"`
				MinHours           interface{}   `json:"min_hours"`
				MaxHours           interface{}   `json:"max_hours"`
				OpenQuestions      []interface{} `json:"open_questions"`
				Title              string        `json:"title"`
				Description        string        `json:"description"`
				Requirements       string        `json:"requirements"`
				Location           string        `json:"location"`
				City               string        `json:"city"`
				Country            string        `json:"country"`
				CareersURL         string        `json:"careers_url"`
				CareersApplyURL    string        `json:"careers_apply_url"`
				MailboxEmail       string        `json:"mailbox_email"`
				CompanyName        string        `json:"company_name"`
				Department         string        `json:"department"`
				CreatedAt          string        `json:"created_at"`
				EmploymentTypeCode string        `json:"employment_type_code"`
				CategoryCode       string        `json:"category_code"`
				ExperienceCode     string        `json:"experience_code"`
				EducationCode      string        `json:"education_code"`
				Tags               []interface{} `json:"tags"`
				Translations       struct {
					En struct {
						Title        string `json:"title"`
						Description  string `json:"description"`
						Requirements string `json:"requirements"`
					} `json:"en"`
				} `json:"translations"`
			} `json:"offers"`
		}

		var jsonJobs Jobs

		c.OnResponse(func(r *colly.Response) {
			var tempJson Jobs
			err := json.Unmarshal(r.Body, &tempJson)
			if err != nil {
				panic(err.Error())
			}

			for _, elem := range tempJson.Offers {

				result_title := elem.Title
				result_url := elem.CareersURL

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

			jsonJobs.Offers = append(jsonJobs.Offers, tempJson.Offers...)
		})

		c.OnRequest(func(r *colly.Request) {
			fmt.Println(Gray(8-1, "Visiting"), Gray(8-1, r.URL.String()))
		})

		c.OnScraped(func(r *colly.Response) {
			jsonJobs_marshal, err := json.Marshal(jsonJobs)
			if err != nil {
				panic(err.Error())
			}
			response = Response{[]byte(jsonJobs_marshal)}
		})

		c.OnError(func(r *colly.Response, err error) {
			fmt.Println(
				Red("Request URL:"), Red(r.Request.URL),
				Red("failed with response:"), Red(r),
				Red("\nError:"), Red(err))
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
			c.Visit(p_start_url)
		}

	case 2:

		c := colly.NewCollector()

		url := "https://boards.greenhouse.io/embed/job_board?for=penta"
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
			fmt.Println(Gray(8-1, "Visiting"), Gray(8-1, r.URL.String()))
		})

		c.OnError(func(r *colly.Response, err error) {
			fmt.Println(
				Red("Request URL:"), Red(r.Request.URL),
				Red("failed with response:"), Red(r),
				Red("\nError:"), Red(err))
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

func (runtime Runtime) Contentful(
	version int, isLocal bool) (
	response Response, results []Result) {
	switch version {
	case 1:

		c := colly.NewCollector()

		url := "https://boards.greenhouse.io/embed/job_board?for=contentful"
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
			fmt.Println(Gray(8-1, "Visiting"), Gray(8-1, r.URL.String()))
		})

		c.OnError(func(r *colly.Response, err error) {
			fmt.Println(
				Red("Request URL:"), Red(r.Request.URL),
				Red("failed with response:"), Red(r),
				Red("\nError:"), Red(err))
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

func (runtime Runtime) Gympass(
	version int, isLocal bool) (
	response Response, results []Result) {
	switch version {
	case 1:

		c := colly.NewCollector()

		url := "https://boards.greenhouse.io/embed/job_board?for=gympass"
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
			fmt.Println(Gray(8-1, "Visiting"), Gray(8-1, r.URL.String()))
		})

		c.OnError(func(r *colly.Response, err error) {
			fmt.Println(
				Red("Request URL:"), Red(r.Request.URL),
				Red("failed with response:"), Red(r),
				Red("\nError:"), Red(err))
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

func (runtime Runtime) Hometogo(
	version int, isLocal bool) (
	response Response, results []Result) {
	switch version {
	case 1:

		c := colly.NewCollector()

		h_start_url := "https://api.heavenhr.com/api/v1/positions/public/vacancies/?companyId=_VBAnjTs72rz0J-zBe1sYtA_"
		h_job_url := "https://hometogo.heavenhr.com/jobs/"

		type Jobs struct {
			Links []interface{} `json:"links"`
			Data  []struct {
				ID                  string      `json:"id"`
				Email               interface{} `json:"email"`
				JobTitle            string      `json:"jobTitle"`
				EmploymentTypes     []string    `json:"employmentTypes"`
				Location            string      `json:"location"`
				Department          string      `json:"department"`
				PublicationDate     string      `json:"publicationDate"`
				Status              string      `json:"status"`
				Industry            interface{} `json:"industry"`
				FieldOfWork         interface{} `json:"fieldOfWork"`
				PositionType        interface{} `json:"positionType"`
				Seniority           interface{} `json:"seniority"`
				EmploymentStartDate interface{} `json:"employmentStartDate"`
				HiringOrganization  string      `json:"hiringOrganization"`
				Qualifications      string      `json:"qualifications"`
				Responsibilities    string      `json:"responsibilities"`
				Incentives          string      `json:"incentives"`
				Contact             string      `json:"contact"`
			} `json:"data"`
			Meta struct {
				Page     int `json:"page"`
				PageSize int `json:"pageSize"`
				Count    int `json:"count"`
			} `json:"meta"`
		}

		var jsonJobs Jobs

		c.OnResponse(func(r *colly.Response) {
			var tempJson Jobs
			err := json.Unmarshal(r.Body, &tempJson)
			if err != nil {
				panic(err.Error())
			}

			for _, elem := range tempJson.Data {

				result_title := elem.JobTitle
				result_url := h_job_url + elem.ID + "/apply"

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

			jsonJobs.Data = append(jsonJobs.Data, tempJson.Data...)
		})

		c.OnRequest(func(r *colly.Request) {
			fmt.Println(Gray(8-1, "Visiting"), Gray(8-1, r.URL.String()))
		})

		c.OnScraped(func(r *colly.Response) {
			jsonJobs_marshal, err := json.Marshal(jsonJobs)
			if err != nil {
				panic(err.Error())
			}
			response = Response{[]byte(jsonJobs_marshal)}
		})

		c.OnError(func(r *colly.Response, err error) {
			fmt.Println(
				Red("Request URL:"), Red(r.Request.URL),
				Red("failed with response:"), Red(r),
				Red("\nError:"), Red(err))
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
			c.Visit(h_start_url)
		}
	}
	return
}

func (runtime Runtime) Amazon(
	version int, isLocal bool) (response Response, results []Result) {
	switch version {
	case 1:

		c := colly.NewCollector()

		a_start_url := "https://www.amazon.jobs/en/search.json?loc_query=Belgium&country=BEL&result_limit=1000&offset="
		a_job_url := "https://www.amazon.jobs"
		number_results_per_page := 1000
		counter := 0

		type JsonJobs struct {
			Error  interface{} `json:"error"`
			Hits   int         `json:"hits"`
			Facets struct {
			} `json:"facets"`
			Jobs []struct {
				BasicQualifications     string      `json:"basic_qualifications"`
				BusinessCategory        string      `json:"business_category"`
				City                    string      `json:"city"`
				CompanyName             string      `json:"company_name"`
				CountryCode             string      `json:"country_code"`
				Description             string      `json:"description"`
				DescriptionShort        string      `json:"description_short"`
				DisplayDistance         interface{} `json:"display_distance"`
				ID                      string      `json:"id"`
				IDIcims                 string      `json:"id_icims"`
				JobCategory             string      `json:"job_category"`
				JobFamily               string      `json:"job_family"`
				JobPath                 string      `json:"job_path"`
				JobScheduleType         string      `json:"job_schedule_type"`
				Location                string      `json:"location"`
				NormalizedLocation      string      `json:"normalized_location"`
				OptionalSearchLabels    []string    `json:"optional_search_labels"`
				PostedDate              string      `json:"posted_date"`
				PreferredQualifications interface{} `json:"preferred_qualifications"`
				PrimarySearchLabel      interface{} `json:"primary_search_label"`
				SourceSystem            string      `json:"source_system"`
				State                   interface{} `json:"state"`
				Title                   string      `json:"title"`
				UniversityJob           interface{} `json:"university_job"`
				UpdatedTime             string      `json:"updated_time"`
				URLNextStep             string      `json:"url_next_step"`
				Team                    struct {
					ID                   interface{} `json:"id"`
					BusinessCategoryID   interface{} `json:"business_category_id"`
					Identifier           interface{} `json:"identifier"`
					Label                interface{} `json:"label"`
					CreatedAt            interface{} `json:"created_at"`
					UpdatedAt            interface{} `json:"updated_at"`
					ImageFileName        interface{} `json:"image_file_name"`
					ImageContentType     interface{} `json:"image_content_type"`
					ImageFileSize        interface{} `json:"image_file_size"`
					ImageUpdatedAt       interface{} `json:"image_updated_at"`
					ThumbnailFileName    interface{} `json:"thumbnail_file_name"`
					ThumbnailContentType interface{} `json:"thumbnail_content_type"`
					ThumbnailFileSize    interface{} `json:"thumbnail_file_size"`
					ThumbnailUpdatedAt   interface{} `json:"thumbnail_updated_at"`
					HideJobs             interface{} `json:"hide_jobs"`
					Title                interface{} `json:"title"`
					Headline             interface{} `json:"headline"`
					Description          interface{} `json:"description"`
				} `json:"team"`
			} `json:"jobs"`
		}

		var jsonJobs JsonJobs

		c.OnResponse(func(r *colly.Response) {
			var tempJsonJobs JsonJobs
			err := json.Unmarshal(r.Body, &tempJsonJobs)
			if err != nil {
				panic(err.Error())
			}

			for _, elem := range tempJsonJobs.Jobs {

				result_title := elem.Title
				result_url := a_job_url + elem.JobPath

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

			jsonJobs.Jobs = append(jsonJobs.Jobs, tempJsonJobs.Jobs...)

			if isLocal {
				return
			} else {
				total_pages := tempJsonJobs.Hits / number_results_per_page
				if counter < total_pages+1 {
					counter = counter + 1
					next_page := a_start_url + strconv.Itoa(counter*1000)
					time.Sleep(SecondsSleep * time.Second)
					c.Visit(next_page)
				}
			}
		})

		c.OnRequest(func(r *colly.Request) {
			fmt.Println(Gray(8-1, "Visiting"), Gray(8-1, r.URL.String()))
		})

		c.OnScraped(func(r *colly.Response) {
			jsonJobs_marshal, err := json.Marshal(jsonJobs)
			if err != nil {
				panic(err.Error())
			}
			response = Response{[]byte(jsonJobs_marshal)}
		})

		c.OnError(func(r *colly.Response, err error) {
			fmt.Println(
				Red("Request URL:"), Red(r.Request.URL),
				Red("failed with response:"), Red(r),
				Red("\nError:"), Red(err))
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
			c.Visit(a_start_url + "0")
		}
	}
	return
}

func (runtime Runtime) Lanalabs(
	version int, isLocal bool) (
	response Response, results []Result) {
	switch version {
	case 1:

		c := colly.NewCollector()

		url := "https://lana-labs.breezy.hr/"
		main_tag := "ul"
		main_tag_attr := "class"
		main_tag_value := "position"
		tag_title := "h2"
		tag_url := "a"
		tag_department := "li[class=department]"
		tag_type := "li[class=type]"
		tag_location := "li[class=location]"

		type Job struct {
			Title      string
			Url        string
			Department string
			Type       string
			Location   string
		}

		c.OnHTML(main_tag, func(e *colly.HTMLElement) {
			if strings.Contains(e.Attr(main_tag_attr), main_tag_value) {
				e.ForEach("li", func(_ int, el *colly.HTMLElement) {
					result_title := el.ChildText(tag_title)
					result_url := url + el.ChildAttr(tag_url, "href")
					result_department := el.ChildText(tag_department)
					result_type := el.ChildText(tag_type)
					result_location := el.ChildText(tag_location)

					_, err := netUrl.ParseRequestURI(result_url)
					if err == nil {

						temp_elem_json := Job{
							result_title,
							result_url,
							result_department,
							result_type,
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
			fmt.Println(Gray(8-1, "Visiting"), Gray(8-1, r.URL.String()))
		})

		c.OnError(func(r *colly.Response, err error) {
			fmt.Println(
				Red("Request URL:"), Red(r.Request.URL),
				Red("failed with response:"), Red(r),
				Red("\nError:"), Red(err))
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

func (runtime Runtime) Slack(
	version int, isLocal bool) (
	response Response, results []Result) {
	switch version {
	case 1:

		c := colly.NewCollector()

		url := "https://slack.com/intl/de-de/careers?eu_nc=1#opening"
		main_tag := ".shadow-table"
		sub_tag := "table"
		tag_division := "th"
		tag_data := "tr"
		sub_tag_data := ".for-desktop-only--table-cell"
		tag_url := "a"
		attr_url := "href"

		type Job struct {
			Title    string
			Url      string
			Location string
			Division string
		}

		c.OnHTML(main_tag, func(e *colly.HTMLElement) {
			e.ForEach(sub_tag, func(_ int, el *colly.HTMLElement) {
				job_division := el.ChildText(tag_division)
				el.ForEach(tag_data, func(_ int, ell *colly.HTMLElement) {
					job_data := ell.ChildTexts(sub_tag_data)
					if len(job_data) > 0 {
						result_title := job_data[0]
						result_url := ell.ChildAttr(tag_url, attr_url)
						result_location := job_data[2]

						temp_elem_json := Job{
							result_title,
							result_url,
							result_location,
							job_division,
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
			})
		})

		c.OnResponse(func(r *colly.Response) {
			response = Response{r.Body}
		})

		c.OnRequest(func(r *colly.Request) {
			fmt.Println(Gray(8-1, "Visiting"), Gray(8-1, r.URL.String()))
		})

		c.OnError(func(r *colly.Response, err error) {
			fmt.Println(
				Red("Request URL:"), Red(r.Request.URL),
				Red("failed with response:"), Red(r),
				Red("\nError:"), Red(err))
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

func (runtime Runtime) Revolut(
	version int, isLocal bool) (
	response Response, results []Result) {
	switch version {
	case 1:

		c := colly.NewCollector()

		c_start_url := "https://api.lever.co/v0/postings/revolut?mode=json"

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

		var jsonJobs Jobs

		c.OnResponse(func(r *colly.Response) {
			var tempJson Jobs
			err := json.Unmarshal(r.Body, &tempJson)
			if err != nil {
				panic(err.Error())
			}

			for _, elem := range tempJson {

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

			jsonJobs = append(jsonJobs, tempJson...)
		})

		c.OnRequest(func(r *colly.Request) {
			fmt.Println(Gray(8-1, "Visiting"), Gray(8-1, r.URL.String()))
		})

		c.OnScraped(func(r *colly.Response) {
			jsonJobs_marshal, err := json.Marshal(jsonJobs)
			if err != nil {
				panic(err.Error())
			}
			response = Response{[]byte(jsonJobs_marshal)}
		})

		c.OnError(func(r *colly.Response, err error) {
			fmt.Println(
				Red("Request URL:"), Red(r.Request.URL),
				Red("failed with response:"), Red(r),
				Red("\nError:"), Red(err))
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
			c.Visit(c_start_url)
		}
	}
	return
}

func (runtime Runtime) Mollie(
	version int, isLocal bool) (
	response Response, results []Result) {
	switch version {
	case 1:

		c := colly.NewCollector()

		c_start_url := "https://api.lever.co/v0/postings/mollie?mode=json"

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

		var jsonJobs Jobs

		c.OnResponse(func(r *colly.Response) {
			var tempJson Jobs
			err := json.Unmarshal(r.Body, &tempJson)
			if err != nil {
				panic(err.Error())
			}

			for _, elem := range tempJson {

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

			jsonJobs = append(jsonJobs, tempJson...)
		})

		c.OnRequest(func(r *colly.Request) {
			fmt.Println(Gray(8-1, "Visiting"), Gray(8-1, r.URL.String()))
		})

		c.OnScraped(func(r *colly.Response) {
			jsonJobs_marshal, err := json.Marshal(jsonJobs)
			if err != nil {
				panic(err.Error())
			}
			response = Response{[]byte(jsonJobs_marshal)}
		})

		c.OnError(func(r *colly.Response, err error) {
			fmt.Println(
				Red("Request URL:"), Red(r.Request.URL),
				Red("failed with response:"), Red(r),
				Red("\nError:"), Red(err))
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
			c.Visit(c_start_url)
		}
	}
	return
}

func (runtime Runtime) Circleci(
	version int, isLocal bool) (
	response Response, results []Result) {
	switch version {
	case 1:

		c := colly.NewCollector()

		url := "https://boards.greenhouse.io/embed/job_board?for=circleci"
		base_url := "https://boards.greenhouse.io/circleci/jobs/"
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
					t_j_url := strings.Split(el.ChildAttr(tag_url, "href"), "=")[1]
					result_url := base_url + t_j_url
					result_location := el.ChildText(tag_location)

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
				})
			}
		})

		c.OnResponse(func(r *colly.Response) {
			response = Response{r.Body}
		})

		c.OnRequest(func(r *colly.Request) {
			fmt.Println(Gray(8-1, "Visiting"), Gray(8-1, r.URL.String()))
		})

		c.OnError(func(r *colly.Response, err error) {
			fmt.Println(
				Red("Request URL:"), Red(r.Request.URL),
				Red("failed with response:"), Red(r),
				Red("\nError:"), Red(err))
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

func (runtime Runtime) Blacklane(
	version int, isLocal bool) (
	response Response, results []Result) {
	switch version {
	case 1:

		c := colly.NewCollector()

		url := "https://boards.greenhouse.io/blacklane"
		base_url := "https://boards.greenhouse.io"
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
                    t_j_url := el.ChildAttr(tag_url, "href")
					result_url := base_url + t_j_url
                    result_location := el.ChildText(tag_location)
                    

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
				})
			}
		})

		c.OnResponse(func(r *colly.Response) {
			response = Response{r.Body}
		})

		c.OnRequest(func(r *colly.Request) {
			fmt.Println(Gray(8-1, "Visiting"), Gray(8-1, r.URL.String()))
		})

		c.OnError(func(r *colly.Response, err error) {
			fmt.Println(
				Red("Request URL:"), Red(r.Request.URL),
				Red("failed with response:"), Red(r),
				Red("\nError:"), Red(err))
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
