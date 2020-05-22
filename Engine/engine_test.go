package main

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
)

func TestUnique(t *testing.T) {
	result_1 := Result{"Test_1", "https://www.gooo_1.com", "Title_1", "Result_1"}
	result_2 := Result{"Test_2", "https://www.gooo_2.com", "Title_2", "Result_2"}
	result_3 := Result{"Test_1", "https://www.gooo_1.com", "Title_1", "Result_1"}
	result_4 := Result{"Test_1", "https://www.gooo_1.com", "Title_1", "Result_1"}
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
	SaveResponseToFile(want)
	RemoveFile()
	dir, err := os.Getwd()
	file, err := ioutil.ReadFile(dir + "/response.html")
	_ = file
	if assert.Error(t, err) {
		assert.NotNil(t, err, "The error should be nil")
	}
}
