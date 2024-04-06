package database

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/Oik17/FamPay-BE/utils"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // Import the PostgreSQL driver
)

// Database instance
type Dbinstance struct {
	Db *sqlx.DB
}

var DB Dbinstance

// Connect function
func Connect() {
	p := utils.Config("DB_PORT")
	port, err := strconv.Atoi(p)
	if err != nil {
		fmt.Println("Error parsing str to int")
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai", utils.Config("DB_HOST"), utils.Config("DB_USER"), utils.Config("DB_PASSWORD"), utils.Config("DB_NAME"), port)

	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Failed to connect to database. \n", err)
		os.Exit(2)
	}

	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping the database. \n", err)
		os.Exit(2)
	}

	log.Println("Connected")

	runMigrations(db)

	DB = Dbinstance{
		Db: db,
	}
}

func runMigrations(db *sqlx.DB) {
	_, err := db.Exec(`
	
		CREATE TABLE IF NOT EXISTS video (
			id UUID PRIMARY KEY,
			prompt VARCHAR(255),
			title VARCHAR(255) NOT NULL,
			description TEXT,
			channelTitle VARCHAR(255),
			publishedAt TIMESTAMP,
			thumbnail VARCHAR(255),
			url VARCHAR(255)
		);

	`)
	if err != nil {
		log.Fatal("Failed to run migrations. \n", err)
		os.Exit(2)
	}

	log.Println("Migrations completed")
}
