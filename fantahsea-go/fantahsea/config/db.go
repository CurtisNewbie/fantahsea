package config

import (
	"database/sql"
	"fmt"
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
func InitDbFromConfig(config *DBConfig) (error) {
	return InitDb(config.User, config.Password, config.Database)
}

/* Init Handle to the database */
func InitDb(user string, password string, dbname string) (error) {

	db, err := sql.Open("mysql", fmt.Sprintf("%v:%v@/%v", user, password, dbname))
	if err != nil {
		log.Printf("Failed to Open DB Handle, err: %v", err)
		return err
	}

	db.SetConnMaxLifetime(ConnMaxLifeTime) 
	db.SetMaxOpenConns(MaxOpenConns)
	db.SetMaxIdleConns(MaxIdleConns)

	err = db.Ping() // make sure the handle is actually connected 
	if err != nil {
		log.Printf("Ping DB Error, %v, connection may not be established", err)
		return err 
	}

	log.Println("DB Handle initialized")

	DB = db

	return nil 
}