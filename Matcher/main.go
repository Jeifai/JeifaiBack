package main

import (
	"fmt"
)

func main() {
	DbConnect()

	matches, err := GetMatches()
	if err != nil {
		panic(err.Error())
	}

	for _, elem := range matches {
		fmt.Println(elem.CreatedAt)
		fmt.Println("\t", elem.CompanyName)
		fmt.Println("\t\t", elem.JobTitle)
		fmt.Println("\t\t\t", elem.JobUrl)
		fmt.Println("\t\t\t\t", elem.Keyword)
	}

	defer Db.Close()
}
