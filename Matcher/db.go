package main

import (
	"database/sql"
	"fmt"
    "os"
    "time"
    "strconv"
    "encoding/json"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Scraper struct {
	Id      int
	Name    string
	Version int
}

type Scraping struct {
	Id        int
	ScraperId int
	CreatedAt time.Time
}

type Result struct {
	Title       string
	ResultUrl   string
	Data        json.RawMessage
}

var Db *sql.DB

func DbConnect() {
	err := godotenv.Load()
	if err != nil {
		panic(err.Error())
	}
	dbhost := os.Getenv("DBHOST")
	dbuser := os.Getenv("DBUSER")
	dbport, err := strconv.ParseInt(os.Getenv("DBPORT"), 10, 64)
	if err != nil {
		panic(err.Error())
	}
	dbname := os.Getenv("DBNAME")
	dbpassword := os.Getenv("DBPASSWORD")

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		dbhost, dbport, dbuser, dbpassword, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err.Error())
	}
	Db = db

	if err = Db.Ping(); err != nil {
		Db.Close()
		fmt.Println("Unsuccessfully connected to the database")
		return
	}
	fmt.Println("Successfully connected to the database")
}

func GetScraperByScraperName(name string) (scraper Scraper, err error) {
	fmt.Println("Starting GetScraperByScraperName...")
	err = Db.QueryRow(`SELECT
                        id,
                        name,
                        version
                      FROM scrapers
                      WHERE name = $1`,
		name,
	).
		Scan(
			&scraper.Id,
			&scraper.Name,
			&scraper.Version,
		)
	return
}

func GetLastScrapingByScraperId(scraper Scraper) (scraping Scraping, err error) {
	fmt.Println("Starting GetLastScraping...")
	err = Db.QueryRow(`SELECT
                        MAX(id)
                      FROM scrapings
                      WHERE scraperid = $1`,
		scraper.Id,
	).
		Scan(
			&scraping.Id,
		)
	return
}

func GetNewResultsByScrapingId(scraping Scraping) (results []Result, err error) {
	fmt.Println("Starting GetNewResultsByScrapingId...")
    rows, err := Db.Query(`SELECT
                                r.title,
                                r.url
                            FROM results r
                            WHERE r.scrapingid = $1`, scraping.Id) // AND r.createdat = r.updatedat
	if err != nil {
		return
	}
	for rows.Next() {
		result := Result{}
		if err = rows.Scan(&result.Title, &result.ResultUrl); err != nil {
			return
		}
		results = append(results, result)
	}
	rows.Close()
	return
}