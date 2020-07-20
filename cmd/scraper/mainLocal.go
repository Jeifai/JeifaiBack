package main

import (
	"fmt"

	. "github.com/logrusorgru/aurora"
)

func mainLocal() {
	DbConnect()
	scraper_name := "Shopify"
	fmt.Println(BrightBlue("Scraping Locally -->"), Bold(BrightBlue(scraper_name)))
	scraper_version := 1
	scraping, err := LastScrapingByNameVersion(scraper_name, scraper_version)
	if err != nil {
		panic(err.Error())
	}
	file_path := GenerateFilePath(scraper_name, scraping)
	fileResponse := GetResponseFromStorage(file_path)
	SaveResponseToFile(fileResponse)
	response, results := Scrape(scraper_name, scraper_version, true)
	RemoveFile()
	_ = response
	_ = results
}
