package cmd

import (
	"os"
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
		ResultUrl string
	}
	m := make(map[key]int)
	for _, v := range results {
		k := key{v.ResultUrl}
		if i, ok := m[k]; ok {
			unique[i] = v
		} else {
			m[k] = len(unique)
			unique = append(unique, v)
		}
	}
	return unique
}

func RemoveDuplicatedFromSliceOfString(slice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range slice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
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

