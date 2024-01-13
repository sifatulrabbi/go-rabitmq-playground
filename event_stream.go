package main

import (
	"context"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type EventStream struct {
	URL    string
	Events []string
}

func (es *EventStream) Send(ev string, data []byte) error {
	conn, err := amqp.Dial(es.URL)
	if err != nil {
		return err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(ev, false, false, false, false, nil)
	if err != nil {
		return err
	}

	msg := amqp.Publishing{
		ContentType: "application/json",
		Body:        data,
	}
	err = ch.PublishWithContext(context.Background(), "", q.Name, false, false, msg)
	return err
}

func (es *EventStream) Consume(ev string, callback func(data []byte)) {
	conn, err := amqp.Dial(es.URL)
	if err != nil {
		log.Fatalln("unable to connect to MQ", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalln("failed to create channel", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(ev, false, false, false, false, nil)
	if err != nil {
		log.Fatalln("unable to declare the queue", err)
	}

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Fatalln("Unable to consume", err)
	}
	for d := range msgs {
		callback(d.Body)
	}
}
