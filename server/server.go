package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/sifatulrabbi/go-rabitmq-playground/eventstream"
)

func StartServer(port string, es *eventstream.EventStream) {
	mux := http.NewServeMux()

	mux.HandleFunc("/hello", enableCORS(handleHello(es)))
	mux.HandleFunc("/crypto-price", enableCORS(handleCryptoPrice))
	mux.HandleFunc("/jobs", enableCORS(handleJobsRoute(es)))

	fmt.Printf("Server starting at: %s\n", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), mux); err != nil {
		panic(err)
	}
}

func handleHello(es *eventstream.EventStream) func(w http.ResponseWriter, r *http.Request) {
	handleGet := func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(map[string]string{"message": "Route is not yet implemented"}, w)
	}

	handlePost := func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		if b, err := io.ReadAll(r.Body); err != nil {
			msg := fmt.Sprintf("unable to post: %s", err)
			jsonResponse(map[string]string{"message": msg}, w)
		} else if err := es.Send("hello", b); err != nil {
			msg := fmt.Sprintf("unable to post: %s", err)
			jsonResponse(map[string]string{"message": msg}, w)
		} else {
			jsonResponse(map[string]string{"message": "Post done"}, w)
		}
	}

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			handlePost(w, r)
		} else {
			handleGet(w, r)
		}
	}
}

func jsonResponse(data interface{}, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	if b, err := json.Marshal(data); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("{\"message\": \"invalid json resonse from server\", \"success\": false}"))
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write(b)
	}
}

func enableCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")                   // Allow any origin
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS") // Allowed methods
		w.Header().Set("Access-Control-Allow-Headers", "*")                  // Allow any headers
		next(w, r)
	}
}
