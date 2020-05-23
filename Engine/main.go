package main

import (
// "fmt"
)

func main() {
	DbConnect()
	scraper_name := "Zalando"
	scrapers, err := GetScrapers()
	if err != nil {
		panic(err.Error())
	}
	for _, elem := range scrapers {
		if elem.Name == scraper_name {
			scraping, err := elem.StartScrapingSession()
			if err != nil {
				panic(err.Error())
			}
			isLocal := false
			response, results := Scrape(elem.Name, elem.Version, isLocal)
			file_path := GenerateFilePath(elem.Name, scraping.Id)
			SaveResponseToStorage(response, file_path)
			SaveResults(elem, scraping, results)
		}
	}
}
