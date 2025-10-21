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

type userDocumentation struct {
	Name     string `bson:"name"`
	Surname  string `bson:"surname"`
	Email    string `bson:"email"`
	Password string `bson:"password"` // bcrypt hash
}

func ctx5(request *http.Request) (context.Context, context.CancelFunc) {
	return context.WithTimeout(request.Context(), 5*time.Second)
}

func CheckEmail(myMongo *db.Mongo) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		email := request.URL.Query().Get("email")
		if email == "" {
			http.Error(writer, "email required", http.StatusBadRequest)
			return
		}

		ctx, cancel := ctx5(request)
		defer cancel()
		cnt, err := myMongo.People().CountDocuments(ctx, bson.M{"email": email})
		if err != nil {
			http.Error(writer, "db error", http.StatusInternalServerError)
			return
		}

		writeJSON(writer, map[string]bool{"exists": cnt > 0})
	}
}

func SignUp(myMongo *db.Mongo) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		var req SignUpRequest
		if err := json.NewDecoder(request.Body).Decode(&req); err != nil {
			http.Error(writer, "bad json", http.StatusBadRequest)
			return
		}

		if req.Email == "" || req.Password == "" {
			http.Error(writer, "email and password required", http.StatusBadRequest)
			return
		}

		documentation := userDocumentation{
			Name:     req.Name,
			Surname:  req.Surname,
			Email:    req.Email,
			Password: req.Password, // store as plain text
		}

		ctx, cancel := ctx5(request)
		defer cancel()

		res, err := myMongo.People().InsertOne(ctx, documentation)
		if err != nil {
			var we mongo.WriteException
			if errors.As(err, &we) {
				for _, e := range we.WriteErrors {
					if e.Code == 11000 {
						http.Error(writer, "email already exists", http.StatusConflict)
						return
					}
				}
			}
			http.Error(writer, "db insert error", http.StatusInternalServerError)
			return
		}

		writeJSON(writer, map[string]any{
			"status": "ok",
			"id":     res.InsertedID,
		})
	}
}

func SignIn(myMongo *db.Mongo) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		var req SignInRequest
		if err := json.NewDecoder(request.Body).Decode(&req); err != nil {
			http.Error(writer, "bad json", http.StatusBadRequest)
			return
		}

		ctx, cancel := ctx5(request)
		defer cancel()

		var u userDocumentation
		err := myMongo.People().FindOne(ctx, bson.M{"email": req.Email}).Decode(&u)
		if err != nil {
			writeJSON(writer, map[string]any{"passwordControl": false})
			return
		}

		// Plain text password check (you said you don't want hashing)
		ok := (u.Password == req.Password)

		// Return exactly the boolean your frontend reads
		writeJSON(writer, map[string]any{"passwordControl": ok})
	}
}

func GetUserDetails(myMongo *db.Mongo) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		email := request.URL.Query().Get("email")
		if email == "" {
			http.Error(writer, "email required", http.StatusBadRequest)
			return
		}

		ctx, cancel := ctx5(request)
		defer cancel()
		var u userDocumentation
		if err := myMongo.People().FindOne(ctx, bson.M{"email": email}).Decode(&u); err != nil {
			http.Error(writer, "not found", http.StatusNotFound)
			return
		}

		writeJSON(writer, map[string]any{"name": u.Name, "surname": u.Surname})
	}
}

func writeJSON(writer http.ResponseWriter, v any) {
	writer.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(writer).Encode(v)
}
