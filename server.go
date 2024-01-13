package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func startHTTPServer(port string, es *EventStream) {
	mux := http.NewServeMux()

	mux.HandleFunc("/hello", handleHello(es))

	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), mux); err != nil {
		panic(err)
	} else {
		fmt.Printf("Server serving at port: %s\n", port)
	}
}

func handleHello(es *EventStream) func(w http.ResponseWriter, r *http.Request) {
	handleGet := func(w http.ResponseWriter, r *http.Request) {
		jsonResponse(
			map[string]string{
				"message": "Route is not yet implemented",
			},
			w,
		)
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
