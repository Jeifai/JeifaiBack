package main

import (
    "github.com/stretchr/testify/assert"
    "testing"
)

func TestUnique(t *testing.T) {
    result_1 := Result{"Test_1", "https://www.gooo_1.com", "Title_1", "Result_1"}
    result_2 := Result{"Test_2", "https://www.gooo_2.com", "Title_2", "Result_2"}
    result_3 := Result{"Test_1", "https://www.gooo_1.com", "Title_1", "Result_1"}
    result_4 := Result{"Test_1", "https://www.gooo_1.com", "Title_1", "Result_1"}
    var results = []Result {result_1, result_2, result_3, result_4}
    got := Unique(results)
    want := []Result{result_1, result_2}
    assert.Equal(t, got, want, "The two []Result should be the same.")
}

func TestGenerateFilePath(t *testing.T) {
    got := GenerateFilePath("scraper_name", 1)
    want := "scraper_name/1/response.html"
    assert.Equal(t, got, want, "The two path should be the same.")    
}