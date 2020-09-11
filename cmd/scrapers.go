package cmd

import (
	// "context"
	"encoding/json"
	"fmt"
	// "io/ioutil"
	// "net/http"
	netUrl "net/url"
	// "os"
	"reflect"
	"strconv"
	"strings"
	"time"

	// "github.com/PuerkitoBio/goquery"
	// xj "github.com/basgys/goxml2json"
	// "github.com/chromedp/chromedp"
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
		fmt.Println(
			Red("Request URL:"), Red(r.Request.URL),
			Red("failed with response:"), Red(r),
			Red("\nError:"), Red(err))
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
		fmt.Println(
			Red("Request URL:"), Red(r.Request.URL),
			Red("failed with response:"), Red(r),
			Red("\nError:"), Red(err))
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
		fmt.Println(
			Red("Request URL:"), Red(r.Request.URL),
			Red("failed with response:"), Red(r),
			Red("\nError:"), Red(err))
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
		fmt.Println(
			Red("Request URL:"), Red(r.Request.URL),
			Red("failed with response:"), Red(r),
			Red("\nError:"), Red(err))
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
		fmt.Println(
			Red("Request URL:"), Red(r.Request.URL),
			Red("failed with response:"), Red(r),
			Red("\nError:"), Red(err))
	})
	c.Visit(start_url)
	return
}

func (runtime Runtime) Zalando(version int) (results Results) {
	switch version {
	case 1:
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
			fmt.Println(
				Red("Request URL:"), Red(r.Request.URL),
				Red("failed with response:"), Red(r),
				Red("\nError:"), Red(err))
		})
		c.Visit(start_url)
	}
	return
}

func (runtime Runtime) Google(version int) (results Results) {
	switch version {
	case 1:
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
			fmt.Println(
				Red("Request URL:"), Red(r.Request.URL),
				Red("failed with response:"), Red(r),
				Red("\nError:"), Red(err))
		})
		c.Visit(start_url)
	}
	return
}

func (runtime Runtime) Soundcloud(version int) (results Results) {
	switch version {
	case 1:
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
			fmt.Println(
				Red("Request URL:"), Red(r.Request.URL),
				Red("failed with response:"), Red(r),
				Red("\nError:"), Red(err))
		})
		c.Visit(start_url)
	}
	return
}

func (runtime Runtime) Microsoft(version int) (results Results) {
	switch version {
	case 1:
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
			fmt.Println(
				Red("Request URL:"), Red(r.Request.URL),
				Red("failed with response:"), Red(r),
				Red("\nError:"), Red(err))
		})
		c.Visit(start_url)
	}
	return
}

func (runtime Runtime) Twitter(version int) (results Results) {
	switch version {
	case 1:
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
			fmt.Println(
				Red("Request URL:"), Red(r.Request.URL),
				Red("failed with response:"), Red(r),
				Red("\nError:"), Red(err))
		})
		c.Visit(start_url)
	}
	return
}

func (runtime Runtime) Shopify(version int) (results Results) {
	switch version {
	case 1:
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
			fmt.Println(
				Red("Request URL:"), Red(r.Request.URL),
				Red("failed with response:"), Red(r),
				Red("\nError:"), Red(err))
		})
		c.Visit(start_url)
	}
	return
}

func (runtime Runtime) Urbansport(version int) (results Results) {
	switch version {
	case 1:
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
			fmt.Println(
				Red("Request URL:"), Red(r.Request.URL),
				Red("failed with response:"), Red(r),
				Red("\nError:"), Red(err))
		})
		c.Visit(start_url)
	}
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