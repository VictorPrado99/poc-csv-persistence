package database

import (
	"log"

	"github.com/VictorPrado99/poc-csv-persistence/pkg/api"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Hold the instance who will be accessed
var Instance *gorm.DB
var err error

// Connect to database and create a gorm instance
func Connect(connectionString string) {
	Instance, err = gorm.Open(mysql.Open(connectionString), &gorm.Config{
		QueryFields: true,
	})
	if err != nil {
		log.Fatal(err)
		panic("Cannot connect to DB")
	}
	log.Println("Connected to Database...")
}

// Auto generate the entities at database
func Migrate() {
	Instance.AutoMigrate(&api.Order{})
	log.Println("Database Migration Completed...")
}
