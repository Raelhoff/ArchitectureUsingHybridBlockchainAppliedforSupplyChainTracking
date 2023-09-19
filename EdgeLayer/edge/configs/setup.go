package initializers

import (
	"log"
	"os"

	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"github.com/wpcodevo/golang-fiber-mysql/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDB() {
	var err error

	DB, err = gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the Database! \n", err.Error())
		os.Exit(1)
	}

	DB.Logger = logger.Default.LogMode(logger.Info)

	log.Println("Running Migrations")
	DB.AutoMigrate(&models.Note{})

	log.Println("🚀 Connected Successfully to the Database")
}

var DB2 *sql.DB = ConnectDB2()

func ConnectDB2() *sql.DB {
	db, err := sql.Open("sqlite3", "./gorm.db")
	if err != nil {
		log.Fatal("Failed to connect to the Database! \n", err.Error())
		os.Exit(1)
	}

	return db
}
