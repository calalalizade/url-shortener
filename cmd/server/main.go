package main

import (
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
	dbConn, err := platform.ConnectDB(cfg.DB)
	if err != nil {
		log.Fatal("database connection failed:", err)
	}
	defer dbConn.Close()

	// ----> Redis setup
	var cacheInstance *platform.RedisCache
	if cfg.Cache.Enabled {
		redisClient, err := platform.ConnectRedis(cfg.Redis)
		if err != nil {
			log.Fatal("redis connection failed: ", err)
		}

		cacheInstance = platform.NewRedisCache(redisClient)
	}

	// ----> Shortener setup
	shortenerRepo := shortener.NewRepository(dbConn)
	shortenerService := shortener.NewService(shortenerRepo, cacheInstance, cfg.Cache)
	shortenerHandler := shortener.NewHandler(shortenerService, cfg.BaseUrl)

	// ----> Gin setup
	r := platform.NewRouter()

	api := r.Group("/api/v1")
	shortener.RegisterRoutes(api, shortenerHandler)

	r.Run(":" + cfg.Port)
}
