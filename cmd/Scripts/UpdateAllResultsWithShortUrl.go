package main

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/teris-io/shortid"
	_ "github.com/lib/pq"
)

var Db *sql.DB

func main() {


	db := DbConnect()
	defer db.Close()
	Db = db
	UpdateShortUrl()
}

func DbConnect() (Db *sql.DB) {
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
	return
}

func UpdateShortUrl() {
	rows, err := Db.Query(`SELECT id FROM results r;`)
	if err != nil {
		panic(err.Error())
	}

	sid, err := shortid.New(1, shortid.DefaultABC, 2342)

	valueStrings := []string{}

	for rows.Next() {
		var id int
		if err = rows.Scan(
			&id); err != nil {
			if err != nil {
				panic(err.Error())
			}
		}
		unique_url_id, err := sid.Generate()
		if err != nil {
			panic(err.Error())
		}

		str_n := fmt.Sprintf("(%d,'%s')", id, unique_url_id)
		valueStrings = append(valueStrings, str_n)

		fmt.Println(id)
	}

	smt := `UPDATE results AS t SET
			    urlshort = c.urlshort
			FROM (VALUES %s ) AS c(id, urlshort) 
			WHERE c.id = t.id;`



	smt = fmt.Sprintf(smt, strings.Join(valueStrings, ","))

	_, err = Db.Exec(smt)
	if err != nil {
		panic(err.Error())
	}

	rows.Close()


	return
}