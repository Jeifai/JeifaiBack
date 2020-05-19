package main

import (
    "fmt"
    "time"
)

var scraper_name = "Mitte"
var scraper_version = 1

type Test struct {
    Name        string
    Version     int
    FilePath    string
    Scraping    int
    Results     []Result
}

type Result struct {
	Title     string
	Url       string
}
func main() {
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
