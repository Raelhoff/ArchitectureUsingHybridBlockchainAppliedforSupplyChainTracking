package configs

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/wpcodevo/golang-fiber-mysql/models"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

func ConnectDB() *mongo.Client {
	client, err := mongo.NewClient(options.Client().ApplyURI(EnvMongoURI()))
	if err != nil {
		log.Fatal(err)
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	//ping the database
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("-------------------------------------------------------------------------------------")
	fmt.Println("Connected to MongoDB")
	return client
}

// Client instance
var DB *mongo.Client = ConnectDB()

// getting database collections
func GetCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	collection := client.Database("golangAPI").Collection(collectionName)
	return collection
}

//---------------------------------------------------------------------------------------------------------------------------
////////////////////////////////////// SQLITE////////////////////////////////////

var DBSQLITE *gorm.DB

func ConnectDB1() {
	var err error

	DBSQLITE, err = gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
	if err != nil {
		fmt.Println("-------------------------------------------------------------------------------------")
		log.Fatal("Failed to connect to the Database (gorm.db)! \n", err.Error())
		os.Exit(1)
	}

	DBSQLITE.Logger = logger.Default.LogMode(logger.Info)

	log.Println("Running Migrations")
	DBSQLITE.AutoMigrate(&models.Note{})
	fmt.Println("-------------------------------------------------------------------------------------")
	log.Println("🚀 Connected Successfully to the Database(gorm.db)")
}

// ---------------------------------------------------------------------------------------------------------------------------
var DB2 *sql.DB = ConnectDB2()

func ConnectDB2() *sql.DB {
	db, err := sql.Open("sqlite3", "./gorm.db")
	if err != nil {
		fmt.Println("-------------------------------------------------------------------------------------")
		log.Fatal("Failed to connect to the Database! \n", err.Error())
		os.Exit(1)
	}

	return db
}
