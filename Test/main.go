package main

import (
    "fmt"
)

var scraper_name = "Mitte"
var scraper_version = 1

type Test struct {
    Name        string
    Version     int
    FilePath    string
    Scraping    int
}

func main() {
    test := Test{Name: scraper_name, Version: scraper_version}
    response := test.GetResponse()
    fmt.Println(response)

    // extract latest results from specific DB table

    // run scraper and get the results

    // compare the results and evaluate accuracy of the test
    

    //test("Mitte")       // Test Mitte with the latest version
}
