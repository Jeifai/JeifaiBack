package main

import (
	"bytes"
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"time"
)

func generate_file_path(scraper_name string, scraper_version int) (file_path string) {
	file_path = filepath.Join(scraper_name, strconv.Itoa(scraper_version), "response.html")
	return
}

func SaveResponseToStorage(scraper Scraper, scraping Scraping, response Response) {

	fmt.Println("Starting SaveResponseToStorage...")

	file_path := generate_file_path(scraper.Name, scraping.Id)

	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()
	wc := client.Bucket("jeifai").Object(file_path).NewWriter(ctx)
	if _, err = io.Copy(wc, bytes.NewReader(response.Html)); err != nil {
		fmt.Println(err)
	}
	if err := wc.Close(); err != nil {
		fmt.Println(err)
	}
}

func (test *Test) GetResponseFromStorage() (response string) {

	fmt.Println("Starting GetResponseFromStorage...")

	err := test.LatestScrapingByNameAndVersion()
	if err != nil {
		fmt.Println(err)
	}

    file_path := generate_file_path(test.Name, test.Scraping)
    fmt.Println("Correctely loaded: " + file_path)

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
