package main

import (
	"fmt"

	. "github.com/logrusorgru/aurora"
)

func main() {
	DbConnect()
	defer Db.Close()

	scrapers, err := GetScrapers()
	if err != nil {
		panic(err.Error())
	}
	for _, elem := range scrapers {

		fmt.Println(Blue("Running --> "), Bold(Blue(elem.Name)))

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
}
