package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)

const (
	host     = "127.0.0.1:9011"
	randSize = 10
)

type response struct {
	Message int    `json:"message"`
	Status  string `json:"status"`
}

const (
	SUCCESS = "success"
	FAIL    = "fail"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		vRand := sRand()
		status := SUCCESS
		//time.Sleep(time.Duration(vRand) * time.Second)
		if vRand > 3 {
			w.WriteHeader(http.StatusInternalServerError)
			status = FAIL
		} else {
			w.WriteHeader(http.StatusOK)
		}
		resp := response{
			Message: vRand,
			Status:  status,
		}

		body, _ := json.Marshal(resp)
		fmt.Println(resp)
		_, _ = w.Write(body)
	})

	log.Fatal(http.ListenAndServe(host, nil))
}

func sRand() int {
	rand.Seed(time.Now().UnixNano())

	return rand.Intn(randSize)
}
