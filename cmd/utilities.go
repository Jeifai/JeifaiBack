package cmd

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"time"

	"cloud.google.com/go/storage"
	"path/filepath"
	"strconv"
)

func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func Unique(results []Result) []Result {
	var unique []Result
	type key struct {
		CompanyName string
		Title       string
		ResultUrl   string
		Data        string
	}
	m := make(map[key]int)
	for _, v := range results {
		k := key{v.CompanyName, v.Title, v.ResultUrl, string(v.Data)}
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

func SaveResponseToFileWithFileName(response string, filename string) {
	dir, err := os.Getwd()
	if err != nil {
		panic(err.Error())
	}
	f, err := os.Create(dir + "/" + filename)
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

func RemoveFileWithFileName(filename string) {
	dir, err := os.Getwd()
	if err != nil {
		panic(err.Error())
	}
	err = os.Remove(dir + "/" + filename)
	if err != nil {
		panic(err.Error())
	}
}

func SaveResponseToStorage(response Response, file_path string) {
	fmt.Println(Gray(8-1, "Starting SaveResponseToStorage..."))

	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		panic(err.Error())
	}
	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	wc := client.Bucket("jeifai").Object(file_path).NewWriter(ctx)
	if _, err = io.Copy(wc, bytes.NewReader(response.Html)); err != nil {
		panic(err.Error())
	}
	if err := wc.Close(); err != nil {
		panic(err.Error())
	}
}

func GetResponseFromStorage(file_path string) (response string) {
	fmt.Println(Gray(8-1, "Starting GetResponseFromStorage..."))

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
	response = string(data)

	return
}
