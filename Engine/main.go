package main

import (
// "fmt"
)

var scraper_name = "Kununu"

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
			results := runner(elem.Name, elem.Version)
			SaveResults(elem, scraping, results)
		}
	}
}
