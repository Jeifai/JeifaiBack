package main

import (
	"bytes"
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"io"
	"path/filepath"
	"strconv"
	"time"
)

func generate_file_path(scraper Scraper, scraping Scraping) (file_path string) {
	file_path = filepath.Join(scraper.Name, strconv.Itoa(scraping.Id), "response.html")
	return
}

func SaveResponse(scraper Scraper, scraping Scraping, response Response) {

	fmt.Println("Starting SaveResponse...")

	file_path := generate_file_path(scraper, scraping)

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
