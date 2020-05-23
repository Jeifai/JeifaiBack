package main

import (
	// "fmt"
)

func main() {
	DbConnect()
	scrape("Kununu")
}
func scrape(scraper_name string) {
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
			response, results := Runner(elem.Name, elem.Version, false)
			file_path := GenerateFilePath(elem.Name, scraping.Id)
			SaveResponseToStorage(response, file_path)
			SaveResults(elem, scraping, results)
		}
	}
}
