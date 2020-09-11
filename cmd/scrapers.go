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

	// "github.com/PuerkitoBio/goquery"
	// xj "github.com/basgys/goxml2json"
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
	Location    string
	Data        json.RawMessage
}

type Results []Result

const SecondsSleep = 2 // Seconds between pagination

func Extract(scraper_name string, scraper_version int) (results Results) {
	fmt.Println(Gray(8-1, "Starting Scrape..."))
	runtime := Runtime{scraper_name}
	strucReflected := reflect.ValueOf(runtime)
	method := strucReflected.MethodByName(scraper_name)
	params := []reflect.Value{}
	function_output := method.Call(params)
	results = function_output[0].Interface().(Results)
	results = Unique(results)
	return
}

func (results *Results) Add(
	scraper_name string,
	job_title string,
	job_url string,
	job_location string,
	job_data interface{}) {

	job_data_json, err := json.Marshal(job_data)
	if err != nil {
		panic(err.Error())
	}
	*results = append(*results, Result{
		scraper_name,
		job_title,
		job_url,
		job_location,
		job_data_json,
	})
}

func (runtime Runtime) Dreamingjobs() (results Results) {
	c := colly.NewCollector()
	start_url := "https://robimalco.github.io/dreamingjobs.github.io/"
	type Job struct {
		Title      string
		Url        string
		Department string
		Type       string
		Location   string
	}
	c.OnHTML("ul", func(e *colly.HTMLElement) {
		if strings.Contains(e.Attr("class"), "position") {
			e.ForEach("li", func(_ int, el *colly.HTMLElement) {
				result_title := el.ChildText("h2")
				result_url := start_url + el.ChildAttr("a", "href")
				result_department := el.ChildText("li[class=department]")
				result_type := el.ChildText("li[class=type]")
				result_location := el.ChildText("li[class=location]")
				results.Add(
					runtime.Name,
					result_title,
					result_url,
					result_location,
					Job{
						result_title,
						result_url,
						result_department,
						result_type,
						result_location,
					},
				)
			})
		}
	})
	c.OnRequest(func(r *colly.Request) {
		fmt.Println(Gray(8-1, "Visiting"), Gray(8-1, r.URL.String()))
	})
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println(Red("Request URL:"), Red(r.Request.URL))
	})
	c.Visit(start_url)
	return
}

func (runtime Runtime) Kununu() (results Results) {
	c := colly.NewCollector()
	start_url := "https://www.kununu.com/at/kununu/jobs"
	type Job struct {
		Title    string
		Url      string
		Location string
	}
	c.OnHTML("div", func(e *colly.HTMLElement) {
		if strings.Contains(e.Attr("class"), "company-profile-job-item") {
			result_title := e.ChildText("a")
			result_url := e.ChildAttr("a", "href")
			result_location := e.ChildText(".item-location")
			results.Add(
				runtime.Name,
				result_title,
				result_url,
				result_location,
				Job{
					result_title,
					result_url,
					result_location,
				},
			)
		}
	})
	c.OnRequest(func(r *colly.Request) {
		fmt.Println(Gray(8-1, "Visiting"), Gray(8-1, r.URL.String()))
	})
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println(Red("Request URL:"), Red(r.Request.URL))
	})
	c.Visit(start_url)
	return
}

func (runtime Runtime) Mitte() (results Results) {
	c := colly.NewCollector()
	start_url := "https://api.lever.co/v0/postings/mitte?&mode=json"
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
	c.OnResponse(func(r *colly.Response) {
		var jsonJobs Jobs
		err := json.Unmarshal(r.Body, &jsonJobs)
		if err != nil {
			panic(err.Error())
		}
		for _, elem := range jsonJobs {
			result_title := elem.Text
			result_url := elem.HostedURL
			result_location := elem.Categories.Location
			results.Add(
				runtime.Name,
				result_title,
				result_url,
				result_location,
				elem,
			)
		}
	})
	c.OnRequest(func(r *colly.Request) {
		fmt.Println(Gray(8-1, "Visiting"), Gray(8-1, r.URL.String()))
	})
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println(Red("Request URL:"), Red(r.Request.URL))
	})
	c.Visit(start_url)
	return
}

func (runtime Runtime) IMusician() (results Results) {
	c := colly.NewCollector()
	start_url := "https://imusician-digital-jobs.personio.de/"
	type Job struct {
		Title       string
		Url         string
		Description string
		Location    string
	}
	c.OnHTML("a", func(e *colly.HTMLElement) {
		if strings.Contains(e.Attr("class"), "job-box-link") {
			result_title := e.ChildText(".jb-title")
			result_url := e.Attr("href")
			result_description := e.ChildTexts("span")[0]
			result_location := e.ChildTexts("span")[2]
			results.Add(
				runtime.Name,
				result_title,
				result_url,
				result_location,
				Job{
					result_title,
					result_url,
					result_description,
					result_location,
				},
			)
		}
	})
	c.OnRequest(func(r *colly.Request) {
		fmt.Println(Gray(8-1, "Visiting"), Gray(8-1, r.URL.String()))
	})
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println(Red("Request URL:"), Red(r.Request.URL))
	})
	c.Visit(start_url)
	return
}

func (runtime Runtime) Babelforce() (results Results) {
	c := colly.NewCollector()
	start_url := "https://www.babelforce.com/jobs/"
	type Job struct {
		Title    string
		Url      string
		Location string
	}
	c.OnHTML("div", func(e *colly.HTMLElement) {
		if strings.Contains(e.Attr("class"), "qodef-portfolio") {
			result_title := e.ChildText("h5")
			result_url := e.ChildAttr("a", "href")
			result_location := "Berlin"
			results.Add(
				runtime.Name,
				result_title,
				result_url,
				result_location,
				Job{
					result_title,
					result_url,
					result_location,
				},
			)
		}
	})
	c.OnRequest(func(r *colly.Request) {
		fmt.Println(Gray(8-1, "Visiting"), Gray(8-1, r.URL.String()))
	})
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println(Red("Request URL:"), Red(r.Request.URL))
	})
	c.Visit(start_url)
	return
}

func (runtime Runtime) Zalando() (results Results) {
	c := colly.NewCollector()
	start_url := "https://jobs.zalando.com/api/jobs/?limit=100&offset=0"
	base_url := "https://jobs.zalando.com"
	base_job_url := "https://jobs.zalando.com/de/jobs/%s"
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
	c.OnResponse(func(r *colly.Response) {
		var jsonJobs Jobs
		err := json.Unmarshal(r.Body, &jsonJobs)
		if err != nil {
			panic(err.Error())
		}
		for _, elem := range jsonJobs.Data {
			result_title := elem.Title
			result_url := fmt.Sprintf(base_job_url, strconv.Itoa(elem.ID))
			result_location := strings.Join(elem.Offices, ",")
			results.Add(
				runtime.Name,
				result_title,
				result_url,
				result_location,
				elem,
			)
		}
		if jsonJobs.Next != "" {
			time.Sleep(SecondsSleep * time.Second)
			c.Visit(base_url + jsonJobs.Next)
		}
	})
	c.OnRequest(func(r *colly.Request) {
		fmt.Println(Gray(8-1, "Visiting"), Gray(8-1, r.URL.String()))
	})
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println(Red("Request URL:"), Red(r.Request.URL))
	})
	c.Visit(start_url)
	return
}

func (runtime Runtime) Google() (results Results) {
	c := colly.NewCollector()
	start_url := "https://careers.google.com/api/v2/jobs/search/?page_size=100&page=1"
	base_url := "https://careers.google.com/api/v2/jobs/search/?page_size=100&page="
	base_result_url := "https://careers.google.com/jobs/results/%s"
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
	c.OnResponse(func(r *colly.Response) {
		var jsonJobs JsonJobs
		err := json.Unmarshal(r.Body, &jsonJobs)
		if err != nil {
			panic(err.Error())
		}
		for _, elem := range jsonJobs.Jobs {
			result_title := elem.JobTitle
			result_url := fmt.Sprintf(base_result_url, strings.Split(elem.JobID, "/")[1])
			result_location := strings.Join(elem.Locations, ",")
			results.Add(
				runtime.Name,
				result_title,
				result_url,
				result_location,
				elem,
			)
		}
		total_pages := jsonJobs.Count/number_results_per_page + 2
		if total_pages <= jsonJobs.NextPage {
			return
		}
		if jsonJobs.NextPage != 0 {
			time.Sleep(SecondsSleep * time.Second)
			c.Visit(base_url + strconv.Itoa(jsonJobs.NextPage))
		}
	})
	c.OnRequest(func(r *colly.Request) {
		fmt.Println(Gray(8-1, "Visiting"), Gray(8-1, r.URL.String()))
	})
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println(Red("Request URL:"), Red(r.Request.URL))
	})
	c.Visit(start_url)
	return
}

