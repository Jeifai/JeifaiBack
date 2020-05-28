package main

import (
	"os"
	"path/filepath"
	"strconv"
)

func Unique(results []Result) []Result {
	var unique []Result
	type key struct{ CompanyName, Title, ResultUrl string }
	m := make(map[key]int)
	for _, v := range results {
		k := key{v.CompanyName, v.Title, v.ResultUrl}
		if i, ok := m[k]; ok {
			unique[i] = v
		} else {
			m[k] = len(unique)
			unique = append(unique, v)
		}
	}
	return unique
}

func GenerateFilePath(
	scraper_name string, scraper_version int) (file_path string) {
	file_path = filepath.Join(
		scraper_name, strconv.Itoa(scraper_version), "response.html")
	return
}

func SaveResponseToFile(response string) {
	dir, err := os.Getwd()
	if err != nil {
		panic(err.Error())
	}
	f, err := os.Create(dir + "/response.html")
	if err != nil {
		panic(err.Error())
	}
	defer f.Close()
	f.WriteString(response)
}

func RemoveFile() {
	dir, err := os.Getwd()
	if err != nil {
		panic(err.Error())
	}
	err = os.Remove(dir + "/response.html")
	if err != nil {
		panic(err.Error())
	}
}
