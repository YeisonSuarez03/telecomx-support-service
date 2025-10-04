package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"

	"telecomx-support-service/internal/application/service"
	"telecomx-support-service/internal/config"
	"telecomx-support-service/internal/infrastructure/adapter/event/listener"
	"telecomx-support-service/internal/infrastructure/adapter/repository"
	"telecomx-support-service/internal/infrastructure/adapter/rest"
)

func main() {
	cfg := config.InstanceConfig()
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		log.Fatal(err)
	}
	db := client.Database("telecomx_support")

	repo := repository.NewMongoRepository(db)
	svc := service.NewSupportService(repo)

	listener.StartKafkaListener(svc, cfg.Brokers, cfg.Topic, cfg.Group, cfg.Client)

	mux := http.NewServeMux()
	rest.NewSupportHandler(svc).RegisterRoutes(mux)

	fmt.Printf("REST server running on :%s\n", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, mux))
}