func (runtime Runtime) Soundcloud() (results Results) {
	c := colly.NewCollector()
	start_url := "https://boards.greenhouse.io/embed/job_board?for=soundcloud71"
	type Job struct {
		Title      string
		Url        string
		Department string
		Location   string
	}
	c.OnHTML("section", func(e *colly.HTMLElement) {
		if strings.Contains(e.Attr("class"), "level-0") {
			result_department := e.ChildText("h3")
			e.ForEach("div", func(_ int, el *colly.HTMLElement) {
				result_title := el.ChildText("a")
				result_url := el.ChildAttr("a", "href")
				result_location := el.ChildText("span")
				results.Add(
					runtime.Name,
					result_title,
					result_url,
					result_location,
					Job{
						result_title,
						result_url,
						result_department,
						result_location,
					},
				)
			})
		}
	})
	c.OnRequest(func(r *colly.Request) {
		fmt.Println(Gray(8-1, "Visiting"), Gray(8-1, r.URL.String()))
	})
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println(Red("Request URL:"), Red(r.Request.URL))
	})
	c.Visit(start_url)
	return
}

func (runtime Runtime) Microsoft() (results Results) {
	c := colly.NewCollector()
	start_url := "https://careers.microsoft.com/us/en/search-results?s=1&from=0"
	base_url := "https://careers.microsoft.com/us/en/search-results?s=1&from=%d"
	base_result_url := "https://careers.microsoft.com/us/en/job/%s"
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
	c.OnResponse(func(r *colly.Response) {
		temp_resultsJson := strings.Split(string(r.Body), `"eagerLoadRefineSearch":`)[1]
		s_resultsJson := strings.Split(temp_resultsJson, `}; phApp.sessionParams`)[0]
		resultsJson := []byte(s_resultsJson)
		var jsonJobs JsonJobs
		err := json.Unmarshal(resultsJson, &jsonJobs)
		if err != nil {
			panic(err.Error())
		}
		for _, elem := range jsonJobs.Data.Jobs {
			result_title := elem.Title
			result_url := fmt.Sprintf(base_result_url, elem.JobID)
			result_location := strings.Join(elem.MultiLocation, ",")
			results.Add(
				runtime.Name,
				result_title,
				result_url,
				result_location,
				elem,
			)
		}
		total_pages := jsonJobs.TotalHits/number_results_per_page + 2
		if counter >= total_pages {
			return
		} else {
			counter++
			time.Sleep(SecondsSleep * time.Second)
			temp_url := fmt.Sprintf(base_url, counter*number_results_per_page)
			c.Visit(temp_url)
		}
	})
	c.OnRequest(func(r *colly.Request) {
		fmt.Println(Gray(8-1, "Visiting"), Gray(8-1, r.URL.String()))
	})
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println(Red("Request URL:"), Red(r.Request.URL))
	})
	c.Visit(start_url)
	return
}

func (runtime Runtime) Twitter() (results Results) {
	c := colly.NewCollector()
	start_url := "https://careers.twitter.com/content/careers-twitter/en/jobs.careers.search.json?limit=100&offset=0"
	base_url := "https://careers.twitter.com/content/careers-twitter/en/jobs.careers.search.json?limit=100&offset=%d"
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
	c.OnResponse(func(r *colly.Response) {
		var jsonJobs Jobs
		err := json.Unmarshal(r.Body, &jsonJobs)
		if err != nil {
			panic(err.Error())
		}
		for _, elem := range jsonJobs.Results {
			result_title := elem.Title
			result_url := elem.URL
			var temp_result_location []string
			for _, location := range elem.Locations {
				temp_result_location = append(temp_result_location, location.Title)
			}
			result_location := strings.Join(temp_result_location, ",")
			results.Add(
				runtime.Name,
				result_title,
				result_url,
				result_location,
				elem,
			)
		}
		total_pages := jsonJobs.TotalCount/number_results_per_page + 1
		if counter >= total_pages {
			return
		} else {
			counter++
			time.Sleep(SecondsSleep * time.Second)
			temp_t_url := fmt.Sprintf(base_url, counter*100)
			c.Visit(temp_t_url)
		}
	})
	c.OnRequest(func(r *colly.Request) {
		fmt.Println(Gray(8-1, "Visiting"), Gray(8-1, r.URL.String()))
	})
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println(Red("Request URL:"), Red(r.Request.URL))
	})
	c.Visit(start_url)
	return
}

func (runtime Runtime) Shopify() (results Results) {
	c := colly.NewCollector()
	start_url := "https://api.lever.co/v0/postings/shopify?mode=json"
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
	c.OnResponse(func(r *colly.Response) {
		var jsonJobs Jobs
		err := json.Unmarshal(r.Body, &jsonJobs)
		if err != nil {
			panic(err.Error())
		}
		for _, elem := range jsonJobs {
			result_title := elem.Text
			result_url := elem.HostedURL
			result_location := elem.Categories.Location
			results.Add(
				runtime.Name,
				result_title,
				result_url,
				result_location,
				elem,
			)
		}
	})
	c.OnRequest(func(r *colly.Request) {
		fmt.Println(Gray(8-1, "Visiting"), Gray(8-1, r.URL.String()))
	})
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println(Red("Request URL:"), Red(r.Request.URL))
	})
	c.Visit(start_url)
	return
}

func (runtime Runtime) Urbansport() (results Results) {
	c := colly.NewCollector()
	start_url := "https://boards.greenhouse.io/urbansportsclub"
	type Job struct {
		Title      string
		Url        string
		Department string
		Location   string
	}
	c.OnHTML("section", func(e *colly.HTMLElement) {
		if strings.Contains(e.Attr("class"), "level-0") {
			result_department := e.ChildText("h3")
			e.ForEach("div", func(_ int, el *colly.HTMLElement) {
				result_title := el.ChildText("a")
				result_url := el.ChildAttr("a", "href")
				result_location := el.ChildText("span")
				results.Add(
					runtime.Name,
					result_title,
					result_url,
					result_location,
					Job{
						result_title,
						result_url,
						result_department,
						result_location,
					},
				)
			})
		}
	})
	c.OnRequest(func(r *colly.Request) {
		fmt.Println(Gray(8-1, "Visiting"), Gray(8-1, r.URL.String()))
	})
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println(Red("Request URL:"), Red(r.Request.URL))
	})
	c.Visit(start_url)
	return
}

