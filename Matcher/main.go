package main

import (
	"fmt"
)

func main() {
	DbConnect()
	defer Db.Close()

	scrapers, err := GetScrapers()
	if err != nil {
		panic(err.Error())
	}
	for _, elem := range scrapers {

		matching := Matching{}
		err = matching.StartMatchingSession(elem.Id)
		if err != nil {
			panic(err.Error())
		}

		matches, err := GetMatches(matching, elem.Id)
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
}
