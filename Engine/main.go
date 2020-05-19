package main

import (
	"fmt"
	"os"
)

func main() {
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
			response, results := runner(elem.Name, elem.Version, false)
			SaveResponseToStorage(elem, scraping, response)
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
	fileResponse := test.GetResponseFromStorage()
	storedResults, err := test.ResultsByScraping()
	if err != nil {
		fmt.Println(err)
	}
	save_response_to_file(fileResponse)
	httpResponse, newResults := runner(scraper_name, scraper_version, true)
	remove_file()
	_ = storedResults
	_ = httpResponse
	_ = newResults

	fmt.Println("storedResults:")
	fmt.Println(storedResults)
	fmt.Println("\n")
	fmt.Println("newResults:")
	fmt.Println(newResults)

	// COMPARE RESULTS
}

func save_response_to_file(response string) {
	f, err := os.Create("response.html")
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()
	f.WriteString(response)
}

func remove_file() {
	err := os.Remove("response.html")
	if err != nil {
		fmt.Println(err)
	}
}
