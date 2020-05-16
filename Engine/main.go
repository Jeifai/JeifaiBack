// 1 --> Extract from database all the website that needs to be scraped
// Get all the targets url
// Associate to the target its name
// 2 -->Run the scraper
// 3 --> Save in database the information scraped

package main

import (
    "fmt"
)

func scrape(scraper Scraper) (jobs []Job, err error) {
	jobs = runner(scraper.Name)
	return
}

func main() {
	scrapers, err := Scrapers()
	if err != nil {
		return
	}
	for _, elem := range scrapers {
        if elem.Name == "Mitte" {
            fmt.Println(elem.Name)
            jobs, err_2 := scrape(elem)
            if err_2 != nil {
                return
            }
            SaveJobs(elem, jobs)
        }
    }
}
