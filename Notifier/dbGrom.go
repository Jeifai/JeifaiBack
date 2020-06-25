package main

import (
    "fmt"
    "encoding/json"
    "io/ioutil"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Users struct {
	ID                          int
    Username                    string
    Email                       string
    Userstargetskeywords        []Userstargetskeywords `gorm:"ForeignKey:Userid"`
}

type Userstargetskeywords struct {
	ID          int
    Userid      int
    Keywordid   int
    Matches     []Matches `gorm:"ForeignKey:Keywordid"`
}

type Matches struct {
    ID          int
    Keywordid   int
    Resultid    int
    Results     []Results `gorm:"ForeignKey:ID"`
}

type Results struct {
    ID          int
    Title       string
    Url         string
    Scraperid   int
    Scrapers    []Scrapers `gorm:"ForeignKey:ID"`
}

type Scrapers struct {
    ID          int
    Name        string
}


func main() {
    db, _ := gorm.Open("postgres", "host=35.223.132.23 port=5432 user=postgres dbname=jeifai password=jeifai")
    db.SingularTable(true)

	var u []Users
    db.
    Preload("Userstargetskeywords").
    Preload("Userstargetskeywords.Matches").
    Preload("Userstargetskeywords.Matches.Results").
    Preload("Userstargetskeywords.Matches.Results.Scrapers").
    Find(&u)

    fmt.Println(u)
    

	file, _ := json.MarshalIndent(u, "", " ")
 
	_ = ioutil.WriteFile("test.json", file, 0644)
}
