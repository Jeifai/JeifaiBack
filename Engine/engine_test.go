package main

import (
	"cloud.google.com/go/storage"
	"context"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"sort"
	"testing"
	"time"
)

func TestUnique(t *testing.T) {
	result_1 := Result{"Test_1", "https://www.g_1.com", "Title_1", "Result_1"}
	result_2 := Result{"Test_2", "https://www.g_2.com", "Title_2", "Result_2"}
	result_3 := Result{"Test_1", "https://www.g_1.com", "Title_1", "Result_1"}
	result_4 := Result{"Test_1", "https://www.g_1.com", "Title_1", "Result_1"}
	var results = []Result{result_1, result_2, result_3, result_4}
	got := Unique(results)
	want := []Result{result_1, result_2}
	assert.Equal(t, got, want, "The two []Result should be the same.")
}

func TestGenerateFilePath(t *testing.T) {
	got := GenerateFilePath("scraper_name", 1)
	want := "scraper_name/1/response.html"
	assert.Equal(t, got, want, "The two path should be the same.")
}

func TestSaveResponseToFile(t *testing.T) {
	want := "this is a test string"
	SaveResponseToFile(want)
	dir, err := os.Getwd()
	got, err := ioutil.ReadFile(dir + "/response.html")
	_ = err
	assert.Equal(t, got, []byte(want), "The two string should be the same.")
}

func TestRemoveFile(t *testing.T) {
	want := "this is a test string"
	dir, err := os.Getwd()
	f, err := os.Create(dir + "/response.html")
	defer f.Close()
	f.WriteString(want)
	RemoveFile()
	file, err := ioutil.ReadFile(dir + "/response.html")
	_ = file
	if assert.Error(t, err) {
		assert.NotNil(t, err, "The error should be nil")
	}
}

func TestSaveResponseToStorage(t *testing.T) {
	want := "this is a test to TestSaveResponseToStorage"
	file_path := "test/TestSaveResponseToStorage.txt"
	response := Response{[]byte(want)}
	SaveResponseToStorage(response, file_path)
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()
	rc, err := client.Bucket("jeifai").Object(file_path).NewReader(ctx)
	defer rc.Close()
	data, err := ioutil.ReadAll(rc)
	_ = err
	got := string(data)
	assert.Equal(t, got, want, "The two string should be the same.")
}

func TestGetResponseFromStorage(t *testing.T) {
	want := "this is a test to TestGetResponseFromStorage"
	file_path := "test/TestGetResponseFromStorage.txt"
	got := GetResponseFromStorage(file_path)
	assert.Equal(t, got, want, "The two string should be the same")
}

func TestDbConnect(t *testing.T) {
	type TempResult struct{ MinUser int }
	temp_result := TempResult{}
	DbConnect()
	err := Db.QueryRow(`SELECT MIN(s.id) FROM users s`).Scan(&temp_result.MinUser)
	_ = err
	assert.Equal(t, temp_result.MinUser, 1, "The two integer should be the same")
}

func TestScrape(t *testing.T) {
	scraper_name := "Kununu"
	scraper_version := 1
	scraping, err := LastScrapingByNameVersion(scraper_name, scraper_version)
	file_path := GenerateFilePath(scraper_name, scraping)
	fileResponse := GetResponseFromStorage(file_path)
	got, err := ResultsByScraping(scraping)
    SaveResponseToFile(fileResponse)
    isLocal := true
	httpResponse, want := Scrape(scraper_name, scraper_version, isLocal)
	RemoveFile()
	_ = httpResponse
	_ = err
	sort.Slice(got, func(i, j int) bool {
		return got[i].ResultUrl < got[j].ResultUrl
	})
	sort.Slice(want, func(i, j int) bool {
		return want[i].ResultUrl < want[j].ResultUrl
	})
	assert.ElementsMatch(t, got, want, "The two []Result should be the same.")
}
