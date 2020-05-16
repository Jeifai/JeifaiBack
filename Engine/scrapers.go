package main

import (
	"github.com/gocolly/colly"
	"reflect"
    "strings"
    "net/http"
    "fmt"
	"io/ioutil"
)

type Runtime struct {
	Name string
}
type Job struct {
	CompanyName    string
	CompanyUrl     string
	Title          string
	JobUrl         string
}

func runner(name string) (job []Job) {
	r := Runtime{name}
	v := reflect.ValueOf(r)
	m := v.MethodByName(r.Name)
	temp_job := m.Call(nil)
	job = temp_job[0].Interface().([]Job)
	return
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
    resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error")
	}

    fmt.Println(string(body))
    return
}
