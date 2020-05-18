package main

import (
	"fmt"
	"path/filepath"
    "strconv"
	"cloud.google.com/go/storage"
    "context"
    "time"
    "io/ioutil"
)

func generate_file_path(scraper_name string, scraper_version int) (file_path string) {
	file_path = filepath.Join(scraper_name, strconv.Itoa(scraper_version), "response.html")
	return
}

func (test *Test) GetResponse() (response string) {

    fmt.Println("Starting GetResponse...")
    
    err := test.LatestScrapingByNameAndVersion()
    if err != nil {
        fmt.Println(err)
    }

    file_path := generate_file_path(test.Name, test.Scraping)

    ctx := context.Background()
    client, err := storage.NewClient(ctx)
    ctx, cancel := context.WithTimeout(ctx, time.Second*50)
    defer cancel()
    rc, err := client.Bucket("jeifai").Object(file_path).NewReader(ctx)
    if err != nil {
        fmt.Println(err)
    }
    defer rc.Close()

    data, err := ioutil.ReadAll(rc)
    if err != nil {
        fmt.Println(err)
    }
    
    response = string(data)

    return 
}
