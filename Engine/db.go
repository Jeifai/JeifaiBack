package main

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"os"
	"strconv"
	"time"
)

type Scraping struct {
	Id        int
	ScraperId int
	CreatedAt time.Time
}

type Scraper struct {
	Id      int
	Name    string
	Version int
}

var Db *sql.DB

func DbConnect() {

	// Load Environmental Variables
	err := godotenv.Load()
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
    // Send a ping to make sure the database connection is alive.
    if err = Db.Ping(); err != nil {
        Db.Close()
	    fmt.Println("Unsuccessfully connected to the database")
        return
    }
	fmt.Println("Successfully connected to the database")
}

// Get all the scrapers and the latests versions
func Scrapers() (scrapers []Scraper, err error) {
	fmt.Println("Starting Scrapers...")
    rows, err := Db.Query(`SELECT
                                s.name,
                                MAX(s.version) AS version,
                                MAX(s.id) AS id 
                           FROM scrapers s GROUP BY 1;`)
	if err != nil {
		return
	}
	for rows.Next() {
		scraper := Scraper{}
		if err = rows.Scan(&scraper.Name, &scraper.Version, &scraper.Id); err != nil {
			return
		}
		scrapers = append(scrapers, scraper)
	}
	rows.Close()
	return
}

// Get the id of a scraper based on name
func (scraper *Scraper) ScraperByName() (err error) {
	fmt.Println("Starting ScraperByName...")
	err = Db.QueryRow(`SELECT s.id
                       FROM scrapers s
                       WHERE s.name=$1`, scraper.Name).Scan(&scraper.Id)
	return
}

// Save in the database a new Scraping session and return its values
func (scraper *Scraper) Scraping() (scraping Scraping, err error) {
    statement := "INSERT INTO scraping
                        (scraper_id, created_at)
                  VALUES ($1, $2) 
                  RETURNING id, scraper_id, created_at"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	err = stmt.QueryRow(scraper.Id, time.Now()).Scan(&scraping.Id, &scraping.ScraperId, &scraping.CreatedAt)
	return
}

// Save all the results
func SaveResults(scraper Scraper, scraping Scraping, results []Result) {
	fmt.Println("Starting SaveResults...")
	for _, elem := range results {
        statement := "INSERT INTO results
                            (scraper_id, scraping_id, title, url, scraping_url, created_at)
                      VALUES ($1, $2, $3, $4, $5, $6)"
		_, err := Db.Exec(statement, scraper.Id, scraping.Id, elem.Title, elem.ResultUrl, elem.ScrapingUrl, time.Now())
		if err != nil {
			return
		}
	}
}

// Get the latest scraping session by scraper name and scraper version
func LatestScrapingByNameAndVersion(scraper_name string, scraper_version int) (scraping int, err error) {
    fmt.Println("Starting LatestScrapingByNameAndVersion...")
	err = Db.QueryRow(`SELECT MAX(s.id) FROM scraping s 
                        LEFT JOIN scrapers ss ON(s.scraper_id = ss.id)
                        LEFT JOIN targets t ON(ss.target_id = t.id)
                        WHERE t.name = $1 AND ss.version = $2;`, scraper_name, scraper_version).Scan(&scraping)
	return
}

// Get all the results by scraping session
func ResultsByScraping(scraping int) (results []Result, err error) {
	fmt.Println("Starting ResultsByScraping...")
	rows, err := Db.Query(`SELECT t.name, r.scraping_url, r.title, r.url
                            FROM results r
                            LEFT JOIN scraping s ON(r.scraping_id = s.id)
                            LEFT JOIN scrapers ss ON(s.scraper_id = ss.id)
                            LEFT JOIN targets t ON(ss.target_id = t.id)
                            WHERE r.scraping_id = $1`, scraping)
	if err != nil {
		return
	}
	for rows.Next() {
		result := Result{}
		if err = rows.Scan(&result.CompanyName, &result.ScrapingUrl, &result.Title, &result.ResultUrl); err != nil {
			return
		}
		results = append(results, result)
	}
	rows.Close()
	fmt.Println("Number of results loaded: " + strconv.Itoa(len(results)))
	return
}
