package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"google.golang.org/api/iterator"

	"cloud.google.com/go/storage"
)

func main() {
	bucket := "jeifai"
	replace_what := ".html"
	replace_with := ".txt"
	RenameFilesInStorage(bucket, replace_what, replace_with)
}

func RenameFilesInStorage(bucket, replace_what, replace_with string) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		panic(err.Error())
	}
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	it := client.Bucket(bucket).Objects(ctx, nil)
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			panic(err.Error())
		}

		old_name := attrs.Name

		if strings.Contains(old_name, ".html") {

			new_name := strings.ReplaceAll(old_name, replace_what, replace_with)

			src := client.Bucket(bucket).Object(old_name)
			dst := client.Bucket(bucket).Object(new_name)

			if _, err := dst.CopierFrom(src).Run(ctx); err != nil {
				if err != nil {
					panic(err.Error())
				}
			}
			if err := src.Delete(ctx); err != nil {
				if err != nil {
					panic(err.Error())
				}
			}

			fmt.Println(old_name, new_name, "\tDONE")
		}
	}
}
