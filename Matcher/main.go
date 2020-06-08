package main

import (
	"fmt"
)

func main() {
	DbConnect()

	scraper_name := "Microsoft"

	matches, err := GetMatches(scraper_name)
	if err != nil {
		panic(err.Error())
	}

	fmt.Println(matches)

	_ = err

	defer Db.Close()
}
