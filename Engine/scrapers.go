package main

import (
	"fmt"
	"github.com/gocolly/colly"
	"io/ioutil"
	"net/http"
	"reflect"
    "strings"
    "encoding/json"
    "time"
)

type Runtime struct {
	Name string
}
type Job struct {
	CompanyName string
	CompanyUrl  string
	Title       string
	JobUrl      string
}

func runner(name string) (job []Job) {
	r := Runtime{name}
	v := reflect.ValueOf(r)
	m := v.MethodByName(r.Name)
	temp_job := m.Call(nil)
	job = temp_job[0].Interface().([]Job)
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

func (runtime Runtime) Kununu() (jobs []Job) {
	url := "https://www.kununu.com/at/kununu/jobs"
	main_tag := "div"
	main_tag_attr := "class"
	main_tag_value := "company-profile-job-item"
	tag_title := "a"
	tag_url := "a"

	c := colly.NewCollector()
	c.OnHTML(main_tag, func(e *colly.HTMLElement) {
		if strings.Contains(e.Attr(main_tag_attr), main_tag_value) {
			job_title := e.ChildText(tag_title)
			job_url := e.ChildAttr(tag_url, "href")
			jobs = append(jobs, Job{
				runtime.Name,
				url,
				job_title,
				job_url})
		}
	})
	c.Visit(url)

	return
}

func (runtime Runtime) Mitte() (jobs []Job) {
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
        Title    string  `json:"text"`
        Url      string  `json:"hostedUrl"`
    }

    type JsonJob struct {
        Positions   []Postings `json:"postings"`
    }
    type JsonJobs []JsonJob
    var jsonJobs JsonJobs
    err = json.Unmarshal(body, &jsonJobs)
    if err != nil {
            fmt.Println(err)
        }
    for _, elem := range jsonJobs {
        fmt.Println("\t")
        job_title := elem.Positions[0].Title
		job_url := elem.Positions[0].Url
        fmt.Println(elem.Positions[0].Title)
        fmt.Println(elem.Positions[0].Url)
        jobs = append(jobs, Job{
            runtime.Name,
			url,
			job_title,
			job_url})
    }
	return
}
