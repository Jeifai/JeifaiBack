package main

import ("fmt")

func main() {
    scraper_name := "Babelforce"
    scraper_version := 1
    response, results := Scrape(scraper_name, scraper_version, true)
    _ = response
    fmt.Println(results)
}