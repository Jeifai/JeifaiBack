package cmd

import (
	"context"
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
	xj "github.com/basgys/goxml2json"
	"github.com/chromedp/chromedp"
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

func Extract(
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
					counter++
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
					counter++
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
		n_base_url := "https://www.n26.com"
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
	case 2:

		c := colly.NewCollector()
		l := c.Clone()

		n_start_url := "https://n26.com/en/careers"
		n_base_url := "https://www.n26.com"
		main_tag := "a"
		main_attr := "href"
		string_location_url := "locations"

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

		l.OnHTML("li", func(e *colly.HTMLElement) {
			e.ForEach("div", func(_ int, el *colly.HTMLElement) {
				if strings.Contains(el.ChildAttr("a", "href"), "positions") {
					temp_result_url := el.ChildAttr("a", "href")

					result_url := n_base_url + temp_result_url

					result_title := el.ChildText("a")

					goquerySelection := el.DOM

					details_nodes := goquerySelection.Find("dd").Nodes
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

		start := 0
		d_start_url := "https://karriere.deutschebahn.com/service/search/karriere-de/2653760?sort=pubExternalDate_td&pageNum=" + strconv.Itoa(start)
		d_base_url := "https://karriere.deutschebahn.com/service/search/karriere-de/2653760?sort=pubExternalDate_td&pageNum="
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

		c.OnHTML("a[class=active]", func(e *colly.HTMLElement) {
			next_page_url := d_base_url + e.Text
			time.Sleep(SecondsSleep * time.Second)
			e.Request.Visit(next_page_url)
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

		a_start_url := "https://www.amazon.jobs/en/search.json?loc_query=Germany&country=DEU&result_limit=1000&offset="
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
					counter++
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

func (runtime Runtime) Auto1(
	version int, isLocal bool) (response Response, results []Result) {
	switch version {
	case 1:

		c := colly.NewCollector()

		a_base_url := "https://www.auto1-group.com/smart-recruiters/jobs/search/?page="
		a_base_job_url := "https://www.auto1-group.com/de/jobs/"
		current_page := 1
		number_results_per_page := 15

		type Auto1Jobs struct {
			Jobs struct {
				Hits struct {
					Total    int         `json:"total"`
					MaxScore interface{} `json:"max_score"`
					Hits     []struct {
						Index  string      `json:"_index"`
						Type   string      `json:"_type"`
						ID     string      `json:"_id"`
						Score  interface{} `json:"_score"`
						Source struct {
							Title string `json:"title"`
							JobAd struct {
								Sections struct {
									CompanyDescription struct {
										Title string `json:"title"`
										Text  string `json:"text"`
									} `json:"companyDescription"`
									JobDescription struct {
										Title string `json:"title"`
										Text  string `json:"text"`
									} `json:"jobDescription"`
									Qualifications struct {
										Title string `json:"title"`
										Text  string `json:"text"`
									} `json:"qualifications"`
									AdditionalInformation struct {
										Title string `json:"title"`
										Text  string `json:"text"`
									} `json:"additionalInformation"`
								} `json:"Jobssections"`
							} `json:"jobAd"`
							LocationCity     string    `json:"locationCity"`
							LocationCountry  string    `json:"locationCountry"`
							Brand            string    `json:"brand"`
							Company          string    `json:"company"`
							ExperienceLevel  string    `json:"experienceLevel"`
							Department       string    `json:"department"`
							TypeOfEmployment string    `json:"typeOfEmployment"`
							CreatedOn        time.Time `json:"createdOn"`
							IsActive         int       `json:"isActive"`
							URL              string    `json:"url"`
						} `json:"_source"`
						Sort []int64 `json:"sort"`
					} `json:"hits"`
				} `json:"hits"`
			} `json:"jobs"`
		}

		var jsonJobs Auto1Jobs

		c.OnResponse(func(r *colly.Response) {
			var tempJsonJobs Auto1Jobs
			err := json.Unmarshal(r.Body, &tempJsonJobs)
			if err != nil {
				panic(err.Error())
			}

			for _, elem := range tempJsonJobs.Jobs.Hits.Hits {

				result_title := elem.Source.Title
				result_url := a_base_job_url + elem.Source.URL

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

			jsonJobs.Jobs.Hits.Hits = append(jsonJobs.Jobs.Hits.Hits, tempJsonJobs.Jobs.Hits.Hits...)

			total_pages := tempJsonJobs.Jobs.Hits.Total/number_results_per_page + 2

			if current_page > total_pages {
				return
			} else {
				time.Sleep(SecondsSleep * time.Second)
				current_page++
				c.Visit(a_base_url + strconv.Itoa(current_page))
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
			c.Visit(a_base_url + strconv.Itoa(current_page))
		}
	}
	return
}

func (runtime Runtime) Flixbus(
	version int, isLocal bool) (
	response Response, results []Result) {
	switch version {
	case 1:

		c := colly.NewCollector()

		c_start_url := "https://flix.careers/api/jobs"

		type FlixbusJobs struct {
			Jobs []struct {
				AbsoluteURL    string `json:"absolute_url"`
				DataCompliance []struct {
					Type            string      `json:"type"`
					RequiresConsent bool        `json:"requires_consent"`
					RetentionPeriod interface{} `json:"retention_period"`
				} `json:"data_compliance"`
				InternalJobID int64 `json:"internal_job_id"`
				Location      struct {
					Name string `json:"name"`
				} `json:"location"`
				Metadata []struct {
					ID        int64  `json:"id"`
					Name      string `json:"name"`
					Value     string `json:"value"`
					ValueType string `json:"value_type"`
				} `json:"metadata"`
				ID            int64  `json:"id"`
				UpdatedAt     string `json:"updated_at"`
				RequisitionID string `json:"requisition_id"`
				Title         string `json:"title"`
				Departments   []struct {
					ID       int64         `json:"id"`
					Name     string        `json:"name"`
					ChildIds []interface{} `json:"child_ids"`
					ParentID interface{}   `json:"parent_id"`
				} `json:"departments"`
				Offices []struct {
					ID       int64         `json:"id"`
					Name     string        `json:"name"`
					Location interface{}   `json:"location"`
					ChildIds []interface{} `json:"child_ids"`
					ParentID int64         `json:"parent_id"`
				} `json:"offices"`
			} `json:"jobs"`
			Meta struct {
				Total int `json:"total"`
			} `json:"meta"`
		}

		var jsonJobs FlixbusJobs

		c.OnResponse(func(r *colly.Response) {
			var tempJson FlixbusJobs
			err := json.Unmarshal(r.Body, &tempJson)
			if err != nil {
				panic(err.Error())
			}

			for _, elem := range tempJson.Jobs {

				result_title := elem.Title
				result_url := elem.AbsoluteURL

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

			jsonJobs.Jobs = append(jsonJobs.Jobs, tempJson.Jobs...)
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

func (runtime Runtime) Quora(
	version int, isLocal bool) (
	response Response, results []Result) {
	switch version {
	case 1:

		c := colly.NewCollector()

		url := "https://boards.greenhouse.io/quora"
		q_base_job_url := "https://boards.greenhouse.io/"
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
					result_url := q_base_job_url + el.ChildAttr(tag_url, "href")
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

func (runtime Runtime) Greenhouse(
	version int, isLocal bool) (
	response Response, results []Result) {
	switch version {
	case 1:

		c := colly.NewCollector()

		url := "https://boards.greenhouse.io/embed/job_board?for=greenhouse"
		main_tag := "section"
		main_tag_attr := "class"
		main_tag_value := "level-0"
		tag_title := "a"
		tag_url := "a"
		tag_department := "h2"
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
					result_url := t_j_url
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

func (runtime Runtime) Docker(
	version int, isLocal bool) (
	response Response, results []Result) {
	switch version {
	case 1:

		c := colly.NewCollector()

		d_start_url := "https://newton.newtonsoftware.com/career/CareerHome.action?clientId=8a7883c6708df1d40170a6df29950b39"
		main_tag := ".gnewtonCareerGroupRowClass"
		tag_title := "a"
		tag_location := ".gnewtonCareerGroupJobDescriptionClass"

		type Job struct {
			Title    string
			Url      string
			Location string
		}

		c.OnHTML(main_tag, func(e *colly.HTMLElement) {
			result_title := e.ChildText(tag_title)
			result_url := e.ChildAttr(tag_title, "href")
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
			c.Visit(d_start_url)
		}
	}
	return
}

func (runtime Runtime) Zapier(
	version int, isLocal bool) (
	response Response, results []Result) {
	switch version {
	case 1:

		c := colly.NewCollector()

		z_start_url := "https://zapier.com/jobs"
		z_base_url := "https://zapier.com"
		main_tag := "section"
		main_tag_attr := "id"
		main_tag_value := "job-openings"
		tag_section_job := "li"
		tag_info := "a"

		type Job struct {
			Title      string
			Url        string
			Department string
		}

		c.OnHTML(main_tag, func(e *colly.HTMLElement) {
			if strings.Contains(e.Attr(main_tag_attr), main_tag_value) {
				e.ForEach(tag_section_job, func(_ int, el *colly.HTMLElement) {
					result_info := el.ChildText(tag_info)
					result_temp_url := el.ChildAttr(tag_info, "href")

					if !strings.Contains(result_temp_url, "https") {

						result_url := z_base_url + result_temp_url

						info_split := strings.Split(result_info, " - ")
						result_department := info_split[0]
						result_title := info_split[1]

						temp_elem_json := Job{
							result_title,
							result_url,
							result_department,
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
			c.Visit(z_start_url)
		}
	}
	return
}

func (runtime Runtime) Datadog(
	version int, isLocal bool) (
	response Response, results []Result) {
	switch version {
	case 1:

		c := colly.NewCollector()

		start_url := "https://api.greenhouse.io/v1/boards/datadog/jobs/"

		type JsonJobs struct {
			Jobs []struct {
				AbsoluteURL    string `json:"absolute_url"`
				DataCompliance []struct {
					Type            string      `json:"type"`
					RequiresConsent bool        `json:"requires_consent"`
					RetentionPeriod interface{} `json:"retention_period"`
				} `json:"data_compliance"`
				InternalJobID int `json:"internal_job_id"`
				Location      struct {
					Name string `json:"name"`
				} `json:"location"`
				Metadata []struct {
					ID        int    `json:"id"`
					Name      string `json:"name"`
					Value     string `json:"value"`
					ValueType string `json:"value_type"`
				} `json:"metadata"`
				ID            int    `json:"id"`
				UpdatedAt     string `json:"updated_at"`
				RequisitionID string `json:"requisition_id"`
				Title         string `json:"title"`
				Education     string `json:"education,omitempty"`
			} `json:"jobs"`
			Meta struct {
				Total int `json:"total"`
			} `json:"meta"`
		}

		var jsonJobs JsonJobs

		c.OnResponse(func(r *colly.Response) {
			var tempJson JsonJobs
			err := json.Unmarshal(r.Body, &tempJson)
			if err != nil {
				panic(err.Error())
			}

			for _, elem := range tempJson.Jobs {

				result_title := elem.Title
				result_url := elem.AbsoluteURL

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

			jsonJobs.Jobs = append(jsonJobs.Jobs, tempJson.Jobs...)
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
			c.Visit(start_url)
		}
	}
	return
}

func (runtime Runtime) Stripe(
	version int, isLocal bool) (
	response Response, results []Result) {
	switch version {
	case 1:

		c := colly.NewCollector()

		start_url := "https://api.greenhouse.io/v1/boards/stripe/jobs/"

		type JsonJobs struct {
			Jobs []struct {
				AbsoluteURL    string `json:"absolute_url"`
				DataCompliance []struct {
					Type            string      `json:"type"`
					RequiresConsent bool        `json:"requires_consent"`
					RetentionPeriod interface{} `json:"retention_period"`
				} `json:"data_compliance"`
				InternalJobID int `json:"internal_job_id"`
				Location      struct {
					Name string `json:"name"`
				} `json:"location"`
				Metadata []struct {
					ID        int    `json:"id"`
					Name      string `json:"name"`
					Value     string `json:"value"`
					ValueType string `json:"value_type"`
				} `json:"metadata"`
				ID            int    `json:"id"`
				UpdatedAt     string `json:"updated_at"`
				RequisitionID string `json:"requisition_id"`
				Title         string `json:"title"`
				Education     string `json:"education,omitempty"`
			} `json:"jobs"`
			Meta struct {
				Total int `json:"total"`
			} `json:"meta"`
		}

		var jsonJobs JsonJobs

		c.OnResponse(func(r *colly.Response) {
			var tempJson JsonJobs
			err := json.Unmarshal(r.Body, &tempJson)
			if err != nil {
				panic(err.Error())
			}

			for _, elem := range tempJson.Jobs {

				result_title := elem.Title
				result_url := elem.AbsoluteURL

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

			jsonJobs.Jobs = append(jsonJobs.Jobs, tempJson.Jobs...)
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
			c.Visit(start_url)
		}
	}
	return
}

func (runtime Runtime) Github(
	version int, isLocal bool) (
	response Response, results []Result) {
	switch version {
	case 1:

		c := colly.NewCollector()

		start_url := "https://api.greenhouse.io/v1/boards/github/jobs/"

		type JsonJobs struct {
			Jobs []struct {
				AbsoluteURL    string `json:"absolute_url"`
				DataCompliance []struct {
					Type            string      `json:"type"`
					RequiresConsent bool        `json:"requires_consent"`
					RetentionPeriod interface{} `json:"retention_period"`
				} `json:"data_compliance"`
				InternalJobID int `json:"internal_job_id"`
				Location      struct {
					Name string `json:"name"`
				} `json:"location"`
				Metadata []struct {
					ID        int    `json:"id"`
					Name      string `json:"name"`
					Value     string `json:"value"`
					ValueType string `json:"value_type"`
				} `json:"metadata"`
				ID            int    `json:"id"`
				UpdatedAt     string `json:"updated_at"`
				RequisitionID string `json:"requisition_id"`
				Title         string `json:"title"`
				Education     string `json:"education,omitempty"`
			} `json:"jobs"`
			Meta struct {
				Total int `json:"total"`
			} `json:"meta"`
		}

		var jsonJobs JsonJobs

		c.OnResponse(func(r *colly.Response) {
			var tempJson JsonJobs
			err := json.Unmarshal(r.Body, &tempJson)
			if err != nil {
				panic(err.Error())
			}

			for _, elem := range tempJson.Jobs {

				result_title := elem.Title
				result_url := elem.AbsoluteURL

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

			jsonJobs.Jobs = append(jsonJobs.Jobs, tempJson.Jobs...)
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
			c.Visit(start_url)
		}
	}
	return
}

func (runtime Runtime) Getyourguide(
	version int, isLocal bool) (
	response Response, results []Result) {
	switch version {
	case 1:

		c := colly.NewCollector()

		start_url := "https://api.greenhouse.io/v1/boards/getyourguide/jobs"

		type JsonJobs struct {
			Jobs []struct {
				AbsoluteURL    string `json:"absolute_url"`
				DataCompliance []struct {
					Type            string      `json:"type"`
					RequiresConsent bool        `json:"requires_consent"`
					RetentionPeriod interface{} `json:"retention_period"`
				} `json:"data_compliance"`
				Education     string `json:"education,omitempty"`
				InternalJobID int    `json:"internal_job_id"`
				Location      struct {
					Name string `json:"name"`
				} `json:"location"`
				Metadata []struct {
					ID        int         `json:"id"`
					Name      string      `json:"name"`
					Value     interface{} `json:"value"`
					ValueType string      `json:"value_type"`
				} `json:"metadata"`
				ID            int    `json:"id"`
				UpdatedAt     string `json:"updated_at"`
				RequisitionID string `json:"requisition_id"`
				Title         string `json:"title"`
			} `json:"jobs"`
			Meta struct {
				Total int `json:"total"`
			} `json:"meta"`
		}

		var jsonJobs JsonJobs

		c.OnResponse(func(r *colly.Response) {
			var tempJson JsonJobs
			err := json.Unmarshal(r.Body, &tempJson)
			if err != nil {
				panic(err.Error())
			}

			for _, elem := range tempJson.Jobs {

				result_title := elem.Title
				result_url := elem.AbsoluteURL

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

			jsonJobs.Jobs = append(jsonJobs.Jobs, tempJson.Jobs...)
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
			c.Visit(start_url)
		}
	}
	return
}

func (runtime Runtime) Wefox(
	version int, isLocal bool) (
	response Response, results []Result) {
	switch version {
	case 1:

		c := colly.NewCollector()

		start_url := "https://api.greenhouse.io/v1/boards/wefoxgroup/jobs"

		type JsonJobs struct {
			Jobs []struct {
				AbsoluteURL    string `json:"absolute_url"`
				DataCompliance []struct {
					Type            string      `json:"type"`
					RequiresConsent bool        `json:"requires_consent"`
					RetentionPeriod interface{} `json:"retention_period"`
				} `json:"data_compliance"`
				Education     string `json:"education,omitempty"`
				InternalJobID int    `json:"internal_job_id"`
				Location      struct {
					Name string `json:"name"`
				} `json:"location"`
				Metadata []struct {
					ID        int         `json:"id"`
					Name      string      `json:"name"`
					Value     interface{} `json:"value"`
					ValueType string      `json:"value_type"`
				} `json:"metadata"`
				ID            int    `json:"id"`
				UpdatedAt     string `json:"updated_at"`
				RequisitionID string `json:"requisition_id"`
				Title         string `json:"title"`
			} `json:"jobs"`
			Meta struct {
				Total int `json:"total"`
			} `json:"meta"`
		}

		var jsonJobs JsonJobs

		c.OnResponse(func(r *colly.Response) {
			var tempJson JsonJobs
			err := json.Unmarshal(r.Body, &tempJson)
			if err != nil {
				panic(err.Error())
			}

			for _, elem := range tempJson.Jobs {

				result_title := elem.Title
				result_url := elem.AbsoluteURL

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

			jsonJobs.Jobs = append(jsonJobs.Jobs, tempJson.Jobs...)
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
			c.Visit(start_url)
		}
	}
	return
}

func (runtime Runtime) Celonis(
	version int, isLocal bool) (
	response Response, results []Result) {
	switch version {
	case 1:

		c := colly.NewCollector()

		start_url := "https://api.greenhouse.io/v1/boards/celonis/jobs"

		type JsonJobs struct {
			Jobs []struct {
				AbsoluteURL    string `json:"absolute_url"`
				DataCompliance []struct {
					Type            string      `json:"type"`
					RequiresConsent bool        `json:"requires_consent"`
					RetentionPeriod interface{} `json:"retention_period"`
				} `json:"data_compliance"`
				Education     string `json:"education,omitempty"`
				InternalJobID int    `json:"internal_job_id"`
				Location      struct {
					Name string `json:"name"`
				} `json:"location"`
				Metadata []struct {
					ID        int         `json:"id"`
					Name      string      `json:"name"`
					Value     interface{} `json:"value"`
					ValueType string      `json:"value_type"`
				} `json:"metadata"`
				ID            int    `json:"id"`
				UpdatedAt     string `json:"updated_at"`
				RequisitionID string `json:"requisition_id"`
				Title         string `json:"title"`
			} `json:"jobs"`
			Meta struct {
				Total int `json:"total"`
			} `json:"meta"`
		}

		var jsonJobs JsonJobs

		c.OnResponse(func(r *colly.Response) {
			var tempJson JsonJobs
			err := json.Unmarshal(r.Body, &tempJson)
			if err != nil {
				panic(err.Error())
			}

			for _, elem := range tempJson.Jobs {

				result_title := elem.Title
				result_url := elem.AbsoluteURL

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

			jsonJobs.Jobs = append(jsonJobs.Jobs, tempJson.Jobs...)
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
			c.Visit(start_url)
		}
	}
	return
}

func (runtime Runtime) Omio(
	version int, isLocal bool) (
	response Response, results []Result) {
	switch version {
	case 1:

		c := colly.NewCollector()

		start_url := "https://api.smartrecruiters.com/v1/companies/Omio1/postings"
		base_job_url := "https://www.omio.com/jobs/#"

		type JsonJobs struct {
			Offset     int `json:"offset"`
			Limit      int `json:"limit"`
			TotalFound int `json:"totalFound"`
			Content    []struct {
				ID        string `json:"id"`
				Name      string `json:"name"`
				UUID      string `json:"uuid"`
				RefNumber string `json:"refNumber"`
				Company   struct {
					Identifier string `json:"identifier"`
					Name       string `json:"name"`
				} `json:"company"`
				ReleasedDate time.Time `json:"releasedDate"`
				Location     struct {
					City       string `json:"city"`
					Region     string `json:"region"`
					Country    string `json:"country"`
					Address    string `json:"address"`
					PostalCode string `json:"postalCode"`
					Remote     bool   `json:"remote"`
				} `json:"location,omitempty"`
				Industry struct {
					ID    string `json:"id"`
					Label string `json:"label"`
				} `json:"industry"`
				Department struct {
					ID    string `json:"id"`
					Label string `json:"label"`
				} `json:"department,omitempty"`
				Function struct {
					ID    string `json:"id"`
					Label string `json:"label"`
				} `json:"function"`
				TypeOfEmployment struct {
					Label string `json:"label"`
				} `json:"typeOfEmployment"`
				ExperienceLevel struct {
					ID    string `json:"id"`
					Label string `json:"label"`
				} `json:"experienceLevel"`
				CustomField []struct {
					FieldID    string `json:"fieldId"`
					FieldLabel string `json:"fieldLabel"`
					ValueID    string `json:"valueId"`
					ValueLabel string `json:"valueLabel"`
				} `json:"customField"`
				Ref     string `json:"ref"`
				Creator struct {
					Name string `json:"name"`
				} `json:"creator"`
				Language struct {
					Code        string `json:"code"`
					Label       string `json:"label"`
					LabelNative string `json:"labelNative"`
				} `json:"language"`
			} `json:"content"`
		}

		var jsonJobs JsonJobs

		c.OnResponse(func(r *colly.Response) {
			var tempJson JsonJobs
			err := json.Unmarshal(r.Body, &tempJson)
			if err != nil {
				panic(err.Error())
			}

			for _, elem := range tempJson.Content {

				result_title := elem.Name
				result_url := base_job_url + elem.ID

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

			jsonJobs.Content = append(jsonJobs.Content, tempJson.Content...)
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
			c.Visit(start_url)
		}
	}
	return
}

func (runtime Runtime) Aboutyou(
	version int, isLocal bool) (
	response Response, results []Result) {
	switch version {
	case 1:

		c := colly.NewCollector()

		start_url := "https://corporate.aboutyou.de/app/api/openpositions.php?posts_per_page=500"

		type JsonJobs struct {
			Posts []struct {
				ID         int    `json:"id"`
				Title      string `json:"title"`
				Department string `json:"department"`
				Location   string `json:"location"`
				URL        string `json:"url"`
				Type       struct {
					ID   int    `json:"id"`
					Name string `json:"name"`
				} `json:"type"`
			} `json:"posts"`
			TotalCount int `json:"totalCount"`
		}

		var jsonJobs JsonJobs

		c.OnResponse(func(r *colly.Response) {
			var tempJson JsonJobs
			err := json.Unmarshal(r.Body, &tempJson)
			if err != nil {
				panic(err.Error())
			}

			for _, elem := range tempJson.Posts {

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

			jsonJobs.Posts = append(jsonJobs.Posts, tempJson.Posts...)
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
			c.Visit(start_url)
		}
	}
	return
}

func (runtime Runtime) Depositsolutions(
	version int, isLocal bool) (
	response Response, results []Result) {
	switch version {
	case 1:

		c := colly.NewCollector()

		start_url := "https://careers-page.workable.com/api/v3/accounts/deposit-solutions/jobs"
		base_job_url := "https://apply.workable.com/deposit-solutions/j/"

		type JsonJobs struct {
			Total   int `json:"total"`
			Results []struct {
				ID           int    `json:"id"`
				Shortcode    string `json:"shortcode"`
				Title        string `json:"title"`
				Description  string `json:"description"`
				Requirements string `json:"requirements"`
				Benefits     string `json:"benefits"`
				Remote       bool   `json:"remote"`
				Location     struct {
					Country     string `json:"country"`
					CountryCode string `json:"countryCode"`
					City        string `json:"city"`
					Region      string `json:"region"`
				} `json:"location"`
				State          string      `json:"state"`
				IsInternal     bool        `json:"isInternal"`
				Code           interface{} `json:"code"`
				Published      time.Time   `json:"published"`
				Type           string      `json:"type"`
				Language       string      `json:"language"`
				Department     []string    `json:"department"`
				AccountUID     string      `json:"accountUid"`
				ApprovalStatus string      `json:"approvalStatus"`
			} `json:"results"`
		}

		var jsonJobs JsonJobs

		c.OnResponse(func(r *colly.Response) {
			var tempJson JsonJobs
			err := json.Unmarshal(r.Body, &tempJson)
			if err != nil {
				panic(err.Error())
			}

			for _, elem := range tempJson.Results {

				result_title := elem.Title
				result_url := base_job_url + elem.Shortcode

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

			jsonJobs.Results = append(jsonJobs.Results, tempJson.Results...)
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
			c.Request(
				"POST",
				start_url,
				strings.NewReader(""),
				nil,
				http.Header{"Content-Type": []string{"application/x-www-form-urlencoded"}},
			)
		}
	}
	return
}

func (runtime Runtime) Taxfix(
	version int, isLocal bool) (
	response Response, results []Result) {
	switch version {
	case 1:

		c := colly.NewCollector()

		start_url := "https://api.greenhouse.io/v1/boards/taxfix/jobs"

		type JsonJobs struct {
			Jobs []struct {
				AbsoluteURL    string `json:"absolute_url"`
				DataCompliance []struct {
					Type            string      `json:"type"`
					RequiresConsent bool        `json:"requires_consent"`
					RetentionPeriod interface{} `json:"retention_period"`
				} `json:"data_compliance"`
				Education     string `json:"education,omitempty"`
				InternalJobID int    `json:"internal_job_id"`
				Location      struct {
					Name string `json:"name"`
				} `json:"location"`
				Metadata []struct {
					ID        int         `json:"id"`
					Name      string      `json:"name"`
					Value     interface{} `json:"value"`
					ValueType string      `json:"value_type"`
				} `json:"metadata"`
				ID            int    `json:"id"`
				UpdatedAt     string `json:"updated_at"`
				RequisitionID string `json:"requisition_id"`
				Title         string `json:"title"`
			} `json:"jobs"`
			Meta struct {
				Total int `json:"total"`
			} `json:"meta"`
		}

		var jsonJobs JsonJobs

		c.OnResponse(func(r *colly.Response) {
			var tempJson JsonJobs
			err := json.Unmarshal(r.Body, &tempJson)
			if err != nil {
				panic(err.Error())
			}

			for _, elem := range tempJson.Jobs {

				result_title := elem.Title
				result_url := elem.AbsoluteURL

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

			jsonJobs.Jobs = append(jsonJobs.Jobs, tempJson.Jobs...)
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
			c.Visit(start_url)
		}
	}
	return
}

func (runtime Runtime) Moonfare(
	version int, isLocal bool) (
	response Response, results []Result) {
	switch version {
	case 1:

		c := colly.NewCollector()

		start_url := "https://api.greenhouse.io/v1/boards/moonfare/jobs"

		type JsonJobs struct {
			Jobs []struct {
				AbsoluteURL    string `json:"absolute_url"`
				DataCompliance []struct {
					Type            string      `json:"type"`
					RequiresConsent bool        `json:"requires_consent"`
					RetentionPeriod interface{} `json:"retention_period"`
				} `json:"data_compliance"`
				Education     string `json:"education,omitempty"`
				InternalJobID int    `json:"internal_job_id"`
				Location      struct {
					Name string `json:"name"`
				} `json:"location"`
				Metadata []struct {
					ID        int         `json:"id"`
					Name      string      `json:"name"`
					Value     interface{} `json:"value"`
					ValueType string      `json:"value_type"`
				} `json:"metadata"`
				ID            int    `json:"id"`
				UpdatedAt     string `json:"updated_at"`
				RequisitionID string `json:"requisition_id"`
				Title         string `json:"title"`
			} `json:"jobs"`
			Meta struct {
				Total int `json:"total"`
			} `json:"meta"`
		}

		var jsonJobs JsonJobs

		c.OnResponse(func(r *colly.Response) {
			var tempJson JsonJobs
			err := json.Unmarshal(r.Body, &tempJson)
			if err != nil {
				panic(err.Error())
			}

			for _, elem := range tempJson.Jobs {

				result_title := elem.Title
				result_url := elem.AbsoluteURL

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

			jsonJobs.Jobs = append(jsonJobs.Jobs, tempJson.Jobs...)
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
			c.Visit(start_url)
		}
	}
	return
}

func (runtime Runtime) Fincompare(
	version int, isLocal bool) (
	response Response, results []Result) {
	switch version {
	case 1:

		c := colly.NewCollector()

		c_start_url := "https://api.lever.co/v0/postings/fincompare?mode=json"

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

func (runtime Runtime) Billie(
	version int, isLocal bool) (
	response Response, results []Result) {
	switch version {
	case 1:

		c := colly.NewCollector()

		start_url := "https://api.greenhouse.io/v1/boards/billie/jobs"

		type JsonJobs struct {
			Jobs []struct {
				AbsoluteURL    string `json:"absolute_url"`
				DataCompliance []struct {
					Type            string      `json:"type"`
					RequiresConsent bool        `json:"requires_consent"`
					RetentionPeriod interface{} `json:"retention_period"`
				} `json:"data_compliance"`
				Education     string `json:"education,omitempty"`
				InternalJobID int    `json:"internal_job_id"`
				Location      struct {
					Name string `json:"name"`
				} `json:"location"`
				Metadata []struct {
					ID        int         `json:"id"`
					Name      string      `json:"name"`
					Value     interface{} `json:"value"`
					ValueType string      `json:"value_type"`
				} `json:"metadata"`
				ID            int    `json:"id"`
				UpdatedAt     string `json:"updated_at"`
				RequisitionID string `json:"requisition_id"`
				Title         string `json:"title"`
			} `json:"jobs"`
			Meta struct {
				Total int `json:"total"`
			} `json:"meta"`
		}

		var jsonJobs JsonJobs

		c.OnResponse(func(r *colly.Response) {
			var tempJson JsonJobs
			err := json.Unmarshal(r.Body, &tempJson)
			if err != nil {
				panic(err.Error())
			}

			for _, elem := range tempJson.Jobs {

				result_title := elem.Title
				result_url := elem.AbsoluteURL

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

			jsonJobs.Jobs = append(jsonJobs.Jobs, tempJson.Jobs...)
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
			c.Visit(start_url)
		}
	}
	return
}

func (runtime Runtime) Pairfinance(
	version int, isLocal bool) (
	response Response, results []Result) {
	switch version {
	case 1:

		c := colly.NewCollector()

		start_url := "https://api.greenhouse.io/v1/boards/pairfinance/jobs"

		type JsonJobs struct {
			Jobs []struct {
				AbsoluteURL    string `json:"absolute_url"`
				DataCompliance []struct {
					Type            string      `json:"type"`
					RequiresConsent bool        `json:"requires_consent"`
					RetentionPeriod interface{} `json:"retention_period"`
				} `json:"data_compliance"`
				Education     string `json:"education,omitempty"`
				InternalJobID int    `json:"internal_job_id"`
				Location      struct {
					Name string `json:"name"`
				} `json:"location"`
				Metadata []struct {
					ID        int         `json:"id"`
					Name      string      `json:"name"`
					Value     interface{} `json:"value"`
					ValueType string      `json:"value_type"`
				} `json:"metadata"`
				ID            int    `json:"id"`
				UpdatedAt     string `json:"updated_at"`
				RequisitionID string `json:"requisition_id"`
				Title         string `json:"title"`
			} `json:"jobs"`
			Meta struct {
				Total int `json:"total"`
			} `json:"meta"`
		}

		var jsonJobs JsonJobs

		c.OnResponse(func(r *colly.Response) {
			var tempJson JsonJobs
			err := json.Unmarshal(r.Body, &tempJson)
			if err != nil {
				panic(err.Error())
			}

			for _, elem := range tempJson.Jobs {

				result_title := elem.Title
				result_url := elem.AbsoluteURL

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

			jsonJobs.Jobs = append(jsonJobs.Jobs, tempJson.Jobs...)
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
			c.Visit(start_url)
		}
	}
	return
}

func (runtime Runtime) Getsafe(
	version int, isLocal bool) (
	response Response, results []Result) {
	switch version {
	case 1:

		c := colly.NewCollector()

		url := "https://getsafe-jobs.personio.de"
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

func (runtime Runtime) Liqid(
	version int, isLocal bool) (
	response Response, results []Result) {
	switch version {
	case 1:

		c := colly.NewCollector()

		url := "https://liqid-jobs.personio.de"
		main_tag := "div"
		main_tag_attr := "class"
		main_tag_value := "job-list-desc"
		tag_title := "a"
		tag_info := "p"
		separator := "·"

		type Job struct {
			Title    string
			Url      string
			Type     string
			Location string
		}

		c.OnHTML(main_tag, func(e *colly.HTMLElement) {
			if strings.Contains(e.Attr(main_tag_attr), main_tag_value) {
				result_title := e.ChildText(tag_title)
				result_url := e.ChildAttr(tag_title, "href")
				result_info := strings.Split(e.ChildText(tag_info), separator)
				result_type := strings.Join(strings.Fields(strings.TrimSpace(result_info[0])), " ")
				result_location := strings.Join(strings.Fields(strings.TrimSpace(result_info[1])), " ")

				_, err := netUrl.ParseRequestURI(result_url)
				if err == nil {

					temp_elem_json := Job{
						result_title,
						result_url,
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

func (runtime Runtime) Elementinsurance(
	version int, isLocal bool) (
	response Response, results []Result) {
	switch version {
	case 1:

		c := colly.NewCollector()

		p_start_url := "https://elementinsuranceag.recruitee.com/api/offers"

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
				MinHours           int           `json:"min_hours"`
				MaxHours           int           `json:"max_hours"`
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
	}
	return
}

func (runtime Runtime) Freeda(
	version int, isLocal bool) (
	response Response, results []Result) {
	switch version {
	case 1:

		c := colly.NewCollector()

		start_url := "https://api.greenhouse.io/v1/boards/freedamedia/jobs"

		type JsonJobs struct {
			Jobs []struct {
				AbsoluteURL    string `json:"absolute_url"`
				DataCompliance []struct {
					Type            string      `json:"type"`
					RequiresConsent bool        `json:"requires_consent"`
					RetentionPeriod interface{} `json:"retention_period"`
				} `json:"data_compliance"`
				Education     string `json:"education,omitempty"`
				InternalJobID int    `json:"internal_job_id"`
				Location      struct {
					Name string `json:"name"`
				} `json:"location"`
				Metadata []struct {
					ID        int         `json:"id"`
					Name      string      `json:"name"`
					Value     interface{} `json:"value"`
					ValueType string      `json:"value_type"`
				} `json:"metadata"`
				ID            int    `json:"id"`
				UpdatedAt     string `json:"updated_at"`
				RequisitionID string `json:"requisition_id"`
				Title         string `json:"title"`
			} `json:"jobs"`
			Meta struct {
				Total int `json:"total"`
			} `json:"meta"`
		}

		var jsonJobs JsonJobs

		c.OnResponse(func(r *colly.Response) {
			var tempJson JsonJobs
			err := json.Unmarshal(r.Body, &tempJson)
			if err != nil {
				panic(err.Error())
			}

			for _, elem := range tempJson.Jobs {

				result_title := elem.Title
				result_url := elem.AbsoluteURL

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

			jsonJobs.Jobs = append(jsonJobs.Jobs, tempJson.Jobs...)
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
			c.Visit(start_url)
		}
	}
	return
}

func (runtime Runtime) Talentgarden(
	version int, isLocal bool) (
	response Response, results []Result) {
	switch version {
	case 1:

		c := colly.NewCollector()

		url := "https://talentgarden.bamboohr.com/jobs/embed2.php?departmentId=0"
		main_tag := "div"
		main_tag_attr := "class"
		main_tag_value := "BambooHR-ATS-board"
		tag_department_section := "li[class=BambooHR-ATS-Department-Item]"
		tag_department := "div[class=BambooHR-ATS-Department-Header]"
		tag_job_section := "ul[class=BambooHR-ATS-Jobs-List]"
		tag_title := "a"
		tag_location := "span"

		type Job struct {
			Department string
			Title      string
			Url        string
			Location   string
		}

		c.OnHTML(main_tag, func(e *colly.HTMLElement) {
			if strings.Contains(e.Attr(main_tag_attr), main_tag_value) {
				e.ForEach(tag_department_section, func(_ int, el *colly.HTMLElement) {
					result_department := strings.TrimSpace(el.ChildText(tag_department))
					el.ForEach(tag_job_section, func(_ int, ell *colly.HTMLElement) {
						result_title := ell.ChildText(tag_title)
						result_url := "https:" + ell.ChildAttr(tag_title, "href")
						result_location := ell.ChildText(tag_location)

						_, err := netUrl.ParseRequestURI(result_url)
						if err == nil {

							temp_elem_json := Job{
								result_department,
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

func (runtime Runtime) Facileit(
	version int, isLocal bool) (
	response Response, results []Result) {
	switch version {
	case 1:

		c := colly.NewCollector()
		l := c.Clone()

		url := "https://jobs.facile.it/chi-cerchiamo.html"
		tag_main := "div[id=JB_central]"

		type Job struct {
			Title       string
			Url         string
			Location    string
			Description string
		}

		if !isLocal {

			c.OnHTML(tag_main, func(e *colly.HTMLElement) {
				script_url := e.ChildAttr("script", "src")
				k := strings.Split(strings.Split(script_url, "&k=")[1], "&LAC")[0]
				base_url := "https://inrecruiting.intervieweb.it/app.php?module=iframeAnnunci&k=" + k + "&LAC=Facileit&act1=23"
				l.Visit(base_url)
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

			var jsonJobs []Job

			l.OnResponse(func(r *colly.Response) {
				responseText := string(r.Body)
				url := strings.Split(strings.Split(responseText, "$.post('")[1], "',")[0]
				cookie := c.Cookies(r.Request.URL.String())[0].Raw

				client := &http.Client{}
				data := strings.NewReader(`orderBy=byfunction&descEn=1`)
				req, _ := http.NewRequest("POST", url, data)
				req.Header.Set("content-type", "application/x-www-form-urlencoded")
				req.Header.Set("cookie", cookie)
				resp, _ := client.Do(req)
				bodyText, _ := ioutil.ReadAll(resp.Body)
				doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(bodyText)))
				if err != nil {
					panic(err.Error())
				}

				var titles []string
				var urls []string
				var locations []string
				doc.Find("dt").Each(func(i int, s *goquery.Selection) {
					response_title := strings.TrimSpace(s.Find("a").Text())
					temp_response_url, _ := s.Find("a").Attr("href")
					response_url := strings.ReplaceAll(temp_response_url, "defgroup=", "defgroup=function") + "400&d=jobs.facile.it"
					response_location := strings.TrimSpace(s.Find("span[class=location_annuncio]").Text())
					titles = append(titles, response_title)
					urls = append(urls, response_url)
					locations = append(locations, response_location)
				})

				var descriptions []string
				doc.Find("dd").Each(func(i int, s *goquery.Selection) {
					response_description := strings.TrimSpace(s.Find("p[class=description]").Text())
					descriptions = append(descriptions, response_description)
				})

				for i := range titles {
					temp_elem_json := Job{
						titles[i],
						urls[i],
						locations[i],
						descriptions[i],
					}

					elem_json, err := json.Marshal(temp_elem_json)
					if err != nil {
						panic(err.Error())
					}

					results = append(results, Result{
						runtime.Name,
						titles[i],
						urls[i],
						elem_json,
					})

					jsonJobs = append(jsonJobs, temp_elem_json)
				}
			})

			l.OnScraped(func(r *colly.Response) {
				jsonJobs_marshal, err := json.Marshal(jsonJobs)
				if err != nil {
					panic(err.Error())
				}
				response = Response{[]byte(jsonJobs_marshal)}
			})

			c.Visit(url)
		} else {

			var jsonJobs []Job

			c.OnResponse(func(r *colly.Response) {
				err := json.Unmarshal(r.Body, &jsonJobs)
				if err != nil {
					panic(err.Error())
				}

				for _, elem := range jsonJobs {

					elem_json, err := json.Marshal(elem)
					if err != nil {
						panic(err.Error())
					}

					results = append(results, Result{
						runtime.Name,
						elem.Url,
						elem.Title,
						elem_json,
					})
				}
			})

			c.OnScraped(func(r *colly.Response) {
				jsonJobs_marshal, err := json.Marshal(jsonJobs)
				if err != nil {
					panic(err.Error())
				}
				response = Response{[]byte(jsonJobs_marshal)}
			})

			t := &http.Transport{}
			t.RegisterProtocol("file", http.NewFileTransport(http.Dir("/")))
			c.WithTransport(t)
			dir, err := os.Getwd()
			if err != nil {
				panic(err.Error())
			}
			c.Visit("file:" + dir + "/response.html")
		}
	}
	return
}

func (runtime Runtime) Vodafone(
	version int, isLocal bool) (response Response, results []Result) {
	switch version {
	case 1:

		c := colly.NewCollector()

		// today_date := "&date=" + strings.ReplaceAll(time.Now().Format("02/01/06"), "/", "%2F")

		v_start_url := "https://careers.vodafone.com/search/?startrow=%d"
		v_base_url := "https://careers.vodafone.com"
		number_results_per_page := 25
		counter := 0
		tag_result := "a"
		tag_location := "span[class=jobLocation]"
		tag_date := "span[class=jobDate]"
		tag_total_results := ".paginationLabel"

		type Job struct {
			Title    string
			Url      string
			Location string
			Date     string
		}

		var jsonJobs []Job

		c.OnHTML(".html5", func(e *colly.HTMLElement) {
			e.ForEach(".data-row", func(_ int, el *colly.HTMLElement) {
				result_title := strings.Join(strings.Fields(strings.TrimSpace(el.ChildTexts(tag_result)[0])), " ")
				result_url := v_base_url + strings.Join(strings.Fields(strings.TrimSpace(el.ChildAttr(tag_result, "href"))), " ")
				result_location := strings.Join(strings.Fields(strings.TrimSpace(el.ChildText(tag_location))), " ")
				result_date := strings.Join(strings.Fields(strings.TrimSpace(el.ChildText(tag_date))), " ")

				_, err := netUrl.ParseRequestURI(result_url)
				if err == nil {

					temp_elem_json := Job{
						result_title,
						result_url,
						result_location,
						result_date,
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

			temp_total_results := strings.Split(e.ChildText(tag_total_results), " ")
			string_total_results := temp_total_results[len(temp_total_results)-1]
			total_results, err := strconv.Atoi(string_total_results)
			if err != nil {
				panic(err.Error())
			}

			total_pages := total_results/number_results_per_page + 2

			if isLocal {
				return
			} else {
				if counter >= total_pages {
					return
				} else {
					counter++
					time.Sleep(SecondsSleep * time.Second)
					temp_v_url := fmt.Sprintf(v_start_url, counter*number_results_per_page)
					c.Visit(temp_v_url)
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
			c.Visit(fmt.Sprintf(v_start_url, 0))
		}
	}
	return
}

func (runtime Runtime) Glovo(
	version int, isLocal bool) (
	response Response, results []Result) {
	switch version {
	case 1:

		c := colly.NewCollector()

		start_url := "https://api.greenhouse.io/v1/boards/glovo/jobs"

		type JsonJobs struct {
			Jobs []struct {
				AbsoluteURL    string `json:"absolute_url"`
				DataCompliance []struct {
					Type            string      `json:"type"`
					RequiresConsent bool        `json:"requires_consent"`
					RetentionPeriod interface{} `json:"retention_period"`
				} `json:"data_compliance"`
				Education     string `json:"education,omitempty"`
				InternalJobID int    `json:"internal_job_id"`
				Location      struct {
					Name string `json:"name"`
				} `json:"location"`
				Metadata []struct {
					ID        int         `json:"id"`
					Name      string      `json:"name"`
					Value     interface{} `json:"value"`
					ValueType string      `json:"value_type"`
				} `json:"metadata"`
				ID            int    `json:"id"`
				UpdatedAt     string `json:"updated_at"`
				RequisitionID string `json:"requisition_id"`
				Title         string `json:"title"`
			} `json:"jobs"`
			Meta struct {
				Total int `json:"total"`
			} `json:"meta"`
		}

		var jsonJobs JsonJobs

		c.OnResponse(func(r *colly.Response) {
			var tempJson JsonJobs
			err := json.Unmarshal(r.Body, &tempJson)
			if err != nil {
				panic(err.Error())
			}

			for _, elem := range tempJson.Jobs {

				result_title := elem.Title
				result_url := elem.AbsoluteURL

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

			jsonJobs.Jobs = append(jsonJobs.Jobs, tempJson.Jobs...)
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
			c.Visit(start_url)
		}
	}
	return
}

func (runtime Runtime) Glickon(
	version int, isLocal bool) (
	response Response, results []Result) {
	switch version {
	case 1:

		c := colly.NewCollector()
		l := c.Clone()

		section_url := "https://core.glickon.com/api/candidate/latest/companies/glickon"
		department_url := "https://core.glickon.com/api/candidate/latest/sections/%s?from_www=true"
		job_base_url := "https://www.glickon.com/en/challenges/"

		type Departments struct {
			Sections []struct {
				ID   int    `json:"id"`
				Name string `json:"name"`
			} `json:"sections"`
		}

		type Jobs struct {
			ID                  int           `json:"id"`
			Name                string        `json:"name"`
			ShortDescription    string        `json:"short_description"`
			Description         string        `json:"description"`
			Color               string        `json:"color"`
			IconURL             string        `json:"icon_url"`
			BackgroundURL       string        `json:"background_url"`
			IsPublic            bool          `json:"is_public"`
			ForEmployees        bool          `json:"for_employees"`
			CompanyCareersName  string        `json:"company_careers_name"`
			ShowLeaderboard     bool          `json:"show_leaderboard"`
			ShowTeamLeaderboard bool          `json:"show_team_leaderboard"`
			Images              []interface{} `json:"images"`
			Videos              []interface{} `json:"videos"`
			Files               []interface{} `json:"files"`
			Challenges          []struct {
				Hash                            string `json:"hash"`
				Name                            string `json:"name"`
				Description                     string `json:"description"`
				ShortDescription                string `json:"short_description"`
				HasPassword                     bool   `json:"has_password"`
				SponsoredImageURL               string `json:"sponsored_image_url"`
				Color                           string `json:"color"`
				NameForExternalPage             string `json:"name_for_external_page"`
				ShortDescriptionForExternalPage string `json:"short_description_for_external_page"`
				PlayButtonForExternalPage       string `json:"play_button_for_external_page"`
				NumberOfQuestions               int    `json:"number_of_questions"`
				EstimatedCompletionTime         int    `json:"estimated_completion_time"`
			} `json:"challenges"`
		}

		var jsonJobs Jobs

		c.OnResponse(func(r *colly.Response) {
			var departments Departments
			err := json.Unmarshal(r.Body, &departments)
			if err != nil {
				panic(err.Error())
			}
			for _, elem := range departments.Sections {
				department_id := elem.ID
				department_url := fmt.Sprintf(department_url, strconv.Itoa(department_id))
				l.Visit(department_url)
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

		l.OnResponse(func(r *colly.Response) {
			var tempJsonJobs Jobs
			err := json.Unmarshal(r.Body, &tempJsonJobs)
			if err != nil {
				panic(err.Error())
			}
			jsonJobs.Challenges = append(jsonJobs.Challenges, tempJsonJobs.Challenges...)

			for _, elem := range tempJsonJobs.Challenges {
				result_title := elem.Name
				result_url := job_base_url + elem.Hash
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
			c.Visit(section_url)
		}
	}
	return
}

func (runtime Runtime) Satispay(
	version int, isLocal bool) (
	response Response, results []Result) {
	switch version {
	case 1:

		c := colly.NewCollector()

		url := "https://satispay.breezy.hr"
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

func (runtime Runtime) Medtronic(
	version int, isLocal bool) (
	response Response, results []Result) {
	switch version {
	case 1:

		c := colly.NewCollector()
		l := c.Clone()
		x := l.Clone()

		start_url := "https://jobs.medtronic.com"
		id_url := "https://jobs.medtronic.com/ajax/jobs/search/create?uid=661"
		temp_jobs_url := "https://jobs.medtronic.com/jobs/search/%d/page/%d"
		var session_id int
		var tsstoken string
		var ORA_OTSS_SESSION_ID string
		var cookies []string

		type Session struct {
			Status      string `json:"Status"`
			UserMessage string `json:"UserMessage"`
			Result      struct {
				JobSearchID int `json:"JobSearch.id"`
			} `json:"Result"`
		}

		type Job struct {
			Title       string
			Url         string
			Location    string
			Category    string
			Description string
		}

		c.OnHTML("body", func(e *colly.HTMLElement) {
			tsstoken = e.ChildAttr("input[name=tsstoken]", "value")
			l.Visit(id_url)
		})

		c.OnResponse(func(r *colly.Response) {
			responseData := string(r.Body)
			ORA_OTSS_SESSION_ID = strings.Split(strings.Split(responseData, `session_id":"`)[2], `","`)[0]
			cookies = append(cookies, "ORA_OTSS_SESSION_ID="+ORA_OTSS_SESSION_ID)
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

		l.OnResponse(func(r *colly.Response) {
			var session Session
			err := json.Unmarshal(r.Body, &session)
			if err != nil {
				panic(err.Error())
			}
			session_id = session.Result.JobSearchID
			jobs_url := fmt.Sprintf(temp_jobs_url, session_id, 1)
			x.Visit(jobs_url)
		})

		x.OnHTML("body", func(e *colly.HTMLElement) {
			e.ForEach(".job_list_row", func(_ int, el *colly.HTMLElement) {
				result_title := el.ChildText("a")
				result_url := el.ChildAttr("a", "href")
				result_location := el.ChildText("span[class=location]")
				result_category := el.ChildText("span[class=category]")
				result_description := el.ChildText("p[class=jlr_description]")
				_, err := netUrl.ParseRequestURI(result_url)
				if err == nil {

					temp_elem_json := Job{
						result_title,
						result_url,
						result_location,
						result_category,
						result_description,
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

			string_number_pages := e.ChildText("div[id=jPaginateNumPages]")
			number_pages, _ := strconv.Atoi(strings.Split(string_number_pages, ".")[0])

			for counter := 2; counter <= number_pages; counter++ {
				time.Sleep(SecondsSleep * time.Second)

				temp_url := "https://jobs.medtronic.com/ajax/content/job_results?JobSearch.id=%d&page_index=%d"
				temp_temp_url := fmt.Sprintf(temp_url, session_id, counter)
				fmt.Println(Gray(8-1, "Visiting"), Gray(8-1, temp_temp_url))

				client := &http.Client{}
				req, err := http.NewRequest("POST", temp_temp_url, nil)
				if err != nil {
					panic(err.Error())
				}

				req.Header.Set("tss-token", tsstoken)
				req.Header.Set("Cookie", "ORA_OTSS_SESSION_ID="+ORA_OTSS_SESSION_ID)
				resp, err := client.Do(req)
				if err != nil {
					panic(err.Error())
				}
				bodyText, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					panic(err.Error())
				}

				body := strings.ReplaceAll(
					strings.ReplaceAll(
						strings.ReplaceAll(string(bodyText), "\t", ""), "\n", ""), `\`, "")

				results_sections := strings.Split(body, "job_list_row")
				for i := 1; i < len(results_sections); i++ {
					elem := results_sections[i]
					result_title := strings.Split(strings.Split(elem, `class="job_link font_bold">`)[1], `</a>`)[0]
					result_url := strings.Split(strings.Split(elem, `<p><a href="`)[1], `"`)[0]
					result_location := strings.Split(strings.Split(elem, `<span class="location">`)[1], `</span>`)[0]
					result_category := strings.Split(strings.Split(elem, `<span class="category">`)[1], `</span>`)[0]
					result_description := strings.Split(strings.Split(elem, `<p class="jlr_description">`)[1], `</p>`)[0]

					temp_elem_json := Job{
						result_title,
						result_url,
						result_location,
						result_category,
						result_description,
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

		x.OnRequest(func(r *colly.Request) {
			fmt.Println(Gray(8-1, "Visiting"), Gray(8-1, r.URL.String()))
		})

		x.OnScraped(func(r *colly.Response) {
			results_marshal, err := json.Marshal(results)
			if err != nil {
				panic(err.Error())
			}
			response = Response{[]byte(results_marshal)}
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
			c.Visit(start_url)
		}
	}
	return
}

func (runtime Runtime) Bendingspoons(
	version int, isLocal bool) (
	response Response, results []Result) {
	switch version {
	case 1:

		c := colly.NewCollector()

		start_url := "https://website.rolemodel.bendingspoons.com/roles.json"
		job_url := "https://bendingspoons.com/careers.html?x="

		type Jobs []struct {
			Salary            string `json:"salary,omitempty"`
			Title             string `json:"title"`
			Photo             string `json:"photo"`
			ID                string `json:"id"`
			Area              string `json:"area"`
			ApplicationFields []struct {
				MaxFileSize   int           `json:"max_file_size"`
				Subtitle      string        `json:"subtitle"`
				Title         string        `json:"title"`
				Optional      bool          `json:"optional"`
				Choices       []interface{} `json:"choices"`
				Extensions    []interface{} `json:"extensions"`
				FileTypes     []interface{} `json:"file_types"`
				Type          string        `json:"type"`
				ID            string        `json:"id"`
				MaxCharacters int           `json:"max_characters"`
			} `json:"application_fields"`
			Contract string `json:"contract,omitempty"`
			WeOffer  []struct {
				Text  string `json:"text"`
				Title string `json:"title"`
			} `json:"we_offer"`
			ShortDescription string `json:"short_description"`
			Version          int    `json:"version"`
			Location         string `json:"location,omitempty"`
			JobVisible       bool   `json:"job_visible"`
			JobActive        bool   `json:"job_active"`
			WeLookFor        []struct {
				Text  string `json:"text"`
				Title string `json:"title"`
			} `json:"we_look_for"`
			LongDescription string `json:"long_description"`
		}

		var jsonJobs Jobs

		c.OnResponse(func(r *colly.Response) {
			var tempJson Jobs
			err := json.Unmarshal(r.Body, &tempJson)
			if err != nil {
				panic(err.Error())
			}

			for _, elem := range tempJson {

				result_title := elem.Title
				result_url := job_url + elem.ID

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
			c.Visit(start_url)
		}
	}
	return
}

func (runtime Runtime) Bcg(
	version int, isLocal bool) (response Response, results []Result) {
	switch version {
	case 1:

		type Job struct {
			Title       string
			Url         string
			Location    string
			Date        string
			Description string
		}

		if !isLocal {

			ctx, cancel := chromedp.NewContext(context.Background())
			defer cancel()

			b_start_url := "https://talent.bcg.com/en_US/apply/SearchJobs/?folderOffset=%d"
			start_offset := 0
			number_results_per_page := 20
			_ = number_results_per_page

			var initialResponse string
			if err := chromedp.Run(ctx,
				chromedp.Navigate(fmt.Sprintf(b_start_url, start_offset)),
				chromedp.OuterHTML(".body_Chrome", &initialResponse),
			); err != nil {
				panic(err)
			}

			temp_total_results := strings.Split(
				strings.Split(
					strings.Split(initialResponse, `jobPaginationLegend`)[1], "</span>")[0], " ")
			total_results, _ := strconv.Atoi(temp_total_results[len(temp_total_results)-1])

			for i := 0; i <= total_results; i += number_results_per_page {
				fmt.Println(Gray(8-1, "Visiting"), Gray(8-1, fmt.Sprintf(b_start_url, i)))
				var pageResponse string
				if err := chromedp.Run(ctx,
					chromedp.Navigate(fmt.Sprintf(b_start_url, i)),
					chromedp.OuterHTML(".body_Chrome", &pageResponse),
				); err != nil {
					panic(err)
				}

				results_sections := strings.Split(pageResponse, `<li class="listSingleColumnItem">`)
				for q := 1; q < len(results_sections); q++ {
					elem := results_sections[q]
					result_title := strings.Split(strings.Split(strings.Split(elem, `<a href="`)[1], `">`)[1], `</a>`)[0]
					result_url := strings.Split(strings.Split(elem, `<a href="`)[1], `"`)[0]
					result_location := strings.Split(strings.Split(elem, `<span class="listSingleColumnItemMiscDataItem">`)[1], `</span>`)[0]
					result_date := strings.Split(strings.Split(elem, `Posted `)[1], `</span>`)[0]
					result_description := strings.Split(strings.Split(elem, `<div class="listSingleColumnItemDescription">`)[1], `<a`)[0]

					temp_elem_json := Job{
						result_title,
						result_url,
						result_location,
						result_date,
						result_description,
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
			results_marshal, err := json.Marshal(results)
			if err != nil {
				panic(err.Error())
			}
			response = Response{[]byte(results_marshal)}

		} else {
			file, _ := os.Open("response.html")
			pageResponse, _ := ioutil.ReadAll(file)
			var jobs []Job
			err := json.Unmarshal(pageResponse, &jobs)
			if err != nil {
				panic(err.Error())
			}

			for _, elem := range jobs {

				result_title := elem.Title
				result_url := elem.Url

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
	}
	return
}

func (runtime Runtime) Deloitte(
	version int, isLocal bool) (response Response, results []Result) {
	switch version {
	case 1:

		initial_file_name := "deloitteDepartments.html"

		type Job struct {
			Url         string
			Title       string
			Company     string
			Entity      string
			Department  string
			Id          string
			Type        string
			Date        string
			Description string
		}

		if !isLocal {

			t := &http.Transport{}
			t.RegisterProtocol("file", http.NewFileTransport(http.Dir("/")))
			dir, err := os.Getwd()
			if err != nil {
				panic(err.Error())
			}

			ctx, cancel := chromedp.NewContext(context.Background())
			defer cancel()

			var res []byte
			var initialResponse string
			if err := chromedp.Run(ctx,
				chromedp.Navigate("https://jobs2.deloitte.com/us/en/c/analytics-jobs"),
				chromedp.WaitReady(`.jobs-list-item`, chromedp.ByQuery),
				chromedp.EvaluateAsDevTools(`document.getElementsByClassName("clearall")[0].click()`, &res),
				chromedp.Sleep(SecondsSleep*time.Second),
				chromedp.WaitReady(`.phs-jobs-block`, chromedp.ByQuery),
				chromedp.OuterHTML("html", &initialResponse),
			); err != nil {
				panic(err)
			}
			SaveResponseToFileWithFileName(initialResponse, initial_file_name)

			c := colly.NewCollector()
			c.WithTransport(t)
			x := c.Clone()
			x.WithTransport(t)

			c.OnHTML("html", func(e *colly.HTMLElement) {
				e.ForEach(".jobs-list-item", func(_ int, el *colly.HTMLElement) {
					result_url := strings.Join(strings.Fields(strings.TrimSpace(el.ChildAttr("a", "href"))), " ")
					result_title := strings.Join(strings.Fields(strings.TrimSpace(el.ChildText("h4"))), " ")
					result_company := strings.Join(strings.Fields(strings.TrimSpace(el.ChildText(".memberfirm"))), " ")
					result_entity := strings.Join(strings.Fields(strings.TrimSpace(el.ChildText(".memberentity"))), " ")
					result_department := strings.Join(strings.Fields(strings.TrimSpace(el.ChildText(".job-category"))), " ")
					result_id := strings.Join(strings.Fields(strings.TrimSpace(el.ChildText(".job-id"))), " ")
					result_type := strings.Join(strings.Fields(strings.TrimSpace(el.ChildText(".job-type"))), " ")
					result_date := strings.Join(strings.Fields(strings.TrimSpace(el.ChildText(".job-postdate"))), " ")
					result_description := strings.Join(strings.Fields(strings.TrimSpace(el.ChildText(".job-description"))), " ")

					temp_elem_json := Job{
						result_url,
						result_title,
						result_company,
						result_entity,
						result_department,
						result_id,
						result_type,
						result_date,
						result_description,
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

				temp_number_of_jobs := e.ChildAttr(".search-bottom-count", "data-ph-at-total-jobs-text")
				number_of_jobs, _ := strconv.Atoi(temp_number_of_jobs)

				number_results_per_page := 50
				jobs_base_url := e.ChildAttr(`meta[property="og:url"]`, "content") + "?s=1&from=%d"

				for i := number_results_per_page; i <= number_of_jobs; i += number_results_per_page {

					sub_department_url := fmt.Sprintf(jobs_base_url, i)

					var departmentSubPageResponse string
					if err := chromedp.Run(ctx,
						chromedp.Navigate(sub_department_url),
						chromedp.WaitReady(`.jobs-list-item`, chromedp.ByQuery),
						chromedp.OuterHTML("html", &departmentSubPageResponse),
					); err != nil {
						panic(err)
					}

					sub_file_name := fmt.Sprintf("sub_department_url%d.html", i)
					SaveResponseToFileWithFileName(departmentSubPageResponse, sub_file_name)
					x.Visit("file:" + dir + "/" + sub_file_name)
					time.Sleep(SecondsSleep * time.Second)

					RemoveFileWithFileName(sub_file_name)
				}
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

			c.OnScraped(func(r *colly.Response) {
				results_marshal, err := json.Marshal(results)
				if err != nil {
					panic(err.Error())
				}
				response = Response{[]byte(results_marshal)}

				RemoveFileWithFileName(initial_file_name)
			})

			x.OnHTML("html", func(e *colly.HTMLElement) {
				e.ForEach(".jobs-list-item", func(_ int, el *colly.HTMLElement) {
					result_url := strings.Join(strings.Fields(strings.TrimSpace(el.ChildAttr("a", "href"))), " ")
					result_title := strings.Join(strings.Fields(strings.TrimSpace(el.ChildText("h4"))), " ")
					result_company := strings.Join(strings.Fields(strings.TrimSpace(el.ChildText(".memberfirm"))), " ")
					result_entity := strings.Join(strings.Fields(strings.TrimSpace(el.ChildText(".memberentity"))), " ")
					result_department := strings.Join(strings.Fields(strings.TrimSpace(el.ChildText(".job-category"))), " ")
					result_id := strings.Join(strings.Fields(strings.TrimSpace(el.ChildText(".job-id"))), " ")
					result_type := strings.Join(strings.Fields(strings.TrimSpace(el.ChildText(".job-type"))), " ")
					result_date := strings.Join(strings.Fields(strings.TrimSpace(el.ChildText(".job-postdate"))), " ")
					result_description := strings.Join(strings.Fields(strings.TrimSpace(el.ChildText(".job-description"))), " ")

					temp_elem_json := Job{
						result_url,
						result_title,
						result_company,
						result_entity,
						result_department,
						result_id,
						result_type,
						result_date,
						result_description,
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
			})

			x.OnRequest(func(r *colly.Request) {
				fmt.Println(Gray(8-1, "Visiting"), Gray(8-1, r.URL.String()))
			})

			x.OnError(func(r *colly.Response, err error) {
				fmt.Println(
					Red("Request URL:"), Red(r.Request.URL),
					Red("failed with response:"), Red(r),
					Red("\nError:"), Red(err))
			})

			c.WithTransport(t)
			c.Visit("file:" + dir + "/" + initial_file_name)
		} else {
			file, _ := os.Open("response.html")
			pageResponse, _ := ioutil.ReadAll(file)
			var jobs []Job
			err := json.Unmarshal(pageResponse, &jobs)
			if err != nil {
				panic(err.Error())
			}

			for _, elem := range jobs {

				result_title := elem.Title
				result_url := elem.Url

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
	}

	return
}

func (runtime Runtime) Bayer(
	version int, isLocal bool) (response Response, results []Result) {
	switch version {
	case 1:

		c := colly.NewCollector()

		start_url := "https://career.bayer.com/en/jobs-search?page=%d"
		base_job_url := "https://career.bayer.com%s"
		tag_title := "a"
		tag_date := ".views-field-field-job-last-modify-time"
		tag_country := ".views-field-field-job-country"
		tag_location := ".views-field-field-job-location"
		tag_last_page := ".pager__item--last"
		counter := 0

		type Job struct {
			Title    string
			Url      string
			Date     string
			Country  string
			Location string
		}

		var jsonJobs []Job

		if !isLocal {

			c.OnHTML(".content", func(e *colly.HTMLElement) {
				e.ForEach("tr", func(_ int, el *colly.HTMLElement) {
					result_title := el.ChildText(tag_title)
					result_url := fmt.Sprintf(base_job_url, el.ChildAttr(tag_title, "href"))
					result_date := el.ChildText(tag_date)
					result_country := el.ChildText(tag_country)
					result_location := el.ChildText(tag_location)

					_, err := netUrl.ParseRequestURI(result_url)
					if err == nil {

						temp_elem_json := Job{
							result_title,
							result_url,
							result_date,
							result_country,
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

				goqueryselect := e.DOM
				temp_last_page, _ := goqueryselect.Find(tag_last_page).Find("a").Attr("href")
				split_temp_last_page := strings.Split(temp_last_page, "=")
				last_page, _ := strconv.Atoi(split_temp_last_page[len(split_temp_last_page)-1])
				if counter <= last_page {
					counter++
					time.Sleep(SecondsSleep * time.Second)
					c.Visit(fmt.Sprintf(start_url, counter))
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
			c.Visit(fmt.Sprintf(start_url, 0))
		} else {
			var jsonJobs []Job
			c.OnResponse(func(r *colly.Response) {
				err := json.Unmarshal(r.Body, &jsonJobs)
				if err != nil {
					panic(err.Error())
				}
				for _, elem := range jsonJobs {
					elem_json, err := json.Marshal(elem)
					if err != nil {
						panic(err.Error())
					}
					results = append(results, Result{
						runtime.Name,
						elem.Url,
						elem.Title,
						elem_json,
					})
				}
			})
			c.OnScraped(func(r *colly.Response) {
				jsonJobs_marshal, err := json.Marshal(jsonJobs)
				if err != nil {
					panic(err.Error())
				}
				response = Response{[]byte(jsonJobs_marshal)}
			})

			t := &http.Transport{}
			t.RegisterProtocol("file", http.NewFileTransport(http.Dir("/")))
			c.WithTransport(t)
			dir, err := os.Getwd()
			if err != nil {
				panic(err.Error())
			}
			c.Visit("file:" + dir + "/response.html")
		}
	}
	return
}

func (runtime Runtime) Roche(
	version int, isLocal bool) (response Response, results []Result) {
	switch version {
	case 1:

		c := colly.NewCollector()

		start_url := "https://www.roche.com/toolbox/jobSearch.json?type=json&api=jobs&pageLength=%d&offset=%d"
		base_url := "https://www.roche.com%s"
		number_results_per_page := 300

		type JsonJobs struct {
			Jobs struct {
				Status       string `json:"status"`
				TotalMatches int    `json:"totalMatches"`
				Items        []struct {
					Title           string `json:"title"`
					DetailsURL      string `json:"detailsUrl"`
					OpenDate        string `json:"openDate"`
					JobLevel        string `json:"jobLevel"`
					PrimaryLocation struct {
						Country string `json:"country"`
						State   string `json:"state"`
						City    string `json:"city"`
					} `json:"primaryLocation"`
					PrimaryLocationCode struct {
						CountryCode string `json:"countryCode"`
						StateCode   string `json:"stateCode"`
						CityCode    string `json:"cityCode"`
					} `json:"primaryLocationCode"`
					OtherLocations     []interface{} `json:"otherLocations"`
					OtherLocationCodes []interface{} `json:"otherLocationCodes"`
					ReqID              string        `json:"reqId"`
					JobBoard           string        `json:"jobBoard"`
				} `json:"items"`
			} `json:"jobs"`
		}

		var jsonJobs JsonJobs

		c.OnResponse(func(r *colly.Response) {
			var tempJsonJobs JsonJobs
			err := json.Unmarshal(r.Body, &tempJsonJobs)
			if err != nil {
				panic(err.Error())
			}

			for _, elem := range tempJsonJobs.Jobs.Items {

				result_title := elem.Title
				result_url := fmt.Sprintf(base_url, elem.DetailsURL)

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

			jsonJobs.Jobs.Items = append(jsonJobs.Jobs.Items, tempJsonJobs.Jobs.Items...)

			total_matches := tempJsonJobs.Jobs.TotalMatches
			total_pages := total_matches / number_results_per_page
			for i := 1; i <= total_pages; i++ {
				time.Sleep(SecondsSleep * time.Second)
				c.Visit(fmt.Sprintf(start_url, number_results_per_page, i))
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
			c.Visit(fmt.Sprintf(start_url, number_results_per_page, 0))
		}
	}
	return
}

func (runtime Runtime) Msd(
	version int, isLocal bool) (response Response, results []Result) {
	switch version {
	case 1:

		type JsonJobs struct {
			EagerLoadRefineSearch struct {
				Status    int `json:"status"`
				Hits      int `json:"hits"`
				TotalHits int `json:"totalHits"`
				Data      struct {
					Jobs []struct {
						Country           string    `json:"country"`
						CityState         string    `json:"cityState"`
						City              string    `json:"city"`
						MlSkills          []string  `json:"ml_skills"`
						Type              string    `json:"type"`
						Experience        string    `json:"experience,omitempty"`
						Locale            string    `json:"locale"`
						Title             string    `json:"title"`
						MultiLocation     []string  `json:"multi_location"`
						PostedDate        time.Time `json:"postedDate"`
						JobSeqNo          string    `json:"jobSeqNo"`
						DescriptionTeaser string    `json:"descriptionTeaser"`
						DateCreated       time.Time `json:"dateCreated"`
						State             string    `json:"state"`
						CityStateCountry  string    `json:"cityStateCountry"`
						Department        string    `json:"department,omitempty"`
						VisibilityType    string    `json:"visibilityType"`
						SiteType          string    `json:"siteType"`
						IsMultiCategory   bool      `json:"isMultiCategory"`
						ReqID             string    `json:"reqId"`
						JobID             string    `json:"jobId"`
						Badge             string    `json:"badge"`
						JobVisibility     []string  `json:"jobVisibility"`
						IsMultiLocation   bool      `json:"isMultiLocation"`
						Location          string    `json:"location"`
						Category          string    `json:"category"`
						ExternalApply     bool      `json:"externalApply"`
						SubCategory       string    `json:"subCategory,omitempty"`
						Industry          string    `json:"industry,omitempty"`
						WorkLocation      string    `json:"workLocation,omitempty"`
						Address           string    `json:"address,omitempty"`
						MultiCategory     []string  `json:"multi_category,omitempty"`
						ApplyURL          string    `json:"applyUrl,omitempty"`
					} `json:"jobs"`
				} `json:"data"`
			} `json:"eagerLoadRefineSearch"`
		}

		start_url := "https://jobs.msd.com/gb/en/search-results?s=1&from=%d"
		base_job_url := "https://jobs.msd.com/gb/en/job/%s"

		var jsonJobs JsonJobs

		if !isLocal {

			ctx, cancel := chromedp.NewContext(context.Background())
			defer cancel()

			var initialResponse string
			if err := chromedp.Run(ctx,
				chromedp.Navigate(fmt.Sprintf(start_url, 0)),
				chromedp.OuterHTML(".desktop", &initialResponse),
			); err != nil {
				panic(err)
			}

			temp_jsonjob_section := strings.Split(
				strings.Split(
					initialResponse, `"eagerLoadRefineSearch":`)[1], `,"jobwidgetsettings`)[0]
			jsonjobs_sections := `{"eagerLoadRefineSearch":` + temp_jsonjob_section + "}"

			var tempJsonJobs JsonJobs
			err := json.Unmarshal([]byte(jsonjobs_sections), &tempJsonJobs)
			if err != nil {
				panic(err.Error())
			}

			items_per_page := tempJsonJobs.EagerLoadRefineSearch.Hits
			total_matches := tempJsonJobs.EagerLoadRefineSearch.TotalHits
			total_pages := total_matches / items_per_page
			for i := 1; i <= total_pages+1; i++ {

				jobs_url := fmt.Sprintf(start_url, i*items_per_page)
				fmt.Println(Gray(8-1, "Visiting"), Gray(8-1, jobs_url))

				var jobResponse string
				if err := chromedp.Run(ctx,
					chromedp.Navigate(jobs_url),
					chromedp.OuterHTML(".desktop", &jobResponse),
				); err != nil {
					panic(err)
				}

				temp_jsonjob_section := strings.Split(
					strings.Split(
						jobResponse, `"eagerLoadRefineSearch":`)[1], `,"jobwidgetsettings`)[0]
				jsonjobs_sections := `{"eagerLoadRefineSearch":` + temp_jsonjob_section + "}"

				var tempJson JsonJobs
				err := json.Unmarshal([]byte(jsonjobs_sections), &tempJson)
				if err != nil {
					panic(err.Error())
				}

				for _, elem := range tempJson.EagerLoadRefineSearch.Data.Jobs {

					result_title := elem.Title
					result_url := fmt.Sprintf(base_job_url, elem.JobID)

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

				jsonJobs.EagerLoadRefineSearch.Data.Jobs = append(
					jsonJobs.EagerLoadRefineSearch.Data.Jobs,
					tempJson.EagerLoadRefineSearch.Data.Jobs...)
			}

			results_marshal, err := json.Marshal(jsonJobs)
			if err != nil {
				panic(err.Error())
			}
			response = Response{[]byte(results_marshal)}
		} else {
			file, _ := os.Open("response.html")
			pageResponse, _ := ioutil.ReadAll(file)
			var jsonJobs JsonJobs
			err := json.Unmarshal(pageResponse, &jsonJobs)
			if err != nil {
				panic(err.Error())
			}

			for _, elem := range jsonJobs.EagerLoadRefineSearch.Data.Jobs {

				result_title := elem.Title
				result_url := fmt.Sprintf(base_job_url, elem.JobID)

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
	}
	return
}

func (runtime Runtime) Subitoit(
	version int, isLocal bool) (response Response, results []Result) {
	switch version {
	case 1:

		c := colly.NewCollector()

		start_url := "https://info.subito.it/lavora-con-noi.htm"
		tag_section := ".work-openings"
		tag_result := ".list-box-item"
		tag_title := "a"
		tag_department := "h4"

		type Job struct {
			Url        string
			Title      string
			Department string
		}

		c.OnHTML(tag_section, func(e *colly.HTMLElement) {
			e.ForEach(tag_result, func(_ int, el *colly.HTMLElement) {
				result_title := el.ChildText(tag_title)
				result_url := el.ChildAttr(tag_title, "href")
				result_department := el.ChildText(tag_department)

				_, err := netUrl.ParseRequestURI(result_url)
				if err == nil {

					temp_elem_json := Job{
						result_title,
						result_url,
						result_department,
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

		c.OnRequest(func(r *colly.Request) {
			fmt.Println(Gray(8-1, "Visiting"), Gray(8-1, r.URL.String()))
		})

		c.OnResponse(func(r *colly.Response) {
			response = Response{r.Body}
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
			c.Visit(start_url)
		}
	}
	return
}

func (runtime Runtime) Square(
	version int, isLocal bool) (response Response, results []Result) {
	switch version {
	case 1:

		c := colly.NewCollector()

		start_url := "https://api.smartrecruiters.com/v1/companies/square/postings?offset=%d"
		base_job_url := "https://www.smartrecruiters.com/Square/%s"
		number_results_per_page := 100

		type Jobs struct {
			Offset     int `json:"offset"`
			Limit      int `json:"limit"`
			TotalFound int `json:"totalFound"`
			Content    []struct {
				ID        string `json:"id"`
				Name      string `json:"name"`
				UUID      string `json:"uuid"`
				RefNumber string `json:"refNumber"`
				Company   struct {
					Identifier string `json:"identifier"`
					Name       string `json:"name"`
				} `json:"company"`
				ReleasedDate time.Time `json:"releasedDate"`
				Location     struct {
					City    string `json:"city"`
					Region  string `json:"region"`
					Country string `json:"country"`
					Remote  bool   `json:"remote"`
				} `json:"location"`
				Industry struct {
					ID    string `json:"id"`
					Label string `json:"label"`
				} `json:"industry"`
				Department struct {
					ID    string `json:"id"`
					Label string `json:"label"`
				} `json:"department"`
				Function struct {
					ID    string `json:"id"`
					Label string `json:"label"`
				} `json:"function"`
				TypeOfEmployment struct {
					Label string `json:"label"`
				} `json:"typeOfEmployment"`
				ExperienceLevel struct {
					ID    string `json:"id"`
					Label string `json:"label"`
				} `json:"experienceLevel"`
				CustomField []struct {
					FieldID    string `json:"fieldId"`
					FieldLabel string `json:"fieldLabel"`
					ValueID    string `json:"valueId"`
					ValueLabel string `json:"valueLabel"`
				} `json:"customField"`
				Ref     string `json:"ref"`
				Creator struct {
					Name string `json:"name"`
				} `json:"creator"`
				Language struct {
					Code        string `json:"code"`
					Label       string `json:"label"`
					LabelNative string `json:"labelNative"`
				} `json:"language"`
			} `json:"content"`
		}

		var jsonJobs Jobs

		c.OnResponse(func(r *colly.Response) {
			var tempJson Jobs
			err := json.Unmarshal(r.Body, &tempJson)
			if err != nil {
				panic(err.Error())
			}

			for _, elem := range tempJson.Content {

				result_title := elem.Name
				result_url := fmt.Sprintf(base_job_url, elem.ID)

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

			jsonJobs.Content = append(jsonJobs.Content, tempJson.Content...)

			if isLocal {
				return
			} else {
				total_matches := tempJson.TotalFound
				total_pages := total_matches / number_results_per_page
				for i := 1; i <= total_pages; i++ {
					time.Sleep(SecondsSleep * time.Second)
					c.Visit(fmt.Sprintf(start_url, number_results_per_page*i))
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
			c.Visit(fmt.Sprintf(start_url, 0))
		}
	}
	return
}

func (runtime Runtime) Facebook(
	version int, isLocal bool) (response Response, results []Result) {
	switch version {
	case 1:

		c := colly.NewCollector()

		start_url := "https://www.facebook.com/careers/jobs?results_per_page=100&page=%d"
		base_job_url := "https://www.facebook.com%s"
		number_results_per_page := 100

		type Job struct {
			Title    string
			Url      string
			Location string
			Info     string
		}

		if !isLocal {

			c.OnHTML("#search_result", func(e *colly.HTMLElement) {
				e.ForEach("a", func(_ int, el *colly.HTMLElement) {
					goqueryselector := el.DOM
					result_url := fmt.Sprintf(base_job_url, el.Attr("href"))
					result_title := el.ChildText("._8sel")
					result_location := goqueryselector.Find("._97fe ._8sen").Find("span").Text()

					var result_info []string
					temp_result_info := el.ChildTexts("._8see")
					for _, elem := range temp_result_info {
						if !strings.Contains(elem, "+") {
							result_info = append(result_info, elem)
						}
					}

					_, err := netUrl.ParseRequestURI(result_url)
					if err == nil {
						temp_elem_json := Job{
							result_title,
							result_url,
							result_location,
							strings.Join(result_info, " - "),
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

				array_number_results := strings.Split(e.ChildText("._6v-m"), " ")
				string_number_results := array_number_results[len(array_number_results)-1]
				number_results, _ := strconv.Atoi(string_number_results)
				total_pages := number_results / number_results_per_page

				for i := 2; i <= total_pages; i++ {
					time.Sleep(SecondsSleep * time.Second)
					c.Visit(fmt.Sprintf(start_url, i))
				}
			})

			c.OnRequest(func(r *colly.Request) {
				fmt.Println(Gray(8-1, "Visiting"), Gray(8-1, r.URL.String()))
			})

			c.OnScraped(func(r *colly.Response) {
				results_marshal, err := json.Marshal(results)
				if err != nil {
					panic(err.Error())
				}
				response = Response{[]byte(results_marshal)}
			})

			c.OnError(func(r *colly.Response, err error) {
				fmt.Println(
					Red("Request URL:"), Red(r.Request.URL),
					Red("failed with response:"), Red(r),
					Red("\nError:"), Red(err))
			})

			c.Visit(fmt.Sprintf(start_url, 1))
		} else {
			file, _ := os.Open("response.html")
			pageResponse, _ := ioutil.ReadAll(file)
			var jsonJobs []Job
			err := json.Unmarshal(pageResponse, &jsonJobs)
			if err != nil {
				panic(err.Error())
			}

			for _, elem := range jsonJobs {

				result_title := elem.Title
				result_url := elem.Url

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
	}
	return
}

func (runtime Runtime) Paintgun(
	version int, isLocal bool) (
	response Response, results []Result) {
	switch version {
	case 1:

		c := colly.NewCollector()

		start_url := "https://join.com/api/public/companies/9628/jobs?page=1&pageSize=100"
		job_base_url := "https://paintgun.join.com/jobs/%s"

		type JsonJobs struct {
			Items []struct {
				ID               int       `json:"id"`
				LastID           int       `json:"lastId"`
				OriginIDParam    string    `json:"originIdParam"`
				IDParam          string    `json:"idParam"`
				Title            string    `json:"title"`
				PlaceID          string    `json:"placeId"`
				Zip              string    `json:"zip"`
				IsRemote         bool      `json:"isRemote"`
				CountryID        int       `json:"countryId"`
				EmploymentTypeID int       `json:"employmentTypeId"`
				LanguageID       int       `json:"languageId"`
				CategoryID       int       `json:"categoryId"`
				CreatedAt        time.Time `json:"createdAt"`
				EmploymentType   struct {
					ID          int       `json:"id"`
					Name        string    `json:"name"`
					Slug        string    `json:"slug"`
					CreatedAt   time.Time `json:"createdAt"`
					UpdatedAt   time.Time `json:"updatedAt"`
					IsNullValue bool      `json:"isNullValue"`
					GoogleType  string    `json:"googleType"`
					NameEn      string    `json:"nameEn"`
					NameDe      string    `json:"nameDe"`
					NameIt      string    `json:"nameIt"`
					NameFr      string    `json:"nameFr"`
					SortOrder   int       `json:"sortOrder"`
				} `json:"employmentType"`
				Language struct {
					ID        int       `json:"id"`
					Name      string    `json:"name"`
					Iso6391   string    `json:"iso6391"`
					IsDefault bool      `json:"isDefault"`
					CreatedAt time.Time `json:"createdAt"`
					UpdatedAt time.Time `json:"updatedAt"`
					Locale    string    `json:"locale"`
				} `json:"language"`
				Country struct {
					ID        int       `json:"id"`
					Name      string    `json:"name"`
					Iso3166   string    `json:"iso3166"`
					CreatedAt time.Time `json:"createdAt"`
					UpdatedAt time.Time `json:"updatedAt"`
				} `json:"country"`
				UnifiedDescription bool `json:"unifiedDescription"`
				PendingDeletion    bool `json:"pendingDeletion"`
				EducationID        int  `json:"educationId,omitempty"`
				Education          struct {
					ID          int       `json:"id"`
					Name        string    `json:"name"`
					Slug        string    `json:"slug"`
					CreatedAt   time.Time `json:"createdAt"`
					UpdatedAt   time.Time `json:"updatedAt"`
					IsNullValue bool      `json:"isNullValue"`
				} `json:"education,omitempty"`
			} `json:"items"`
			Pagination struct {
				RowCount  int `json:"rowCount"`
				PageCount int `json:"pageCount"`
				Page      int `json:"page"`
				PageSize  int `json:"pageSize"`
			} `json:"pagination"`
			Aggregations      []interface{} `json:"aggregations"`
			UsingFallbackData bool          `json:"usingFallbackData"`
		}

		var jsonJobs JsonJobs

		c.OnResponse(func(r *colly.Response) {
			var tempJson JsonJobs
			err := json.Unmarshal(r.Body, &tempJson)
			if err != nil {
				panic(err.Error())
			}

			for _, elem := range tempJson.Items {

				result_title := elem.Title
				result_url := fmt.Sprintf(job_base_url, elem.IDParam)

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

			jsonJobs.Items = append(jsonJobs.Items, tempJson.Items...)
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
			c.Visit(start_url)
		}
	}
	return
}

func (runtime Runtime) Nen(
	version int, isLocal bool) (
	response Response, results []Result) {
	switch version {
	case 1:

		results = append(results, Result{
			runtime.Name,
			"Salesforce Lead",
			"https://www.linkedin.com/jobs/view/1947567619",
			[]byte("{}"),
		})

		results_marshal, err := json.Marshal(results)
		if err != nil {
			panic(err.Error())
		}
		response = Response{[]byte(results_marshal)}
	}
	return
}

func (runtime Runtime) Amboss(
	version int, isLocal bool) (
	response Response, results []Result) {
	switch version {
	case 1:

		c := colly.NewCollector()

		start_url := "https://www.amboss.com/us/career-opportunities"
		job_base_url := "https://www.amboss.com%s"
		tag_main_section := ".jobs-list"
		tag_job_section := "._pwggpq"
		tag_title := "._pulkya"
		tag_location := "._1f1zsnz"

		type Job struct {
			Url      string
			Title    string
			Location string
		}

		c.OnHTML(tag_main_section, func(e *colly.HTMLElement) {
			e.ForEach(tag_job_section, func(_ int, el *colly.HTMLElement) {
				result_url := fmt.Sprintf(job_base_url, el.Attr("href"))
				result_title := el.ChildText(tag_title)
				result_location := el.ChildText(tag_location)

				_, err := netUrl.ParseRequestURI(result_url)
				if err == nil {

					temp_elem_json := Job{
						result_url,
						result_title,
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
			c.Visit(start_url)
		}
	}
	return
}

func (runtime Runtime) Chatterbug(
	version int, isLocal bool) (
	response Response, results []Result) {
	switch version {
	case 1:

		c := colly.NewCollector()

		p_start_url := "https://chatterbug.recruitee.com/api/offers"

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
				MinHours           int           `json:"min_hours"`
				MaxHours           int           `json:"max_hours"`
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
				Department         interface{}   `json:"department"`
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
				OpenQuestions []interface{} `json:"open_questions"`
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
	}
	return
}

func (runtime Runtime) Infarm(
	version int, isLocal bool) (
	response Response, results []Result) {
	switch version {
	case 1:

		c := colly.NewCollector()

		start_url := "https://api.greenhouse.io/v1/boards/infarm/jobs"

		type JsonJobs struct {
			Jobs []struct {
				AbsoluteURL    string `json:"absolute_url"`
				DataCompliance []struct {
					Type            string      `json:"type"`
					RequiresConsent bool        `json:"requires_consent"`
					RetentionPeriod interface{} `json:"retention_period"`
				} `json:"data_compliance"`
				Education     string `json:"education,omitempty"`
				InternalJobID int    `json:"internal_job_id"`
				Location      struct {
					Name string `json:"name"`
				} `json:"location"`
				Metadata []struct {
					ID        int         `json:"id"`
					Name      string      `json:"name"`
					Value     interface{} `json:"value"`
					ValueType string      `json:"value_type"`
				} `json:"metadata"`
				ID            int    `json:"id"`
				UpdatedAt     string `json:"updated_at"`
				RequisitionID string `json:"requisition_id"`
				Title         string `json:"title"`
			} `json:"jobs"`
			Meta struct {
				Total int `json:"total"`
			} `json:"meta"`
		}

		var jsonJobs JsonJobs

		c.OnResponse(func(r *colly.Response) {
			var tempJson JsonJobs
			err := json.Unmarshal(r.Body, &tempJson)
			if err != nil {
				panic(err.Error())
			}

			for _, elem := range tempJson.Jobs {

				result_title := elem.Title
				result_url := elem.AbsoluteURL

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

			jsonJobs.Jobs = append(jsonJobs.Jobs, tempJson.Jobs...)
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
			c.Visit(start_url)
		}
	}
	return
}

func (runtime Runtime) Pitch(
	version int, isLocal bool) (
	response Response, results []Result) {
	switch version {
	case 1:

		c := colly.NewCollector()

		start_url := "https://api.greenhouse.io/v1/boards/pitch/jobs"

		type JsonJobs struct {
			Jobs []struct {
				AbsoluteURL    string `json:"absolute_url"`
				DataCompliance []struct {
					Type            string      `json:"type"`
					RequiresConsent bool        `json:"requires_consent"`
					RetentionPeriod interface{} `json:"retention_period"`
				} `json:"data_compliance"`
				Education     string `json:"education,omitempty"`
				InternalJobID int    `json:"internal_job_id"`
				Location      struct {
					Name string `json:"name"`
				} `json:"location"`
				Metadata []struct {
					ID        int         `json:"id"`
					Name      string      `json:"name"`
					Value     interface{} `json:"value"`
					ValueType string      `json:"value_type"`
				} `json:"metadata"`
				ID            int    `json:"id"`
				UpdatedAt     string `json:"updated_at"`
				RequisitionID string `json:"requisition_id"`
				Title         string `json:"title"`
			} `json:"jobs"`
			Meta struct {
				Total int `json:"total"`
			} `json:"meta"`
		}

		var jsonJobs JsonJobs

		c.OnResponse(func(r *colly.Response) {
			var tempJson JsonJobs
			err := json.Unmarshal(r.Body, &tempJson)
			if err != nil {
				panic(err.Error())
			}

			for _, elem := range tempJson.Jobs {

				result_title := elem.Title
				result_url := elem.AbsoluteURL

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

			jsonJobs.Jobs = append(jsonJobs.Jobs, tempJson.Jobs...)
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
			c.Visit(start_url)
		}
	}
	return
}

func (runtime Runtime) Beat81(
	version int, isLocal bool) (
	response Response, results []Result) {
	switch version {
	case 1:

		c := colly.NewCollector()

		c_start_url := "https://api.lever.co/v0/postings/beat81?mode=json"

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

func (runtime Runtime) Careerfoundry(
	version int, isLocal bool) (
	response Response, results []Result) {
	switch version {
	case 1:

		c := colly.NewCollector()

		p_start_url := "https://careerfoundry.recruitee.com/api/offers"

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
				MinHours           int           `json:"min_hours"`
				MaxHours           int           `json:"max_hours"`
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
				Department         interface{}   `json:"department"`
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
				OpenQuestions []interface{} `json:"open_questions"`
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
	}
	return
}

func (runtime Runtime) Casparhealth(
	version int, isLocal bool) (
	response Response, results []Result) {
	switch version {
	case 1:

		c := colly.NewCollector()

		url := "https://goreha-jobs.personio.de"
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

func (runtime Runtime) Ecosia(
	version int, isLocal bool) (
	response Response, results []Result) {
	switch version {
	case 1:

		c := colly.NewCollector()

		c_start_url := "https://api.lever.co/v0/postings/ecosia?mode=json"

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

func (runtime Runtime) Forto(
	version int, isLocal bool) (
	response Response, results []Result) {
	switch version {
	case 1:

		c := colly.NewCollector()

		c_start_url := "https://api.lever.co/v0/postings/forto?mode=json"

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

func (runtime Runtime) Idagio(
	version int, isLocal bool) (
	response Response, results []Result) {
	switch version {
	case 1:

		c := colly.NewCollector()

		url := "https://idagio-jobs.personio.de"
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

func (runtime Runtime) Joblift(
	version int, isLocal bool) (
	response Response, results []Result) {
	switch version {
	case 1:

		c := colly.NewCollector()

		start_url := "https://joblift-talent.freshteam.com/jobs"
		tag_main_section := ".job-role-list"
		tag_department_section := "li:not([class])"
		tag_department := ".role-title"
		tag_job_section := ".job-list-info"
		tag_title := ".job-title"
		tag_location := ".location-info"
		job_base_url := "https://joblift-talent.freshteam.com%s"

		type Job struct {
			Url        string
			Title      string
			Location   string
			Department string
		}

		c.OnHTML(tag_main_section, func(e *colly.HTMLElement) {
			e.ForEach(tag_department_section, func(_ int, el *colly.HTMLElement) {
				result_department := strings.Split(el.ChildText(tag_department), "-")[0]
				el.ForEach(tag_job_section, func(_ int, ell *colly.HTMLElement) {
					result_url := fmt.Sprintf(job_base_url, ell.ChildAttr("a", "href"))
					result_title := ell.ChildText(tag_title)
					result_location := strings.Split(ell.ChildText(tag_location), "\n")[0]

					_, err := netUrl.ParseRequestURI(result_url)
					if err == nil {

						temp_elem_json := Job{
							result_url,
							result_title,
							result_location,
							result_department,
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
			c.Visit(start_url)
		}
	}
	return
}

func (runtime Runtime) Kontist(
	version int, isLocal bool) (
	response Response, results []Result) {
	switch version {
	case 1:

		c := colly.NewCollector()

		c_start_url := "https://kontist.com/careers/jobs.json"

		type Jobs struct {
			Data []struct {
				ID         string `json:"id"`
				ActiveTime []struct {
					OpenedAt int         `json:"opened_at"`
					ClosedAt interface{} `json:"closed_at"`
				} `json:"active_time"`
				ApplicationFormURL string `json:"application_form_url"`
				CreatedAt          int    `json:"created_at"`
				JobURL             string `json:"job_url"`
				SeoContent         struct {
				} `json:"seo_content"`
				ShareImageURL       string `json:"share_image_url"`
				Status              string `json:"status"`
				Title               string `json:"title"`
				TmpDepartment       string `json:"tmp_department"`
				TmpLocation         string `json:"tmp_location"`
				TotalCandidateCount int    `json:"total_candidate_count"`
				Type                string `json:"type"`
			} `json:"data"`
			HasMore bool `json:"has_more"`
		}

		var jsonJobs Jobs

		c.OnResponse(func(r *colly.Response) {
			var tempJson Jobs
			err := json.Unmarshal(r.Body, &tempJson)
			if err != nil {
				panic(err.Error())
			}

			for _, elem := range tempJson.Data {

				result_title := elem.Title
				result_url := elem.JobURL

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
			c.Visit(c_start_url)
		}
	}
	return
}

func (runtime Runtime) Medloop(
	version int, isLocal bool) (
	response Response, results []Result) {
	switch version {
	case 1:

		c := colly.NewCollector()

		start_url := "https://api.greenhouse.io/v1/boards/medloop/jobs"

		type JsonJobs struct {
			Jobs []struct {
				AbsoluteURL    string `json:"absolute_url"`
				DataCompliance []struct {
					Type            string      `json:"type"`
					RequiresConsent bool        `json:"requires_consent"`
					RetentionPeriod interface{} `json:"retention_period"`
				} `json:"data_compliance"`
				Education     string `json:"education,omitempty"`
				InternalJobID int    `json:"internal_job_id"`
				Location      struct {
					Name string `json:"name"`
				} `json:"location"`
				Metadata []struct {
					ID        int         `json:"id"`
					Name      string      `json:"name"`
					Value     interface{} `json:"value"`
					ValueType string      `json:"value_type"`
				} `json:"metadata"`
				ID            int    `json:"id"`
				UpdatedAt     string `json:"updated_at"`
				RequisitionID string `json:"requisition_id"`
				Title         string `json:"title"`
			} `json:"jobs"`
			Meta struct {
				Total int `json:"total"`
			} `json:"meta"`
		}

		var jsonJobs JsonJobs

		c.OnResponse(func(r *colly.Response) {
			var tempJson JsonJobs
			err := json.Unmarshal(r.Body, &tempJson)
			if err != nil {
				panic(err.Error())
			}

			for _, elem := range tempJson.Jobs {

				result_title := elem.Title
				result_url := elem.AbsoluteURL

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

			jsonJobs.Jobs = append(jsonJobs.Jobs, tempJson.Jobs...)
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
			c.Visit(start_url)
		}
	}
	return
}

func (runtime Runtime) Medwing(
	version int, isLocal bool) (
	response Response, results []Result) {
	switch version {
	case 1:

		c := colly.NewCollector()

		c_start_url := "https://team.medwing.com/wp-json/wp/v2/jobs"

		type Jobs []struct {
			ID      int    `json:"id"`
			Date    string `json:"date"`
			DateGmt string `json:"date_gmt"`
			GUID    struct {
				Rendered string `json:"rendered"`
			} `json:"guid"`
			Modified    string `json:"modified"`
			ModifiedGmt string `json:"modified_gmt"`
			Slug        string `json:"slug"`
			Status      string `json:"status"`
			Type        string `json:"type"`
			Link        string `json:"link"`
			Title       struct {
				Rendered string `json:"rendered"`
			} `json:"title"`
			Content struct {
				Rendered  string `json:"rendered"`
				Protected bool   `json:"protected"`
			} `json:"content"`
			Excerpt struct {
				Rendered  string `json:"rendered"`
				Protected bool   `json:"protected"`
			} `json:"excerpt"`
			Author        int           `json:"author"`
			FeaturedMedia int           `json:"featured_media"`
			CommentStatus string        `json:"comment_status"`
			PingStatus    string        `json:"ping_status"`
			Template      string        `json:"template"`
			Format        string        `json:"format"`
			Meta          []interface{} `json:"meta"`
			Kategorie     []int         `json:"kategorie"`
			Department    []int         `json:"department"`
			Einstieg      []int         `json:"einstieg"`
			Vertrag       []int         `json:"vertrag"`
			Location      []int         `json:"location"`
		}

		var jsonJobs Jobs

		c.OnResponse(func(r *colly.Response) {
			var tempJson Jobs
			err := json.Unmarshal(r.Body, &tempJson)
			if err != nil {
				panic(err.Error())
			}

			for _, elem := range tempJson {

				result_title := elem.Title.Rendered
				result_url := elem.Link

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

func (runtime Runtime) Merantix(
	version int, isLocal bool) (
	response Response, results []Result) {
	switch version {
	case 1:

		c := colly.NewCollector()

		url := "https://merantix.bamboohr.com/jobs/embed2.php?departmentId=0"
		main_section_tag := ".BambooHR-ATS-Department-List"
		tag_division_section := ".BambooHR-ATS-Department-Item"
		tag_division := ".BambooHR-ATS-Department-Header"
		tag_job_section := ".BambooHR-ATS-Jobs-Item"
		tag_title := "a"
		tag_location := "span"

		type Job struct {
			Division string
			Title    string
			Url      string
			Location string
		}

		c.OnHTML(main_section_tag, func(e *colly.HTMLElement) {
			e.ForEach(tag_division_section, func(_ int, el *colly.HTMLElement) {
				result_division := strings.TrimSpace(el.ChildText(tag_division))
				el.ForEach(tag_job_section, func(_ int, ell *colly.HTMLElement) {
					result_title := ell.ChildText(tag_title)
					result_url := "https:" + ell.ChildAttr(tag_title, "href")
					result_location := ell.ChildText(tag_location)

					_, err := netUrl.ParseRequestURI(result_url)
					if err == nil {

						temp_elem_json := Job{
							result_division,
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

func (runtime Runtime) Ninox(
	version int, isLocal bool) (
	response Response, results []Result) {
	switch version {
	case 1:

		c := colly.NewCollector()

		start_url := "https://ninox.com/en/jobs"
		job_base_url := "https://ninox.com/%s"
		tag_job_section := ".job-new"
		tag_title := "h4"
		tag_location := ".jobs-j-openinglugar"

		type Job struct {
			Url      string
			Title    string
			Location string
		}

		c.OnHTML(tag_job_section, func(e *colly.HTMLElement) {
			result_url := fmt.Sprintf(job_base_url, e.ChildAttr("a", "href"))
			result_title := e.ChildText(tag_title)
			result_location := e.ChildText(tag_location)

			_, err := netUrl.ParseRequestURI(result_url)
			if err == nil {

				temp_elem_json := Job{
					result_url,
					result_title,
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
			c.Visit(start_url)
		}
	}
	return
}

func (runtime Runtime) Zenjob(
	version int, isLocal bool) (
	response Response, results []Result) {
	switch version {
	case 1:

		c := colly.NewCollector()

		start_url := "https://zenjob.teamtailor.com"
		tag_main_section := ".jobs"
		tag_job_section := "li"
		tag_title := ".title"
		tag_location := ".meta"
		job_base_url := "https://zenjob.teamtailor.com%s"

		type Job struct {
			Url      string
			Title    string
			Location string
		}

		c.OnHTML(tag_main_section, func(e *colly.HTMLElement) {
			e.ForEach(tag_job_section, func(_ int, el *colly.HTMLElement) {
				result_url := fmt.Sprintf(job_base_url, el.ChildAttr("a", "href"))
				result_title := el.ChildText(tag_title)
				result_location := el.ChildText(tag_location)

				_, err := netUrl.ParseRequestURI(result_url)
				if err == nil {

					temp_elem_json := Job{
						result_url,
						result_title,
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
			c.Visit(start_url)
		}
	}
	return
}

func (runtime Runtime) Plantix(
	version int, isLocal bool) (
	response Response, results []Result) {
	switch version {
	case 1:

		c := colly.NewCollector()

		p_start_url := "https://plantix.recruitee.com/api/offers"

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
				MinHours           int           `json:"min_hours"`
				MaxHours           int           `json:"max_hours"`
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
				Department         interface{}   `json:"department"`
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
				OpenQuestions []interface{} `json:"open_questions"`
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
	}
	return
}

func (runtime Runtime) Coachhub(
	version int, isLocal bool) (
	response Response, results []Result) {
	switch version {
	case 1:

		c := colly.NewCollector()

		start_url := "https://coachhub-jobs.personio.de/"
		tag_main := ".panel-container"
		tag_job_section := ".recent-job-list"
		tag_title := "h6"
		tag_location := "p"

		type Job struct {
			Url      string
			Title    string
			Location string
		}

		c.OnHTML(tag_main, func(e *colly.HTMLElement) {
			e.ForEach(tag_job_section, func(_ int, el *colly.HTMLElement) {
				result_url := el.ChildAttr("a", "href")
				result_title := el.ChildText(tag_title)
				result_location := strings.Split(el.ChildText(tag_location), "·")[1]

				_, err := netUrl.ParseRequestURI(result_url)
				if err == nil {

					temp_elem_json := Job{
						result_url,
						result_title,
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
			c.Visit(start_url)
		}
	}
	return
}

func (runtime Runtime) Raisin(
	version int, isLocal bool) (
	response Response, results []Result) {
	switch version {
	case 1:

		c := colly.NewCollector()

		url := "https://www.raisin.com/careers/"
		tag_main_section := ".module-jobs-listing"
		tag_title := ".col-sm-9"

		type Job struct {
			Url   string
			Title string
		}

		c.OnHTML(tag_main_section, func(e *colly.HTMLElement) {
			e.ForEach("a", func(_ int, el *colly.HTMLElement) {
				result_url := el.Attr("href")
				result_title := el.ChildText(tag_title)

				_, err := netUrl.ParseRequestURI(result_url)
				if err == nil {

					temp_elem_json := Job{
						result_url,
						result_title,
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

func (runtime Runtime) Acatus(
	version int, isLocal bool) (
	response Response, results []Result) {
	switch version {
	case 1:

		c := colly.NewCollector()

		url := "https://acatus-jobs.personio.de/?language=en#all"
		main_tag_job := ".job-list-desc"
		tag_title := "a"
		tag_info := "p"
		separator := "·"

		type Job struct {
			Title    string
			Url      string
			Type     string
			Location string
		}

		c.OnHTML("body", func(e *colly.HTMLElement) {
			e.ForEach(main_tag_job, func(_ int, el *colly.HTMLElement) {
				result_url := el.ChildAttr(tag_title, "href")
				result_title := el.ChildText(tag_title)
				result_info := strings.Split(el.ChildText(tag_info), separator)
				result_type := strings.Join(strings.Fields(strings.TrimSpace(result_info[0])), " ")
				result_location := strings.Join(strings.Fields(strings.TrimSpace(result_info[1])), " ")

				_, err := netUrl.ParseRequestURI(result_url)
				if err == nil {

					temp_elem_json := Job{
						result_title,
						result_url,
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

func (runtime Runtime) Adjust(
	version int, isLocal bool) (
	response Response, results []Result) {
	switch version {
	case 1:

		c := colly.NewCollector()

		start_url := "https://api.greenhouse.io/v1/boards/adjust/jobs"

		type JsonJobs struct {
			Jobs []struct {
				AbsoluteURL    string `json:"absolute_url"`
				DataCompliance []struct {
					Type            string      `json:"type"`
					RequiresConsent bool        `json:"requires_consent"`
					RetentionPeriod interface{} `json:"retention_period"`
				} `json:"data_compliance"`
				Education     string `json:"education,omitempty"`
				InternalJobID int    `json:"internal_job_id"`
				Location      struct {
					Name string `json:"name"`
				} `json:"location"`
				Metadata []struct {
					ID        int         `json:"id"`
					Name      string      `json:"name"`
					Value     interface{} `json:"value"`
					ValueType string      `json:"value_type"`
				} `json:"metadata"`
				ID            int    `json:"id"`
				UpdatedAt     string `json:"updated_at"`
				RequisitionID string `json:"requisition_id"`
				Title         string `json:"title"`
			} `json:"jobs"`
			Meta struct {
				Total int `json:"total"`
			} `json:"meta"`
		}

		var jsonJobs JsonJobs

		c.OnResponse(func(r *colly.Response) {
			var tempJson JsonJobs
			err := json.Unmarshal(r.Body, &tempJson)
			if err != nil {
				panic(err.Error())
			}

			for _, elem := range tempJson.Jobs {

				result_title := elem.Title
				result_url := elem.AbsoluteURL

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

			jsonJobs.Jobs = append(jsonJobs.Jobs, tempJson.Jobs...)
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
			c.Visit(start_url)
		}
	}
	return
}

func (runtime Runtime) Automationhero(
	version int, isLocal bool) (
	response Response, results []Result) {
	switch version {
	case 1:

		c := colly.NewCollector()

		c_start_url := "https://api.lever.co/v0/postings/automationhero?mode=json"

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

func (runtime Runtime) Bonify(
	version int, isLocal bool) (
	response Response, results []Result) {
	switch version {
	case 1:

		c := colly.NewCollector()

		start_url := "http://www.bonify.de/jobs"
		job_base_url := "https://www.bonify.de/jobs/%s"

		type JsonJobs struct {
			Page             int         `json:"page"`
			ResultsPerPage   int         `json:"results_per_page"`
			ResultsSize      int         `json:"results_size"`
			TotalResultsSize int         `json:"total_results_size"`
			TotalPages       int         `json:"total_pages"`
			NextPage         interface{} `json:"next_page"`
			PrevPage         interface{} `json:"prev_page"`
			Results          []struct {
				ID                   string        `json:"id"`
				UID                  string        `json:"uid"`
				Type                 string        `json:"type"`
				Href                 string        `json:"href"`
				Tags                 []interface{} `json:"tags"`
				FirstPublicationDate string        `json:"first_publication_date"`
				LastPublicationDate  string        `json:"last_publication_date"`
				Slugs                []string      `json:"slugs"`
				LinkedDocuments      []interface{} `json:"linked_documents"`
				Lang                 string        `json:"lang"`
				AlternateLanguages   []interface{} `json:"alternate_languages"`
				Data                 struct {
					Title []struct {
						Type  string        `json:"type"`
						Text  string        `json:"text"`
						Spans []interface{} `json:"spans"`
					} `json:"title"`
					PersonioJobID        string `json:"personio_job_id"`
					JobType              string `json:"job_type"`
					Index                string `json:"index"`
					FocusKeyPhrase       string `json:"focus_key_phrase"`
					BreadcrumbVisibility string `json:"breadcrumb_visibility"`
					Department           string `json:"department"`
				} `json:"data"`
			} `json:"results"`
			Version string `json:"version"`
			License string `json:"license"`
		}

		var jsonJobs JsonJobs

		c.OnResponse(func(r *colly.Response) {
			body := string(r.Body)
			json_body := strings.Split(
				strings.Split(
					body, `resultsAllJobsListingsTrimmed":`)[1], `,"resultsCompanyBenefits`)[0]

			var tempJson JsonJobs
			err := json.Unmarshal([]byte(json_body), &tempJson)
			if err != nil {
				panic(err.Error())
			}

			for _, elem := range tempJson.Results {

				result_title := elem.Data.Title[0].Text
				result_url := fmt.Sprintf(job_base_url, elem.UID)

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

			jsonJobs.Results = append(jsonJobs.Results, tempJson.Results...)
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
			c.Visit(start_url)
		}
	}
	return
}

func (runtime Runtime) Bryter(
	version int, isLocal bool) (
	response Response, results []Result) {
	switch version {
	case 1:

		c := colly.NewCollector()

		url := "https://bryter.io/careers"
		tag_main_section := "#careers-listing"
		tag_title := "h4"

		type Job struct {
			Title string
			Url   string
		}

		c.OnHTML(tag_main_section, func(e *colly.HTMLElement) {
			e.ForEach("a", func(_ int, el *colly.HTMLElement) {
				result_url := el.Attr("href")
				result_title := el.ChildText(tag_title)

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

func (runtime Runtime) Bunch(
	version int, isLocal bool) (
	response Response, results []Result) {
	switch version {
	case 1:

		results = append(results, Result{
			runtime.Name,
			"Freelance/Full-time Product Designer",
			"https://angel.co/company/bunch-hq/jobs/682927-freelance-full-time-product-designer",
			[]byte("{}"),
		})

		results = append(results, Result{
			runtime.Name,
			"Product Launch Intern (Internship)",
			"https://angel.co/company/bunch-hq/jobs/907192-product-launch-intern-internship",
			[]byte("{}"),
		})

		results_marshal, err := json.Marshal(results)
		if err != nil {
			panic(err.Error())
		}
		response = Response{[]byte(results_marshal)}
	}
	return
}

func (runtime Runtime) Candis(
	version int, isLocal bool) (
	response Response, results []Result) {
	switch version {
	case 1:

		c := colly.NewCollector()

		p_start_url := "https://career.recruitee.com/api/c/50731/widget/?widget=true"

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
				OpenQuestions []interface{} `json:"open_questions"`
			} `json:"offers"`
			Terms []interface{} `json:"terms"`
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
	}
	return
}

func (runtime Runtime) Cargoone(
	version int, isLocal bool) (
	response Response, results []Result) {
	switch version {
	case 1:

		c := colly.NewCollector()

		c_start_url := "https://api.lever.co/v0/postings/cargo-2?mode=json"

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

func (runtime Runtime) Construyo(
	version int, isLocal bool) (
	response Response, results []Result) {
	switch version {
	case 1:

		c := colly.NewCollector()

		url := "https://partum-gmbh-jobs.personio.de"
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

func (runtime Runtime) Crosslend(
	version int, isLocal bool) (
	response Response, results []Result) {
	switch version {
	case 1:

		c := colly.NewCollector()

		url := "https://www.crosslend.com/home/careers"

		type Job struct {
			Title string
			Url   string
		}

		c.OnHTML(".tab-content", func(e *colly.HTMLElement) {
			e.ForEach("p", func(_ int, el *colly.HTMLElement) {
				result_title := el.ChildText("a")
				result_url := el.ChildAttr("a", "href")

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

func (runtime Runtime) Bytedance(
	version int, isLocal bool) (response Response, results []Result) {
	switch version {
	case 1:

		base_url := "https://job.bytedance.com/en/position/detail/%s"

		type Jobs struct {
			Code int `json:"code"`
			Data struct {
				JobPostList []struct {
					ID          string `json:"id"`
					Title       string `json:"title"`
					SubTitle    string `json:"sub_title"`
					Description string `json:"description"`
					Requirement string `json:"requirement"`
					JobCategory struct {
						ID       string `json:"id"`
						Name     string `json:"name"`
						EnName   string `json:"en_name"`
						I18NName string `json:"i18n_name"`
						Depth    int    `json:"depth"`
						Parent   struct {
							ID       string      `json:"id"`
							Name     string      `json:"name"`
							EnName   string      `json:"en_name"`
							I18NName string      `json:"i18n_name"`
							Depth    int         `json:"depth"`
							Parent   interface{} `json:"parent"`
							Children interface{} `json:"children"`
						} `json:"parent"`
						Children interface{} `json:"children"`
					} `json:"job_category"`
					CityInfo struct {
						Code         string      `json:"code"`
						Name         string      `json:"name"`
						EnName       string      `json:"en_name"`
						LocationType interface{} `json:"location_type"`
						I18NName     string      `json:"i18n_name"`
						PyName       interface{} `json:"py_name"`
					} `json:"city_info"`
					RecruitType struct {
						ID       string `json:"id"`
						Name     string `json:"name"`
						EnName   string `json:"en_name"`
						I18NName string `json:"i18n_name"`
						Depth    int    `json:"depth"`
						Parent   struct {
							ID       string      `json:"id"`
							Name     string      `json:"name"`
							EnName   string      `json:"en_name"`
							I18NName string      `json:"i18n_name"`
							Depth    int         `json:"depth"`
							Parent   interface{} `json:"parent"`
							Children interface{} `json:"children"`
						} `json:"parent"`
						Children interface{} `json:"children"`
					} `json:"recruit_type"`
					PublishTime int64       `json:"publish_time"`
					JobHotFlag  int         `json:"job_hot_flag"`
					JobSubject  interface{} `json:"job_subject"`
				} `json:"job_post_list"`
				Count int    `json:"count"`
				Extra string `json:"extra"`
			} `json:"data"`
			Message string      `json:"message"`
			Error   interface{} `json:"error"`
		}

		ctx, cancel := chromedp.NewContext(context.Background())
		defer cancel()

		var res string
		if err := chromedp.Run(ctx,
			chromedp.Navigate("https://job.bytedance.com/en/position?limit=10"),
			chromedp.Sleep(5*time.Second),
			chromedp.EvaluateAsDevTools(`document.cookie`, &res),
		); err != nil {
			panic(err.Error())
		}

		token := strings.Split(res, "atsx-csrf-token=")[1]

		client := &http.Client{}
		data := strings.NewReader(`{"limit":1000}`)
		req, err := http.NewRequest("POST", "https://job.bytedance.com/api/v1/search/job/posts", data)
		if err != nil {
			panic(err.Error())
		}
		req.Header.Set("x-csrf-token", strings.ReplaceAll(token, "%3D", "="))
		req.Header.Set("Cookie", "channel=overseas; atsx-csrf-token="+token)
		resp, err := client.Do(req)
		if err != nil {
			panic(err.Error())
		}
		bodyText, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err.Error())
		}

		var tempJson Jobs
		err = json.Unmarshal(bodyText, &tempJson)
		if err != nil {
			panic(err.Error())
		}

		for _, elem := range tempJson.Data.JobPostList {
			result_title := elem.Title
			result_url := fmt.Sprintf(base_url, elem.ID)

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
		response = Response{bodyText}
	}
	return
}

func (runtime Runtime) Bmw(
	version int, isLocal bool) (
	response Response, results []Result) {
	switch version {
	case 1:

		c := colly.NewCollector()

		start_url := "https://www.bmwgroup.jobs/content/grpw/websites/jobfinder.joblist.en.de.json"
		base_job_url := "https://www.bmwgroup.jobs/en/jobfinder/job-description.%s"

		type JsonJobs struct {
			Data []struct {
				PostingDate string `json:"postingDate"`
				Favorite    bool   `json:"favorite"`
				RefNo       string `json:"refNo"`
				ReqTitle    string `json:"reqTitle"`
				JobType     struct {
					Value   string `json:"value"`
					Display string `json:"display"`
				} `json:"jobType"`
				LegalEntity struct {
					Value   string `json:"value"`
					Display string `json:"display"`
				} `json:"legalEntity"`
				JobField struct {
					Value   string `json:"value"`
					Display string `json:"display"`
				} `json:"jobField"`
				Location struct {
					Value   string `json:"value"`
					Display string `json:"display"`
				} `json:"location"`
				JobDescriptionLink string `json:"jobDescriptionLink"`
				JobLevel           string `json:"jobLevel"`
				EmployeeStatus     string `json:"employeeStatus"`
				Schedule           string `json:"schedule"`
				HotJob             bool   `json:"hotJob"`
				Fulltext           string `json:"fulltext"`
			} `json:"data"`
		}

		var jsonJobs JsonJobs

		c.OnResponse(func(r *colly.Response) {
			var tempJson JsonJobs
			err := json.Unmarshal(r.Body, &tempJson)
			if err != nil {
				panic(err.Error())
			}

			for _, elem := range tempJson.Data {

				result_title := elem.ReqTitle
				result_url := fmt.Sprintf(base_job_url, elem.JobDescriptionLink)

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
			c.Visit(start_url)
		}
	}
	return
}

func (runtime Runtime) Infineon(
	version int, isLocal bool) (response Response, results []Result) {
	switch version {
	case 1:

		base_url := "https://www.infineon.com%s"

		type Jobs struct {
			Pages struct {
				Count int `json:"count"`
				Items []struct {
					PublicationLanguageDe bool     `json:"publication_language_de"`
					LocationEn            string   `json:"location_en"`
					ID                    string   `json:"id"`
					PublicationLanguageEn bool     `json:"publication_language_en"`
					CreationDate          string   `json:"creation_date"`
					FieldsOfStudy         []string `json:"fields_of_study,omitempty"`
					FunctionalArea        string   `json:"functional_area"`
					Location              string   `json:"location"`
					Country               string   `json:"country"`
					EntryLevel            string   `json:"entry_level"`
					Division              string   `json:"division"`
					DesiredStartDate      string   `json:"desired_start_date"`
					DetailPageURL         string   `json:"detail_page_url"`
					JobAttributes         []string `json:"job_attributes"`
					Role                  string   `json:"role"`
					Title                 string   `json:"title"`
					Description           string   `json:"description"`
					DetailDataURL         string   `json:"detail_data_url"`
					Icons                 []struct {
						Type string `json:"type"`
						Text string `json:"text"`
					} `json:"icons,omitempty"`
					Tags []string `json:"tags,omitempty"`
				} `json:"items"`
			} `json:"pages"`
			Offset     int `json:"offset"`
			HasResults int `json:"has_results"`
			Count      int `json:"count"`
		}

		client := &http.Client{}
		data := strings.NewReader(`term=&offset=0&max_results=1000&lang=en`)
		req, err := http.NewRequest("POST", "https://www.infineon.com/search/jobs/jobs", data)
		if err != nil {
			panic(err.Error())
		}
		req.Header.Set("content-type", "application/x-www-form-urlencoded; charset=UTF-8")
		resp, err := client.Do(req)
		if err != nil {
			panic(err.Error())
		}
		bodyText, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err.Error())
		}

		var tempJson Jobs
		err = json.Unmarshal(bodyText, &tempJson)
		if err != nil {
			panic(err.Error())
		}

		for _, elem := range tempJson.Pages.Items {
			result_title := elem.Title
			result_url := fmt.Sprintf(base_url, elem.DetailPageURL)

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
		response = Response{bodyText}
	}
	return
}

func (runtime Runtime) Porsche(
	version int, isLocal bool) (
	response Response, results []Result) {
	switch version {
	case 1:

		c := colly.NewCollector()

		start_url := `https://api-jobs.porsche.com/search/?data={"SearchParameters":{"FirstItem":1,"CountItem":1000}}`

		type JsonJobs struct {
			LanguageCode string `json:"LanguageCode"`
			SearchResult struct {
				SearchResultCount    int `json:"SearchResultCount"`
				SearchResultCountAll int `json:"SearchResultCountAll"`
				SearchResultItems    []struct {
					MatchedObjectID         string `json:"MatchedObjectId"`
					MatchedObjectDescriptor struct {
						ID                  string   `json:"ID"`
						PositionID          string   `json:"PositionID"`
						PositionTitle       string   `json:"PositionTitle"`
						PublicationCode     string   `json:"PublicationCode"`
						PositionURI         string   `json:"PositionURI"`
						ApplyURI            []string `json:"ApplyURI"`
						PublicationLanguage struct {
							Code string `json:"Code"`
						} `json:"PublicationLanguage"`
						PublicationChannel []struct {
							ID        int    `json:"Id"`
							StartDate string `json:"StartDate"`
							EndDate   string `json:"EndDate"`
						} `json:"PublicationChannel"`
						PublicationEndDate string `json:"PublicationEndDate"`
						PositionIndustry   []struct {
							Code string `json:"Code"`
							Name string `json:"Name"`
						} `json:"PositionIndustry"`
						JobCategory []struct {
							Code string `json:"Code"`
							Name string `json:"Name"`
						} `json:"JobCategory"`
						CareerLevel []struct {
							Code string `json:"Code"`
							Name string `json:"Name"`
						} `json:"CareerLevel"`
						TargetGroup      []interface{} `json:"TargetGroup"`
						PositionSchedule []struct {
							Code string `json:"Code"`
							Name string `json:"Name"`
						} `json:"PositionSchedule"`
						PositionOfferingType []struct {
							Code string `json:"Code"`
							Name string `json:"Name"`
						} `json:"PositionOfferingType"`
						ParentOrganization     string `json:"ParentOrganization"`
						ParentOrganizationName string `json:"ParentOrganizationName"`
						PositionLocation       []struct {
							Continent              string `json:"Continent"`
							ContinentName          string `json:"ContinentName"`
							Country                string `json:"Country"`
							CountryName            string `json:"CountryName"`
							CountryCode            string `json:"CountryCode"`
							CountrySubDivision     string `json:"CountrySubDivision"`
							CountrySubDivisionName string `json:"CountrySubDivisionName"`
							City                   string `json:"City"`
							CityName               string `json:"CityName"`
						} `json:"PositionLocation"`
						Organization         string   `json:"Organization"`
						OrganizationName     string   `json:"OrganizationName"`
						LogoURI              []string `json:"LogoURI"`
						PublicationStartDate string   `json:"PublicationStartDate"`
					} `json:"MatchedObjectDescriptor"`
					RelevanceScore int `json:"RelevanceScore"`
					RelevanceRank  int `json:"RelevanceRank"`
				} `json:"SearchResultItems"`
				UserArea struct {
					ExecutionError int `json:"ExecutionError"`
				} `json:"UserArea"`
			} `json:"SearchResult"`
		}

		var jsonJobs JsonJobs

		c.OnResponse(func(r *colly.Response) {
			var tempJson JsonJobs
			err := json.Unmarshal(r.Body, &tempJson)
			if err != nil {
				panic(err.Error())
			}

			for _, elem := range tempJson.SearchResult.SearchResultItems {

				result_title := elem.MatchedObjectDescriptor.PositionTitle
				result_url := elem.MatchedObjectDescriptor.PositionURI

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

			jsonJobs.SearchResult.SearchResultItems = append(
				jsonJobs.SearchResult.SearchResultItems,
				tempJson.SearchResult.SearchResultItems...)
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
			c.Visit(start_url)
		}
	}
	return
}

func (runtime Runtime) Bosch(
	version int, isLocal bool) (response Response, results []Result) {
	switch version {
	case 1:

		c := colly.NewCollector()

		start_url := "https://api.smartrecruiters.com/v1/companies/BoschGroup/postings?offset=%d"
		base_job_url := "https://www.smartrecruiters.com/BoschGroup/%s"
		number_results_per_page := 100

		type Jobs struct {
			Offset     int `json:"offset"`
			Limit      int `json:"limit"`
			TotalFound int `json:"totalFound"`
			Content    []struct {
				ID        string `json:"id"`
				Name      string `json:"name"`
				UUID      string `json:"uuid"`
				RefNumber string `json:"refNumber"`
				Company   struct {
					Identifier string `json:"identifier"`
					Name       string `json:"name"`
				} `json:"company"`
				ReleasedDate time.Time `json:"releasedDate"`
			} `json:"content"`
		}

		var jsonJobs Jobs

		c.OnResponse(func(r *colly.Response) {
			var tempJson Jobs
			err := json.Unmarshal(r.Body, &tempJson)
			if err != nil {
				panic(err.Error())
			}

			for _, elem := range tempJson.Content {

				result_title := elem.Name
				result_url := fmt.Sprintf(base_job_url, elem.ID)

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

			jsonJobs.Content = append(jsonJobs.Content, tempJson.Content...)

			if isLocal {
				return
			} else {
				total_matches := tempJson.TotalFound
				total_pages := total_matches / number_results_per_page
				for i := 1; i <= total_pages; i++ {
					time.Sleep(SecondsSleep * time.Second)
					c.Visit(fmt.Sprintf(start_url, number_results_per_page*i))
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
			c.Visit(fmt.Sprintf(start_url, 0))
		}
	}
	return
}

func (runtime Runtime) Mckinsey(
	version int, isLocal bool) (response Response, results []Result) {
	switch version {
	case 1:

		c := colly.NewCollector()

		start_url := "https://mobileservices.mckinsey.com/services/ContentAPI/SearchAPI.svc/jobs/search?&pageSize=100&start=%d"
		base_url := "https://www.mckinsey.com/careers/search-jobs/jobs/%s"
		number_results_per_page := 100
		counter := 0

		type JsonJobs struct {
			Response struct {
				NumFound int `json:"numFound"`
				Start    int `json:"start"`
				Docs     []struct {
					JobID                  string   `json:"jobID"`
					Title                  string   `json:"title"`
					RecordTypeName         []string `json:"recordTypeName"`
					JobSkillGroup          []string `json:"jobSkillGroup"`
					JobSkillCode           []string `json:"jobSkillCode"`
					Interest               string   `json:"interest"`
					InterestCategory       string   `json:"interestCategory"`
					Cities                 []string `json:"cities"`
					Countries              []string `json:"countries"`
					Continents             []string `json:"continents"`
					Functions              []string `json:"functions,omitempty"`
					Industries             []string `json:"industries,omitempty"`
					WhoYouWillWorkWith     string   `json:"whoYouWillWorkWith"`
					WhatYouWillDo          string   `json:"whatYouWillDo"`
					YourBackground         string   `json:"yourBackground"`
					LinkedInSeniorityLevel []string `json:"linkedInSeniorityLevel,omitempty"`
					JobApplyURL            string   `json:"jobApplyURL"`
					FriendlyURL            string   `json:"friendlyURL"`
					ShortJobSummary        string   `json:"shortJobSummary,omitempty"`
				} `json:"docs"`
			} `json:"response"`
		}

		var jsonJobs JsonJobs

		c.OnResponse(func(r *colly.Response) {
			var tempJsonJobs JsonJobs
			err := json.Unmarshal(r.Body, &tempJsonJobs)
			if err != nil {
				panic(err.Error())
			}

			for _, elem := range tempJsonJobs.Response.Docs {

				result_title := elem.Title
				result_url := fmt.Sprintf(base_url, elem.FriendlyURL)

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

			jsonJobs.Response.Docs = append(
				jsonJobs.Response.Docs,
				tempJsonJobs.Response.Docs...)

			total_pages := tempJsonJobs.Response.NumFound / number_results_per_page

			if counter >= total_pages {
				return
			} else {
				counter++
				time.Sleep(SecondsSleep * time.Second)
				c.Visit(fmt.Sprintf(start_url, counter*number_results_per_page))
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
			c.Visit(fmt.Sprintf(start_url, 0))
		}
	}
	return
}

func (runtime Runtime) Sap(
	version int, isLocal bool) (
	response Response, results []Result) {
	switch version {
	case 1:

		c := colly.NewCollector()

		start_url := "https://jobs.sap.com/search/?q=&sortColumn=referencedate&sortDirection=desc&startrow=%d"
		base_url := "https://jobs.sap.com/%s"
		number_results_per_page := 25
		counter := 0

		type Job struct {
			Title    string
			Url      string
			Location string
		}

		c.OnHTML(".html5", func(e *colly.HTMLElement) {
			e.ForEach(".data-row", func(_ int, el *colly.HTMLElement) {
				result_url := fmt.Sprintf(base_url, el.ChildAttr("a", "href"))
				result_title := el.ChildText("a")
				result_location := el.ChildText(".jobLocation.visible-phone")

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
			})

			temp_pages := strings.Split(e.ChildText(".srHelp"), " ")
			s_temp_pages := temp_pages[len(temp_pages)-1]
			total_pages, err := strconv.Atoi(s_temp_pages)
			if err != nil {
				panic(err.Error())
			}

			if counter > total_pages {
				return
			} else {
				counter++
				time.Sleep(SecondsSleep * time.Second)
				c.Visit(fmt.Sprintf(start_url, counter*number_results_per_page))
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
			c.Visit(fmt.Sprintf(start_url, 0))
		}
	}
	return
}

func (runtime Runtime) Puma(
	version int, isLocal bool) (
	response Response, results []Result) {
	switch version {
	case 1:

		c := colly.NewCollector()

		start_url := "https://about.puma.com/api/PUMA/Feature/JobFinder?loadMore=500"
		base_url := "https://about.puma.com%s"

		type JsonJobs struct {
			NumberFound string `json:"numberFound"`
			LoadMoreURL string `json:"loadMoreUrl"`
			Teaser      []struct {
				Jobitemid  string      `json:"jobitemid"`
				URL        string      `json:"url"`
				Title      string      `json:"title"`
				Team       string      `json:"team"`
				Location   string      `json:"location"`
				LocationID interface{} `json:"locationId"`
			} `json:"teaser"`
		}

		var jsonJobs JsonJobs

		c.OnResponse(func(r *colly.Response) {
			var tempJson JsonJobs
			err := json.Unmarshal(r.Body, &tempJson)
			if err != nil {
				panic(err.Error())
			}

			for _, elem := range tempJson.Teaser {

				result_title := elem.Title
				result_url := fmt.Sprintf(base_url, elem.URL)

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

			jsonJobs.Teaser = append(jsonJobs.Teaser, tempJson.Teaser...)
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
			c.Visit(start_url)
		}
	}
	return
}

func (runtime Runtime) Daimler(
	version int, isLocal bool) (
	response Response, results []Result) {
	switch version {
	case 1:

		c := colly.NewCollector()

		start_url := `https://global-jobboard-api.daimler.com/v3/search/{"SearchParameters":{"MatchedObjectDescriptor":["PositionID","PositionTitle","PositionURI","OrganizationName","PositionLocation.CityName","JobCategory.Name","CareerLevel.Name","Facet:PositionLocation.CityName","Facet:PositionLocation.CountryName","PublicationStartDate"],"FirstItem":0,"CountItem":1000000},"SearchCriteria":[{"CriterionName":"PublicationLanguage.Code","CriterionValue":["EN"]}]}`

		type JsonJobs struct {
			SearchResult struct {
				SearchResultCount    int `json:"SearchResultCount"`
				SearchResultCountAll int `json:"SearchResultCountAll"`
				SearchResultItems    []struct {
					MatchedObjectID         string `json:"MatchedObjectId"`
					MatchedObjectDescriptor struct {
						PublicationStartDate string `json:"PublicationStartDate"`
						PositionTitle        string `json:"PositionTitle"`
						PositionURI          string `json:"PositionURI"`
						PositionLocation     []struct {
							CityName string `json:"CityName"`
						} `json:"PositionLocation"`
						OrganizationName string `json:"OrganizationName"`
						JobCategory      []struct {
							Name string `json:"Name"`
						} `json:"JobCategory"`
						CareerLevel []struct {
							Name string `json:"Name"`
						} `json:"CareerLevel"`
						PositionID string `json:"PositionID"`
					} `json:"MatchedObjectDescriptor"`
				} `json:"SearchResultItems"`
			} `json:"SearchResult"`
		}

		var jsonJobs JsonJobs

		c.OnResponse(func(r *colly.Response) {
			var tempJson JsonJobs
			err := json.Unmarshal(r.Body, &tempJson)
			if err != nil {
				panic(err.Error())
			}

			for _, elem := range tempJson.SearchResult.SearchResultItems {

				result_title := elem.MatchedObjectDescriptor.PositionTitle
				result_url := elem.MatchedObjectDescriptor.PositionURI

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

			jsonJobs.SearchResult.SearchResultItems = append(
				jsonJobs.SearchResult.SearchResultItems,
				tempJson.SearchResult.SearchResultItems...)
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
			c.Visit(start_url)
		}
	}
	return
}

func (runtime Runtime) Siemens(
	version int, isLocal bool) (response Response, results []Result) {
	switch version {
	case 1:

		c := colly.NewCollector()

		start_url := "https://jobs.siemens.com/api/jobs?page=%d&limit=100"
		number_results_per_page := 100
		counter := 1

		type JsonJobs struct {
			Jobs []struct {
				Data struct {
					Slug         string   `json:"slug"`
					Language     string   `json:"language"`
					Languages    []string `json:"languages"`
					Title        string   `json:"title"`
					Description  string   `json:"description"`
					City         string   `json:"city"`
					State        string   `json:"state"`
					Country      string   `json:"country"`
					CountryCode  string   `json:"country_code"`
					PostalCode   string   `json:"postal_code"`
					LocationType string   `json:"location_type"`
					Latitude     float64  `json:"latitude"`
					Longitude    float64  `json:"longitude"`
					Categories   []struct {
						Name string `json:"name"`
					} `json:"categories"`
					Tags1            []string `json:"tags1"`
					Brand            string   `json:"brand"`
					PromotionValue   int      `json:"promotion_value"`
					ExperienceLevels []string `json:"experience_levels"`
					Source           string   `json:"source"`
					PostedDate       string   `json:"posted_date"`
					Internal         bool     `json:"internal"`
					Searchable       bool     `json:"searchable"`
					Applyable        bool     `json:"applyable"`
					LiEasyApplyable  bool     `json:"li_easy_applyable"`
					AtsCode          string   `json:"ats_code"`
					MetaData         struct {
						CanonicalURL string `json:"canonical_url"`
					} `json:"meta_data"`
					UpdateDate    string   `json:"update_date"`
					CreateDate    string   `json:"create_date"`
					Category      []string `json:"category"`
					FullLocation  string   `json:"full_location"`
					ShortLocation string   `json:"short_location"`
				} `json:"data"`
			} `json:"jobs"`
			TotalCount int `json:"totalCount"`
			Count      int `json:"count"`
		}

		var jsonJobs JsonJobs

		c.OnResponse(func(r *colly.Response) {
			var tempJsonJobs JsonJobs
			err := json.Unmarshal(r.Body, &tempJsonJobs)
			if err != nil {
				panic(err.Error())
			}

			for _, elem := range tempJsonJobs.Jobs {

				result_title := elem.Data.Title
				result_url := elem.Data.MetaData.CanonicalURL

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

			total_pages := tempJsonJobs.TotalCount / number_results_per_page

			if counter > total_pages {
				return
			} else {
				counter++
				time.Sleep(SecondsSleep * time.Second)
				c.Visit(fmt.Sprintf(start_url, counter))
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
			c.Visit(fmt.Sprintf(start_url, 1))
		}
	}
	return
}

func (runtime Runtime) Continental(
	version int, isLocal bool) (
	response Response, results []Result) {
	switch version {
	case 1:

		c := colly.NewCollector()

		start_url := `https://api.continental-jobs.com/search/?data={"LanguageCode":"EN","SearchParameters":{"FirstItem":1,"CountItem":2000,"Sort":[{"Criterion":"PublicationStartDate","Direction":"DESC"}],"MatchedObjectDescriptor":["ID","PositionID","PositionTitle","PositionURI","PositionLocation.CountryName","PositionLocation.CityName","PositionLocation.Longitude","PositionLocation.Latitude","PositionIndustry.Name","JobCategory.Name","PublicationStartDate","VacancyDivision"]},"SearchCriteria":[{"CriterionName":"PublicationLanguage.Code","CriterionValue":["EN"]},{"CriterionName":"PublicationChannel.Code","CriterionValue":["12"]}]}`

		type JsonJobs struct {
			SearchResult struct {
				SearchResultCount    int `json:"SearchResultCount"`
				SearchResultCountAll int `json:"SearchResultCountAll"`
				SearchResultItems    []struct {
					MatchedObjectID         string `json:"MatchedObjectId"`
					MatchedObjectDescriptor struct {
						PositionIndustry struct {
							Name string `json:"Name"`
						} `json:"PositionIndustry"`
						PublicationStartDate string `json:"PublicationStartDate"`
						PositionTitle        string `json:"PositionTitle"`
						PositionLocation     []struct {
							CityName    string  `json:"CityName"`
							Longitude   float64 `json:"Longitude"`
							Latitude    float64 `json:"Latitude"`
							CountryName string  `json:"CountryName"`
						} `json:"PositionLocation"`
						PositionURI string `json:"PositionURI"`
						ID          int    `json:"ID"`
						JobCategory struct {
							Name string `json:"Name"`
						} `json:"JobCategory"`
						PositionID string `json:"PositionID"`
					} `json:"MatchedObjectDescriptor,omitempty"`
					RelevanceScore int `json:"RelevanceScore"`
					RelevanceRank  int `json:"RelevanceRank"`
				} `json:"SearchResultItems"`
			} `json:"SearchResult"`
		}

		var jsonJobs JsonJobs

		c.OnResponse(func(r *colly.Response) {
			var tempJson JsonJobs
			err := json.Unmarshal(r.Body, &tempJson)
			if err != nil {
				panic(err.Error())
			}

			for _, elem := range tempJson.SearchResult.SearchResultItems {

				result_title := elem.MatchedObjectDescriptor.PositionTitle
				result_url := elem.MatchedObjectDescriptor.PositionURI

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

			jsonJobs.SearchResult.SearchResultItems = append(
				jsonJobs.SearchResult.SearchResultItems,
				tempJson.SearchResult.SearchResultItems...)
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
			c.Visit(start_url)
		}
	}
	return
}

func (runtime Runtime) Deliveryhero(
	version int, isLocal bool) (
	response Response, results []Result) {
	switch version {
	case 1:

		c := colly.NewCollector()

		start_url := "https://careers.deliveryhero.com/global/en/search-results?s=1&from=%d"
		base_job_url := "https://careers.deliveryhero.com/global/en/job/%s"
		number_results_per_page := 50
		counter := 0

		type JsonJobs struct {
			EagerLoadRefineSearch struct {
				Status    int `json:"status"`
				Hits      int `json:"hits"`
				TotalHits int `json:"totalHits"`
				Data      struct {
					Jobs []struct {
						Country            string   `json:"country"`
						CityState          string   `json:"cityState"`
						SubCategory        string   `json:"subCategory"`
						City               string   `json:"city"`
						MlSkills           []string `json:"ml_skills"`
						PostalCode         string   `json:"postalCode"`
						Industry           string   `json:"industry"`
						Type               string   `json:"type"`
						MultiLocation      []string `json:"multi_location"`
						Locale             string   `json:"locale"`
						Title              string   `json:"title"`
						MultiLocationArray []struct {
							Location string `json:"location"`
						} `json:"multi_location_array"`
						JobSeqNo           string    `json:"jobSeqNo"`
						PostedDate         time.Time `json:"postedDate"`
						DescriptionTeaser  string    `json:"descriptionTeaser"`
						DateCreated        time.Time `json:"dateCreated"`
						State              string    `json:"state"`
						CityStateCountry   string    `json:"cityStateCountry"`
						Brand              string    `json:"brand"`
						VisibilityType     string    `json:"visibilityType"`
						SiteType           string    `json:"siteType"`
						Address            string    `json:"address"`
						IsMultiCategory    bool      `json:"isMultiCategory"`
						MultiCategory      []string  `json:"multi_category"`
						ReqID              string    `json:"reqId"`
						JobID              string    `json:"jobId"`
						Badge              string    `json:"badge"`
						JobVisibility      []string  `json:"jobVisibility"`
						IsMultiLocation    bool      `json:"isMultiLocation"`
						ApplyURL           string    `json:"applyUrl"`
						MultiCategoryArray []struct {
							Category string `json:"category"`
						} `json:"multi_category_array"`
						Location        string      `json:"location"`
						Category        string      `json:"category"`
						ExternalApply   bool        `json:"externalApply"`
						LocationLatlong interface{} `json:"locationLatlong"`
					} `json:"jobs"`
				} `json:"data"`
			} `json:"eagerLoadRefineSearch"`
		}

		var jsonJobs JsonJobs

		c.OnResponse(func(r *colly.Response) {
			response = Response{r.Body}
			response_body := string(response.Html)
			response_json := strings.Split(
				strings.Split(
					response_body, "phApp.ddo = ")[1], "; phApp.experimentData")[0]

			var tempJson JsonJobs
			err := json.Unmarshal([]byte(response_json), &tempJson)
			if err != nil {
				panic(err.Error())
			}

			for _, elem := range tempJson.EagerLoadRefineSearch.Data.Jobs {

				result_title := elem.Title
				result_url := fmt.Sprintf(base_job_url, elem.JobID)

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

			jsonJobs.EagerLoadRefineSearch.Data.Jobs = append(
				jsonJobs.EagerLoadRefineSearch.Data.Jobs,
				tempJson.EagerLoadRefineSearch.Data.Jobs...)

			total_pages := tempJson.EagerLoadRefineSearch.TotalHits / number_results_per_page

			if counter > total_pages {
				return
			} else {
				counter++
				time.Sleep(SecondsSleep * time.Second)
				c.Visit(fmt.Sprintf(start_url, counter*number_results_per_page))
			}
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
			c.Visit(fmt.Sprintf(start_url, 0))
		}
	}
	return
}

func (runtime Runtime) Volkswagen(
	version int, isLocal bool) (
	response Response, results []Result) {
	switch version {
	case 1:

		c := colly.NewCollector()

		start_url := "https://karriere.volkswagen.de/sap/opu/odata/sap/zaudi_ui_open_srv/JobSet?sap-client=100"
		job_base_url := "https://karriere.volkswagen.de%s"

		type Jobs struct {
			Feed struct {
				Entry []struct {
					Properties struct {
						ContentType           string `json:"ContentType"`
						FilterContractTypes   string `json:"FilterContractTypes"`
						FilterFunctionalAreas string `json:"FilterFunctionalAreas"`
						ZLanguage             string `json:"ZLanguage"`
						ZRecruiterEmail       string `json:"ZRecruiterEmail"`
						Address               struct {
							Type   string `json:"-type"`
							Line07 string `json:"Line07"`
							Line00 string `json:"Line00"`
							Line01 string `json:"Line01"`
							Line02 string `json:"Line02"`
							Line03 string `json:"Line03"`
							Line04 string `json:"Line04"`
							Line05 string `json:"Line05"`
							Line06 string `json:"Line06"`
							Line08 string `json:"Line08"`
							Line09 string `json:"Line09"`
						} `json:"Address"`
						ApplicationStatusTxt string `json:"ApplicationStatusTxt"`
						Foot1                string `json:"Foot1"`
						ZEmailApplication    string `json:"ZEmailApplication"`
						FilterInterestGroups string `json:"FilterInterestGroups"`
						ZRecruiterPhone      string `json:"ZRecruiterPhone"`
						ZZCompanyWide        string `json:"ZZCompanyWide"`
						ZPublicationDate     string `json:"ZPublicationDate"`
						AdditionalInfo       string `json:"AdditionalInfo"`
						FilterCompanies      string `json:"FilterCompanies"`
						ZIsForeignCompany    string `json:"ZIsForeignCompany"`
						Title                string `json:"Title"`
						PostingAge           string `json:"PostingAge"`
						ProjectDesc          string `json:"ProjectDesc"`
						RefCode              string `json:"RefCode"`
						TravelRatio          string `json:"TravelRatio"`
						IsHotJob             string `json:"IsHotJob"`
						DepartmentDesc       string `json:"DepartmentDesc"`
						TaskDesc             string `json:"TaskDesc"`
						D                    string `json:"-d"`
						Posting              string `json:"Posting"`
						ApplicationExists    string `json:"ApplicationExists"`
						Department           string `json:"Department"`
						Contact              string `json:"Contact"`
						ZEmploymentStartDate struct {
							Null string `json:"-null"`
						} `json:"ZEmploymentStartDate"`
						IsApplicationGroup    string `json:"IsApplicationGroup"`
						Foot2                 string `json:"Foot2"`
						EEOText               string `json:"EEOText"`
						FilterLocations       string `json:"FilterLocations"`
						M                     string `json:"-m"`
						ApplicationURL        string `json:"ApplicationUrl"`
						ReportingTo           string `json:"ReportingTo"`
						JobID                 string `json:"JobID"`
						FilterHierarchyLevels string `json:"FilterHierarchyLevels"`
						ZEmploymentEndDate    struct {
							Null string `json:"-null"`
						} `json:"ZEmploymentEndDate"`
						ZRecruiterFullName string `json:"ZRecruiterFullName"`
						JobDetailsURL      string `json:"JobDetailsUrl"`
						IsExpired          string `json:"IsExpired"`
						ApplicationStatus  string `json:"ApplicationStatus"`
						TravelReq          string `json:"TravelReq"`
						ApplicationEndDate string `json:"ApplicationEndDate"`
						Rank               string `json:"Rank"`
						IsFavorite         string `json:"IsFavorite"`
						CompanyDesc        string `json:"CompanyDesc"`
					} `json:"properties"`
				} `json:"entry"`
			} `json:"feed"`
		}

		var jsonJobs Jobs

		c.OnResponse(func(r *colly.Response) {

			body_xml := strings.NewReader(string(r.Body))
			body_json, err := xj.Convert(body_xml)
			if err != nil {
				panic("That's embarrassing...")
			}

			var tempJson Jobs
			err = json.Unmarshal(body_json.Bytes(), &tempJson)
			if err != nil {
				panic(err.Error())
			}

			for _, elem := range tempJson.Feed.Entry {

				result_title := elem.Properties.Title
				result_url := fmt.Sprintf(job_base_url, elem.Properties.JobDetailsURL)

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

			jsonJobs.Feed.Entry = append(jsonJobs.Feed.Entry, tempJson.Feed.Entry...)
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
			c.Visit(start_url)
		}
	}
	return
}

func (runtime Runtime) Tesla(
	version int, isLocal bool) (
	response Response, results []Result) {
	switch version {
	case 1:

		start_url := "https://www.tesla.com/de_DE/careers/search#/"
		base_url := "https://www.tesla.com/careers/%s"
		file_name := "tesla.html"

		type Job struct {
			Title      string
			Url        string
			Department string
			Location   string
			Date       string
		}

		ctx, cancel := chromedp.NewContext(context.Background())
		defer cancel()
		var initialResponse string
		if err := chromedp.Run(ctx,
			chromedp.Navigate(start_url),
			chromedp.Sleep(5*time.Second),
			chromedp.OuterHTML("html", &initialResponse),
		); err != nil {
			panic(err)
		}
		SaveResponseToFileWithFileName(initialResponse, file_name)

		c := colly.NewCollector()

		c.OnHTML("html", func(e *colly.HTMLElement) {
			e.ForEach(".table-row", func(_ int, el *colly.HTMLElement) {
				result_title := el.ChildText("a")
				result_url := fmt.Sprintf(base_url, el.ChildAttr("a", "href"))
				result_department := el.ChildText(".listing-department")
				result_location := el.ChildText(".listing-location")
				result_date := el.ChildText(".listing-dateposted")

				_, err := netUrl.ParseRequestURI(result_url)
				if err == nil {

					temp_elem_json := Job{
						result_title,
						result_url,
						result_department,
						result_location,
						result_date,
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

		c.OnResponse(func(r *colly.Response) {
		})

		c.OnRequest(func(r *colly.Request) {
			fmt.Println(Gray(8-1, "Visiting"), Gray(8-1, r.URL.String()))
		})

		c.OnScraped(func(r *colly.Response) {
			RemoveFileWithFileName(file_name)
		})

		c.OnError(func(r *colly.Response, err error) {
			fmt.Println(
				Red("Request URL:"), Red(r.Request.URL),
				Red("failed with response:"), Red(r),
				Red("\nError:"), Red(err))
		})

		t := &http.Transport{}
		t.RegisterProtocol("file", http.NewFileTransport(http.Dir("/")))
		dir, err := os.Getwd()
		if err != nil {
			panic(err.Error())
		}
		c.WithTransport(t)
		c.Visit("file:" + dir + "/tesla.html")
	}
	return
}
