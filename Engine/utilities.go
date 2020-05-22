package main

import (
    "os"
    "fmt"
)

func Unique(result []Result) []Result {
    var unique []Result
    type key struct{ CompanyName, CompanyUrl, Title, ResultUrl string }
    m := make(map[key]int)
    for _, v := range result {
        k := key{v.CompanyName, v.CompanyUrl, v.Title, v.ResultUrl}
        if i, ok := m[k]; ok {
            unique[i] = v
        } else {
            m[k] = len(unique)
            unique = append(unique, v)
        }
    }
    return unique
}

func SaveResponseToFile(response string) {
    dir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	f, err := os.Create(dir + "/response.html")
	if err != nil {
		fmt.Println(err)
	}
    defer f.Close()
	f.WriteString(response)
}

func RemoveFile() {
    dir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	err_2 := os.Remove(dir + "/response.html")
	if err_2 != nil {
		fmt.Println(err_2)
	}
}