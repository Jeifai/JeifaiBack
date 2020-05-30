package main

import (
	"fmt"
)

func main() {
    DbConnect()
	fmt.Println("\n\nTestScrape")
	scrapers, err := GetScrapers()
	if err != nil {
		panic(err.Error())
	}
	for _, elem := range scrapers {
        if elem.Name == "Microsoft" {
            fmt.Println("TESTING -> ", elem.Name)
            scraping, err := LastScrapingByNameVersion(elem.Name, elem.Version)
            if err != nil {
                panic(err.Error())
            }
            file_path := GenerateFilePath(elem.Name, scraping)
            fileResponse := GetResponseFromStorage(file_path)
            fmt.Println(scraping)
            got, err := ResultsByScraping(scraping)
            if err != nil {
                panic(err.Error())
            }
            SaveResponseToFile(fileResponse)
            isLocal := true
            httpResponse, want := Scrape(elem.Name, elem.Version, isLocal)
            _ = want
            _ = got
            _ = httpResponse
            _ = err
        }
	}
}