func (runtime Runtime) N26() (results Results) {
	c := colly.NewCollector()
	l := c.Clone()
	start_url := "https://n26.com/en/careers"
	base_url := "https://www.n26.com%s"
	type Job struct {
		Title    string
		Url      string
		Location string
		Contract string
	}
	c.OnHTML("a", func(e *colly.HTMLElement) {
		if strings.Contains(e.Attr("href"), "locations") {
			temp_location_url := e.Attr("href")
			location_url := fmt.Sprintf(base_url, temp_location_url)
			l.Visit(location_url)
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
	l.OnHTML("li", func(e *colly.HTMLElement) {
		e.ForEach("div", func(_ int, el *colly.HTMLElement) {
			if strings.Contains(el.ChildAttr("a", "href"), "positions") {
				temp_result_url := el.ChildAttr("a", "href")
				result_url := fmt.Sprintf(base_url, temp_result_url)
				result_title := el.ChildText("a")
				goquerySelection := el.DOM
				details_nodes := goquerySelection.Find("dd").Nodes
				result_location := details_nodes[0].FirstChild.Data
				result_contract := ""
				if len(details_nodes) > 1 {
					result_contract = details_nodes[1].FirstChild.Data
				}
				results.Add(
					runtime.Name,
					result_title,
					result_url,
					result_location,
					Job{
						result_title,
						result_url,
						result_location,
						result_contract,
					},
				)
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
	c.Visit(start_url)
	return
}

func (runtime Runtime) Blinkist() (results Results) {
	c := colly.NewCollector()
	start_url := "https://api.lever.co/v0/postings/blinkist?&mode=json"
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
	c.OnResponse(func(r *colly.Response) {
		var jsonJobs Jobs
		err := json.Unmarshal(r.Body, &jsonJobs)
		if err != nil {
			panic(err.Error())
		}
		for _, elem := range jsonJobs {
			result_title := elem.Text
			result_url := elem.HostedURL
			result_location := elem.Categories.Location
			results.Add(
				runtime.Name,
				result_title,
				result_url,
				result_location,
				elem,
			)
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
	c.Visit(start_url)
	return
}

func (runtime Runtime) Deutschebahn() (results Results) {
	c := colly.NewCollector()
	start_url := "https://karriere.deutschebahn.com/service/search/karriere-de/2653760?sort=pubExternalDate_td&pageNum=%s"
	base_result_url := "https://karriere.deutschebahn.com/%s"
	type Job struct {
		Url         string
		Title       string
		Location    string
		Entity      string
		Publication string
		Description string
	}
	c.OnHTML("ul", func(e *colly.HTMLElement) {
		if strings.Contains(e.Attr("class"), "result-items") {
			e.ForEach("li", func(_ int, el *colly.HTMLElement) {
				result_title := el.DOM.Find("span[class=title]").Text()
				result_location := strings.TrimSpace(el.DOM.Find("span[class=location]").Text())
				result_entity := strings.TrimSpace(el.DOM.Find("span[class=entity]").Text())
				result_publication := strings.TrimSpace(el.DOM.Find("span[class=publication]").Text())
				result_description := strings.TrimSpace(el.DOM.Find("p[class=responsibilities-text]").Text())
				temp_result_url, _ := el.DOM.Find("div[class=info]").Find("a").Attr("href")
				temp_result_url = fmt.Sprintf(base_result_url, temp_result_url)
				u, err := netUrl.Parse(temp_result_url)
				if err != nil {
					panic(err.Error())
				}
				u.RawQuery = ""
				result_url := u.String()
				results.Add(
					runtime.Name,
					result_title,
					result_url,
					result_location,
					Job{
						result_url,
						result_title,
						result_location,
						result_entity,
						result_publication,
						result_description,
					},
				)
			})
		}
	})
	c.OnHTML("a[class=active]", func(e *colly.HTMLElement) {
		next_page_url := fmt.Sprintf(start_url, e.Text)
		time.Sleep(SecondsSleep * time.Second)
		e.Request.Visit(next_page_url)
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
	c.Visit(fmt.Sprintf(start_url, "0"))
	return
}

func (runtime Runtime) Celo() (results Results) {
	c := colly.NewCollector()
	start_url := "https://api.lever.co/v0/postings/celo?mode=json"
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
	c.OnResponse(func(r *colly.Response) {
		var jsonJobs Jobs
		err := json.Unmarshal(r.Body, &jsonJobs)
		if err != nil {
			panic(err.Error())
		}
		for _, elem := range jsonJobs {
			result_title := elem.Text
			result_url := elem.HostedURL
			result_location := elem.Categories.Location
			results.Add(
				runtime.Name,
				result_title,
				result_url,
				result_location,
				elem,
			)
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
	c.Visit(start_url)
	return
}

func (runtime Runtime) Penta() (results Results) {
	c := colly.NewCollector()
	start_url := "https://boards.greenhouse.io/embed/job_board?for=penta"
	type Job struct {
		Title      string
		Url        string
		Department string
		Location   string
	}
	c.OnHTML("section", func(e *colly.HTMLElement) {
		if strings.Contains(e.Attr("class"), "level-0") {
			result_department := e.ChildText("h3")
			e.ForEach("div", func(_ int, el *colly.HTMLElement) {
				result_title := el.ChildText("a")
				result_url := el.ChildAttr("a", "href")
				result_location := el.ChildText("span")
				results.Add(
					runtime.Name,
					result_title,
					result_url,
					result_location,
					Job{
						result_title,
						result_url,
						result_department,
						result_location,
					},
				)
			})
		}
	})
	c.OnRequest(func(r *colly.Request) {
		fmt.Println(Gray(8-1, "Visiting"), Gray(8-1, r.URL.String()))
	})
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println(Red("Request URL:"), Red(r.Request.URL))
	})
	c.Visit(start_url)
	return
}

func (runtime Runtime) Contentful() (results Results) {
	c := colly.NewCollector()
	start_url := "https://boards.greenhouse.io/embed/job_board?for=contentful"
	type Job struct {
		Title      string
		Url        string
		Department string
		Location   string
	}
	c.OnHTML("section", func(e *colly.HTMLElement) {
		if strings.Contains(e.Attr("class"), "level-0") {
			result_department := e.ChildText("h3")
			e.ForEach("div", func(_ int, el *colly.HTMLElement) {
				result_title := el.ChildText("a")
				result_url := el.ChildAttr("a", "href")
				result_location := el.ChildText("span")
				results.Add(
					runtime.Name,
					result_title,
					result_url,
					result_location,
					Job{
						result_title,
						result_url,
						result_department,
						result_location,
					},
				)
			})
		}
	})
	c.OnRequest(func(r *colly.Request) {
		fmt.Println(Gray(8-1, "Visiting"), Gray(8-1, r.URL.String()))
	})
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println(Red("Request URL:"), Red(r.Request.URL))
	})
	c.Visit(start_url)
	return
}

func (runtime Runtime) Gympass() (results Results) {
	c := colly.NewCollector()
	start_url := "https://boards.greenhouse.io/embed/job_board?for=gympass"
	type Job struct {
		Title      string
		Url        string
		Department string
		Location   string
	}
	c.OnHTML("section", func(e *colly.HTMLElement) {
		if strings.Contains(e.Attr("class"), "level-0") {
			result_department := e.ChildText("h3")
			e.ForEach("div", func(_ int, el *colly.HTMLElement) {
				result_title := el.ChildText("a")
				result_url := el.ChildAttr("a", "href")
				result_location := el.ChildText("span")
				results.Add(
					runtime.Name,
					result_title,
					result_url,
					result_location,
					Job{
						result_title,
						result_url,
						result_department,
						result_location,
					},
				)
			})
		}
	})
	c.OnRequest(func(r *colly.Request) {
		fmt.Println(Gray(8-1, "Visiting"), Gray(8-1, r.URL.String()))
	})
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println(Red("Request URL:"), Red(r.Request.URL))
	})
	c.Visit(start_url)
	return
}

func (runtime Runtime) Hometogo() (results Results) {
	c := colly.NewCollector()
	start_url := "https://api.heavenhr.com/api/v1/positions/public/vacancies/?companyId=_VBAnjTs72rz0J-zBe1sYtA_"
	base_job_url := "https://hometogo.heavenhr.com/jobs/%s%s"
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
	c.OnResponse(func(r *colly.Response) {
		var jsonJobs Jobs
		err := json.Unmarshal(r.Body, &jsonJobs)
		if err != nil {
			panic(err.Error())
		}
		for _, elem := range jsonJobs.Data {
			result_title := elem.JobTitle
			result_url := fmt.Sprintf(base_job_url, elem.ID, "/apply")
			result_location := elem.Location
			results.Add(
				runtime.Name,
				result_title,
				result_url,
				result_location,
				elem,
			)
		}
	})
	c.OnRequest(func(r *colly.Request) {
		fmt.Println(Gray(8-1, "Visiting"), Gray(8-1, r.URL.String()))
	})
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println(Red("Request URL:"), Red(r.Request.URL))
	})
	c.Visit(start_url)
	return
}

func (runtime Runtime) Amazon() (results Results) {
	c := colly.NewCollector()
	start_url := "https://www.amazon.jobs/en/search.json?loc_query=Germany&country=DEU&result_limit=1000&offset=%d"
	base_job_url := "https://www.amazon.jobs%s"
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
	c.OnResponse(func(r *colly.Response) {
		var jsonJobs JsonJobs
		err := json.Unmarshal(r.Body, &jsonJobs)
		if err != nil {
			panic(err.Error())
		}
		for _, elem := range jsonJobs.Jobs {
			result_title := elem.Title
			result_url := fmt.Sprintf(base_job_url, elem.JobPath)
			result_location := elem.Location
			results.Add(
				runtime.Name,
				result_title,
				result_url,
				result_location,
				elem,
			)
		}
		total_pages := jsonJobs.Hits / number_results_per_page
		if counter < total_pages+1 {
			counter++
			next_page := fmt.Sprintf(start_url, counter*1000)
			time.Sleep(SecondsSleep * time.Second)
			c.Visit(next_page)
		}
	})
	c.OnRequest(func(r *colly.Request) {
		fmt.Println(Gray(8-1, "Visiting"), Gray(8-1, r.URL.String()))
	})
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println(Red("Request URL:"), Red(r.Request.URL))
	})
	c.Visit(fmt.Sprintf(start_url, 0))
	return
}

func (runtime Runtime) Lanalabs() (results Results) {
	c := colly.NewCollector()
	start_url := "https://lana-labs.breezy.hr/%s"
	type Job struct {
		Title      string
		Url        string
		Department string
		Type       string
		Location   string
	}
	c.OnHTML("ul", func(e *colly.HTMLElement) {
		if strings.Contains(e.Attr("class"), "position") {
			e.ForEach("li", func(_ int, el *colly.HTMLElement) {
				result_title := el.ChildText("h2")
				result_url := fmt.Sprintf(start_url, el.ChildAttr("a", "href"))
				result_department := el.ChildText("li[class=department]")
				result_type := el.ChildText("li[class=type]")
				result_location := el.ChildText("li[class=location]")
				results.Add(
					runtime.Name,
					result_title,
					result_url,
					result_location,
					Job{
						result_title,
						result_url,
						result_department,
						result_type,
						result_location,
					},
				)
			})
		}
	})
	c.OnRequest(func(r *colly.Request) {
		fmt.Println(Gray(8-1, "Visiting"), Gray(8-1, r.URL.String()))
	})
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println(Red("Request URL:"), Red(r.Request.URL))
	})
	c.Visit(fmt.Sprintf(start_url, ""))
	return
}

func (runtime Runtime) Slack() (results Results) {
	c := colly.NewCollector()
	start_url := "https://slack.com/intl/de-de/careers?eu_nc=1#opening"
	type Job struct {
		Title    string
		Url      string
		Location string
		Division string
	}
	c.OnHTML(".shadow-table", func(e *colly.HTMLElement) {
		e.ForEach("table", func(_ int, el *colly.HTMLElement) {
			result_division := el.ChildText("th")
			el.ForEach("tr", func(_ int, ell *colly.HTMLElement) {
				job_data := ell.ChildTexts(".for-desktop-only--table-cell")
				if len(job_data) > 0 {
					result_title := job_data[0]
					result_url := ell.ChildAttr("a", "href")
					result_location := job_data[2]
					results.Add(
						runtime.Name,
						result_title,
						result_url,
						result_location,
						Job{
							result_title,
							result_url,
							result_location,
							result_division,
						},
					)
				}
			})
		})
	})
	c.OnRequest(func(r *colly.Request) {
		fmt.Println(Gray(8-1, "Visiting"), Gray(8-1, r.URL.String()))
	})
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println(Red("Request URL:"), Red(r.Request.URL))
	})
	c.Visit(start_url)
	return
}

func (runtime Runtime) Revolut() (results Results) {
	c := colly.NewCollector()
	start_url := "https://api.lever.co/v0/postings/revolut?mode=json"
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
	c.OnResponse(func(r *colly.Response) {
		var jsonJobs Jobs
		err := json.Unmarshal(r.Body, &jsonJobs)
		if err != nil {
			panic(err.Error())
		}
		for _, elem := range jsonJobs {
			result_title := elem.Text
			result_url := elem.HostedURL
			result_location := elem.Categories.Location
			results.Add(
				runtime.Name,
				result_title,
				result_url,
				result_location,
				elem,
			)
		}
	})
	c.OnRequest(func(r *colly.Request) {
		fmt.Println(Gray(8-1, "Visiting"), Gray(8-1, r.URL.String()))
	})
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println(Red("Request URL:"), Red(r.Request.URL))
	})
	c.Visit(start_url)
	return
}

