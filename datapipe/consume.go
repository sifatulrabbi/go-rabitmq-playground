package datapipe

import (
	"log"
)

func Subscribe(sessions chan chan Session, queue, routingKey string) {
	for session := range sessions {
		sub := <-session

		if _, err := sub.Ch.QueueDeclare(queue, false, true, true, false, nil); err != nil {
			log.Printf("cannot consume from exclusive queue: %s, %v\n", queue, err)
			return
		}

		if err := sub.Ch.QueueBind(queue, routingKey, sub.ExchangeKey, false, nil); err != nil {
			log.Printf("cannot consume without a binding to exchange: %s, %v\n", sub.ExchangeKey, err)
			return
		}

		deliveries, err := sub.Ch.Consume(queue, "", false, true, false, false, nil)
		if err != nil {
			log.Printf("unable to consume from: %s, %v\n", queue, err)
		}

		log.Println("subscribing...")
		for msg := range deliveries {
			sub.Ch.Ack(msg.DeliveryTag, false)
		}
	}
}
