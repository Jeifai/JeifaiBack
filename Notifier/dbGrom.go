package main

import (
    "fmt"
    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/postgres"
)

type Email struct {
    UserId		int
    UserName	string
    Company struct {
        CompanyName string
        Job struct {
            JobTitle 	string
            JobUrl		string
        }
    }
}

func main() {
    db, err := gorm.Open("postgres", "host=35.223.132.23 port=5432 user=ostgres dbname=jeifai password=jeifai")

    var emails []Email

    db.Raw(`
            SELECT DISTINCT
                u.id AS UserId,
                u.username AS UserName,
                s.name AS CompanyName,
                r.title AS JobTitle,
                r.url AS JobUrl
            FROM results r
            INNER JOIN matches m ON(r.id = m.resultid)
            LEFT JOIN scrapers s ON(r.scraperid = s.id)
            LEFT JOIN notifications n ON(m.id = n.matchid)
            LEFT JOIN userstargetskeywords utk ON(m.keywordid = utk.keywordid)
            LEFT JOIN users u ON(utk.userid = u.id)
            WHERE m.createdat > current_date - interval '0' day
            AND s.id = 10
            AND n.id IS NULL
            ORDER BY 1 DESC;`).Scan(&emails)

    fmt.Println(emails)


    _ = err
    defer db.Close()
}