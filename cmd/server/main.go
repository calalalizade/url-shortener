package main

import (
	"fmt"
	"log"

	"github.com/calalalizade/url-shortener/internal/platform"
	"github.com/calalalizade/url-shortener/internal/shortener"
)

type Shortener struct {
	ShortUrl string `json:"shortUrl"`
	LongUrl  string `json:"longUrl"`
}

func main() {
	// ----> Config setup
	cfg := platform.LoadConfig()

	// ----> DB setup
	db, err := platform.ConnectDB(cfg.DB)
	if err != nil {
		log.Fatal("database connection failed:", err)
	}
	defer db.Close()

	// ----> Shortener setup
	shortenerRepo := shortener.NewRepository(db)
	shortenerService := shortener.NewService(shortenerRepo)
	shortenerHandler := shortener.NewHandler(shortenerService, cfg.BaseUrl)

	// ----> Gin setup
	r := platform.NewRouter()

	api := r.Group("/api/v1")
	shortener.RegisterRoutes(api, shortenerHandler)

	fmt.Println("port: ---->", cfg.Port)
	r.Run(":" + cfg.Port)
}
