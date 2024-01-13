package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
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

	fmt.Println("starting the server...")
	startHTTPServer(PORT, &eventStream)
}
