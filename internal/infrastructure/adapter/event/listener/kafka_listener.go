package listener

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"log"
	"time"

	"telecomx-support-service/internal/application/service"
	"telecomx-support-service/internal/domain/model"
)

type CustomerEvent struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type UserPayload struct {
	UserID    string `json:"userId"`
	Email     string `json:"email,omitempty"`
	Suspended bool   `json:"suspended,omitempty"`
	Deleted   bool   `json:"deleted,omitempty"`
}

func StartKafkaListener(svc *service.SupportService, brokers []string, topic, group, client string) error {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   topic,
		GroupID: group,
		Dialer: &kafka.Dialer{
			ClientID: client,
		},
	})
	defer reader.Close()

	log.Printf("[Kafka] Listening on topic: %s", topic)

	for {
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Println("Kafka error:", err)
			return err
		}

		var event CustomerEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Println("Invalid event:", err)
			continue
		}

		var payload UserPayload
		_ = json.Unmarshal(event.Payload, &payload)

		switch event.Type {
		case "Customer.Created":
			err := svc.Create(context.Background(), &model.Support{
				UserID:     payload.UserID,
				Issue:      "Customer created from Telecomx",
				Status:     "finish",
				CreatedAt:  time.Now(),
				ResolvedAt: time.Now(),
			})
			if err != nil {
				log.Println("Error creating customer:", err)
				return err
			}
		case "Customer.Updated":
			if payload.Deleted || payload.Suspended {
				err := svc.Create(context.Background(), &model.Support{
					UserID:     payload.UserID,
					Issue:      "Customer canceled support tickets from Telecomx",
					Status:     "finish",
					CreatedAt:  time.Now(),
					ResolvedAt: time.Now(),
				})
				if err != nil {
					log.Println("Error creating customer:", err)
					return err
				}
			}
		case "Customer.Deleted":
			err := svc.Delete(context.Background(), payload.UserID)
			if err != nil {
				log.Println("Error deleting customer:", err)
				return err
			}
		}
	}
}
