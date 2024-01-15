package datapipe

import (
	"context"
	"encoding/json"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func Publish(sessions chan chan Session, routingKey string, body interface{}) error {
	// convert the body in to bytes
	b, err := json.Marshal(body)
	if err != nil {
		return err
	}
	msg := amqp.Publishing{
		ContentType: "application/json",
		Body:        b,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for session := range sessions {
		confirm := make(chan amqp.Confirmation, 1)

		pub := <-session

		if err := pub.Ch.Confirm(true); err != nil {
			log.Printf("publishers confirm not supported: %v\n", err)
			close(confirm)
		} else {
			pub.Ch.NotifyPublish(confirm)
		}

		log.Println("publishing...")

		confirmed, ok := <-confirm
		if !ok {
			break
		}
		if !confirmed.Ack {
			log.Printf("nack message %d, body: %v\n", confirmed.DeliveryTag, body)
		}

		if err := pub.Ch.PublishWithContext(ctx, pub.ExchangeKey, routingKey, false, false, msg); err != nil {
			log.Printf("unable to publish the message: %v\n", err)
		} else {
			log.Println("message published")
		}
	}

	return nil
}
