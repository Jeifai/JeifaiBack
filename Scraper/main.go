package main

import (
	"fmt"

	. "github.com/logrusorgru/aurora"
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
		"Kununu", "IMusician", "Mitte", "Revolut",
		"Soundcloud", "Penta", "Celo", "N26", "Mollie",
		"Shopify", "Twitter", "Zalando", "Slack", "Circleci",
		"Google", "Hometogo", "Contentful", "Gympass", "Lanalabs", "Dreamingjobs",
    }

	scrapers, err := GetScrapers()
	if err != nil {
		panic(err.Error())
	}

	for _, elem := range scrapers {
		if Contains(want_scrapers, elem.Name) {
			fmt.Println(BrightBlue("Scraping -->"), Bold(BrightBlue(elem.Name)))
			response, results := Scrape(elem.Name, elem.Version, false)
			if len(results) > 0 {
				scraping, err := elem.StartScrapingSession()
				if err != nil {
					panic(err.Error())
				}
				file_path := GenerateFilePath(elem.Name, scraping.Id)
				SaveResults(elem, scraping, results)
				SaveResponseToStorage(response, file_path)
				if err != nil {
					panic(err.Error())
				}
			}
		}
	}
	if err != nil {
		panic(err.Error())
	}
}
