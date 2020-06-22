package main

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Scraper struct {
	Id      int
	Name    string
	Version int
}

type Matching struct {
	Id        int
	CreatedAt time.Time
}

type Match struct {
	Id              int
	ResultCreatedAt time.Time
	CompanyName     string
	JobTitle        string
	JobUrl          string
	KeywordText     string
	KeywordId       int
	ResultId        int
	CreatedAt       time.Time
	MatchingId      int
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

func GetScrapers() (scrapers []Scraper, err error) {
	fmt.Println("Starting GetScrapers...")
	rows, err := Db.Query(`SELECT
                                s.name, 
                                MAX(s.version) AS version, 
                                MAX(s.id) AS id 
                           FROM scrapers s
                           GROUP BY 1;`)
	if err != nil {
		panic(err.Error())
	}
	for rows.Next() {
		scraper := Scraper{}
		if err = rows.Scan(
			&scraper.Name,
			&scraper.Version,
			&scraper.Id); err != nil {
			return
		}
		scrapers = append(scrapers, scraper)
	}
	rows.Close()
	return
}

func (matching *Matching) StartMatchingSession(scraper_id int) (err error) {
	fmt.Println("Starting StartMatchingSession...")
	statement := `INSERT INTO matchings (scraperid, createdat)
                  VALUES ($1, $2)
                  RETURNING id`
	stmt, err := Db.Prepare(statement)
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()
	err = stmt.QueryRow(
		scraper_id,
		time.Now()).Scan(
		&matching.Id)
	if err != nil {
		panic(err.Error())
	}
	return
}

func GetMatches(matching Matching, scraper_id int) (matches []Match, err error) {
	fmt.Println("Starting GetMatches...")
	rows, err := Db.Query(`
                        SELECT
                            r.createdat AS created_at,
                            s.name AS company,
                            r.title AS job_title,
                            r.url AS job_url,
                            k.text AS keyword_text,
                            k.id AS keyword_id,
                            r.id AS result_id
                        FROM targets t
                        INNER JOIN scrapers s ON(t.id = s.targetid)
                        INNER JOIN results r ON(s.id = r.scraperid)
                        LEFT JOIN userstargetskeywords utk ON(t.id = utk.targetid)
                        LEFT JOIN keywords k ON(utk.keywordid = k.id)
                        WHERE r.createdat::date = date_trunc('day', now()) 
                        AND r.createdat > current_date
                        AND s.id = $1
                        AND REPLACE(LOWER(r.title), ' ', '') LIKE '%' || REPLACE(LOWER(k.text), ' ', '') || '%'`, scraper_id)
	if err != nil {
		return
	}
	for rows.Next() {
		match := Match{}
		if err = rows.Scan(
			&match.ResultCreatedAt,
			&match.CompanyName,
			&match.JobTitle,
			&match.JobUrl,
			&match.KeywordText,
			&match.KeywordId,
			&match.ResultId); err != nil {
			return
		}
		match.CreatedAt = time.Now()
		match.MatchingId = matching.Id
		matches = append(matches, match)
	}
	rows.Close()
	return
}

func SaveMatches(matching Matching, matches []Match) {
	fmt.Println("Starting SaveMatches...")
	valueStrings := []string{}
	valueArgs := []interface{}{}
	timeNow := time.Now() // updatedAt and createdAt will be identical
	for i, elem := range matches {
		str1 := "$" + strconv.Itoa(1+i*4) + ","
		str2 := "$" + strconv.Itoa(2+i*4) + ","
		str3 := "$" + strconv.Itoa(3+i*4) + ","
		str4 := "$" + strconv.Itoa(4+i*4)
		str_n := "(" + str1 + str2 + str3 + str4 + ")"
		valueStrings = append(valueStrings, str_n)
		valueArgs = append(valueArgs, elem.MatchingId)
		valueArgs = append(valueArgs, elem.ResultId)
		valueArgs = append(valueArgs, elem.KeywordId)
		valueArgs = append(valueArgs, timeNow)
	}
	smt := `INSERT INTO matches (
                matchingid, resultid, keywordid, createdat) 
            VALUES %s ON CONFLICT DO NOTHING` //
	smt = fmt.Sprintf(smt, strings.Join(valueStrings, ","))

	_, err := Db.Exec(smt, valueArgs...)
	if err != nil {
		panic(err.Error())
	}
}
