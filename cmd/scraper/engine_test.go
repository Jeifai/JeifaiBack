package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"testing"
	"time"

	"cloud.google.com/go/storage"
	. "github.com/logrusorgru/aurora"
	"github.com/stretchr/testify/assert"
	// "log"
)

/**
func TestMain(m *testing.M) {
	// log.SetOutput(ioutil.Discard)
}
*/

func TestUnique(t *testing.T) {
	fmt.Println(Blue("Running --> "), Bold(Blue("TestUnique")))
	testJson, err := json.Marshal("test")
	if err != nil {
		panic(err.Error())
	}
	result_1 := Result{"Test_1", "https://www.g_1.com", "Title_1", testJson}
	result_2 := Result{"Test_2", "https://www.g_2.com", "Title_2", testJson}
	result_3 := Result{"Test_1", "https://www.g_1.com", "Title_1", testJson}
	result_4 := Result{"Test_1", "https://www.g_1.com", "Title_1", testJson}
	results := []Result{result_1, result_2, result_3, result_4}
	got := Unique(results)
	want := []Result{result_1, result_2}
	assert.Equal(t, got, want, "The two []Result should be the same.")
}

func TestGenerateFilePath(t *testing.T) {
	fmt.Println(Blue("Running --> "), Bold(Blue("TestGenerateFilePath")))
	got := GenerateFilePath("scraper_name", 1)
	want := "scraper_name/1/response.html"
	assert.Equal(t, got, want, "The two path should be the same.")
}

func TestSaveResponseToFile(t *testing.T) {
	fmt.Println(Blue("Running --> "), Bold(Blue("TestSaveResponseToFile")))
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
	fmt.Println(Blue("Running --> "), Bold(Blue("TestRemoveFile")))
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
	fmt.Println(Blue("Running --> "), Bold(Blue("TestSaveResponseToStorage")))
	want := "this is a test to TestSaveResponseToStorage"
	file_path := "test/TestSaveResponseToStorage.html"
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
	fmt.Println(Blue("Running --> "), Bold(Blue("TestGetResponseFromStorage")))
	want := "this is a test to TestGetResponseFromStorage"
	file_path := "test/TestGetResponseFromStorage.html"
	got := GetResponseFromStorage(file_path)
	assert.Equal(t, got, want, "The two string should be the same")
}

func TestDbConnect(t *testing.T) {
	fmt.Println(Blue("Running --> "), Bold(Blue("TestDbConnect")))
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
	fmt.Println(Blue("Running --> "), Bold(Blue("TestScrape")))

	exclude_scrapers := []string{"Mitte", "Microsoft", "Amazon", "Deutschebahn"}
	scrapers, err := GetScrapers()
	if err != nil {
		panic(err.Error())
	}
	for _, elem := range scrapers {
		if !Contains(exclude_scrapers, elem.Name) {
			fmt.Println(BrightBlue("TESTING -> "), Bold(BrightBlue(elem.Name)))
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
			if err != nil {
				fmt.Println(err.Error())
			}
			sort.Slice(got, func(i, j int) bool {
				return got[i].ResultUrl < got[j].ResultUrl
			})
			sort.Slice(want, func(i, j int) bool {
				return want[i].ResultUrl < want[j].ResultUrl
			})
			assert.ElementsMatch(t, got, want, "The two []Result should be the same.")
		}
	}
}

func TestGetScrapers(t *testing.T) {
	fmt.Println(Blue("Running --> "), Bold(Blue("TestGetScrapers")))
	scrapers, err := GetScrapers()
	if err != nil {
		panic(err.Error())
	}
	assert.Equal(
		t, scrapers[0].Name, "Soundcloud", "The two string should be the same")
}

func TestLastScrapingByNameVersion(t *testing.T) {
	fmt.Println(Blue("Running --> "), Bold(Blue("TestLastScrapingByNameVersion")))
	last_scraping_version, err := LastScrapingByNameVersion("Mitte", 1)
	if err != nil {
		panic(err.Error())
	}
	assert.Greater(
		t, last_scraping_version, 1, "The last version should be bigger than 0")
}

func TestResultsByScraping(t *testing.T) {
	fmt.Println(Blue("Running --> "), Bold(Blue("TestResultsByScraping")))
	results, err := ResultsByScraping(245)
	if err != nil {
		panic(err.Error())
	}
	assert.Greater(
		t, len(results), 1, "The numbers of results should be bigger than 0")
}