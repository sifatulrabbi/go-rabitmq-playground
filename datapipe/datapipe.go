package datapipe

import (
	"context"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Session struct {
	ExchangeKey string
	Conn        *amqp.Connection
	Ch          *amqp.Channel
}

func RedialMQ(ctx context.Context, url, exchange string) chan chan Session {
	sessions := make(chan chan Session)

	go func() {
		session := make(chan Session)
		defer close(session)

		for {
			select {
			case sessions <- session:
			case <-ctx.Done():
				log.Println("shutting down sessions factory")
				return
			}

			conn, err := amqp.Dial(url)
			if err != nil {
				log.Fatalf("cannot (re)dial: %s, %v\n", url, err)
			}

			ch, err := conn.Channel()
			if err != nil {
				log.Fatalf("cannot create channel: %v\n", err)
			}

			if err := ch.ExchangeDeclare(exchange, "fanout", false, true, false, false, nil); err != nil {
				log.Fatalf("unable to declare fanout exchange: %v\n", err)
			}

			select {
			case session <- Session{Conn: conn, Ch: ch, ExchangeKey: exchange}:
			case <-ctx.Done():
				log.Println("shutting down the sessions factory")
				return
			}
		}
	}()

	return sessions
}