func (runtime Runtime) Mollie() (results Results) {
	c := colly.NewCollector()
	start_url := "https://api.lever.co/v0/postings/mollie?mode=json"
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
	c.OnResponse(func(r *colly.Response) {
		var jsonJobs Jobs
		err := json.Unmarshal(r.Body, &jsonJobs)
		if err != nil {
			panic(err.Error())
		}
		for _, elem := range jsonJobs {
			result_title := elem.Text
			result_url := elem.HostedURL
			result_location := elem.Categories.Location
			results.Add(
				runtime.Name,
				result_title,
				result_url,
				result_location,
				elem,
			)
		}
	})
	c.OnRequest(func(r *colly.Request) {
		fmt.Println(Gray(8-1, "Visiting"), Gray(8-1, r.URL.String()))
	})
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println(Red("Request URL:"), Red(r.Request.URL))
	})
	c.Visit(start_url)
	return
}

func (runtime Runtime) Circleci() (results Results) {
	c := colly.NewCollector()
	start_url := "https://boards.greenhouse.io/embed/job_board?for=circleci"
	base_job_url := "https://boards.greenhouse.io/circleci/jobs/%s"
	type Job struct {
		Title      string
		Url        string
		Department string
		Location   string
	}
	c.OnHTML("section", func(e *colly.HTMLElement) {
		if strings.Contains(e.Attr("class"), "level-0") {
			result_department := e.ChildText("h3")
			e.ForEach("div", func(_ int, el *colly.HTMLElement) {
				result_title := el.ChildText("a")
				result_url := fmt.Sprintf(base_job_url, strings.Split(el.ChildAttr("a", "href"), "=")[1])
				result_location := el.ChildText("span")
				results.Add(
					runtime.Name,
					result_title,
					result_url,
					result_location,
					Job{
						result_title,
						result_url,
						result_department,
						result_location,
					},
				)
			})
		}
	})
	c.OnRequest(func(r *colly.Request) {
		fmt.Println(Gray(8-1, "Visiting"), Gray(8-1, r.URL.String()))
	})
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println(Red("Request URL:"), Red(r.Request.URL))
	})
	c.Visit(start_url)
	return
}

func (runtime Runtime) Blacklane() (results Results) {
	c := colly.NewCollector()
	start_url := "https://boards.greenhouse.io/blacklane"
	base_job_url := "https://boards.greenhouse.io"
	type Job struct {
		Title      string
		Url        string
		Department string
		Location   string
	}
	c.OnHTML("section", func(e *colly.HTMLElement) {
		if strings.Contains(e.Attr("class"), "level-0") {
			result_department := e.ChildText("h3")
			e.ForEach("div", func(_ int, el *colly.HTMLElement) {
				result_title := el.ChildText("a")
				result_url := fmt.Sprintf(base_job_url, el.ChildAttr("a", "href"))
				result_location := el.ChildText("span")
				results.Add(
					runtime.Name,
					result_title,
					result_url,
					result_location,
					Job{
						result_title,
						result_url,
						result_department,
						result_location,
					},
				)
			})
		}
	})
	c.OnRequest(func(r *colly.Request) {
		fmt.Println(Gray(8-1, "Visiting"), Gray(8-1, r.URL.String()))
	})
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println(Red("Request URL:"), Red(r.Request.URL))
	})
	c.Visit(start_url)
	return
}

func (runtime Runtime) Auto1() (results Results) {
	c := colly.NewCollector()
	start_url := "https://www.auto1-group.com/smart-recruiters/jobs/search/?page=%d"
	base_job_url := "https://www.auto1-group.com/de/jobs/%s"
	counter := 1
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
	c.OnResponse(func(r *colly.Response) {
		var jsonJobs Auto1Jobs
		err := json.Unmarshal(r.Body, &jsonJobs)
		if err != nil {
			panic(err.Error())
		}
		for _, elem := range jsonJobs.Jobs.Hits.Hits {
			result_title := elem.Source.Title
			result_url := fmt.Sprintf(base_job_url, elem.Source.URL)
			result_location := elem.Source.LocationCity + "," + elem.Source.LocationCountry
			results.Add(
				runtime.Name,
				result_title,
				result_url,
				result_location,
				elem,
			)
		}
		total_pages := jsonJobs.Jobs.Hits.Total/number_results_per_page + 2
		if counter > total_pages {
			return
		} else {
			time.Sleep(SecondsSleep * time.Second)
			counter++
			c.Visit(fmt.Sprintf(start_url, counter))
		}
	})
	c.OnRequest(func(r *colly.Request) {
		fmt.Println(Gray(8-1, "Visiting"), Gray(8-1, r.URL.String()))
	})
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println(Red("Request URL:"), Red(r.Request.URL))
	})
	c.Visit(fmt.Sprintf(start_url, counter))
	return
}

