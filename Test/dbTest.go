package main

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"os"
	"strconv"
)

var Db *sql.DB

func init() {

	// Load Environmental Variables
	err := godotenv.Load("/home/robimalco/jeifai/.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	dbhost := os.Getenv("DBHOST")
	dbuser := os.Getenv("DBUSER")
	dbport, err := strconv.ParseInt(os.Getenv("DBPORT"), 10, 64)
	dbname := os.Getenv("DBNAME")
	dbpassword := os.Getenv("DBPASSWORD")

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		dbhost, dbport, dbuser, dbpassword, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	Db = db
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected to the database")
}

// Get the latest scraping session by scraper name and scraper version
func (test *Test) LatestScrapingByNameAndVersion() (err error) {
	fmt.Println("Starting LatestScrapingByNameAndVersion...")
    err = Db.QueryRow(`SELECT MAX(s.id) FROM scraping s 
                        LEFT JOIN scrapers ss ON(s.scraper_id = ss.id)
                        LEFT JOIN targets t ON(ss.target_id = t.id)
                        WHERE t.name = $1 AND ss.version = $2;`, test.Name, test.Version).Scan(&test.Scraping)
    return
}

// Get all the results belonging to a specific scraping session
func (test *Test) ResultsByScraping() (results []Result, err error) {
	fmt.Println("Starting ResultsByUser...")
	rows, err := Db.Query(`SELECT r.title, r.url
                            FROM results r
                            WHERE r.scraping_id = $1`, test.Scraping)
	if err != nil {
		return
	}
	for rows.Next() {
		result := Result{}
		if err = rows.Scan(&result.Title, &result.Url); err != nil {
			return
		}
		results = append(results, result)
	}
	rows.Close()
	return
}