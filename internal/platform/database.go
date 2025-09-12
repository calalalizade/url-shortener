package platform

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)


func ConnectDB(cfg DbConfig) (*sql.DB, error) {
	psqlInfo := fmt.Sprintf(
		"host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		cfg.Host, 
		cfg.Port, 
		cfg.User, 
		cfg.Pass, 
		cfg.Name,
	)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(30 * time.Minute)

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}