func (runtime Runtime) Flixbus() (results Results) {
	c := colly.NewCollector()
	start_url := "https://flix.careers/api/jobs"
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
	c.OnResponse(func(r *colly.Response) {
		var jsonJobs FlixbusJobs
		err := json.Unmarshal(r.Body, &jsonJobs)
		if err != nil {
			panic(err.Error())
		}
		for _, elem := range jsonJobs.Jobs {
			result_title := elem.Title
			result_url := elem.AbsoluteURL
			result_location := elem.Location.Name
			results.Add(
				runtime.Name,
				result_title,
				result_url,
				result_location,
				elem,
			)
		}
	})
	c.OnRequest(func(r *colly.Request) {
		fmt.Println(Gray(8-1, "Visiting"), Gray(8-1, r.URL.String()))
	})
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println(Red("Request URL:"), Red(r.Request.URL))
	})
	c.Visit(start_url)
	return
}

func (runtime Runtime) Quora() (results Results) {
	c := colly.NewCollector()
	start_url := "https://boards.greenhouse.io/quora"
	base_job_url := "https://boards.greenhouse.io/%s"
	type Job struct {
		Title      string
		Url        string
		Department string
		Location   string
	}
	c.OnHTML("section", func(e *colly.HTMLElement) {
		if strings.Contains(e.Attr("class"), "level-0") {
			result_department := e.ChildText("h3")
			e.ForEach("div", func(_ int, el *colly.HTMLElement) {
				result_title := el.ChildText("a")
				result_url := fmt.Sprintf(base_job_url, el.ChildAttr("a", "href"))
				result_location := el.ChildText("span")
				results.Add(
					runtime.Name,
					result_title,
					result_url,
					result_location,
					Job{
						result_title,
						result_url,
						result_department,
						result_location,
					},
				)
			})
		}
	})
	c.OnRequest(func(r *colly.Request) {
		fmt.Println(Gray(8-1, "Visiting"), Gray(8-1, r.URL.String()))
	})
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println(Red("Request URL:"), Red(r.Request.URL))
	})
	c.Visit(start_url)
	return
}

func (runtime Runtime) Greenhouse() (results Results) {
	c := colly.NewCollector()
	start_url := "https://boards.greenhouse.io/embed/job_board?for=greenhouse"
	type Job struct {
		Title      string
		Url        string
		Department string
		Location   string
	}
	c.OnHTML("section", func(e *colly.HTMLElement) {
		if strings.Contains(e.Attr("class"), "level-0") {
			result_department := e.ChildText("h2")
			e.ForEach("div", func(_ int, el *colly.HTMLElement) {
				result_title := el.ChildText("a")
				t_j_url := strings.Split(el.ChildAttr("a", "href"), "=")[1]
				result_url := t_j_url
				result_location := el.ChildText("span")
				results.Add(
					runtime.Name,
					result_title,
					result_url,
					result_location,
					Job{
						result_title,
						result_url,
						result_department,
						result_location,
					},
				)
			})
		}
	})
	c.OnRequest(func(r *colly.Request) {
		fmt.Println(Gray(8-1, "Visiting"), Gray(8-1, r.URL.String()))
	})
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println(Red("Request URL:"), Red(r.Request.URL))
	})
	c.Visit(start_url)
	return
}

func (runtime Runtime) Docker() (results Results) {
	c := colly.NewCollector()
	start_url := "https://newton.newtonsoftware.com/career/CareerHome.action?clientId=8a7883c6708df1d40170a6df29950b39"
	type Job struct {
		Title    string
		Url      string
		Location string
	}
	c.OnHTML(".gnewtonCareerGroupRowClass", func(e *colly.HTMLElement) {
		result_title := e.ChildText("a")
		result_url := e.ChildAttr("a", "href")
		result_location := e.ChildText(".gnewtonCareerGroupJobDescriptionClass")
		results.Add(
			runtime.Name,
			result_title,
			result_url,
			result_location,
			Job{
				result_title,
				result_url,
				result_location,
			},
		)
	})
	c.OnRequest(func(r *colly.Request) {
		fmt.Println(Gray(8-1, "Visiting"), Gray(8-1, r.URL.String()))
	})
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println(Red("Request URL:"), Red(r.Request.URL))
	})
	c.Visit(start_url)
	return
}

func (runtime Runtime) Zapier() (results Results) {
	c := colly.NewCollector()
	start_url := "https://zapier.com/jobs"
	base_job_url := "https://zapier.com%s"
	type Job struct {
		Title      string
		Url        string
		Location   string
		Department string
	}
	c.OnHTML("section", func(e *colly.HTMLElement) {
		if strings.Contains(e.Attr("id"), "job-openings") {
			e.ForEach("li", func(_ int, el *colly.HTMLElement) {
				result_info := el.ChildText("a")
				result_temp_url := el.ChildAttr("a", "href")
				if !strings.Contains(result_temp_url, "https") {
					result_url := fmt.Sprintf(base_job_url, result_temp_url)
					info_split := strings.Split(result_info, " - ")
					result_department := info_split[0]
					result_title := info_split[1]
					result_location := "Remote"
					results.Add(
						runtime.Name,
						result_title,
						result_url,
						result_location,
						Job{
							result_title,
							result_url,
							result_location,
							result_department,
						},
					)
				}
			})
		}
	})
	c.OnRequest(func(r *colly.Request) {
		fmt.Println(Gray(8-1, "Visiting"), Gray(8-1, r.URL.String()))
	})
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println(Red("Request URL:"), Red(r.Request.URL))
	})
	c.Visit(start_url)
	return
}

func (runtime Runtime) Datadog() (results Results) {
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
	c.OnResponse(func(r *colly.Response) {
		var jsonJobs JsonJobs
		err := json.Unmarshal(r.Body, &jsonJobs)
		if err != nil {
			panic(err.Error())
		}
		for _, elem := range jsonJobs.Jobs {
			result_title := elem.Title
			result_url := elem.AbsoluteURL
			result_location := elem.Location.Name
			results.Add(
				runtime.Name,
				result_title,
				result_url,
				result_location,
				elem,
			)
		}
	})
	c.OnRequest(func(r *colly.Request) {
		fmt.Println(Gray(8-1, "Visiting"), Gray(8-1, r.URL.String()))
	})
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println(Red("Request URL:"), Red(r.Request.URL))
	})
	c.Visit(start_url)
	return
}

func (runtime Runtime) Stripe() (results Results) {
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
	c.OnResponse(func(r *colly.Response) {
		var jsonJobs JsonJobs
		err := json.Unmarshal(r.Body, &jsonJobs)
		if err != nil {
			panic(err.Error())
		}
		for _, elem := range jsonJobs.Jobs {
			result_title := elem.Title
			result_url := elem.AbsoluteURL
			result_location := elem.Location.Name
			results.Add(
				runtime.Name,
				result_title,
				result_url,
				result_location,
				elem,
			)
		}
	})
	c.OnRequest(func(r *colly.Request) {
		fmt.Println(Gray(8-1, "Visiting"), Gray(8-1, r.URL.String()))
	})
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println(Red("Request URL:"), Red(r.Request.URL))
	})
	c.Visit(start_url)
	return
}

