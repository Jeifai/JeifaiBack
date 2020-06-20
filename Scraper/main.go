package main

import (
	"fmt"
)

func main() {
	DbConnect()
	defer Db.Close()

	/**
		want_scrapers := []string{"Amazon"}
	    // want_scrapers := []string{"Deutschebahn"}
	    // want_scrapers := []string{"Microsoft"}
	*/
	want_scrapers := []string{
		"Blinkist", "Urbansport", "Babelforce",
		"Kununu", "IMusician", "Mitte",
		"Soundcloud", "Penta", "Celo", "N26",
		"Shopify", "Twitter", "Zalando", "Slack",
		"Google", "Hometogo", "Contentful", "Gympass", "Lanalabs", "Dreamingjobs",
	}

	scrapers, err := GetScrapers()

	for _, elem := range scrapers {
		if Contains(want_scrapers, elem.Name) {
			fmt.Println(elem.Name)
			response, results := Scrape(elem.Name, elem.Version, false)
			if len(results) > 0 {
				scraping, err := elem.StartScrapingSession()
				file_path := GenerateFilePath(elem.Name, scraping.Id)
				SaveResults(elem, scraping, results)
				SaveResponseToStorage(response, file_path)
				_ = err
			}
		}
	}
	_ = err
}
