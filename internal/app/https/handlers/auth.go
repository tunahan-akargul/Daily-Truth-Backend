package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"main-services/internal/db"
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

type userDoc struct {
	Name     string `bson:"name"`
	Surname  string `bson:"surname"`
	Email    string `bson:"email"`
	Password string `bson:"password"` // bcrypt hash
}

func ctx5(r *http.Request) (context.Context, context.CancelFunc) {
	return context.WithTimeout(r.Context(), 5*time.Second)
}

func CheckEmail(m *db.Mongo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email := r.URL.Query().Get("email")
		if email == "" {
			http.Error(w, "email required", http.StatusBadRequest)
			return
		}

		ctx, cancel := ctx5(r)
		defer cancel()
		cnt, err := m.People().CountDocuments(ctx, bson.M{"email": email})
		if err != nil {
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}

		writeJSON(w, map[string]bool{"exists": cnt > 0})
	}
}

func SignUp(m *db.Mongo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req SignUpRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "bad json", http.StatusBadRequest)
			return
		}

		if req.Email == "" || req.Password == "" {
			http.Error(w, "email and password required", http.StatusBadRequest)
			return
		}

		doc := userDoc{
			Name:     req.Name,
			Surname:  req.Surname,
			Email:    req.Email,
			Password: req.Password, // store as plain text
		}

		ctx, cancel := ctx5(r)
		defer cancel()

		res, err := m.People().InsertOne(ctx, doc)
		if err != nil {
			var we mongo.WriteException
			if errors.As(err, &we) {
				for _, e := range we.WriteErrors {
					if e.Code == 11000 {
						http.Error(w, "email already exists", http.StatusConflict)
						return
					}
				}
			}
			http.Error(w, "db insert error", http.StatusInternalServerError)
			return
		}

		writeJSON(w, map[string]any{
			"status": "ok",
			"id":     res.InsertedID,
		})
	}
}

func SignIn(m *db.Mongo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req SignInRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "bad json", http.StatusBadRequest)
			return
		}

		ctx, cancel := ctx5(r)
		defer cancel()

		var u userDoc
		err := m.People().FindOne(ctx, bson.M{"email": req.Email}).Decode(&u)
		if err != nil {
			writeJSON(w, map[string]any{"passwordControl": false})
			return
		}

		// Plain text password check (you said you don't want hashing)
		ok := (u.Password == req.Password)

		// Return exactly the boolean your frontend reads
		writeJSON(w, map[string]any{"passwordControl": ok})
	}
}

func GetUserDetails(m *db.Mongo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email := r.URL.Query().Get("email")
		if email == "" {
			http.Error(w, "email required", http.StatusBadRequest)
			return
		}

		ctx, cancel := ctx5(r)
		defer cancel()
		var u userDoc
		if err := m.People().FindOne(ctx, bson.M{"email": email}).Decode(&u); err != nil {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}

		writeJSON(w, map[string]any{"name": u.Name, "surname": u.Surname})
	}
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}
