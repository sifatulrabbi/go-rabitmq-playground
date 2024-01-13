package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"

	"github.com/sifatulrabbi/go-rabitmq-playground/scheduler"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	PORT := os.Getenv("PORT")
	RABITMQ_URL := os.Getenv("RABITMQ_URL")
	if PORT == "" || RABITMQ_URL == "" {
		panic(errors.New("`PORT` and `RABITMQ_URL` envs are required to run the server."))
	}

	eventStream := EventStream{
		URL:    RABITMQ_URL,
		Events: []string{},
	}

	go eventStream.Consume("hello", func(d []byte) {
		data := map[string]string{}
		if err := json.Unmarshal(d, &data); err != nil {
			log.Println(err)
		} else {
			log.Println(data)
		}
	})

	s := scheduler.Scheduler{
		LastTaskIndex: 0,
		Tasks:         []scheduler.Task{},
		Executioners:  map[string]func(body map[string]string) error{},
	}
	s.AddNewExecutioner("printHello", func(body map[string]string) error {
		fmt.Println("Hello world")
		return nil
	})
	s.AddNewTask(scheduler.NewSchedulerTask("print hello", "printHello", time.Now().Add(time.Second*20), map[string]string{}))
	go s.StartInBg(s.LastTaskIndex)

	forever := make(chan string)
	<-forever

	// fmt.Println("starting the server...")
	// startHTTPServer(PORT, &eventStream)
}
