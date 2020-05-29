package main

import (
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
    "io/ioutil"
	"encoding/json"
	"sort"
	"testing"
	"time"
	// "log"
	"os"
)

/**
func TestMain(m *testing.M) {
	// log.SetOutput(ioutil.Discard)
}
*/

func TestUnique(t *testing.T) {
    fmt.Println("TestUnique")
    testJson, err := json.Marshal("test")
    _ = err
	result_1 := Result{"Test_1", "https://www.g_1.com", "Title_1", testJson}
	result_2 := Result{"Test_2", "https://www.g_2.com", "Title_2", testJson}
	result_3 := Result{"Test_1", "https://www.g_1.com", "Title_1", testJson}
	result_4 := Result{"Test_1", "https://www.g_1.com", "Title_1", testJson}
	var results = []Result{result_1, result_2, result_3, result_4}
	got := Unique(results)
	want := []Result{result_1, result_2}
	assert.Equal(t, got, want, "The two []Result should be the same.")
}

func TestGenerateFilePath(t *testing.T) {
	fmt.Println("\n\nTestGenerateFilePath")
	got := GenerateFilePath("scraper_name", 1)
	want := "scraper_name/1/response.html"
	assert.Equal(t, got, want, "The two path should be the same.")
}

func TestSaveResponseToFile(t *testing.T) {
	fmt.Println("\n\nTestSaveResponseToFile")
	want := "this is a test string"
	SaveResponseToFile(want)
	dir, err := os.Getwd()
	if err != nil {
		panic(err.Error())
	}
	got, err := ioutil.ReadFile(dir + "/response.html")
	if err != nil {
		panic(err.Error())
	}
	assert.Equal(t, got, []byte(want), "The two string should be the same.")
}

func TestRemoveFile(t *testing.T) {
	fmt.Println("\n\nTestRemoveFile")
	want := "this is a test string"
	dir, err := os.Getwd()
	if err != nil {
		panic(err.Error())
	}
	f, err := os.Create(dir + "/response.html")
	if err != nil {
		panic(err.Error())
	}
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
	fmt.Println("\n\nTestSaveResponseToStorage")
	want := "this is a test to TestSaveResponseToStorage"
	file_path := "test/TestSaveResponseToStorage.txt"
	response := Response{[]byte(want)}
	SaveResponseToStorage(response, file_path)
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		panic(err.Error())
	}
	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()
	rc, err := client.Bucket("jeifai").Object(file_path).NewReader(ctx)
	if err != nil {
		panic(err.Error())
	}
	defer rc.Close()
	data, err := ioutil.ReadAll(rc)
	if err != nil {
		panic(err.Error())
	}
	got := string(data)
	assert.Equal(t, got, want, "The two string should be the same.")
}

func TestGetResponseFromStorage(t *testing.T) {
	fmt.Println("\n\nTestGetResponseFromStorage")
	want := "this is a test to TestGetResponseFromStorage"
	file_path := "test/TestGetResponseFromStorage.txt"
	got := GetResponseFromStorage(file_path)
	assert.Equal(t, got, want, "The two string should be the same")
}

func TestDbConnect(t *testing.T) {
	fmt.Println("\n\nTestDbConnect")
	type TempResult struct{ MinUser int }
	temp_result := TempResult{}
	DbConnect()
	err := Db.QueryRow(`SELECT MIN(s.id) FROM users s`).Scan(&temp_result.MinUser)
	if err != nil {
		panic(err.Error())
	}
	assert.Equal(t, temp_result.MinUser, 1, "The two integer should be the same")
}

func TestScrape(t *testing.T) {
	fmt.Println("\n\nTestScrape")
	scrapers, err := GetScrapers()
	if err != nil {
		panic(err.Error())
	}
	for _, elem := range scrapers {
		fmt.Println("TESTING -> ", elem.Name)
		scraping, err := LastScrapingByNameVersion(elem.Name, elem.Version)
		if err != nil {
			panic(err.Error())
		}
		file_path := GenerateFilePath(elem.Name, scraping)
		fileResponse := GetResponseFromStorage(file_path)
		got, err := ResultsByScraping(scraping)
		if err != nil {
			panic(err.Error())
		}
		SaveResponseToFile(fileResponse)
		isLocal := true
		httpResponse, want := Scrape(elem.Name, elem.Version, isLocal)
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
}

func TestGetScrapers(t *testing.T) {
	fmt.Println("\n\nTestGetScrapers")
	scrapers, err := GetScrapers()
	if err != nil {
		panic(err.Error())
	}
	assert.Equal(
		t, scrapers[0].Name, "IMusician", "The two string should be the same")
}

func TestLastScrapingByNameVersion(t *testing.T) {
	fmt.Println("\n\nTestLastScrapingByNameVersion")
	last_scraping_version, err := LastScrapingByNameVersion("Mitte", 1)
	if err != nil {
		panic(err.Error())
	}
	assert.Greater(
		t, last_scraping_version, 1, "The last version should be bigger than 0")
}

func TestResultsByScraping(t *testing.T) {
	fmt.Println("\n\nTestResultsByScraping")
	results, err := ResultsByScraping(136)
	if err != nil {
		panic(err.Error())
	}
	assert.Greater(
		t, len(results), 1, "The numbers of results should be bigger than 0")
}
