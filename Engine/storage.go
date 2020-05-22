package main

import (
	"bytes"
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"time"
)

func SaveResponseToStorage(response Response, file_path string) {

	fmt.Println("Starting SaveResponseToStorage...")

	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	wc := client.Bucket("jeifai").Object(file_path).NewWriter(ctx)
	if _, err = io.Copy(wc, bytes.NewReader(response.Html)); err != nil {
		fmt.Println(err)
	}
	if err := wc.Close(); err != nil {
		fmt.Println(err)
	}
}

func (test *Test) GetResponseFromStorage(file_path string) (response string) {

	fmt.Println("Starting GetResponseFromStorage...")

	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	rc, err := client.Bucket("jeifai").Object(file_path).NewReader(ctx)
	if err != nil {
		fmt.Println(err)
	}
	defer rc.Close()

	data, err := ioutil.ReadAll(rc)
	if err != nil {
		fmt.Println(err)
	}
	response = string(data)

	return
}
