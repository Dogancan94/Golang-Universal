package main

import (
	"log"
	"net/http"

	"github.com/Dogancan94/toolkit"
)

type RequestPayload struct {
	Action  string `json:"action"`
	Message string `json:"message"`
}

type ResponsePayload struct {
	Message    string `json:"message"`
	StatusCode int    `json:"status_code,omitempty"`
}

func main() {
	mux := http.NewServeMux()
	mux.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("."))))
	mux.HandleFunc("/receive-post", receivePost)
	mux.HandleFunc("/remote-service", remoteService)
	log.Println("Starting service....")

	err := http.ListenAndServe(":8081", mux)
	if err != nil {
		log.Fatal(err)
	}

}

func receivePost(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload
	var t toolkit.Tools

	err := t.ReadJson(w, r, &requestPayload)
	if err != nil {
		t.ErrorJSON(w, err)
		return
	}

	responsePayload := ResponsePayload{
		Message: "hit the handler okay, and sending response",
	}

	err = t.WriteJSON(w, http.StatusAccepted, responsePayload)
	if err != nil {
		log.Println(err)
	}
}

func remoteService(w http.ResponseWriter, r *http.Request) {

}
