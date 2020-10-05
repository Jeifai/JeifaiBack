package cmd

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	. "github.com/logrusorgru/aurora"
	"github.com/teris-io/shortid"
)

type Scraping struct {
	Id        int
	ScraperId int
	CreatedAt time.Time
}

type Scraper struct {
	Id   int
	Name string
}

type Company struct {
	Name string
	Job  []Job
}

type Job struct {
	Title 		string
	Url   	 	string
	Location 	string
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
		fmt.Println(Bold(Red("Unsuccessfully connected to the database")))
		return
	}
	fmt.Println(Bold(Green("Successfully connected to the database")))
}

func GetScraper(company string) (scraper Scraper) {
	fmt.Println(Gray(8-1, "Starting GetScraper..."))
	err := Db.QueryRow(`SELECT
                            s.name, 
                            s.id AS id 
                        FROM scrapers s
                        WHERE s.name=$1;`,
		company).Scan(
		&scraper.Name,
		&scraper.Id)
	if err != nil {
		panic(err.Error())
	}
	return
}

func GetScrapers() (scrapers []Scraper) {
	fmt.Println(Gray(8-1, "Starting GetScrapers..."))
	rows, err := Db.Query(`SELECT
                                s.name, 
                                s.id AS id 
                           FROM scrapers s;`)
	if err != nil {
		panic(err.Error())
	}
	for rows.Next() {
		scraper := Scraper{}
		if err = rows.Scan(
			&scraper.Name,
			&scraper.Id); err != nil {
			return
		}
		scrapers = append(scrapers, scraper)
	}
	rows.Close()
	return
}

func (scraper *Scraper) StartScrapingSession(count_results int) (scraping Scraping) {
	fmt.Println(Gray(8-1, "Starting StartScrapingSession..."))
	statement := `INSERT INTO scrapings (scraperid, countresults, createdat)
                  VALUES ($1, $2, $3) 
                  RETURNING id, scraperid, createdat`
	stmt, err := Db.Prepare(statement)
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()
	err = stmt.QueryRow(scraper.Id, count_results, time.Now()).Scan(
		&scraping.Id, &scraping.ScraperId, &scraping.CreatedAt)
	if err != nil {
		panic(err.Error())
	}
	return
}

func SaveResults(scraper Scraper, scraping Scraping, results []Result) {
	fmt.Println(Gray(8-1, "Starting SaveResults..."))
	valueStrings := []string{}
	valueArgs := []interface{}{}
	timeNow := time.Now() // updatedAt and createdAt will be identical
	sid, err := shortid.New(1, shortid.DefaultABC, 2342)

	for i, elem := range results {
		unique_url_id, err := sid.Generate()
		if err != nil {
			panic(err.Error())
		}
		str1 := "$" + strconv.Itoa(1+i*9) + ","
		str2 := "$" + strconv.Itoa(2+i*9) + ","
		str3 := "$" + strconv.Itoa(3+i*9) + ","
		str4 := "$" + strconv.Itoa(4+i*9) + ","
		str5 := "$" + strconv.Itoa(5+i*9) + ","
		str6 := "$" + strconv.Itoa(6+i*9) + ","
		str7 := "$" + strconv.Itoa(7+i*9) + ","
		str8 := "$" + strconv.Itoa(8+i*9) + ","
		str9 := "$" + strconv.Itoa(9+i*9)
		str_n := "(" + str1 + str2 + str3 + str4 + str5 + str6 + str7 + str8 + str9 + ")"
		valueStrings = append(valueStrings, str_n)
		valueArgs = append(valueArgs, scraper.Id)
		valueArgs = append(valueArgs, scraping.Id)
		valueArgs = append(valueArgs, elem.Title)
		valueArgs = append(valueArgs, elem.ResultUrl)
		valueArgs = append(valueArgs, unique_url_id)
		valueArgs = append(valueArgs, elem.Location)
		valueArgs = append(valueArgs, elem.Data)
		valueArgs = append(valueArgs, timeNow)
		valueArgs = append(valueArgs, timeNow)
	}
	smt := `INSERT INTO results (
                scraperid, scrapingid, title, url, urlshort, location, data, createdat, updatedat)
            VALUES %s ON CONFLICT (url) DO UPDATE
            SET scrapingid = EXCLUDED.scrapingid,
                title = EXCLUDED.title,
                location = EXCLUDED.location,
                updatedat = EXCLUDED.updatedat,
                data = EXCLUDED.data`
	smt = fmt.Sprintf(smt, strings.Join(valueStrings, ","))

	_, err = Db.Exec(smt, valueArgs...)
	if err != nil {
		panic(err.Error())
	}
}