func (runtime Runtime) Github() (results Results) {
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
	c.OnResponse(func(r *colly.Response) {
		var jsonJobs JsonJobs
		err := json.Unmarshal(r.Body, &jsonJobs)
		if err != nil {
			panic(err.Error())
		}
		for _, elem := range jsonJobs.Jobs {
			result_title := elem.Title
			result_url := elem.AbsoluteURL
			result_location := elem.Location.Name
			results.Add(
				runtime.Name,
				result_title,
				result_url,
				result_location,
				elem,
			)
		}
	})
	c.OnRequest(func(r *colly.Request) {
		fmt.Println(Gray(8-1, "Visiting"), Gray(8-1, r.URL.String()))
	})
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println(Red("Request URL:"), Red(r.Request.URL))
	})
	c.Visit(start_url)
	return
}

func (runtime Runtime) Getyourguide() (results Results) {
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
	c.OnResponse(func(r *colly.Response) {
		var jsonJobs JsonJobs
		err := json.Unmarshal(r.Body, &jsonJobs)
		if err != nil {
			panic(err.Error())
		}
		for _, elem := range jsonJobs.Jobs {
			result_title := elem.Title
			result_url := elem.AbsoluteURL
			result_location := elem.Location.Name
			results.Add(
				runtime.Name,
				result_title,
				result_url,
				result_location,
				elem,
			)
		}
	})
	c.OnRequest(func(r *colly.Request) {
		fmt.Println(Gray(8-1, "Visiting"), Gray(8-1, r.URL.String()))
	})
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println(Red("Request URL:"), Red(r.Request.URL))
	})
	c.Visit(start_url)
	return
}

func (runtime Runtime) Wefox() (results Results) {
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
	c.OnResponse(func(r *colly.Response) {
		var jsonJobs JsonJobs
		err := json.Unmarshal(r.Body, &jsonJobs)
		if err != nil {
			panic(err.Error())
		}
		for _, elem := range jsonJobs.Jobs {
			result_title := elem.Title
			result_url := elem.AbsoluteURL
			result_location := elem.Location.Name
			results.Add(
				runtime.Name,
				result_title,
				result_url,
				result_location,
				elem,
			)
		}
	})
	c.OnRequest(func(r *colly.Request) {
		fmt.Println(Gray(8-1, "Visiting"), Gray(8-1, r.URL.String()))
	})
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println(Red("Request URL:"), Red(r.Request.URL))
	})
	c.Visit(start_url)
	return
}

func (runtime Runtime) Celonis() (results Results) {
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
	c.OnResponse(func(r *colly.Response) {
		var jsonJobs JsonJobs
		err := json.Unmarshal(r.Body, &jsonJobs)
		if err != nil {
			panic(err.Error())
		}
		for _, elem := range jsonJobs.Jobs {
			result_title := elem.Title
			result_url := elem.AbsoluteURL
			result_location := elem.Location.Name
			results.Add(
				runtime.Name,
				result_title,
				result_url,
				result_location,
				elem,
			)
		}
	})
	c.OnRequest(func(r *colly.Request) {
		fmt.Println(Gray(8-1, "Visiting"), Gray(8-1, r.URL.String()))
	})
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println(Red("Request URL:"), Red(r.Request.URL))
	})
	c.Visit(start_url)
	return
}

func (runtime Runtime) Omio() (results Results) {
	c := colly.NewCollector()
	start_url := "https://api.smartrecruiters.com/v1/companies/Omio1/postings"
	base_job_url := "https://www.omio.com/jobs/#%s"
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
	c.OnResponse(func(r *colly.Response) {
		var jsonJobs JsonJobs
		err := json.Unmarshal(r.Body, &jsonJobs)
		if err != nil {
			panic(err.Error())
		}
		for _, elem := range jsonJobs.Content {
			result_title := elem.Name
			result_url := fmt.Sprintf(base_job_url, elem.ID)
			result_location := elem.Location.City + "," + elem.Location.Country
			if elem.Location.Remote {
				result_location = result_location + "," + "Remote"
			}
			results.Add(
				runtime.Name,
				result_title,
				result_url,
				result_location,
				elem,
			)
		}
	})
	c.OnRequest(func(r *colly.Request) {
		fmt.Println(Gray(8-1, "Visiting"), Gray(8-1, r.URL.String()))
	})
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println(Red("Request URL:"), Red(r.Request.URL))
	})
	c.Visit(start_url)
	return
}

func (runtime Runtime) Aboutyou() (results Results) {
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
	c.OnResponse(func(r *colly.Response) {
		var jsonJobs JsonJobs
		err := json.Unmarshal(r.Body, &jsonJobs)
		if err != nil {
			panic(err.Error())
		}
		for _, elem := range jsonJobs.Posts {
			result_title := elem.Title
			result_url := elem.URL
			result_location := elem.Location
			results.Add(
				runtime.Name,
				result_title,
				result_url,
				result_location,
				elem,
			)
		}
	})
	c.OnRequest(func(r *colly.Request) {
		fmt.Println(Gray(8-1, "Visiting"), Gray(8-1, r.URL.String()))
	})
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println(Red("Request URL:"), Red(r.Request.URL))
	})
	c.Visit(start_url)
	return
}

func (runtime Runtime) Depositsolutions() (results Results) {
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
	c.OnResponse(func(r *colly.Response) {
		var jsonJobs JsonJobs
		err := json.Unmarshal(r.Body, &jsonJobs)
		if err != nil {
			panic(err.Error())
		}
		for _, elem := range jsonJobs.Results {
			result_title := elem.Title
			result_url := fmt.Sprintf(base_job_url, elem.Shortcode)
			result_location := elem.Location.City + "," + elem.Location.Country
			if elem.Remote {
				result_location = result_location + "," + "Remote"
			}
			results.Add(
				runtime.Name,
				result_title,
				result_url,
				result_location,
				elem,
			)
		}
	})
	c.OnRequest(func(r *colly.Request) {
		fmt.Println(Gray(8-1, "Visiting"), Gray(8-1, r.URL.String()))
	})
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println(Red("Request URL:"), Red(r.Request.URL))
	})
	c.Request(
		"POST",
		start_url,
		strings.NewReader(""),
		nil,
		http.Header{"Content-Type": []string{"application/x-www-form-urlencoded"}},
	)
	return
}

func (runtime Runtime) Taxfix() (results Results) {
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
	c.OnResponse(func(r *colly.Response) {
		var jsonJobs JsonJobs
		err := json.Unmarshal(r.Body, &jsonJobs)
		if err != nil {
			panic(err.Error())
		}
		for _, elem := range jsonJobs.Jobs {
			result_title := elem.Title
			result_url := elem.AbsoluteURL
			result_location := elem.Location.Name
			results.Add(
				runtime.Name,
				result_title,
				result_url,
				result_location,
				elem,
			)
		}
	})
	c.OnRequest(func(r *colly.Request) {
		fmt.Println(Gray(8-1, "Visiting"), Gray(8-1, r.URL.String()))
	})
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println(Red("Request URL:"), Red(r.Request.URL))
	})
	c.Visit(start_url)
	return
}

func (runtime Runtime) Moonfare() (results Results) {
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
	c.OnResponse(func(r *colly.Response) {
		var jsonJobs JsonJobs
		err := json.Unmarshal(r.Body, &jsonJobs)
		if err != nil {
			panic(err.Error())
		}
		for _, elem := range jsonJobs.Jobs {
			result_title := elem.Title
			result_url := elem.AbsoluteURL
			result_location := elem.Location.Name
			results.Add(
				runtime.Name,
				result_title,
				result_url,
				result_location,
				elem,
			)
		}
	})
	c.OnRequest(func(r *colly.Request) {
		fmt.Println(Gray(8-1, "Visiting"), Gray(8-1, r.URL.String()))
	})
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println(Red("Request URL:"), Red(r.Request.URL))
	})
	c.Visit(start_url)
	return
}

func (runtime Runtime) Fincompare() (results Results) {
	c := colly.NewCollector()
	start_url := "https://api.lever.co/v0/postings/fincompare?mode=json"
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
	c.OnResponse(func(r *colly.Response) {
		var jsonJobs Jobs
		err := json.Unmarshal(r.Body, &jsonJobs)
		if err != nil {
			panic(err.Error())
		}
		for _, elem := range jsonJobs {
			result_title := elem.Text
			result_url := elem.HostedURL
			result_location := elem.Categories.Location
			results.Add(
				runtime.Name,
				result_title,
				result_url,
				result_location,
				elem,
			)
		}
	})
	c.OnRequest(func(r *colly.Request) {
		fmt.Println(Gray(8-1, "Visiting"), Gray(8-1, r.URL.String()))
	})
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println(Red("Request URL:"), Red(r.Request.URL))
	})
	c.Visit(start_url)
	return
}

