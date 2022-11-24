package main

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/practice2311/api"
	"github.com/practice2311/config"
	"github.com/practice2311/storage"

	"github.com/go-redis/redis/v9"
)

func main() {
	cfg := config.Load(".")

	psqlUrl := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.PostConfig.Host,
		cfg.PostConfig.Port,
		cfg.PostConfig.User,
		cfg.PostConfig.Password,
		cfg.PostConfig.Database,
	)

	psqlConn, err := sqlx.Connect("postgres", psqlUrl)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	strg := storage.NewStoragePg(psqlConn)

	rdb := redis.NewClient(&redis.Options{
		Addr: cfg.RedisConfig.RedisHost + ":" + cfg.RedisConfig.RedisPort,
	})

	strg = storage.NewStoragePg(psqlConn)
	inMemory := storage.NewInMemoryStorage(rdb)

	apiServer := api.New(api.RouterOptions{
		Cfg:      &cfg,
		Storage:  strg,
		InMemory: inMemory,
	})

	err = apiServer.Run(cfg.HttpPort)
	if err != nil {
		log.Fatalf("failed to run server: %v", err)
	}

	// log.Print("Server stopped")
}
