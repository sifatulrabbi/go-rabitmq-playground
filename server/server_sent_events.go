package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/sifatulrabbi/go-rabitmq-playground/eventstream"
)

func mockCryptoPrice(ctx context.Context, priceCh chan<- int) {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	ticker := time.NewTicker(time.Second)

outerloop:
	for {
		select {
		case <-ctx.Done():
			break outerloop
		case <-ticker.C:
			priceCh <- r.Intn(1000)
		}
	}

	ticker.Stop()
	close(priceCh)
}

func handleCryptoPrice(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		msg := "SSE not supported by the client"
		http.Error(w, msg, http.StatusInternalServerError)
		log.Println(msg)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	priceCh := make(chan int)
	go mockCryptoPrice(r.Context(), priceCh)

	for price := range priceCh {
		evPayload, err := formatSSEData("crypto_price", price)
		if err != nil {
			log.Println("Error while formating event data", err)
			break
		}

		_, err = fmt.Fprint(w, evPayload)
		if err != nil {
			log.Println("Error while sending data to the client:", err)
			break
		}

		flusher.Flush()
	}
}

func handleJobsRoute(es *eventstream.EventStream) func(w http.ResponseWriter, r *http.Request) {
	getJobs := func(w http.ResponseWriter, r *http.Request) {
		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "SSE not supported", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Connection", "keep-alive")

		es.Consume("jobs", func(b []byte) {
			data, err := formatSSEData("jobs", string(b))
			if err != nil {
				// TODO: close the http connection
				return
			}

			_, err = fmt.Fprint(w, data)
			if err != nil {
				return
			}

			fmt.Println("sending job info to client", string(data))
			flusher.Flush()
		})
	}

	postJob := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		defer r.Body.Close()

		b, err := io.ReadAll(r.Body)
		if err != nil {
			jsonResponse(map[string]string{"error": err.Error()}, w)
			return
		}

		if err := es.Send("jobs", b); err != nil {
			jsonResponse(map[string]string{"error": err.Error()}, w)
			return
		}

		fmt.Println("New job created", string(b))
		jsonResponse(map[string]string{"message": "Job created"}, w)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("%s: %s\n", r.Method, r.URL.Path)

		switch r.Method {
		case http.MethodGet:
			getJobs(w, r)
			break
		case http.MethodPost:
			postJob(w, r)
			break
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			fmt.Fprint(w, "Method not allowed")
			break
		}
	}
}

// formatSSEData takes name of an event and any kind of data and transforms
// into a server sent event payload structure.
// Data is sent as a json object, { "data": <your_data> }.
//
// Example:
//
//	Input:
//		event="price-update"
//		data=10
//	Output:
//		event: price-update\n
//		data: "{\"data\":10}"\n\n
func formatSSEData(ev string, data any) (string, error) {
	m := map[string]any{
		"data": data,
	}

	b, err := json.Marshal(m)
	if err != nil {
		return "", err
	}

	sp := strings.Builder{}
	sp.WriteString(fmt.Sprintf("event:%s\n", ev))
	sp.WriteString(fmt.Sprintf("data:%s\n\n", b))

	return sp.String(), nil
}
