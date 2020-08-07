package main

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	. "github.com/logrusorgru/aurora"
)

var Db *sql.DB

func main() {
	DbConnect()
	linkedinDict := GenerateLinkedinDict()
	saveStringToFile(linkedinDict)
}

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

func GenerateLinkedinDict() (linkedinDict string) {
	err := Db.QueryRow(`
    SELECT
        (json_agg(t))::text
    FROM (
        SELECT
            id,
            linkedinurl
        FROM targets
        WHERE linkedinurl IS NOT NULL
    ) t;`).Scan(&linkedinDict)
	if err != nil {
		panic(err.Error())
	}
	return strings.ReplaceAll(
		strings.ReplaceAll(
			linkedinDict, " ", ""), "\n", "")
}

func saveStringToFile(linkedinDict string) {
	f, err := os.Create("linkedinDict.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = f.WriteString(linkedinDict)
	if err != nil {
		fmt.Println(err)
		f.Close()
		return
	}
	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
}
