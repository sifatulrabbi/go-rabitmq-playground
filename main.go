package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/sifatulrabbi/go-rabitmq-playground/datapipe"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	consumers := []*datapipe.Consumer{}

	go func() {
		ch := make(chan datapipe.Bytes)

		consumer := datapipe.Consumer{
			Name:      "test_consumer_1",
			EventName: "test",
			Chan:      ch,
		}
		consumers = append(consumers, &consumer)

		for {
			select {
			case data := <-ch:
				fmt.Println(data)
			}
		}
	}()

	go func() {
		ch := make(chan datapipe.Bytes)

		consumer := datapipe.Consumer{
			Name:      "test_consumer_2",
			EventName: "ping",
			Chan:      ch,
		}
		consumers = append(consumers, &consumer)

		for {
			select {
			case data := <-ch:
				fmt.Println(data)
			}
		}
	}()

	go datapipe.NewCustomDataPipe(ctx, consumers)

	for {
		log.Println("publishing info")

		if err := datapipe.PublishToCustom("test", map[string]string{"message": "Testing"}); err != nil {
			log.Fatalln(err)
		}
		if err := datapipe.PublishToCustom("ping", map[string]string{"message": "pinging"}); err != nil {
			log.Fatalln(err)
		}

		time.Sleep(2 * time.Second)
	}
}

// func main() {
// 	if err := godotenv.Load(); err != nil {
// 		panic(err)
// 	}
//
// 	PORT := os.Getenv("PORT")
// 	RABITMQ_URL := os.Getenv("RABITMQ_URL")
// 	if PORT == "" || RABITMQ_URL == "" {
// 		panic(errors.New("`PORT` and `RABITMQ_URL` envs are required to run the server."))
// 	}
//
// 	eventStream := eventstream.EventStream{
// 		URL:    RABITMQ_URL,
// 		Events: []string{},
// 	}
//
// 	go eventStream.Consume("hello", func(d []byte) {
// 		data := map[string]string{}
// 		if err := json.Unmarshal(d, &data); err != nil {
// 			log.Println(err)
// 		} else {
// 			log.Println(data)
// 		}
// 	})
//
// 	s := scheduler.Scheduler{
// 		LastTaskIndex: 0,
// 		Tasks:         []scheduler.Task{},
// 		Executioners:  map[string]func(body map[string]string) error{},
// 	}
//
// 	s.AddNewExecutioner("printHello", func(body map[string]string) error {
// 		fmt.Println("Hello world")
// 		return nil
// 	})
// 	s.AddNewExecutioner("scheduleTask", func(body map[string]string) error {
// 		fmt.Println("Scheduled task is running with body:", body)
// 		return nil
// 	})
//
// 	// s.AddNewTask(scheduler.NewSchedulerTask("print hello", "printHello", time.Now().Add(time.Second*20), map[string]string{}))
//
// 	go s.StartInBg(s.LastTaskIndex)
// 	server.StartServer(PORT, &eventStream)
// }
