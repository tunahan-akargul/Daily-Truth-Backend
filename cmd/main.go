package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	// chi -- 
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SignUpRequest struct {
	Name     string `json:"name"`
	Surname  string `json:"surname"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignInRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

var peopleColl *mongo.Collection

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil { log.Fatal(err) }
	if err := client.Ping(ctx, nil); err != nil { log.Fatal("mongo ping failed:", err) }

	peopleColl = client.Database("testdb").Collection("people")
	if peopleColl == nil { log.Fatal("peopleColl is nil") }
	fmt.Println("Mongo OK, server on :8080")

	http.HandleFunc("/check-email", CheckEmailHandler)
	http.HandleFunc("/signup", SignUpHandler)
	http.HandleFunc("/signin", SignInHandler)
	http.HandleFunc("/get-details", GetUserDetailsHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func CheckEmailHandler(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	if email == "" {
		http.Error(w, "email required", http.StatusBadRequest)
		return
	}
	if peopleColl == nil {
		http.Error(w, "db not ready", http.StatusInternalServerError)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	count, err := peopleColl.CountDocuments(ctx, bson.M{"email": email})
	if err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]bool{"exists": count > 0})
}

func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	var req SignUpRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad json", http.StatusBadRequest)
		return
	}
	if peopleColl == nil {
		http.Error(w, "db not ready", http.StatusInternalServerError)
		return
	}



	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := peopleColl.InsertOne(ctx, req)
	if err != nil {
		http.Error(w, "db insert error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{"status": "ok", "id": res.InsertedID})
}

func SignInHandler(w http.ResponseWriter, r *http.Request) {
	var req SignInRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad json", http.StatusBadRequest)
		return
	}
	if peopleColl == nil {
		http.Error(w, "db not ready", http.StatusInternalServerError)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user SignUpRequest
	response := peopleColl.FindOne(ctx, bson.M{"email": req.Email})
	if err := response.Decode(&user); err != nil {
		http.Error(w, "user not found", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]bool{"passwordControl": req.Password == user.Password})
}

func GetUserDetailsHandler(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	if email == "" {
		http.Error(w, "email required", http.StatusBadRequest)
		return
	}
	if peopleColl == nil {
		http.Error(w, "db not ready", http.StatusInternalServerError)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user SignUpRequest
	response := peopleColl.FindOne(ctx, bson.M{"email": email})
	if err := response.Decode(&user); err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{"name": user.Name, "surname": user.Surname})
}