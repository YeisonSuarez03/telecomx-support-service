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
            // Try fallback shape {"event": "...", "data": {...}}
            var alt AltCustomerEvent
            if err2 := json.Unmarshal(msg.Value, &alt); err2 != nil {
                log.Println("Invalid event (both shapes failed):", err, "|", err2)
                continue
            }
            event.Type = alt.Event
            event.Payload = alt.Data
        } else if event.Type == "" && len(event.Payload) == 0 {
            // Some producers may use different keys; attempt fallback even if first unmarshal succeeded but empty
            var alt AltCustomerEvent
            if err2 := json.Unmarshal(msg.Value, &alt); err2 == nil {
                event.Type = alt.Event
                event.Payload = alt.Data
            }
        }

        // Log parsed event type and a trimmed payload preview
        log.Printf("[Kafka] Parsed event type=%s payload_len=%d", event.Type, len(event.Payload))
        // Log full event JSON for debugging
        log.Printf("[Kafka] Event JSON=%s", toJSON(event))

		var payload UserPayload
		_ = json.Unmarshal(event.Payload, &payload)
        // Log full payload JSON for debugging
        log.Printf("[Kafka] Payload JSON=%s", toJSON(payload))

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