func (runtime Runtime) Billie() (results Results) {
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
	c.OnResponse(func(r *colly.Response) {
		var jsonJobs JsonJobs
		err := json.Unmarshal(r.Body, &jsonJobs)
		if err != nil {
			panic(err.Error())
		}
		for _, elem := range jsonJobs.Jobs {
			result_title := elem.Title
			result_url := elem.AbsoluteURL
			result_location := elem.Location.Name
			results.Add(
				runtime.Name,
				result_title,
				result_url,
				result_location,
				elem,
			)
		}
	})
	c.OnRequest(func(r *colly.Request) {
		fmt.Println(Gray(8-1, "Visiting"), Gray(8-1, r.URL.String()))
	})
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println(Red("Request URL:"), Red(r.Request.URL))
	})
	c.Visit(start_url)
	return
}

func (runtime Runtime) Pairfinance() (results Results) {
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
	c.OnResponse(func(r *colly.Response) {
		var jsonJobs JsonJobs
		err := json.Unmarshal(r.Body, &jsonJobs)
		if err != nil {
			panic(err.Error())
		}
		for _, elem := range jsonJobs.Jobs {
			result_title := elem.Title
			result_url := elem.AbsoluteURL
			result_location := elem.Location.Name
			results.Add(
				runtime.Name,
				result_title,
				result_url,
				result_location,
				elem,
			)
		}
	})
	c.OnRequest(func(r *colly.Request) {
		fmt.Println(Gray(8-1, "Visiting"), Gray(8-1, r.URL.String()))
	})
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println(Red("Request URL:"), Red(r.Request.URL))
	})
	c.Visit(start_url)
	return
}

func (runtime Runtime) Getsafe() (results Results) {
	c := colly.NewCollector()
	start_url := "https://getsafe-jobs.personio.de"
	type Job struct {
		Title       string
		Url         string
		Description string
		Location    string
	}
	c.OnHTML("a", func(e *colly.HTMLElement) {
		if strings.Contains(e.Attr("class"), "job-box-link") {
			result_title := e.ChildText(".jb-title")
			result_url := e.Attr("href")
			result_description := e.ChildTexts("span")[0]
			result_location := e.ChildTexts("span")[2]
			results.Add(
				runtime.Name,
				result_title,
				result_url,
				result_location,
				Job{
					result_title,
					result_url,
					result_description,
					result_location,
				},
			)
		}
	})
	c.OnRequest(func(r *colly.Request) {
		fmt.Println(Gray(8-1, "Visiting"), Gray(8-1, r.URL.String()))
	})
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println(Red("Request URL:"), Red(r.Request.URL))
	})
	c.Visit(start_url)
	return
}

func (runtime Runtime) Liqid() (results Results) {
	c := colly.NewCollector()
	url := "https://liqid-jobs.personio.de"
	type Job struct {
		Title    string
		Url      string
		Type     string
		Location string
	}
	c.OnHTML("div", func(e *colly.HTMLElement) {
		if strings.Contains(e.Attr("class"), "job-list-desc") {
			result_title := e.ChildText("a")
			result_url := e.ChildAttr("a", "href")
			result_info := strings.Split(e.ChildText("p"), "·")
			result_type := strings.Join(strings.Fields(strings.TrimSpace(result_info[0])), " ")
			result_location := strings.Join(strings.Fields(strings.TrimSpace(result_info[1])), " ")
			results.Add(
				runtime.Name,
				result_title,
				result_url,
				result_location,
				Job{
					result_title,
					result_url,
					result_type,
					result_location,
				},
			)
		}
	})
	c.OnRequest(func(r *colly.Request) {
		fmt.Println(Gray(8-1, "Visiting"), Gray(8-1, r.URL.String()))
	})
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println(Red("Request URL:"), Red(r.Request.URL))
	})
	c.Visit(url)
	return
}

func (runtime Runtime) Elementinsurance() (results Results) {
	c := colly.NewCollector()
	start_url := "https://elementinsuranceag.recruitee.com/api/offers"
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
	c.OnResponse(func(r *colly.Response) {
		var jsonJobs Jobs
		err := json.Unmarshal(r.Body, &jsonJobs)
		if err != nil {
			panic(err.Error())
		}
		for _, elem := range jsonJobs.Offers {
			result_title := elem.Title
			result_url := elem.CareersURL
			result_location := elem.City + "," + elem.Country
			results.Add(
				runtime.Name,
				result_title,
				result_url,
				result_location,
				elem,
			)
		}
	})
	c.OnRequest(func(r *colly.Request) {
		fmt.Println(Gray(8-1, "Visiting"), Gray(8-1, r.URL.String()))
	})
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println(Red("Request URL:"), Red(r.Request.URL))
	})
	c.Visit(start_url)
	return
}

func (runtime Runtime) Freeda() (results Results) {
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
	c.OnResponse(func(r *colly.Response) {
		var jsonJobs JsonJobs
		err := json.Unmarshal(r.Body, &jsonJobs)
		if err != nil {
			panic(err.Error())
		}
		for _, elem := range jsonJobs.Jobs {
			result_title := elem.Title
			result_url := elem.AbsoluteURL
			result_location := elem.Location.Name
			results.Add(
				runtime.Name,
				result_title,
				result_url,
				result_location,
				elem,
			)
		}
	})
	c.OnRequest(func(r *colly.Request) {
		fmt.Println(Gray(8-1, "Visiting"), Gray(8-1, r.URL.String()))
	})
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println(Red("Request URL:"), Red(r.Request.URL))
	})
	c.Visit(start_url)
	return
}

func (runtime Runtime) Talentgarden() (results Results) {
	c := colly.NewCollector()
	start_url := "https://talentgarden.bamboohr.com/jobs/embed2.php?departmentId=0"
	type Job struct {
		Title      string
		Url        string
		Location   string
		Department string
	}
	c.OnHTML("div", func(e *colly.HTMLElement) {
		if strings.Contains(e.Attr("class"), "BambooHR-ATS-board") {
			e.ForEach("li[class=BambooHR-ATS-Department-Item]", func(_ int, el *colly.HTMLElement) {
				result_department := strings.TrimSpace(el.ChildText("div[class=BambooHR-ATS-Department-Header]"))
				el.ForEach("ul[class=BambooHR-ATS-Jobs-List]", func(_ int, ell *colly.HTMLElement) {
					result_title := ell.ChildText("a")
					result_url := "https:" + ell.ChildAttr("a", "href")
					result_location := ell.ChildText("span")
					results.Add(
						runtime.Name,
						result_title,
						result_url,
						result_location,
						Job{
							result_title,
							result_url,
							result_location,
							result_department,
						},
					)
				})
			})
		}
	})
	c.OnRequest(func(r *colly.Request) {
		fmt.Println(Gray(8-1, "Visiting"), Gray(8-1, r.URL.String()))
	})
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println(Red("Request URL:"), Red(r.Request.URL))
	})
	c.Visit(start_url)
	return
}

func (runtime Runtime) Facileit() (results Results) {
	start_url := "https://inrecruiting.intervieweb.it/app.php?module=iframeAnnunci&k=1382636f10340a4ca6713ef6df70205a&LAC=Facileit&act1=23"
	file_name := "facileit.html"
	type Job struct {
		Title       string
		Url         string
		Location    string
		Description string
	}
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()
	var initialResponse string
	if err := chromedp.Run(ctx,
		chromedp.Navigate(start_url),
		// chromedp.Sleep(30*time.Second),
		chromedp.WaitVisible(".titolo_annuncio"),
		chromedp.OuterHTML("html", &initialResponse),
	); err != nil {
		panic(err.Error())
	}
	SaveResponseToFileWithFileName(initialResponse, file_name)
	c := colly.NewCollector()
	c.OnHTML("dt", func(e *colly.HTMLElement) {
		result_infos := e.ChildTexts("span")
		result_title := result_infos[0]
		result_location := result_infos[1]
		result_url := e.ChildAttr("a", "href")
		result_description := e.ChildText(".description")
		results.Add(
			runtime.Name,
			result_title,
			result_url,
			result_location,
			Job{
				result_title,
				result_url,
				result_location,
				result_description,
			},
		)
	})
	c.OnRequest(func(r *colly.Request) {
		fmt.Println(Gray(8-1, "Visiting"), Gray(8-1, r.URL.String()))
	})
	c.OnScraped(func(r *colly.Response) {
		RemoveFileWithFileName(file_name)
	})
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println(Red("Request URL:"), Red(r.Request.URL))
	})
	t := &http.Transport{}
	t.RegisterProtocol("file", http.NewFileTransport(http.Dir("/")))
	dir, err := os.Getwd()
	if err != nil {
		panic(err.Error())
	}
	c.WithTransport(t)
	c.Visit("file:" + dir + "/" + file_name)
	return
}

