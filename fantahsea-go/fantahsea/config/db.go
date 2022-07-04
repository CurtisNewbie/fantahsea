package config

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const (
	// hikari recommends 1800000, so we do the same thing 
	ConnMaxLifeTime = time.Minute * 30
	MaxOpenConns = 10
	MaxIdleConns = MaxOpenConns // recommended to be the same as the maxOpenConns
)

var (
	// Global handle to the database
	DB *sql.DB
)

/* Init Handle to the database */
func InitDb() (error) {
	db, err := sql.Open("mysql", "user:password@/dbname")
	if err != nil {
		log.Printf("Failed to Open DB Handle, err: %v", err)
		panic(err)
	}
	db.SetConnMaxLifetime(ConnMaxLifeTime) 
	db.SetMaxOpenConns(MaxOpenConns)
	db.SetMaxIdleConns(MaxIdleConns)

	ping_err := db.Ping()
	if ping_err != nil {
		log.Printf("Ping DB Error, %v, connection may not be established\n", ping_err)
		return ping_err 
	}

	DB = db
	return nil 
}