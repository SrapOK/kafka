package main

import (
	"flag"
	"log"
	"net"
	"net/http"
	"post-service/internal/config"
	"post-service/internal/handler"

	"github.com/gin-gonic/gin"
	"github.com/segmentio/kafka-go"
)

func main() {
	configPath := flag.String("config", "config.yaml", "path to config")
	flag.Parse()

	cfg := config.Load(*configPath)

	log.Println(*cfg)

	addr, err := net.ResolveTCPAddr("tcp", cfg.Kafka.Addr)

	if err != nil {
		log.Fatalln(err)
	}

	kfk := &kafka.Writer{
		Addr:                   addr,
		Topic:                  cfg.Kafka.Topic,
		AllowAutoTopicCreation: true,
	}

	router := gin.Default()
	h := handler.Handler{Client: kfk}

	router.POST("/", h.PostPerson)

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
