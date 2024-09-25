package main

import (
	"context"
	"flag"
	"get-service/internal/client"
	"get-service/internal/config"
	"get-service/internal/handler"
	"get-service/internal/model/person"
	storage "get-service/internal/storage/redis"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func main() {
	configPath := flag.String("config", "config.yaml", "path to config")
	flag.Parse()

	cfg := config.Load(*configPath)
	log.Println(*cfg)

	db, err := storage.New(&redis.Options{Addr: cfg.Redis.Addr})
	if err != nil {
		log.Fatalln(err)
	}
	repo := person.PersonRepo{Db: db}

	go func() {
		reader := client.GetKafkaReader(cfg.Kafka.Addr, cfg.Kafka.Topic)
		defer reader.Close()

		for {
			m, err := reader.ReadMessage(context.Background())
			if err != nil {
				log.Fatalln(err)
			}

			ps, err := person.UnserializePersonDto(m.Value)
			if err != nil {
				log.Printf("cannot unserialize msg value to person: %s", err.Error())
			}

			ctx := context.Background()
			if err := repo.SavePerson(&ctx, ps); err != nil {
				log.Println(err.Error())
			}
		}
	}()

	router := gin.Default()
	h := handler.Handler{Repo: &repo}

	router.GET("/:id", h.GetPerson)
	router.GET("/", h.GetAllPersons)
	router.GET("/random", h.GetRandomPerson)

	server := &http.Server{
		Addr:         cfg.HttpServer.Addr,
		Handler:      router,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Failed to start server: %s", err.Error())
	}

	log.Printf("Server stopped")
}
