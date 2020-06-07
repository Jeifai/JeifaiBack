package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"time"

	"cloud.google.com/go/storage"
)

func SaveResponseToStorage(response Response, file_path string) {
	fmt.Println("Starting SaveResponseToStorage...")

	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		panic(err.Error())
	}
	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	wc := client.Bucket("jeifai").Object(file_path).NewWriter(ctx)
	if _, err = io.Copy(wc, bytes.NewReader(response.Txt)); err != nil {
		panic(err.Error())
	}
	if err := wc.Close(); err != nil {
		panic(err.Error())
	}
}

func GetResponseFromStorage(file_path string) (response string) {
	fmt.Println("Starting GetResponseFromStorage...")

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
