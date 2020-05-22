package main

import (
	"fmt"
)

func main() {
	DbConnect()
	// scrape("Mitte")
	test("Mitte", 1)
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

type Test struct {
	Name     string
	Version  int
	FilePath string
	Scraping int
}

func test(scraper_name string, scraper_version int) {
	test := Test{Name: scraper_name, Version: scraper_version}

	err := test.LatestScrapingByNameAndVersion()
	if err != nil {
		fmt.Println(err)
	}

	file_path := GenerateFilePath(scraper_name, test.Scraping)
	fileResponse := GetResponseFromStorage(file_path)
	storedResults, err := test.ResultsByScraping()
	if err != nil {
		fmt.Println(err)
	}
	SaveResponseToFile(fileResponse)
	httpResponse, newResults := Runner(scraper_name, scraper_version, true)
	RemoveFile()
	_ = httpResponse

	var bool_array []bool
	for _, stored_element := range storedResults {
		for _, new_element := range newResults {
			if stored_element.ResultUrl == new_element.ResultUrl {
				bool_url := new_element.ResultUrl == stored_element.ResultUrl
				bool_title := new_element.Title == stored_element.Title
				bool_comparison := bool_url == bool_title
				bool_array = append(bool_array, bool_comparison)
			}
		}
	}
	if len(storedResults) == len(bool_array) {
		fmt.Println("TEST STATUS: OK")
	} else {
		fmt.Println("TEST STATUS: KO")
	}
}
