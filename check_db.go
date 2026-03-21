package main

import (
	"fmt"
	"log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := "host=14.225.217.105 user=postgres password=Abcd1234 dbname=dev_code_practice port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	var columns []string
	db.Raw("SELECT column_name FROM information_schema.columns WHERE table_name = 'subjects'").Scan(&columns)
	fmt.Println("Columns in subjects table:", columns)
}
