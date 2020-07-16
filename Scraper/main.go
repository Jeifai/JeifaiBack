package main

import (
	"fmt"

	. "github.com/logrusorgru/aurora"
)

func main() {
	DbConnect()
	defer Db.Close()

	// Microsoft, Deutschebahn, Amazon
	want_scrapers := []string{
		"Blinkist", "Urbansport", "Babelforce", "Blacklane", "Kununu", "Docker",
		"IMusician", "Mitte", "Revolut", "Auto1", "Soundcloud", "Penta", "Zapier",
		"Celo", "N26", "Mollie", "Flixbus", "Shopify", "Twitter", "Zalando",
		"Slack", "Circleci", "Quora", "Google", "Hometogo", "Contentful", "Github",
		"Gympass", "Lanalabs", "Dreamingjobs", "Greenhouse", "Datadog", "Stripe",
		"Getyourguide", "Wefox", "Celonis", "Omio", "Aboutyou", "Depositsolutions",
		"Taxfix", "Moonfare", "Fincompare",
	}

	scrapers, err := GetScrapers()
	if err != nil {
		panic(err.Error())
	}

	for _, elem := range scrapers {
		if Contains(want_scrapers, elem.Name) {
			fmt.Println(BrightBlue("Scraping -->"), Bold(BrightBlue(elem.Name)))
			response, results := Scrape(elem.Name, elem.Version, false)
			n_results := len(results)
			if n_results > 0 {
				fmt.Println(Green("Number of results scraped: "), Bold(Green(n_results)))
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
			} else {
				fmt.Println(Bold(Red("DANGER, NO RESULTS FOUND")))
			}
		}
	}
	if err != nil {
		panic(err.Error())
	}
}
