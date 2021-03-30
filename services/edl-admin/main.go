package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
)

var ctx = context.Background()

type Request struct {
	Id     string `json:"id"`
	User   string `json:"user" redis:"user"`
	Url    string `json:"url" redis:"url"`
	Reason string `json:"reason" redis:"reason"`
}

func main() {
	fmt.Println("Web Server Starting...")
	r := mux.NewRouter()
	r.HandleFunc("/health", Health).Methods("GET")
	r.HandleFunc("/requests", Requests).Methods("GET")
	r.HandleFunc("/accept", Accept).Methods("POST")
	r.HandleFunc("/reject", Reject).Methods("POST")
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./assets/")))
	http.ListenAndServe(":80", r)
}

func Health(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func Requests(w http.ResponseWriter, r *http.Request) {
	// Connect to Redis
	redis_host := os.Getenv("REDIS_HOST")
	rdb := redis.NewClient(&redis.Options{
		Addr:     redis_host + ":6379",
		Password: "",
		DB:       0,
	})

	// Get all requests keys
	keys, err := rdb.Keys(ctx, "request:*").Result()
	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// For each request, get all fields and map to Request struct
	var requests []Request
	for _, key := range keys {
		request := Request{}
		reqres := rdb.HGetAll(ctx, key)
		if reqres.Err() != nil {
			fmt.Println(fmt.Errorf("Error: %v", err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if err = reqres.Scan(&request); err != nil {
			fmt.Println(fmt.Errorf("Error: %v", err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		request.Id = strings.TrimPrefix(key, "request:")
		requests = append(requests, request)
	}

	// Build & Send JSON response with array of request objects
	requestListBytes, err := json.Marshal(requests)
	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(requestListBytes)
}

func Accept(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	requestid := r.Form.Get("requestid")

	// Connect to Redis
	redis_host := os.Getenv("REDIS_HOST")
	rdb := redis.NewClient(&redis.Options{
		Addr:     redis_host + ":6379",
		Password: "",
		DB:       0,
	})

	// Get URL to add to EDL
	url, err := rdb.HGet(ctx, "request:"+requestid, "url").Result()
	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Add URL to EDL Set
	addres, err := rdb.SAdd(ctx, "edl", url).Result()
	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if addres == 0 {
		fmt.Printf("Error: URL already exists in EDL Set [%v]\n", requestid)
	}

	// Delete Request
	delres, err := rdb.HDel(ctx, "request:"+requestid, "user", "url", "reason").Result()
	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if delres != 3 {
		fmt.Printf("Error: Delete failed [%v]\n", requestid)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Go back to Home
	http.Redirect(w, r, "/", http.StatusFound)
}

func Reject(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	requestid := r.Form.Get("requestid")

	// Connect to Redis
	redis_host := os.Getenv("REDIS_HOST")
	rdb := redis.NewClient(&redis.Options{
		Addr:     redis_host + ":6379",
		Password: "",
		DB:       0,
	})

	// Delete Request
	delres, err := rdb.HDel(ctx, "request:"+requestid, "user", "url", "reason").Result()
	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if delres != 3 {
		fmt.Printf("Error: Delete failed [%v]\n", requestid)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Go back to Home
	http.Redirect(w, r, "/", http.StatusFound)
}
