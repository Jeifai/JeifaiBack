package main

import (
	"fmt"
)

func main() {
	DbConnect()

	matching := Matching{}
	err := matching.StartMatchingSession()
	if err != nil {
		panic(err.Error())
	}

	matches, err := GetMatches(matching)
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

	SaveMatches(matching, matches)

	defer Db.Close()
}
