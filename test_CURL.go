package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func mainCurl() {
	client := &http.Client{}
	data := strings.NewReader(`{"limit":30}`)
	req, err := http.NewRequest("POST", "https://job.bytedance.com/api/v1/search/job/posts", data)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("x-csrf-token", "uUVTqdJlG1S0uWVuXkkTguFfzsTSyg2p6KdqlqbMIus=")
	req.Header.Set("Cookie", "channel=overseas; atsx-csrf-token=uUVTqdJlG1S0uWVuXkkTguFfzsTSyg2p6KdqlqbMIus%3D")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", bodyText)
}
