// Initiate a new Scraping struct
// Using the scraper_id (or scraper name?) get Scraper data

package main

import (
	"fmt"
)

func main() {
	scraper := Scraper{Name: "Mitte"}
	scraper.ScraperByName()
	scraping, err := scraper.Scraping()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(scraping)
    jobs := runner(scraper.Name)
	SaveJobs(scraper, scraping, jobs)
}