func (runtime Runtime) Vodafone() (results Results) {
	c := colly.NewCollector()
	start_url := "https://careers.vodafone.com/search/?startrow=%d"
	base_job_url := "https://careers.vodafone.com"
	number_results_per_page := 25
	counter := 0
	type Job struct {
		Title    string
		Url      string
		Location string
		Date     string
	}
	c.OnHTML(".html5", func(e *colly.HTMLElement) {
		e.ForEach(".data-row", func(_ int, el *colly.HTMLElement) {
			result_title := strings.Join(strings.Fields(strings.TrimSpace(el.ChildTexts("a")[0])), " ")
			result_url := fmt.Sprintf(base_job_url, strings.Join(strings.Fields(strings.TrimSpace(el.ChildAttr("a", "href"))), " "))
			result_location := strings.Join(strings.Fields(strings.TrimSpace(el.ChildText("span[class=jobLocation]"))), " ")
			result_date := strings.Join(strings.Fields(strings.TrimSpace(el.ChildText("span[class=jobDate]"))), " ")
			results.Add(
				runtime.Name,
				result_title,
				result_url,
				result_location,
				Job{
					result_title,
					result_url,
					result_location,
					result_date,
				},
			)
		})
		temp_total_results := strings.Split(e.ChildText(".paginationLabel"), " ")
		string_total_results := temp_total_results[len(temp_total_results)-1]
		total_results, err := strconv.Atoi(string_total_results)
		if err != nil {
			panic(err.Error())
		}
		total_pages := total_results/number_results_per_page + 2
		if counter >= total_pages {
			return
		} else {
			counter++
			time.Sleep(SecondsSleep * time.Second)
			temp_v_url := fmt.Sprintf(start_url, counter*number_results_per_page)
			c.Visit(temp_v_url)
		}
	})
	c.OnRequest(func(r *colly.Request) {
		fmt.Println(Gray(8-1, "Visiting"), Gray(8-1, r.URL.String()))
	})
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println(Red("Request URL:"), Red(r.Request.URL))
	})
	c.Visit(fmt.Sprintf(start_url, 0))
	return
}

func (runtime Runtime) Glovo() (results Results) {
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
	c.OnResponse(func(r *colly.Response) {
		var jsonJobs JsonJobs
		err := json.Unmarshal(r.Body, &jsonJobs)
		if err != nil {
			panic(err.Error())
		}
		for _, elem := range jsonJobs.Jobs {
			result_title := elem.Title
			result_url := elem.AbsoluteURL
			result_location := elem.Location.Name
			results.Add(
				runtime.Name,
				result_title,
				result_url,
				result_location,
				elem,
			)
		}
	})
	c.OnRequest(func(r *colly.Request) {
		fmt.Println(Gray(8-1, "Visiting"), Gray(8-1, r.URL.String()))
	})
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println(Red("Request URL:"), Red(r.Request.URL))
	})
	c.Visit(start_url)
	return
}

func (runtime Runtime) Glickon() (results Results) {
	c := colly.NewCollector()
	l := c.Clone()
	section_url := "https://core.glickon.com/api/candidate/latest/companies/glickon"
	department_url := "https://core.glickon.com/api/candidate/latest/sections/%s?from_www=true"
	job_api_url := "https://core.glickon.com/api/candidate/latest/company_challenges/%s"
	job_base_url := "https://www.glickon.com/en/challenges/"
	type JsonJobs struct {
		Title       string
		Url         string
		Location    string
		Description string
	}
	type Departments struct {
		Sections []struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"sections"`
	}
	type Jobs struct {
		Challenges []struct {
			Hash        string `json:"hash"`
			Name        string `json:"name"`
			Description string `json:"description"`
		} `json:"challenges"`
	}
	type Job struct {
		Location string `json:"location"`
	}
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
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println(Red("Request URL:"), Red(r.Request.URL))
	})
	l.OnResponse(func(r *colly.Response) {
		var tempJsonJobs Jobs
		err := json.Unmarshal(r.Body, &tempJsonJobs)
		if err != nil {
			panic(err.Error())
		}
		for _, elem := range tempJsonJobs.Challenges {
			result_title := elem.Name
			result_url := job_base_url + elem.Hash
			result_description := elem.Description
			location_req, err := http.NewRequest(http.MethodGet, fmt.Sprintf(job_api_url, elem.Hash), nil)
			spaceClient := http.Client{}
			res, err := spaceClient.Do(location_req)
			body, err := ioutil.ReadAll(res.Body)
			temp_location := Job{}
			err = json.Unmarshal(body, &temp_location)
			if err != nil {
				panic(err.Error())
			}
			result_location := temp_location.Location
			results.Add(
				runtime.Name,
				result_title,
				result_url,
				result_location,
				JsonJobs{
					result_title,
					result_url,
					result_location,
					result_description,
				},
			)
		}
	})
	c.Visit(section_url)
	return
}

func (runtime Runtime) Satispay() (results Results) {
	c := colly.NewCollector()
	start_url := "https://satispay.breezy.hr%s"
	type Job struct {
		Title      string
		Url        string
		Location   string
		Department string
		Type       string
	}
	c.OnHTML("ul", func(e *colly.HTMLElement) {
		if strings.Contains(e.Attr("class"), "position") {
			e.ForEach("li", func(_ int, el *colly.HTMLElement) {
				result_title := el.ChildText("h2")
				result_url := fmt.Sprintf(start_url, el.ChildAttr("a", "href"))
				result_department := el.ChildText("li[class=department]")
				result_type := el.ChildText("li[class=type]")
				result_location := el.ChildText("li[class=location]")
				results.Add(
					runtime.Name,
					result_title,
					result_url,
					result_location,
					Job{
						result_title,
						result_url,
						result_location,
						result_department,
						result_type,
					},
				)
			})
		}
	})
	c.OnRequest(func(r *colly.Request) {
		fmt.Println(Gray(8-1, "Visiting"), Gray(8-1, r.URL.String()))
	})
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println(Red("Request URL:"), Red(r.Request.URL))
	})
	c.Visit(fmt.Sprintf(start_url, ""))
	return
}

func (runtime Runtime) Medtronic() (results Results) {
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
		fmt.Println(Red("Request URL:"), Red(r.Request.URL))
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
			results.Add(
				runtime.Name,
				result_title,
				result_url,
				result_location,
				Job{
					result_title,
					result_url,
					result_location,
					result_category,
					result_description,
				},
			)
		})
		// string_number_pages := e.ChildText("div[id=jPaginateNumPages]")
		// number_pages, _ := strconv.Atoi(strings.Split(string_number_pages, ".")[0])
		for counter := 2; counter <= 4; counter++ {
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
				result_location := strings.ReplaceAll(
					strings.ReplaceAll(
						strings.Split(strings.Split(elem, `<span class="location">`)[1], `</span>`)[0],
						"nttttttttttt", ""),
					"ntttttttttttt", "")
				result_category := strings.Split(strings.Split(elem, `<span class="category">`)[1], `</span>`)[0]
				result_description := strings.Split(strings.Split(elem, `<p class="jlr_description">`)[1], `</p>`)[0]
				results.Add(
					runtime.Name,
					result_title,
					result_url,
					result_location,
					Job{
						result_title,
						result_url,
						result_location,
						result_category,
						result_description,
					},
				)
			}
		}
	})
	x.OnRequest(func(r *colly.Request) {
		fmt.Println(Gray(8-1, "Visiting"), Gray(8-1, r.URL.String()))
	})
	c.Visit(start_url)
	return
}