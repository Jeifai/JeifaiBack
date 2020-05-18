package main

import (
// "fmt"
)

var scraper_name = "Mitte"

func main() {
	scrapers, err := Scrapers()
	if err != nil {
		return
	}
	for _, elem := range scrapers {
		if elem.Name == scraper_name {
			scraping, err := elem.Scraping()
			if err != nil {
				return
			}
			response, results := runner(elem.Name, elem.Version)
			SaveResponse(elem, scraping, response)
			SaveResults(elem, scraping, results)
		}
	}
}
