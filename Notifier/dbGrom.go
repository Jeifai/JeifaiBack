package main

import (
    "fmt"
    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/postgres"
)

type User struct {
    ID            int
    Username      string
}
type Session struct {
    ID          int
    Email       string
    Userid      int
    User        User `gorm:"ForeignKey:Userid;AssociationForeignKey:ID"`
}

func main() {
    db, _ := gorm.Open("postgres", "host=35.223.132.23 port=5432 user=postgres dbname=jeifai password=jeifai")
    var s []Session
    db.Preload("User").Find(&s)

    fmt.Println(s)
}