package main

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type User struct {
	ID       int
    Username string
    Session []Session `gorm:"ForeignKey:Userid"`
}

type Session struct {
	ID     int
	Email  string
	Userid int
}

func main() {
	db, _ := gorm.Open("postgres", "host=35.223.132.23 port=5432 user=postgres dbname=jeifai password=jeifai")
	var u []User
	db.Preload("Session").Find(&u)

	fmt.Println(u)
}
