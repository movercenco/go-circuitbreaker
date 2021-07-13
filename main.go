package main

import (
	"fmt"
	"github.com/sony/gobreaker"
	"io"
	"log"
	"net/http"
	"time"
)

var (
	host       = "127.0.0.1:9010"
	downstream = "http://127.0.0.1:9011/"
)

func main() {
	breaker := breaker()
	client := client()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		body, err := do(client, breaker)
		if err != nil {
			fmt.Println(err.Error())
		}
		_, _ = w.Write(body)
	})

	log.Fatal(http.ListenAndServe(host, nil))
}

func client() http.Client {
	return http.Client{
		Timeout: 1 * time.Second,
	}
}

func breaker() *gobreaker.CircuitBreaker {
	settings := gobreaker.Settings{
		Name:    "HTTP GET",
		Timeout: 10 * time.Second,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
			fmt.Println("--------")
			fmt.Println("TotalFailures: ", counts.TotalFailures)
			fmt.Println("TotalSuccesses: ", counts.TotalSuccesses)
			fmt.Println("Requests: ", counts.Requests)
			fmt.Println("failureRatio: ", failureRatio)
			fmt.Println("--------")
			return counts.Requests >= 3 && failureRatio >= 0.5
		},
	}

	return gobreaker.NewCircuitBreaker(settings)
}

func do(client http.Client, breaker *gobreaker.CircuitBreaker) ([]byte, error) {
	body, err := breaker.Execute(func() (interface{}, error) {
		resp, err := client.Get(downstream)
		if err != nil {
			return []byte(`{"message": "timeout"}`), fmt.Errorf("timeout: %v", err.Error())
		}

		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)

		if resp != nil && resp.StatusCode >= http.StatusInternalServerError {
			return body, fmt.Errorf("http response error: %v", resp.StatusCode)
		}

		return body, nil
	})

	if err != nil && body == nil {
		return []byte(`{"message": "circuit breaker is open"}`), err
	}

	return body.([]byte), err
}
