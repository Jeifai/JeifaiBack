package main

import "fmt"

func main() {
    DbConnect()
    fmt.Println("TEST")
    defer Db.Close()
}