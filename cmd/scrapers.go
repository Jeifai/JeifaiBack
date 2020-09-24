package cmd

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	netUrl "net/url"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

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

const SecondsSleep = 2 // Seconds between each pagination

func Extract(scraper_name string) (results Results) {
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

/**
██      ███████ ██    ██ ███████ ██████
██      ██      ██    ██ ██      ██   ██
██      █████   ██    ██ █████   ██████
██      ██       ██  ██  ██      ██   ██
███████ ███████   ████   ███████ ██   ██
*/
func Lever(start_url string, runtime_name string, results *Results) {
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
	c := colly.NewCollector()
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
				runtime_name,
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

func (runtime Runtime) Mitte() (results Results) {
	start_url := "https://api.lever.co/v0/postings/mitte?&mode=json"
	Lever(start_url, runtime.Name, &results)
	return
}

func (runtime Runtime) Shopify() (results Results) {
	start_url := "https://api.lever.co/v0/postings/shopify?mode=json"
	Lever(start_url, runtime.Name, &results)
	return
}

func (runtime Runtime) Blinkist() (results Results) {
	start_url := "https://api.lever.co/v0/postings/blinkist?&mode=json"
	Lever(start_url, runtime.Name, &results)
	return
}

func (runtime Runtime) Celo() (results Results) {
	start_url := "https://api.lever.co/v0/postings/celo?mode=json"
	Lever(start_url, runtime.Name, &results)
	return
}

func (runtime Runtime) Revolut() (results Results) {
	start_url := "https://api.lever.co/v0/postings/revolut?mode=json"
	Lever(start_url, runtime.Name, &results)
	return
}

func (runtime Runtime) Mollie() (results Results) {
	start_url := "https://api.lever.co/v0/postings/mollie?mode=json"
	Lever(start_url, runtime.Name, &results)
	return
}

func (runtime Runtime) Fincompare() (results Results) {
	start_url := "https://api.lever.co/v0/postings/fincompare?mode=json"
	Lever(start_url, runtime.Name, &results)
	return
}

func (runtime Runtime) Beat81() (results Results) {
	start_url := "https://api.lever.co/v0/postings/beat81?mode=json"
	Lever(start_url, runtime.Name, &results)
	return
}

func (runtime Runtime) Forto() (results Results) {
	start_url := "https://api.lever.co/v0/postings/forto?mode=json"
	Lever(start_url, runtime.Name, &results)
	return
}

func (runtime Runtime) Ecosia() (results Results) {
	start_url := "https://api.lever.co/v0/postings/ecosia?mode=json"
	Lever(start_url, runtime.Name, &results)
	return
}

func (runtime Runtime) Automationhero() (results Results) {
	start_url := "https://api.lever.co/v0/postings/automationhero?mode=json"
	Lever(start_url, runtime.Name, &results)
	return
}

func (runtime Runtime) Cargoone() (results Results) {
	start_url := "https://api.lever.co/v0/postings/cargo-2?mode=json"
	Lever(start_url, runtime.Name, &results)
	return
}

func (runtime Runtime) Klarna() (results Results) {
	start_url := "https://api.lever.co/v0/postings/klarna?mode=json"
	Lever(start_url, runtime.Name, &results)
	return
}

func (runtime Runtime) Kaiahealth() (results Results) {
	start_url := "https://api.lever.co/v0/postings/kaiahealth?mode=json"
	Lever(start_url, runtime.Name, &results)
	return
}

func (runtime Runtime) Improbable() (results Results) {
	start_url := "https://api.lever.co/v0/postings/improbable?&mode=json"
	Lever(start_url, runtime.Name, &results)
	return
}

func (runtime Runtime) Klarx() (results Results) {
	start_url := "https://api.lever.co/v0/postings/klarx?&mode=json"
	Lever(start_url, runtime.Name, &results)
	return
}

/**
██████  ███████ ██████  ███████  ██████  ███    ██ ██  ██████       ██
██   ██ ██      ██   ██ ██      ██    ██ ████   ██ ██ ██    ██     ███
██████  █████   ██████  ███████ ██    ██ ██ ██  ██ ██ ██    ██      ██
██      ██      ██   ██      ██ ██    ██ ██  ██ ██ ██ ██    ██      ██
██      ███████ ██   ██ ███████  ██████  ██   ████ ██  ██████       ██
*/
func Personio1(start_url string, runtime_name string, results *Results) {
	type Job struct {
		Title       string
		Url         string
		Location    string
		Description string
	}
	c := colly.NewCollector()
	c.OnHTML("a", func(e *colly.HTMLElement) {
		if strings.Contains(e.Attr("class"), "job-box-link") {
			result_title := e.ChildText(".jb-title")
			result_url := e.Attr("href")
			result_description := e.ChildTexts("span")[0]
			result_location := e.ChildTexts("span")[2]
			results.Add(
				runtime_name,
				result_title,
				result_url,
				result_location,
				Job{
					result_url,
					result_title,
					result_location,
					result_description,
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

func (runtime Runtime) Casparhealth() (results Results) {
	start_url := "https://goreha-jobs.personio.de"
	Personio1(start_url, runtime.Name, &results)
	return
}

func (runtime Runtime) IMusician() (results Results) {
	start_url := "https://imusician-digital-jobs.personio.de/"
	Personio1(start_url, runtime.Name, &results)
	return
}

func (runtime Runtime) Getsafe() (results Results) {
	start_url := "https://getsafe-jobs.personio.de"
	Personio1(start_url, runtime.Name, &results)
	return
}

func (runtime Runtime) Idagio() (results Results) {
	start_url := "https://idagio-jobs.personio.de"
	Personio1(start_url, runtime.Name, &results)
	return
}

func (runtime Runtime) Raisin() (results Results) {
	start_url := "https://raisin-jobs.personio.de"
	Personio1(start_url, runtime.Name, &results)
	return
}

func (runtime Runtime) Construyo() (results Results) {
	start_url := "https://partum-gmbh-jobs.personio.de"
	Personio1(start_url, runtime.Name, &results)
	return
}

func (runtime Runtime) Morressier() (results Results) {
	start_url := "https://morressier-jobs.personio.de/"
	Personio1(start_url, runtime.Name, &results)
	return
}

func (runtime Runtime) Enmacc() (results Results) {
	start_url := "https://enmacc-jobs.personio.de/"
	Personio1(start_url, runtime.Name, &results)
	return
}

func (runtime Runtime) Yfoodlabs() (results Results) {
	start_url := "https://yfoodlabs-jobs.personio.de/"
	Personio1(start_url, runtime.Name, &results)
	return
}

func (runtime Runtime) Personio() (results Results) {
	start_url := "https://personio-jobs.personio.de/"
	Personio1(start_url, runtime.Name, &results)
	return
}

func (runtime Runtime) Egym() (results Results) {
	start_url := "https://egym-jobs.personio.de/"
	Personio1(start_url, runtime.Name, &results)
	return
}

func (runtime Runtime) Buildingradar() (results Results) {
	start_url := "https://building-radar-jobs.personio.de/"
	Personio1(start_url, runtime.Name, &results)
	return
}

func (runtime Runtime) Capmo() (results Results) {
	start_url := "https://capmo-jobs.personio.de/"
	Personio1(start_url, runtime.Name, &results)
	return
}

func (runtime Runtime) Agrilution() (results Results) {
	start_url := "https://agrilution-jobs.personio.de/"
	Personio1(start_url, runtime.Name, &results)
	return
}

func (runtime Runtime) Remberg() (results Results) {
	start_url := "https://remberg-jobs.personio.de/"
	Personio1(start_url, runtime.Name, &results)
	return
}

/**
██████  ███████ ██████  ███████  ██████  ███    ██ ██  ██████      ██████
██   ██ ██      ██   ██ ██      ██    ██ ████   ██ ██ ██    ██          ██
██████  █████   ██████  ███████ ██    ██ ██ ██  ██ ██ ██    ██      █████
██      ██      ██   ██      ██ ██    ██ ██  ██ ██ ██ ██    ██     ██
██      ███████ ██   ██ ███████  ██████  ██   ████ ██  ██████      ███████
*/
func Personio2(start_url string, runtime_name string, results *Results) {
	type Job struct {
		Url      string
		Title    string
		Location string
		Type     string
	}
	c := colly.NewCollector()
	c.OnHTML("body", func(e *colly.HTMLElement) {
		e.ForEach(".job-list-desc", func(_ int, el *colly.HTMLElement) {
			result_url := el.ChildAttr("a", "href")
			result_title := el.ChildText("a")
			result_info := strings.Split(el.ChildText("p"), "·")
			result_type := strings.Join(strings.Fields(strings.TrimSpace(result_info[0])), " ")
			result_location := strings.Join(strings.Fields(strings.TrimSpace(result_info[1])), " ")
			results.Add(
				runtime_name,
				result_title,
				result_url,
				result_location,
				Job{
					result_url,
					result_title,
					result_location,
					result_type,
				},
			)
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

func (runtime Runtime) Crosslend() (results Results) {
	start_url := "https://crosslend-jobs.personio.de/"
	Personio2(start_url, runtime.Name, &results)
	return
}

func (runtime Runtime) Liqid() (results Results) {
	start_url := "https://liqid-jobs.personio.de"
	Personio2(start_url, runtime.Name, &results)
	return
}

func (runtime Runtime) Acatus() (results Results) {
	start_url := "https://acatus-jobs.personio.de/?language=en#all"
	Personio2(start_url, runtime.Name, &results)
	return
}

func (runtime Runtime) Juniqe() (results Results) {
	start_url := "https://juniqe-jobs.personio.de/"
	Personio2(start_url, runtime.Name, &results)
	return
}

func (runtime Runtime) Skoove() (results Results) {
	start_url := "https://skoove-jobs.personio.de/"
	Personio2(start_url, runtime.Name, &results)
	return
}

func (runtime Runtime) Cluno() (results Results) {
	start_url := "https://cluno-jobs.personio.de/"
	Personio2(start_url, runtime.Name, &results)
	return
}

func (runtime Runtime) Navvis() (results Results) {
	start_url := "https://navvis-jobs.personio.de/"
	Personio2(start_url, runtime.Name, &results)
	return
}

func (runtime Runtime) Magazino() (results Results) {
	start_url := "https://magazino-jobs.personio.de/"
	Personio2(start_url, runtime.Name, &results)
	return
}

func (runtime Runtime) Reflekt() (results Results) {
	start_url := "https://reflekt-gmbh-jobs.personio.de"
	Personio2(start_url, runtime.Name, &results)
	return
}

func (runtime Runtime) Ndgit() (results Results) {
	start_url := "https://ndgit-jobs.personio.de/"
	Personio2(start_url, runtime.Name, &results)
	return
}

func (runtime Runtime) Dreamingjobs() (results Results) {
	start_url := "https://robimalco.github.io/dreamingjobs.github.io/"
	type Job struct {
		Title      string
		Url        string
		Department string
		Type       string
		Location   string
	}
	c := colly.NewCollector()
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
	start_url := "https://www.kununu.com/at/kununu/jobs"
	type Job struct {
		Title    string
		Url      string
		Location string
	}
	c := colly.NewCollector()
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

func (runtime Runtime) Babelforce() (results Results) {
	start_url := "https://www.babelforce.com/jobs/"
	type Job struct {
		Title    string
		Url      string
		Location string
	}
	c := colly.NewCollector()
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
	c := colly.NewCollector()
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
	c := colly.NewCollector()
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
	start_url := "https://boards.greenhouse.io/embed/job_board?for=soundcloud71"
	type Job struct {
		Title      string
		Url        string
		Department string
		Location   string
	}
	c := colly.NewCollector()
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
	c := colly.NewCollector()
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
	c := colly.NewCollector()
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

func (runtime Runtime) Urbansport() (results Results) {
	start_url := "https://boards.greenhouse.io/urbansportsclub"
	type Job struct {
		Title      string
		Url        string
		Department string
		Location   string
	}
	c := colly.NewCollector()
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
	start_url := "https://n26.com/en/careers"
	base_url := "https://www.n26.com%s"
	type Job struct {
		Title    string
		Url      string
		Location string
		Contract string
	}
	c := colly.NewCollector()
	l := c.Clone()
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

func (runtime Runtime) Deutschebahn() (results Results) {
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
	c := colly.NewCollector()
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

func (runtime Runtime) Penta() (results Results) {
	start_url := "https://boards.greenhouse.io/embed/job_board?for=penta"
	type Job struct {
		Title      string
		Url        string
		Department string
		Location   string
	}
	c := colly.NewCollector()
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
	start_url := "https://boards.greenhouse.io/embed/job_board?for=contentful"
	type Job struct {
		Title      string
		Url        string
		Department string
		Location   string
	}
	c := colly.NewCollector()
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
	start_url := "https://boards.greenhouse.io/embed/job_board?for=gympass"
	type Job struct {
		Title      string
		Url        string
		Department string
		Location   string
	}
	c := colly.NewCollector()
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
	c := colly.NewCollector()
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
	c := colly.NewCollector()
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
	start_url := "https://lana-labs.breezy.hr/%s"
	type Job struct {
		Title      string
		Url        string
		Department string
		Type       string
		Location   string
	}
	c := colly.NewCollector()
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
	start_url := "https://slack.com/intl/de-de/careers?eu_nc=1#opening"
	type Job struct {
		Title    string
		Url      string
		Location string
		Division string
	}
	c := colly.NewCollector()
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

func (runtime Runtime) Circleci() (results Results) {
	start_url := "https://boards.greenhouse.io/embed/job_board?for=circleci"
	base_job_url := "https://boards.greenhouse.io/circleci/jobs/%s"
	type Job struct {
		Title      string
		Url        string
		Department string
		Location   string
	}
	c := colly.NewCollector()
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
	start_url := "https://boards.greenhouse.io/blacklane"
	base_job_url := "https://boards.greenhouse.io"
	type Job struct {
		Title      string
		Url        string
		Department string
		Location   string
	}
	c := colly.NewCollector()
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
	c := colly.NewCollector()
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
	c := colly.NewCollector()
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
	start_url := "https://boards.greenhouse.io/quora"
	base_job_url := "https://boards.greenhouse.io/%s"
	type Job struct {
		Title      string
		Url        string
		Department string
		Location   string
	}
	c := colly.NewCollector()
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
	start_url := "https://boards.greenhouse.io/embed/job_board?for=greenhouse"
	type Job struct {
		Title      string
		Url        string
		Department string
		Location   string
	}
	c := colly.NewCollector()
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
	start_url := "https://newton.newtonsoftware.com/career/CareerHome.action?clientId=8a7883c6708df1d40170a6df29950b39"
	type Job struct {
		Title    string
		Url      string
		Location string
	}
	c := colly.NewCollector()
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
	start_url := "https://zapier.com/jobs"
	base_job_url := "https://zapier.com%s"
	type Job struct {
		Title      string
		Url        string
		Location   string
		Department string
	}
	c := colly.NewCollector()
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
	c := colly.NewCollector()
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
	c := colly.NewCollector()
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
	c := colly.NewCollector()
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
	c := colly.NewCollector()
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
	c := colly.NewCollector()
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
	c := colly.NewCollector()
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
	c := colly.NewCollector()
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
	c := colly.NewCollector()
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
	c := colly.NewCollector()
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
	c := colly.NewCollector()
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
	c := colly.NewCollector()
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

func (runtime Runtime) Billie() (results Results) {
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
	c := colly.NewCollector()
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
	c := colly.NewCollector()
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

func (runtime Runtime) Elementinsurance() (results Results) {
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
	c := colly.NewCollector()
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
			if fmt.Sprint(elem.Remote) == "true" {
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

func (runtime Runtime) Freeda() (results Results) {
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
	c := colly.NewCollector()
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
	start_url := "https://talentgarden.bamboohr.com/jobs/embed2.php?departmentId=0"
	type Job struct {
		Title      string
		Url        string
		Location   string
		Department string
	}
	c := colly.NewCollector()
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
	c := colly.NewCollector()
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
	c := colly.NewCollector()
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
	c := colly.NewCollector()
	l := c.Clone()
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
	start_url := "https://satispay.breezy.hr%s"
	type Job struct {
		Title      string
		Url        string
		Location   string
		Department string
		Type       string
	}
	c := colly.NewCollector()
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
	c := colly.NewCollector()
	l := c.Clone()
	x := l.Clone()
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

func (runtime Runtime) Bendingspoons() (results Results) {
	start_url := "https://website.rolemodel.bendingspoons.com/roles.json"
	base_job_url := "https://bendingspoons.com/careers.html?x=%s"
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
	c := colly.NewCollector()
	c.OnResponse(func(r *colly.Response) {
		var jsonJobs Jobs
		err := json.Unmarshal(r.Body, &jsonJobs)
		if err != nil {
			panic(err.Error())
		}
		for _, elem := range jsonJobs {
			result_title := elem.Title
			result_url := fmt.Sprintf(base_job_url, elem.ID)
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

func (runtime Runtime) Bcg() (results Results) {
	type Job struct {
		Title       string
		Url         string
		Location    string
		Date        string
		Description string
	}
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()
	start_url := "https://talent.bcg.com/en_US/apply/SearchJobs/?folderOffset=%d"
	start_offset := 0
	number_results_per_page := 20
	_ = number_results_per_page
	var initialResponse string
	if err := chromedp.Run(ctx,
		chromedp.Navigate(fmt.Sprintf(start_url, start_offset)),
		chromedp.OuterHTML(".body_Chrome", &initialResponse),
	); err != nil {
		panic(err)
	}
	temp_total_results := strings.Split(
		strings.Split(
			strings.Split(initialResponse, `jobPaginationLegend`)[1], "</span>")[0], " ")
	total_results, _ := strconv.Atoi(temp_total_results[len(temp_total_results)-1])
	for i := 0; i <= total_results; i += number_results_per_page {
		fmt.Println(Gray(8-1, "Visiting"), Gray(8-1, fmt.Sprintf(start_url, i)))
		var pageResponse string
		if err := chromedp.Run(ctx,
			chromedp.Navigate(fmt.Sprintf(start_url, i)),
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
					result_description,
				},
			)
		}
	}
	return
}

func (runtime Runtime) Deloitte() (results Results) {
	t := &http.Transport{}
	t.RegisterProtocol("file", http.NewFileTransport(http.Dir("/")))
	dir, err := os.Getwd()
	if err != nil {
		panic(err.Error())
	}
	file_name := "deloitteDepartments.html"
	type Job struct {
		Url         string
		Title       string
		Location    string
		Company     string
		Entity      string
		Department  string
		Id          string
		Type        string
		Date        string
		Description string
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
	SaveResponseToFileWithFileName(initialResponse, file_name)
	c := colly.NewCollector()
	c.WithTransport(t)
	x := c.Clone()
	x.WithTransport(t)
	c.OnHTML("html", func(e *colly.HTMLElement) {
		e.ForEach(".jobs-list-item", func(_ int, el *colly.HTMLElement) {
			result_url := strings.Join(strings.Fields(strings.TrimSpace(el.ChildAttr("a", "href"))), " ")
			result_title := strings.Join(strings.Fields(strings.TrimSpace(el.ChildText("h4"))), " ")
			result_location := strings.Join(strings.Fields(strings.TrimSpace(el.ChildText(".job-location"))), " ")
			result_company := strings.Join(strings.Fields(strings.TrimSpace(el.ChildText(".memberfirm"))), " ")
			result_entity := strings.Join(strings.Fields(strings.TrimSpace(el.ChildText(".memberentity"))), " ")
			result_department := strings.Join(strings.Fields(strings.TrimSpace(el.ChildText(".job-category"))), " ")
			result_id := strings.Join(strings.Fields(strings.TrimSpace(el.ChildText(".job-id"))), " ")
			result_type := strings.Join(strings.Fields(strings.TrimSpace(el.ChildText(".job-type"))), " ")
			result_date := strings.Join(strings.Fields(strings.TrimSpace(el.ChildText(".job-postdate"))), " ")
			result_description := strings.Join(strings.Fields(strings.TrimSpace(el.ChildText(".job-description"))), " ")
			results.Add(
				runtime.Name,
				result_title,
				result_url,
				result_location,
				Job{
					result_title,
					result_url,
					result_location,
					result_company,
					result_entity,
					result_department,
					result_id,
					result_type,
					result_date,
					result_description,
				},
			)
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
	c.OnScraped(func(r *colly.Response) {
		RemoveFileWithFileName(file_name)
	})
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println(Red("Request URL:"), Red(r.Request.URL))
	})
	x.OnHTML("html", func(e *colly.HTMLElement) {
		e.ForEach(".jobs-list-item", func(_ int, el *colly.HTMLElement) {
			result_url := strings.Join(strings.Fields(strings.TrimSpace(el.ChildAttr("a", "href"))), " ")
			result_title := strings.Join(strings.Fields(strings.TrimSpace(el.ChildText("h4"))), " ")
			result_location := strings.Join(strings.Fields(strings.TrimSpace(el.ChildText(".job-location"))), " ")
			result_company := strings.Join(strings.Fields(strings.TrimSpace(el.ChildText(".memberfirm"))), " ")
			result_entity := strings.Join(strings.Fields(strings.TrimSpace(el.ChildText(".memberentity"))), " ")
			result_department := strings.Join(strings.Fields(strings.TrimSpace(el.ChildText(".job-category"))), " ")
			result_id := strings.Join(strings.Fields(strings.TrimSpace(el.ChildText(".job-id"))), " ")
			result_type := strings.Join(strings.Fields(strings.TrimSpace(el.ChildText(".job-type"))), " ")
			result_date := strings.Join(strings.Fields(strings.TrimSpace(el.ChildText(".job-postdate"))), " ")
			result_description := strings.Join(strings.Fields(strings.TrimSpace(el.ChildText(".job-description"))), " ")
			results.Add(
				runtime.Name,
				result_title,
				result_url,
				result_location,
				Job{
					result_title,
					result_url,
					result_location,
					result_company,
					result_entity,
					result_department,
					result_id,
					result_type,
					result_date,
					result_description,
				},
			)
		})
	})
	x.OnRequest(func(r *colly.Request) {
		fmt.Println(Gray(8-1, "Visiting"), Gray(8-1, r.URL.String()))
	})
	x.OnError(func(r *colly.Response, err error) {
		fmt.Println(Red("Request URL:"), Red(r.Request.URL))
	})
	c.WithTransport(t)
	c.Visit("file:" + dir + "/" + file_name)
	return
}

func (runtime Runtime) Bayer() (results Results) {
	start_url := "https://career.bayer.com/en/jobs-search?page=%d"
	base_job_url := "https://career.bayer.com%s"
	counter := 0
	type Job struct {
		Title    string
		Url      string
		Location string
		Date     string
		Country  string
	}
	c := colly.NewCollector()
	c.OnHTML(".content", func(e *colly.HTMLElement) {
		e.ForEach("tr", func(_ int, el *colly.HTMLElement) {
			result_title := el.ChildText("a")
			result_url := fmt.Sprintf(base_job_url, el.ChildAttr("a", "href"))
			result_date := el.ChildText(".views-field-field-job-last-modify-time")
			result_country := el.ChildText(".views-field-field-job-country")
			result_location := el.ChildText(".views-field-field-job-location")
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
					result_country,
				},
			)
		})
		goqueryselect := e.DOM
		temp_last_page, _ := goqueryselect.Find(".pager__item--last").Find("a").Attr("href")
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
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println(Red("Request URL:"), Red(r.Request.URL))
	})
	c.Visit(fmt.Sprintf(start_url, 0))
	return
}

func (runtime Runtime) Roche() (results Results) {
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
	c := colly.NewCollector()
	c.OnResponse(func(r *colly.Response) {
		var jsonJobs JsonJobs
		err := json.Unmarshal(r.Body, &jsonJobs)
		if err != nil {
			panic(err.Error())
		}
		for _, elem := range jsonJobs.Jobs.Items {
			result_title := elem.Title
			result_url := fmt.Sprintf(base_url, elem.DetailsURL)
			result_location := elem.PrimaryLocation.City + "," + elem.PrimaryLocation.Country
			results.Add(
				runtime.Name,
				result_title,
				result_url,
				result_location,
				elem,
			)
		}
		total_matches := jsonJobs.Jobs.TotalMatches
		total_pages := total_matches / number_results_per_page
		for i := 1; i <= total_pages; i++ {
			time.Sleep(SecondsSleep * time.Second)
			c.Visit(fmt.Sprintf(start_url, number_results_per_page, i))
		}
	})
	c.OnRequest(func(r *colly.Request) {
		fmt.Println(Gray(8-1, "Visiting"), Gray(8-1, r.URL.String()))
	})
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println(Red("Request URL:"), Red(r.Request.URL))
	})
	c.Visit(fmt.Sprintf(start_url, number_results_per_page, 0))
	return
}

func (runtime Runtime) Msd() (results Results) {
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
	var jsonJobs JsonJobs
	err := json.Unmarshal([]byte(jsonjobs_sections), &jsonJobs)
	if err != nil {
		panic(err.Error())
	}
	items_per_page := jsonJobs.EagerLoadRefineSearch.Hits
	total_matches := jsonJobs.EagerLoadRefineSearch.TotalHits
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
			result_location := elem.Location
			results.Add(
				runtime.Name,
				result_title,
				result_url,
				result_location,
				elem,
			)
		}
	}
	return
}

func (runtime Runtime) Subitoit() (results Results) {
	start_url := "https://info.subito.it/lavora-con-noi.htm"
	type Job struct {
		Url        string
		Title      string
		Location   string
		Department string
	}
	c := colly.NewCollector()
	c.OnHTML(".work-openings", func(e *colly.HTMLElement) {
		e.ForEach(".list-box-item", func(_ int, el *colly.HTMLElement) {
			result_title := el.ChildText("a")
			result_url := el.ChildAttr("a", "href")
			result_department := el.ChildText("h4")
			results.Add(
				runtime.Name,
				result_title,
				result_url,
				"Milano",
				Job{
					result_url,
					result_title,
					"Milano",
					result_department,
				},
			)
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

func (runtime Runtime) Square() (results Results) {
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
	c := colly.NewCollector()
	c.OnResponse(func(r *colly.Response) {
		var jsonJobs Jobs
		err := json.Unmarshal(r.Body, &jsonJobs)
		if err != nil {
			panic(err.Error())
		}
		for _, elem := range jsonJobs.Content {
			result_title := elem.Name
			result_url := fmt.Sprintf(base_job_url, elem.ID)
			result_location := elem.Location.City + "," + elem.Location.Country
			if elem.Location.Remote {
				result_location = result_location + ", Remote"
			}
			results.Add(
				runtime.Name,
				result_title,
				result_url,
				result_location,
				elem,
			)
		}
		total_matches := jsonJobs.TotalFound
		total_pages := total_matches / number_results_per_page
		for i := 1; i <= total_pages; i++ {
			time.Sleep(SecondsSleep * time.Second)
			c.Visit(fmt.Sprintf(start_url, number_results_per_page*i))
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

func (runtime Runtime) Facebook() (results Results) {
	start_url := "https://www.facebook.com/careers/jobs?results_per_page=100&page=%d"
	base_job_url := "https://www.facebook.com%s"
	number_results_per_page := 100
	type Job struct {
		Title    string
		Url      string
		Location string
		Info     string
	}
	c := colly.NewCollector()
	c.OnHTML("#search_result", func(e *colly.HTMLElement) {
		e.ForEach("a", func(_ int, el *colly.HTMLElement) {
			result_url := fmt.Sprintf(base_job_url, el.Attr("href"))
			result_title := el.ChildText("._8sel")
			result_location := el.ChildText("._8sen")
			var result_info []string
			temp_result_info := el.ChildTexts("._8see")
			for _, elem := range temp_result_info {
				if !strings.Contains(elem, "+") {
					result_info = append(result_info, elem)
				}
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
					strings.Join(result_info, " - "),
				},
			)
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
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println(Red("Request URL:"), Red(r.Request.URL))
	})
	c.Visit(fmt.Sprintf(start_url, 1))
	return
}

func (runtime Runtime) Paintgun() (results Results) {
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
	c := colly.NewCollector()
	c.OnResponse(func(r *colly.Response) {
		var jsonJobs JsonJobs
		err := json.Unmarshal(r.Body, &jsonJobs)
		if err != nil {
			panic(err.Error())
		}
		for _, elem := range jsonJobs.Items {
			result_title := elem.Title
			result_url := fmt.Sprintf(job_base_url, elem.IDParam)
			result_location := elem.Country.Name
			if elem.IsRemote {
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

func (runtime Runtime) Nen() (results Results) {
	/*
		results = append(results, Result{
			runtime.Name,
			"Salesforce Lead",
			"https://www.linkedin.com/jobs/view/1947567619",
			"Milano",
			[]byte("{}"),
		})
	*/
	return
}

func (runtime Runtime) Amboss() (results Results) {
	start_url := "https://www.amboss.com/us/career-opportunities"
	base_job_url := "https://www.amboss.com%s"
	type Job struct {
		Url      string
		Title    string
		Location string
	}
	c := colly.NewCollector()
	c.OnHTML(".jobs-list", func(e *colly.HTMLElement) {
		e.ForEach("._pwggpq", func(_ int, el *colly.HTMLElement) {
			result_url := fmt.Sprintf(base_job_url, el.Attr("href"))
			result_title := el.ChildText("._pulkya")
			result_location := el.ChildText("._1f1zsnz")
			results.Add(
				runtime.Name,
				result_title,
				result_url,
				result_location,
				Job{
					result_url,
					result_title,
					result_location,
				},
			)
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

func (runtime Runtime) Chatterbug() (results Results) {
	start_url := "https://chatterbug.recruitee.com/api/offers"
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
	c := colly.NewCollector()
	c.OnResponse(func(r *colly.Response) {
		var jsonJobs Jobs
		err := json.Unmarshal(r.Body, &jsonJobs)
		if err != nil {
			panic(err.Error())
		}
		for _, elem := range jsonJobs.Offers {
			result_title := elem.Title
			result_url := elem.CareersURL
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

func (runtime Runtime) Infarm() (results Results) {
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
	c := colly.NewCollector()
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

func (runtime Runtime) Pitch() (results Results) {
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
	c := colly.NewCollector()
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

func (runtime Runtime) Careerfoundry() (results Results) {
	start_url := "https://careerfoundry.recruitee.com/api/offers"
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
	c := colly.NewCollector()
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
			if fmt.Sprint(elem.Remote) == "true" {
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

func (runtime Runtime) Joblift() (results Results) {
	start_url := "https://joblift-talent.freshteam.com/jobs"
	base_job_url := "https://joblift-talent.freshteam.com%s"
	type Job struct {
		Url        string
		Title      string
		Location   string
		Department string
	}
	c := colly.NewCollector()
	c.OnHTML(".job-role-list", func(e *colly.HTMLElement) {
		e.ForEach("li:not([class])", func(_ int, el *colly.HTMLElement) {
			result_department := strings.Split(el.ChildText(".role-title"), "-")[0]
			el.ForEach(".job-list-info", func(_ int, ell *colly.HTMLElement) {
				result_url := fmt.Sprintf(base_job_url, ell.ChildAttr("a", "href"))
				result_title := ell.ChildText(".job-title")
				result_location := strings.Split(ell.ChildText(".location-info"), "\n")[0]
				results.Add(
					runtime.Name,
					result_title,
					result_url,
					result_location,
					Job{
						result_url,
						result_title,
						result_location,
						result_department,
					},
				)
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

func (runtime Runtime) Kontist() (results Results) {
	start_url := "https://kontist.com/careers/jobs.json"
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
	c := colly.NewCollector()
	c.OnResponse(func(r *colly.Response) {
		var jsonJobs Jobs
		err := json.Unmarshal(r.Body, &jsonJobs)
		if err != nil {
			panic(err.Error())
		}
		for _, elem := range jsonJobs.Data {
			result_title := elem.Title
			result_url := elem.JobURL
			result_location := elem.TmpLocation
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

func (runtime Runtime) Medloop() (results Results) {
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
	c := colly.NewCollector()
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

func (runtime Runtime) Medwing() (results Results) {
	start_url := "https://team.medwing.com/wp-json/wp/v2/jobs"
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
		Links         struct {
			WpTerm []struct {
				Taxonomy   string `json:"taxonomy"`
				Embeddable bool   `json:"embeddable"`
				Href       string `json:"href"`
			} `json:"wp:term"`
		} `json:"_links"`
	}
	type Location []struct {
		Name string `json:"name"`
	}
	c := colly.NewCollector()
	c.OnResponse(func(r *colly.Response) {
		var jsonJobs Jobs
		err := json.Unmarshal(r.Body, &jsonJobs)
		if err != nil {
			panic(err.Error())
		}
		for _, elem := range jsonJobs {
			result_title := elem.Title.Rendered
			result_url := elem.Link
			var location_api_url string
			for _, taxonomy := range elem.Links.WpTerm {
				if taxonomy.Taxonomy == "location" {
					location_api_url = taxonomy.Href
				}
			}
			location_req, err := http.NewRequest(http.MethodGet, location_api_url, nil)
			spaceClient := http.Client{}
			res, err := spaceClient.Do(location_req)
			body, err := ioutil.ReadAll(res.Body)
			temp_location := Location{}
			err = json.Unmarshal(body, &temp_location)
			if err != nil {
				panic(err.Error())
			}
			result_location := temp_location[0].Name
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

func (runtime Runtime) Merantix() (results Results) {
	start_url := "https://merantix.bamboohr.com/jobs/embed2.php?departmentId=0"
	type Job struct {
		Url      string
		Title    string
		Location string
		Division string
	}
	c := colly.NewCollector()
	c.OnHTML(".BambooHR-ATS-Department-List", func(e *colly.HTMLElement) {
		e.ForEach(".BambooHR-ATS-Department-Item", func(_ int, el *colly.HTMLElement) {
			result_division := strings.TrimSpace(el.ChildText(".BambooHR-ATS-Department-Header"))
			el.ForEach(".BambooHR-ATS-Jobs-Item", func(_ int, ell *colly.HTMLElement) {
				result_title := ell.ChildText("a")
				result_url := "https:" + ell.ChildAttr("a", "href")
				result_location := ell.ChildText("span")
				results.Add(
					runtime.Name,
					result_title,
					result_url,
					result_location,
					Job{
						result_url,
						result_title,
						result_location,
						result_division,
					},
				)
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

func (runtime Runtime) Ninox() (results Results) {
	start_url := "https://ninox.com/en/jobs"
	base_job_url := "https://ninox.com/%s"
	type Job struct {
		Url      string
		Title    string
		Location string
	}
	c := colly.NewCollector()
	c.OnHTML(".job-new", func(e *colly.HTMLElement) {
		result_url := fmt.Sprintf(base_job_url, e.ChildAttr("a", "href"))
		result_title := e.ChildText("h4")
		result_location := e.ChildText(".jobs-j-openinglugar")
		results.Add(
			runtime.Name,
			result_title,
			result_url,
			result_location,
			Job{
				result_url,
				result_title,
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

func (runtime Runtime) Zenjob() (results Results) {
	start_url := "https://zenjob.teamtailor.com"
	base_job_url := "https://zenjob.teamtailor.com%s"
	type Job struct {
		Url      string
		Title    string
		Location string
	}
	c := colly.NewCollector()
	c.OnHTML(".jobs", func(e *colly.HTMLElement) {
		e.ForEach("li", func(_ int, el *colly.HTMLElement) {
			result_url := fmt.Sprintf(base_job_url, el.ChildAttr("a", "href"))
			result_title := el.ChildText(".title")
			result_location := el.ChildText(".meta")
			results.Add(
				runtime.Name,
				result_title,
				result_url,
				result_location,
				Job{
					result_url,
					result_title,
					result_location,
				},
			)
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

func (runtime Runtime) Plantix() (results Results) {
	start_url := "https://plantix.recruitee.com/api/offers"
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
	c := colly.NewCollector()
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
			if fmt.Sprint(elem.Remote) == "true" {
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

func (runtime Runtime) Coachhub() (results Results) {
	start_url := "https://coachhub-jobs.personio.de/"
	type Job struct {
		Url      string
		Title    string
		Location string
	}
	c := colly.NewCollector()
	c.OnHTML(".panel-container", func(e *colly.HTMLElement) {
		e.ForEach(".recent-job-list", func(_ int, el *colly.HTMLElement) {
			result_url := el.ChildAttr("a", "href")
			result_title := el.ChildText("h6")
			result_location := strings.Split(el.ChildText("p"), "·")[1]
			results.Add(
				runtime.Name,
				result_title,
				result_url,
				result_location,
				Job{
					result_url,
					result_title,
					result_location,
				},
			)
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

func (runtime Runtime) Adjust() (results Results) {
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
	c := colly.NewCollector()
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

func (runtime Runtime) Bonify() (results Results) {
	start_url := "http://www.bonify.de/jobs"
	base_job_url := "https://www.bonify.de/jobs/%s"
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
	c := colly.NewCollector()
	c.OnResponse(func(r *colly.Response) {
		body := string(r.Body)
		json_body := strings.Split(
			strings.Split(
				body, `resultsAllJobsListingsTrimmed":`)[1], `,"resultsCompanyBenefits`)[0]
		var jsonJobs JsonJobs
		err := json.Unmarshal([]byte(json_body), &jsonJobs)
		if err != nil {
			panic(err.Error())
		}
		for _, elem := range jsonJobs.Results {
			result_title := elem.Data.Title[0].Text
			result_url := fmt.Sprintf(base_job_url, elem.UID)
			result_location := "Berlin"
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

func (runtime Runtime) Bryter() (results Results) {
	start_url := "https://bryter.io/careers"
	type Job struct {
		Title    string
		Url      string
		Location string
	}
	c := colly.NewCollector()
	c.OnHTML("#careers-listing", func(e *colly.HTMLElement) {
		e.ForEach("a", func(_ int, el *colly.HTMLElement) {
			result_url := el.Attr("href")
			result_title := el.ChildText("h4")
			result_location := "Berlin, Frankfurt, London"
			results.Add(
				runtime.Name,
				result_title,
				result_url,
				result_location,
				Job{
					result_url,
					result_title,
					result_location,
				},
			)
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

func (runtime Runtime) Bunch() (results Results) {
	results = append(results, Result{
		runtime.Name,
		"Freelance/Full-time Product Designer",
		"https://angel.co/company/bunch-hq/jobs/682927-freelance-full-time-product-designer",
		"New York City • Berlin • Remote",
		[]byte("{}"),
	})

	results = append(results, Result{
		runtime.Name,
		"Product Launch Intern (Internship)",
		"https://angel.co/company/bunch-hq/jobs/907192-product-launch-intern-internship",
		"Berlin • Remote",
		[]byte("{}"),
	})
	return
}

func (runtime Runtime) Candis() (results Results) {
	start_url := "https://career.recruitee.com/api/c/50731/widget/?widget=true"
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
	c := colly.NewCollector()
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

func (runtime Runtime) Bytedance() (results Results) {
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
	var jsonJobs Jobs
	err = json.Unmarshal(bodyText, &jsonJobs)
	if err != nil {
		panic(err.Error())
	}
	for _, elem := range jsonJobs.Data.JobPostList {
		result_title := elem.Title
		result_url := fmt.Sprintf(base_url, elem.ID)
		result_location := elem.CityInfo.EnName
		results.Add(
			runtime.Name,
			result_title,
			result_url,
			result_location,
			elem,
		)
	}
	return
}

func (runtime Runtime) Bmw() (results Results) {
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
	c := colly.NewCollector()
	c.OnResponse(func(r *colly.Response) {
		var jsonJobs JsonJobs
		err := json.Unmarshal(r.Body, &jsonJobs)
		if err != nil {
			panic(err.Error())
		}
		for _, elem := range jsonJobs.Data {
			result_title := elem.ReqTitle
			result_url := fmt.Sprintf(base_job_url, elem.JobDescriptionLink)
			result_location := elem.Location.Value
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

func (runtime Runtime) Infineon() (results Results) {
	start_url := "https://www.infineon.com%s"
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
	var jsonJobs Jobs
	err = json.Unmarshal(bodyText, &jsonJobs)
	if err != nil {
		panic(err.Error())
	}
	for _, elem := range jsonJobs.Pages.Items {
		result_title := elem.Title
		result_url := fmt.Sprintf(start_url, elem.DetailPageURL)
		result_location := elem.Location
		results.Add(
			runtime.Name,
			result_title,
			result_url,
			result_location,
			elem,
		)
	}
	return
}

func (runtime Runtime) Porsche() (results Results) {
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
	c := colly.NewCollector()
	c.OnResponse(func(r *colly.Response) {
		var jsonJobs JsonJobs
		err := json.Unmarshal(r.Body, &jsonJobs)
		if err != nil {
			panic(err.Error())
		}
		for _, elem := range jsonJobs.SearchResult.SearchResultItems {
			result_title := elem.MatchedObjectDescriptor.PositionTitle
			result_url := elem.MatchedObjectDescriptor.PositionURI
			result_location := elem.MatchedObjectDescriptor.PositionLocation[0].CityName + "," + elem.MatchedObjectDescriptor.PositionLocation[0].CountryName
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

func (runtime Runtime) Bosch() (results Results) {
	start_url := "https://api.smartrecruiters.com/v1/companies/BoschGroup/postings?offset=%d"
	base_job_url := "https://www.smartrecruiters.com/BoschGroup/%s"
	number_results_per_page := 100
	type Jobs struct {
		Content []struct {
			Company struct {
				Identifier string `json:"identifier"`
				Name       string `json:"name"`
			} `json:"company"`
			Creator struct {
				Name string `json:"name"`
			} `json:"creator"`
			CustomField []struct {
				FieldID    string `json:"fieldId"`
				FieldLabel string `json:"fieldLabel"`
				ValueID    string `json:"valueId"`
				ValueLabel string `json:"valueLabel"`
			} `json:"customField"`
			Department      struct{} `json:"department"`
			ExperienceLevel struct {
				ID    string `json:"id"`
				Label string `json:"label"`
			} `json:"experienceLevel"`
			Function struct {
				ID    string `json:"id"`
				Label string `json:"label"`
			} `json:"function"`
			ID       string `json:"id"`
			Industry struct {
				ID    string `json:"id"`
				Label string `json:"label"`
			} `json:"industry"`
			Language struct {
				Code        string `json:"code"`
				Label       string `json:"label"`
				LabelNative string `json:"labelNative"`
			} `json:"language"`
			Location struct {
				Address    string `json:"address"`
				City       string `json:"city"`
				Country    string `json:"country"`
				PostalCode string `json:"postalCode"`
				Region     string `json:"region"`
				Remote     bool   `json:"remote"`
			} `json:"location"`
			Name             string `json:"name"`
			Ref              string `json:"ref"`
			RefNumber        string `json:"refNumber"`
			ReleasedDate     string `json:"releasedDate"`
			TypeOfEmployment struct {
				Label string `json:"label"`
			} `json:"typeOfEmployment"`
			UUID string `json:"uuid"`
		} `json:"content"`
		Limit      int64 `json:"limit"`
		Offset     int64 `json:"offset"`
		TotalFound int64 `json:"totalFound"`
	}
	c := colly.NewCollector()
	c.OnResponse(func(r *colly.Response) {
		var jsonJobs Jobs
		err := json.Unmarshal(r.Body, &jsonJobs)
		if err != nil {
			panic(err.Error())
		}
		for _, elem := range jsonJobs.Content {
			result_title := elem.Name
			result_url := fmt.Sprintf(base_job_url, elem.ID)
			result_location := elem.Location.City + "," + elem.Location.Country
			results.Add(
				runtime.Name,
				result_title,
				result_url,
				result_location,
				elem,
			)
		}
		total_matches := int(jsonJobs.TotalFound)
		total_pages := total_matches / number_results_per_page
		for i := 1; i <= total_pages; i++ {
			time.Sleep(SecondsSleep * time.Second)
			c.Visit(fmt.Sprintf(start_url, number_results_per_page*i))
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

func (runtime Runtime) Mckinsey() (results Results) {
	start_url := "https://mobileservices.mckinsey.com/services/ContentAPI/SearchAPI.svc/jobs/search?&pageSize=100&start=%d"
	base_job_url := "https://www.mckinsey.com/careers/search-jobs/jobs/%s"
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
	c := colly.NewCollector()
	c.OnResponse(func(r *colly.Response) {
		var jsonJobs JsonJobs
		err := json.Unmarshal(r.Body, &jsonJobs)
		if err != nil {
			panic(err.Error())
		}
		for _, elem := range jsonJobs.Response.Docs {
			result_title := elem.Title
			result_url := fmt.Sprintf(base_job_url, elem.FriendlyURL)
			cities := RemoveDuplicatedFromSliceOfString(elem.Cities)
			countries := RemoveDuplicatedFromSliceOfString(elem.Countries)
			result_location := strings.Join(cities, ",") + "-" + strings.Join(countries, ",")
			results.Add(
				runtime.Name,
				result_title,
				result_url,
				result_location,
				elem,
			)
		}
		total_pages := jsonJobs.Response.NumFound / number_results_per_page
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
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println(Red("Request URL:"), Red(r.Request.URL))
	})
	c.Visit(fmt.Sprintf(start_url, 0))
	return
}

func (runtime Runtime) Sap() (results Results) {
	start_url := "https://jobs.sap.com/search/?q=&sortColumn=referencedate&sortDirection=desc&startrow=%d"
	base_job_url := "https://jobs.sap.com/%s"
	number_results_per_page := 25
	counter := 0
	type Job struct {
		Title    string
		Url      string
		Location string
	}
	c := colly.NewCollector()
	c.OnHTML(".html5", func(e *colly.HTMLElement) {
		e.ForEach(".data-row", func(_ int, el *colly.HTMLElement) {
			result_url := fmt.Sprintf(base_job_url, el.ChildAttr("a", "href"))
			result_title := el.ChildText("a")
			result_location := el.ChildText(".jobLocation.visible-phone")
			results.Add(
				runtime.Name,
				result_title,
				result_url,
				result_location,
				Job{
					result_url,
					result_title,
					result_location,
				},
			)
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
	c.OnRequest(func(r *colly.Request) {
		fmt.Println(Gray(8-1, "Visiting"), Gray(8-1, r.URL.String()))
	})
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println(Red("Request URL:"), Red(r.Request.URL))
	})
	c.Visit(fmt.Sprintf(start_url, 0))
	return
}

func (runtime Runtime) Puma() (results Results) {
	start_url := "https://about.puma.com/api/PUMA/Feature/JobFinder?loadMore=500"
	base_job_url := "https://about.puma.com%s"
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
	c := colly.NewCollector()
	c.OnResponse(func(r *colly.Response) {
		var jsonJobs JsonJobs
		err := json.Unmarshal(r.Body, &jsonJobs)
		if err != nil {
			panic(err.Error())
		}
		for _, elem := range jsonJobs.Teaser {
			result_title := elem.Title
			result_url := fmt.Sprintf(base_job_url, elem.URL)
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

func (runtime Runtime) Daimler() (results Results) {
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
	c := colly.NewCollector()
	c.OnResponse(func(r *colly.Response) {
		var jsonJobs JsonJobs
		err := json.Unmarshal(r.Body, &jsonJobs)
		if err != nil {
			panic(err.Error())
		}
		for _, elem := range jsonJobs.SearchResult.SearchResultItems {
			result_title := elem.MatchedObjectDescriptor.PositionTitle
			result_url := elem.MatchedObjectDescriptor.PositionURI
			result_location := elem.MatchedObjectDescriptor.PositionLocation[0].CityName
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

func (runtime Runtime) Siemens() (results Results) {
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
	c := colly.NewCollector()
	c.OnResponse(func(r *colly.Response) {
		var jsonJobs JsonJobs
		err := json.Unmarshal(r.Body, &jsonJobs)
		if err != nil {
			panic(err.Error())
		}

		for _, elem := range jsonJobs.Jobs {

			result_title := elem.Data.Title
			result_url := elem.Data.MetaData.CanonicalURL
			result_location := elem.Data.City + "," + elem.Data.Country
			results.Add(
				runtime.Name,
				result_title,
				result_url,
				result_location,
				elem,
			)
		}
		total_pages := jsonJobs.TotalCount / number_results_per_page
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
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println(Red("Request URL:"), Red(r.Request.URL))
	})
	c.Visit(fmt.Sprintf(start_url, 1))
	return
}

func (runtime Runtime) Continental() (results Results) {
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
	c := colly.NewCollector()
	c.OnResponse(func(r *colly.Response) {
		var jsonJobs JsonJobs
		err := json.Unmarshal(r.Body, &jsonJobs)
		if err != nil {
			panic(err.Error())
		}
		for _, elem := range jsonJobs.SearchResult.SearchResultItems {
			result_title := elem.MatchedObjectDescriptor.PositionTitle
			result_url := elem.MatchedObjectDescriptor.PositionURI
			var result_location string
			if len(elem.MatchedObjectDescriptor.PositionLocation) == 0 {
				result_location = ""
			} else {
				result_location = elem.MatchedObjectDescriptor.PositionLocation[0].CityName + "," + elem.MatchedObjectDescriptor.PositionLocation[0].CountryName
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

func (runtime Runtime) Deliveryhero() (results Results) {
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
	c := colly.NewCollector()
	c.OnResponse(func(r *colly.Response) {
		response := Response{r.Body}
		response_body := string(response.Html)
		response_json := strings.Split(
			strings.Split(
				response_body, "phApp.ddo = ")[1], "; phApp.experimentData")[0]
		var jsonJobs JsonJobs
		err := json.Unmarshal([]byte(response_json), &jsonJobs)
		if err != nil {
			panic(err.Error())
		}
		for _, elem := range jsonJobs.EagerLoadRefineSearch.Data.Jobs {
			result_title := elem.Title
			result_url := fmt.Sprintf(base_job_url, elem.JobID)
			result_location := elem.MultiLocationArray[0].Location
			results.Add(
				runtime.Name,
				result_title,
				result_url,
				result_location,
				elem,
			)
		}
		total_pages := jsonJobs.EagerLoadRefineSearch.TotalHits / number_results_per_page
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
		fmt.Println(Red("Request URL:"), Red(r.Request.URL))
	})
	c.Visit(fmt.Sprintf(start_url, 0))
	return
}

func (runtime Runtime) Volkswagen() (results Results) {
	start_url := "https://karriere.volkswagen.de/sap/bc/bsp/sap/zvw_hcmx_ui_ext/desktop.html#/SEARCH/RESULTS"
	base_url := "https://karriere.volkswagen.de/sap/bc/bsp/sap/zvw_hcmx_ui_ext/?jobId=%s"
	file_name := "volkswagen.html"
	type Job struct {
		Title      string
		Url        string
		Location   string
		Department string
	}
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()
	var initialResponse string
	var res []byte
	if err := chromedp.Run(ctx,
		chromedp.Navigate(start_url),
		chromedp.WaitVisible(".details"),
		chromedp.EvaluateAsDevTools(`for (var i = 0; i < 20; i++) {document.getElementsByClassName("button more showMore")[document.getElementsByClassName("button more showMore").length -1].click();}`, &res),
		chromedp.Sleep(SecondsSleep*time.Second),
		chromedp.OuterHTML("html", &initialResponse),
	); err != nil {
		panic(err)
	}
	SaveResponseToFileWithFileName(initialResponse, file_name)
	c := colly.NewCollector()
	c.OnHTML("html", func(e *colly.HTMLElement) {
		e.ForEach(".listItem", func(_ int, el *colly.HTMLElement) {
			result_title := el.ChildText(".title")
			result_url := fmt.Sprintf(base_url, el.ChildAttr("div", "data-id"))
			result_location := el.ChildText(".locationPrimary")
			result_department := el.ChildText(".functionalArea")
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

func (runtime Runtime) Tesla() (results Results) {
	start_url := "https://www.tesla.com/de_DE/careers/search#/"
	base_job_url := "https://www.tesla.com/careers/%s"
	file_name := "tesla.html"
	type Job struct {
		Title      string
		Url        string
		Location   string
		Department string
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
			result_url := fmt.Sprintf(base_job_url, el.ChildAttr("a", "href"))
			result_department := el.ChildText(".listing-department")
			result_location := el.ChildText(".listing-location")
			result_date := el.ChildText(".listing-dateposted")
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
					result_date,
				},
			)
		})
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

func (runtime Runtime) Researchgate() (results Results) {
	start_url := "https://www.researchgate.net/jobs?page=%d"
	base_job_url := "https://www.researchgate.net/%s"
	counter := 1
	type Job struct {
		Title     string
		Url       string
		Location  string
		Institute string
	}
	c := colly.NewCollector()
	c.OnHTML("html", func(e *colly.HTMLElement) {
		e.ForEach(".jobs-list-item-nova", func(_ int, el *colly.HTMLElement) {
			result_title := el.ChildText(".nova-v-job-item__title")
			result_url := fmt.Sprintf(base_job_url, el.ChildAttr("a", "href"))
			result_infos := el.ChildTexts(".nova-v-job-item__info-section-list-item")
			result_institute := result_infos[0]
			result_location := result_infos[1]
			results.Add(
				runtime.Name,
				result_title,
				result_url,
				result_location,
				Job{
					result_title,
					result_url,
					result_location,
					result_institute,
				},
			)
		})
		page_links := e.ChildAttrs(".pager-link", "data-target-page")
		temp_total_pages := page_links[len(page_links)-2]
		total_pages, _ := strconv.Atoi(temp_total_pages)
		if counter <= total_pages {
			counter++
			time.Sleep(SecondsSleep * time.Second)
			c.Visit(fmt.Sprintf(start_url, counter))
		}
	})
	c.OnRequest(func(r *colly.Request) {
		fmt.Println(Gray(8-1, "Visiting"), Gray(8-1, r.URL.String()))
	})
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println(Red("Request URL:"), Red(r.Request.URL))
	})
	c.Visit(fmt.Sprintf(start_url, 1))
	return
}

func (runtime Runtime) Clue() (results Results) {
	start_url := "https://api.greenhouse.io/v1/boards/clue/jobs"
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
	c := colly.NewCollector()
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

func (runtime Runtime) Adahealth() (results Results) {
	start_url := "https://api.greenhouse.io/v1/boards/adahealth/jobs"
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
	c := colly.NewCollector()
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

func (runtime Runtime) Babbel() (results Results) {
	start_url := "https://api.greenhouse.io/v1/boards/babbel/jobs"
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
	c := colly.NewCollector()
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

func (runtime Runtime) Eyeem() (results Results) {
	start_url := "https://www.eyeem.com/jobs"
	file_name := "eyeem.html"
	type Job struct {
		Title    string
		Url      string
		Location string
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
		e.ForEach(".collection-item-job", func(_ int, el *colly.HTMLElement) {
			result_title := el.ChildText(".bold-s-18")
			result_url := el.ChildAttr("a", "href")
			result_location := el.ChildText(".jobs")
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

func (runtime Runtime) Rocketinternet() (results Results) {
	start_url := "https://www.rocket-internet.com/careers/rocket"
	base_job_url := "https://www.rocket-internet.com%s"
	type Job struct {
		Url      string
		Title    string
		Location string
		Type     string
	}
	c := colly.NewCollector()
	l := c.Clone()
	c.OnHTML(".department", func(e *colly.HTMLElement) {
		department_url := fmt.Sprintf(base_job_url, e.ChildAttr("a", "href"))
		l.Visit(department_url)
	})
	c.OnRequest(func(r *colly.Request) {
		fmt.Println(Gray(8-1, "Visiting"), Gray(8-1, r.URL.String()))
	})
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println(Red("Request URL:"), Red(r.Request.URL))
	})
	l.OnHTML("#careers-listing", func(e *colly.HTMLElement) {
		e.ForEach("div[role=listitem]", func(_ int, el *colly.HTMLElement) {
			result_url := fmt.Sprintf(base_job_url, el.ChildAttr("a", "href"))
			result_info := el.ChildTexts(".text-sans-serif")
			result_title := result_info[0]
			result_type := result_info[1]
			result_location := result_info[2]
			results.Add(
				runtime.Name,
				result_title,
				result_url,
				result_location,
				Job{
					result_url,
					result_title,
					result_location,
					result_type,
				},
			)
		})
	})
	l.OnRequest(func(r *colly.Request) {
		fmt.Println(Gray(8-1, "Visiting"), Gray(8-1, r.URL.String()))
	})
	l.OnError(func(r *colly.Response, err error) {
		fmt.Println(Red("Request URL:"), Red(r.Request.URL))
	})
	c.Visit(start_url)
	return
}

func (runtime Runtime) Here() (results Results) {
	start_url := "https://careers-here.icims.com/jobs/search?in_iframe=1&pr=%d"
	counter := 0
	type Job struct {
		Url         string
		Title       string
		Location    string
		Description string
	}
	c := colly.NewCollector()
	c.OnHTML("html", func(e *colly.HTMLElement) {
		e.ForEach(".row", func(_ int, el *colly.HTMLElement) {
			result_title := strings.TrimSpace(strings.ReplaceAll(el.ChildText(".title"), "Requisition Title", ""))
			result_url := el.ChildAttr("a", "href")
			result_location := strings.TrimSpace(strings.ReplaceAll(el.ChildText(".left"), "Job Locations", ""))
			result_description := el.ChildText(".description")
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
		result_pages := e.ChildTexts(".iCIMS_PagingBatch .sr-only")
		re := regexp.MustCompile("[0-9]+")
		temp_total_pages := re.FindAllString(result_pages[len(result_pages)-1], -1)[0]
		total_pages, _ := strconv.Atoi(temp_total_pages)
		if counter < (total_pages - 1) {
			counter++
			time.Sleep(SecondsSleep * time.Second)
			temp_url := fmt.Sprintf(start_url, counter)
			c.Visit(temp_url)
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

func (runtime Runtime) Fineway() (results Results) {
	start_url := "https://fineway.recruitee.com/api/offers"
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
	c := colly.NewCollector()
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
			if fmt.Sprint(elem.Remote) == "true" {
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

func (runtime Runtime) Limehome() (results Results) {
	start_url := "https://career.limehome.de/"
	base_job_url := "https://career.limehome.de%s"
	type Job struct {
		Url        string
		Title      string
		Location   string
		Department string
	}
	c := colly.NewCollector()
	c.OnHTML(".col-md-6", func(e *colly.HTMLElement) {
		result_url := fmt.Sprintf(base_job_url, e.Attr("href"))
		result_title := e.ChildText(".title")
		result_location := e.ChildText(".location")
		result_department := e.ChildTexts(".department")[0]
		results.Add(
			runtime.Name,
			result_title,
			result_url,
			result_location,
			Job{
				result_url,
				result_title,
				result_location,
				result_department,
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

func (runtime Runtime) Sonomotors() (results Results) {
	start_url := "https://sonomotors.com/de/career/"
	base_job_url := "https://sonomotors.com%s"
	type Job struct {
		Url      string
		Title    string
		Location string
	}
	c := colly.NewCollector()
	c.OnHTML(".job", func(e *colly.HTMLElement) {
		result_url := fmt.Sprintf(base_job_url, e.Attr("href"))
		result_title := e.ChildText("h3")
		result_location := e.ChildText("p")
		results.Add(
			runtime.Name,
			result_title,
			result_url,
			result_location,
			Job{
				result_url,
				result_title,
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

func (runtime Runtime) Twaice() (results Results) {
	start_url := "https://twaice.com/jobs"
	type Job struct {
		Url      string
		Title    string
		Location string
		Type     string
	}
	c := colly.NewCollector()
	c.OnHTML(".shadow-md", func(e *colly.HTMLElement) {
		result_url := e.Attr("href")
		result_title := e.ChildText("h3")
		result_location := "Munich"
		result_type := e.ChildText(".mb-0")
		results.Add(
			runtime.Name,
			result_title,
			result_url,
			result_location,
			Job{
				result_url,
				result_title,
				result_location,
				result_type,
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

func (runtime Runtime) Lilium() (results Results) {
	start_url := "https://api.greenhouse.io/v1/boards/lilium/jobs"
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
	c := colly.NewCollector()
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

func (runtime Runtime) Volocopter() (results Results) {
	start_url := "https://api.smartrecruiters.com/v1/companies/VolocopterGmbH/postings"
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
	c := colly.NewCollector()
	c.OnResponse(func(r *colly.Response) {
		var jsonJobs JsonJobs
		err := json.Unmarshal(r.Body, &jsonJobs)
		if err != nil {
			panic(err.Error())
		}
		for _, elem := range jsonJobs.Content {
			result_title := elem.Name
			result_url := elem.Ref
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

func (runtime Runtime) Idnow() (results Results) {
	start_url := "https://idnow.jobbase.io"
	base_url := "https://idnow.jobbase.io%s"
	file_name := "idnow.html"
	type Job struct {
		Url      string
		Title    string
		Location string
		Date     string
	}
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()
	var initialResponse string
	if err := chromedp.Run(ctx,
		chromedp.Navigate(start_url),
		chromedp.Sleep(SecondsSleep*time.Second),
		chromedp.OuterHTML("html", &initialResponse),
	); err != nil {
		panic(err)
	}
	SaveResponseToFileWithFileName(initialResponse, file_name)
	c := colly.NewCollector()
	c.OnHTML(".row-table-condensed-md", func(e *colly.HTMLElement) {
		result_title := e.ChildText("a")
		result_url := fmt.Sprintf(base_url, e.ChildAttr("a", "href"))
		result_location := e.ChildTexts(".cell-table")[1]
		result_date := e.ChildTexts(".cell-table")[2]
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

func (runtime Runtime) Demodesk() (results Results) {
	start_url := "https://wordpress.demodesk.com/wp-json/wp/v2/jobs"
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
		FeaturedMedia       int         `json:"featured_media"`
		Template            string      `json:"template"`
		BetterFeaturedImage interface{} `json:"better_featured_image"`
		Acf                 struct {
			Department  string `json:"department"`
			JobType     string `json:"job_type"`
			WorkingTime string `json:"working_time"`
			Location    string `json:"location"`
			Profile     string `json:"profile"`
			Blurb       string `json:"blurb"`
			Order       string `json:"order"`
		} `json:"acf"`
		Links struct {
			Self []struct {
				Href string `json:"href"`
			} `json:"self"`
			Collection []struct {
				Href string `json:"href"`
			} `json:"collection"`
			About []struct {
				Href string `json:"href"`
			} `json:"about"`
			WpAttachment []struct {
				Href string `json:"href"`
			} `json:"wp:attachment"`
			Curies []struct {
				Name      string `json:"name"`
				Href      string `json:"href"`
				Templated bool   `json:"templated"`
			} `json:"curies"`
		} `json:"_links"`
	}
	c := colly.NewCollector()
	c.OnResponse(func(r *colly.Response) {
		var jobs Jobs
		err := json.Unmarshal(r.Body, &jobs)
		if err != nil {
			panic(err.Error())
		}
		for _, elem := range jobs {
			result_title := elem.Title.Rendered
			result_url := elem.Link
			result_location := elem.Acf.Location
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

func (runtime Runtime) Tado() (results Results) {
	start_url := "https://apply.workable.com/api/v3/accounts/tado/jobs"
	base_job_url := "https://apply.workable.com/tado/j/%s"
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
	c := colly.NewCollector()
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

func (runtime Runtime) Holidu() (results Results) {
	start_url := "https://api.holidu.com/api/careers"
	type JsonJobs struct {
		Positions []struct {
			ID                 string `json:"id"`
			Subcompany         string `json:"subcompany,omitempty"`
			Office             string `json:"office"`
			Department         string `json:"department"`
			RecruitingCategory string `json:"recruitingCategory"`
			Name               string `json:"name"`
			JobDescriptions    struct {
				JobDescription []struct {
					Name  string `json:"name"`
					Value string `json:"value"`
				} `json:"jobDescription"`
			} `json:"jobDescriptions"`
			EmploymentType     string    `json:"employmentType"`
			Seniority          string    `json:"seniority"`
			Schedule           string    `json:"schedule"`
			YearsOfExperience  string    `json:"yearsOfExperience"`
			Occupation         string    `json:"occupation"`
			OccupationCategory string    `json:"occupationCategory"`
			CreatedAt          time.Time `json:"createdAt"`
		} `json:"positions"`
	}
	c := colly.NewCollector()
	c.OnResponse(func(r *colly.Response) {
		var jsonJobs JsonJobs
		err := json.Unmarshal(r.Body, &jsonJobs)
		if err != nil {
			panic(err.Error())
		}
		for _, elem := range jsonJobs.Positions {
			result_title := elem.Name
			fmt.Println(result_title)
			result_url := "https://www.holidu.com/careers?" + elem.ID
			result_location := elem.Office
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

func (runtime Runtime) Westwing() (results Results) {
	start_url := "https://www.westwing.com/de/career/"
	base_job_url := "https://www.westwing.com%s"
	type Job struct {
		Url      string
		Title    string
		Location string
		Type     string
	}
	c := colly.NewCollector()
	c.OnHTML(".js-job", func(e *colly.HTMLElement) {
		result_title := e.ChildText(".block__jobs-title")
		result_url := fmt.Sprintf(base_job_url, e.ChildAttr("a", "href"))
		result_location := e.ChildTexts(".block__jobs-detail")[0]
		result_type := e.ChildTexts(".block__jobs-detail")[1]
		results.Add(
			runtime.Name,
			result_title,
			result_url,
			result_location,
			Job{
				result_title,
				result_url,
				result_location,
				result_type,
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

func (runtime Runtime) Combyne() (results Results) {
	start_url := "https://combyne.recruitee.com/api/offers"
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
	c := colly.NewCollector()
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
			if fmt.Sprint(elem.Remote) == "true" {
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

func (runtime Runtime) Globalsavingsgroup() (results Results) {
	start_url := "https://api.greenhouse.io/v1/boards/globalsavingsgroup/jobs"
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
	c := colly.NewCollector()
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

func (runtime Runtime) Riskmethods() (results Results) {
	start_url := "https://apply.workable.com/api/v3/accounts/riskmethods/jobs"
	base_job_url := "https://apply.workable.com/riskmethods/j/%s"
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
	c := colly.NewCollector()
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

func (runtime Runtime) Mylivn() (results Results) {
	start_url := "https://mylivn.com/jobs"
	base_job_url := "https://mylivn.com%s"
	type Job struct {
		Title    string
		Url      string
		Location string
		Type     string
	}
	c := colly.NewCollector()
	c.OnHTML(".positions-module--position--23wcR", func(e *colly.HTMLElement) {
		result_title := e.ChildText(".positions-module--positionTitle--i2U3O")
		result_url := fmt.Sprintf(base_job_url, e.ChildAttr("a", "href"))
		result_type := e.ChildText(".positions-module--positionTime--38eR8")
		result_location := e.ChildText(".positions-module--positionLocation--mtUS9")
		results.Add(
			runtime.Name,
			result_title,
			result_url,
			result_location,
			Job{
				result_title,
				result_url,
				result_location,
				result_type,
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

func (runtime Runtime) Teleclinic() (results Results) {
	start_url := "https://teleclinic-gmbh.breezy.hr%s"
	type Job struct {
		Title      string
		Url        string
		Location   string
		Department string
		Type       string
	}
	c := colly.NewCollector()
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

func (runtime Runtime) Shore() (results Results) {
	start_url := "https://www.shore.com/en/career/#on-apply"
	type Job struct {
		Title      string
		Url        string
		Location   string
		Department string
	}
	c := colly.NewCollector()
	c.OnHTML(".job-opening-list-element", func(e *colly.HTMLElement) {
		result_title := strings.TrimSpace(e.ChildText(".job-title"))
		result_url := strings.TrimSpace(e.ChildAttr("a", "href"))
		result_location := strings.Join(strings.Fields(strings.TrimSpace(e.ChildText(".job-location-and-type"))), " ")
		result_department := e.ChildText(".job-department column")
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
	c.OnRequest(func(r *colly.Request) {
		fmt.Println(Gray(8-1, "Visiting"), Gray(8-1, r.URL.String()))
	})
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println(Red("Request URL:"), Red(r.Request.URL))
	})
	c.Visit(start_url)
	return
}

func (runtime Runtime) Inmindcloud() (results Results) {
	url := "https://www.inmindcloud.com/about-us/career/"
	type Job struct {
		Title       string
		Url         string
		Type        string
		Description string
	}
	c := colly.NewCollector()
	c.OnHTML(".post", func(e *colly.HTMLElement) {
		result_url := e.ChildAttr("a", "href")
		result_info := e.ChildText("h4")
		result_title := strings.Split(result_info, ",")[0]
		result_location := strings.Split(result_info, ",")[1]
		result_description := e.ChildText("p")
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
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println(Red("Request URL:"), Red(r.Request.URL))
	})
	c.Visit(url)
	return
}

func (runtime Runtime) Alyne() (results Results) {
	start_url := "https://api.greenhouse.io/v1/boards/Alyne/jobs"
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
	c := colly.NewCollector()
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

func (runtime Runtime) Smartreporting() (results Results) {
	start_url := "https://apply.workable.com/api/v3/accounts/smartreporting/jobs"
	base_job_url := "https://apply.workable.com/smartreporting/j/%s"
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
	c := colly.NewCollector()
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

func (runtime Runtime) Censhare() (results Results) {
	start_url := "https://www.censhare.com/company/careers"
	base_job_url := "https://www.censhare.com%s"
	type Job struct {
		Title      string
		Url        string
		Department string
		Location   string
	}
	c := colly.NewCollector()
	c.OnHTML("div[class=csCard]", func(e *colly.HTMLElement) {
		if e.Attr("data-cid") != "" {
			result_title := e.ChildText(".csCard__title")
			result_url := fmt.Sprintf(base_job_url, e.ChildAttr("a", "href"))
			result_location := e.ChildTexts("p")[1]
			result_department := e.ChildTexts("p")[0]
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

func (runtime Runtime) Stylight() (results Results) {
	start_url := "https://about.stylight.com/jobs"
	type Job struct {
		Title    string
		Url      string
		Location string
	}
	c := colly.NewCollector()
	c.OnHTML(".post-list", func(e *colly.HTMLElement) {
		result_title := e.ChildAttr("a", "title")
		result_url := e.ChildAttr("a", "href")
		result_location := "Munich"
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

func (runtime Runtime) Ryte() (results Results) {
	start_url := "https://en.ryte.com/jobs"
	base_job_url := "https://en.ryte.com/jobs/%s"
	type Job struct {
		Title    string
		Url      string
		Location string
		Type     string
	}
	c := colly.NewCollector()
	c.OnHTML(".card", func(e *colly.HTMLElement) {
		result_title := e.ChildText(".job_name")
		result_url := fmt.Sprintf(base_job_url, e.ChildAttr("a", "href"))
		result_location := e.ChildTexts(".details_item")[2]
		result_type := e.ChildTexts(".details_item")[0]
		results.Add(
			runtime.Name,
			result_title,
			result_url,
			result_location,
			Job{
				result_title,
				result_url,
				result_location,
				result_type,
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

func (runtime Runtime) Codasip() (results Results) {
	start_url := "https://codasip.bamboohr.com/jobs/embed2.php?departmentId=0"
	type Job struct {
		Url      string
		Title    string
		Location string
		Division string
	}
	c := colly.NewCollector()
	c.OnHTML(".BambooHR-ATS-Department-List", func(e *colly.HTMLElement) {
		e.ForEach(".BambooHR-ATS-Department-Item", func(_ int, el *colly.HTMLElement) {
			result_division := strings.TrimSpace(el.ChildText(".BambooHR-ATS-Department-Header"))
			el.ForEach(".BambooHR-ATS-Jobs-Item", func(_ int, ell *colly.HTMLElement) {
				result_title := ell.ChildText("a")
				result_url := "https:" + ell.ChildAttr("a", "href")
				result_location := ell.ChildText("span")
				results.Add(
					runtime.Name,
					result_title,
					result_url,
					result_location,
					Job{
						result_url,
						result_title,
						result_location,
						result_division,
					},
				)
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

func (runtime Runtime) Speexx() (results Results) {
	start_url := "https://apply.workable.com/api/v3/accounts/speexx/jobs"
	base_job_url := "https://apply.workable.com/speexx/j/%s"
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
	c := colly.NewCollector()
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

func (runtime Runtime) Spendit() (results Results) {
	start_url := "https://spendit.recruitee.com/api/offers"
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
	c := colly.NewCollector()
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
			if fmt.Sprint(elem.Remote) == "true" {
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

func (runtime Runtime) Allianz() (results Results) {
	start_url := "https://jobs.allianz.com/sap/hcmx/hitlist_na"
	base_job_url := "https://jobs.allianz.com/sap/bc/bsp/sap/zhcmx_erc_ui_ex/desktop.html?jobId=%s"
	location_url := "https://jobs.allianz.com/sap/opu/odata/hcmx/erc_ui_open_srv/LocationSet"
	type JsonJobs struct {
		D struct {
			Results []struct {
				ApplicationEndDate string `json:"ApplicationEndDate"`
				JobID              string `json:"JobID"`
				Posting            int    `json:"Posting"`
				Title              string `json:"Title"`
				PostingAge         int    `json:"PostingAge"`
				JobDetailsURL      string `json:"JobDetailsUrl"`
				TravelRatio        int    `json:"TravelRatio"`
				Company            struct {
					CompanyID int    `json:"CompanyID"`
					Text      string `json:"Text"`
					LogoURL   string `json:"LogoURL"`
				} `json:"Company"`
				FunctionalArea struct {
					FunctionalAreaID int    `json:"FunctionalAreaID"`
					Text             string `json:"Text"`
					Tooltip          string `json:"Tooltip"`
				} `json:"FunctionalArea"`
				ContractType struct {
					ContractTypeID int    `json:"ContractTypeID"`
					Text           string `json:"Text"`
					Tooltip        string `json:"Tooltip"`
				} `json:"ContractType"`
				HierarchyLevel struct {
					HierarchyLevelID int    `json:"HierarchyLevelID"`
					Text             string `json:"Text"`
					Tooltip          string `json:"Tooltip"`
				} `json:"HierarchyLevel"`
				Location struct {
					LocationID       int    `json:"LocationID"`
					Text             string `json:"Text"`
					Latitude         string `json:"Latitude"`
					Longitude        string `json:"Longitude"`
					ParentLocationID int    `json:"ParentLocationID"`
					Type             int    `json:"Type"`
					Adm1Code         string `json:"Adm1Code"`
				} `json:"Location"`
			} `json:"results"`
		} `json:"d"`
	}
	type Location struct {
		Entry []struct {
			Content struct {
				Properties struct {
					Text       string `xml:"Text"`
					LocationID string `xml:"LocationID"`
				} `xml:"properties"`
			} `xml:"content"`
		} `xml:"entry"`
	}
	c := colly.NewCollector()
	l := colly.NewCollector()
	var xmlLocation Location
	c.OnResponse(func(r *colly.Response) {
		var jsonJobs JsonJobs
		err := json.Unmarshal(r.Body, &jsonJobs)
		if err != nil {
			panic(err.Error())
		}
		l.Visit(location_url)
		for _, elem := range jsonJobs.D.Results {
			result_title := elem.Title
			result_url := fmt.Sprintf(base_job_url, elem.JobID)
			var result_location string
			for _, v := range xmlLocation.Entry {
				if v.Content.Properties.LocationID == strconv.Itoa(elem.Location.LocationID) {
					result_location = v.Content.Properties.Text
				}
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
	l.OnResponse(func(r *colly.Response) {
		err := xml.Unmarshal(r.Body, &xmlLocation)
		if err != nil {
			panic(err.Error())
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

func (runtime Runtime) Tier() (results Results) {
	start_url := "https://tier-mobility-jobs.personio.de/"
	type Job struct {
		Url      string
		Title    string
		Location string
	}
	c := colly.NewCollector()
	c.OnHTML(".panel-container", func(e *colly.HTMLElement) {
		e.ForEach(".recent-job-list", func(_ int, el *colly.HTMLElement) {
			result_url := el.ChildAttr("a", "href")
			result_title := el.ChildText("h6")
			result_location := strings.Split(el.ChildText("p"), "·")[1]
			results.Add(
				runtime.Name,
				result_title,
				result_url,
				result_location,
				Job{
					result_url,
					result_title,
					result_location,
				},
			)
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

func (runtime Runtime) Uniper() (results Results) {
	start_url := "https://jobs.uniper.energy/search"
	new_page_url := "https://jobs.uniper.energy/tile-search-results/?startrow=%d"
	base_job_url := "https://jobs.uniper.energy"
	number_results_per_page := 10
	total_pages := 0
	counter := 0
	type Job struct {
		Title    string
		Url      string
		Location string
	}
	c := colly.NewCollector()
	c.OnHTML("html", func(e *colly.HTMLElement) {
		e.ForEach(".sub-section", func(_ int, el *colly.HTMLElement) {
			result_title := el.ChildText("a")
			result_url := fmt.Sprintf(base_job_url, el.ChildAttr("a", "href"))
			result_location := strings.TrimSpace(strings.ReplaceAll(el.ChildText(".location"), "Location", ""))
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

		if counter == 0 {
			temp_total_results := strings.Split(e.ChildText("#tile-search-results-label"), " ")
			string_total_results := temp_total_results[len(temp_total_results)-2]
			total_results, err := strconv.Atoi(string_total_results)
			if err != nil {
				panic(err.Error())
			}
			total_pages = total_results / number_results_per_page
		}
		if counter >= total_pages {
			return
		} else {
			time.Sleep(SecondsSleep * time.Second)
			temp_v_url := fmt.Sprintf(new_page_url, 10+counter*number_results_per_page)
			counter++
			c.Visit(temp_v_url)
		}
	})

	c.OnResponse(func(r *colly.Response) {
		SaveResponseToFileWithFileName(string(r.Body), "Uniper.html")
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
