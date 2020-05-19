package main

import (
    "fmt"
)

func main() {
    test()
}

func scrape() {
    var scraper_name = "IMusician"
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

type Test struct {
    Name            string
    Version         int
    FilePath        string
    Scraping        int
}
func test() {
    var scraper_name = "IMusician"
    var scraper_version = 1
    test := Test{Name: scraper_name, Version: scraper_version}
    response := test.GetResponse()
    results, err := test.ResultsByScraping()
    if err != nil {
        fmt.Println(err)
    }
    fmt.Println(results)
    _ = response
    // run scraper and get the results
    // compare the results and evaluate accuracy of the test
    //test("Mitte")       // Test Mitte with the latest version
}