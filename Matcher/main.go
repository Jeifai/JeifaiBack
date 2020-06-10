package main

import (
	"fmt"
)

func main() {
	DbConnect()
	defer Db.Close()

	scraper_name := "Google"
	scraper_id, err := GetScraper(scraper_name)
	if err != nil {
		panic(err.Error())
	}

	matching := Matching{}
	err = matching.StartMatchingSession(scraper_id)
	if err != nil {
		panic(err.Error())
	}

	matches, err := GetMatches(matching, scraper_id)
	if err != nil {
		panic(err.Error())
	}

	for _, elem := range matches {
		fmt.Println(elem.CreatedAt)
		fmt.Println("\t", elem.CompanyName)
		fmt.Println("\t\t", elem.JobTitle)
		fmt.Println("\t\t\t", elem.JobUrl)
		fmt.Println("\t\t\t\t", elem.KeywordText)
	}

    if len(matches) > 0 {
	    SaveMatches(matching, matches)
    }
}
