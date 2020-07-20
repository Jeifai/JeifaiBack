package cmd

import (
	"fmt"

	. "github.com/logrusorgru/aurora"
)

func TempMatch(company string) {

	DbConnect()
    defer Db.Close()
    
    if company == "all" {
        scrapers := GetScrapers()
	    for _, elem := range scrapers {
            RunMatcher(elem)
        }
    } else {
        scraper := GetScraper(company)
        RunMatcher(scraper)
    }
}


func RunMatcher(scraper Scraper) {
    fmt.Println(Blue("Running --> "), Bold(Blue(scraper.Name)))

    matching := Matching{}
    matching.StartMatchingSession(scraper.Id)

    matches, err := GetMatches(matching, scraper.Id)
    if err != nil {
        panic(err.Error())
    }

    for _, elem := range matches {
        fmt.Println(
            Bold(Green("\tNew Match -->")),
            Faint(Green(elem.KeywordText)),
            Bold(Green(elem.JobTitle)),
            Faint(Underline(BrightGreen(elem.JobUrl))))
    }

    if len(matches) > 0 {
        SaveMatches(matching, matches)
    }
}