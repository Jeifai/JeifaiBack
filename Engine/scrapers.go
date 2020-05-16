package main

import (
	"github.com/gocolly/colly"
	"reflect"
	"strings"
)

type Runtime struct {
	CompanyName string
}
type Job struct {
	ScraperVersion int
	CompanyName    string
	CompanyUrl     string
	JobTitle       string
	JobUrl         string
}

func runner(companyName string) (job []Job) {
	r := Runtime{companyName}
	v := reflect.ValueOf(r)
	m := v.MethodByName(r.CompanyName)
	temp_job := m.Call(nil)
	job = temp_job[0].Interface().([]Job)
	return
}

func (runtime Runtime) Kununu() (jobs []Job) {
	version := 1

	company_url := "https://www.kununu.com/at/kununu/jobs"
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
				version,
				runtime.CompanyName,
				company_url,
				job_title,
				job_url})
		}
	})
	c.Visit(company_url)

	return
}
