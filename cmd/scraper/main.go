package cmd

import (
	"fmt"

	. "github.com/logrusorgru/aurora"
)

func Scrape(company string) {

	DbConnect()
	defer Db.Close()
    
    if company == "all" {
        scrapers := GetScrapers()
	    for _, elem := range scrapers {
            RunScraper(elem)
        }
    } else {
        scraper := GetScraper(company)
        RunScraper(scraper)
    }
}

func RunScraper(scraper Scraper) {
    if scraper.Name != "Microsoft" || scraper.Name != "Deutschebahn" || scraper.Name != "Amazon" {
        fmt.Println(BrightBlue("Scraping -->"), Bold(BrightBlue(scraper.Name)))
        response, results := Extract(scraper.Name, scraper.Version, false)
        n_results := len(results)
        if n_results > 0 {
            fmt.Println(Green("Number of results scraped: "), Bold(Green(n_results)))
            scraping, err := scraper.StartScrapingSession()
            if err != nil {
                panic(err.Error())
            }
            file_path := GenerateFilePath(scraper.Name, scraping.Id)
            SaveResults(scraper, scraping, results)
            SaveResponseToStorage(response, file_path)
            if err != nil {
                panic(err.Error())
            }
        } else {
            fmt.Println(Bold(Red("DANGER, NO RESULTS FOUND")))
        }
    }
}