package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/OttoRoming/kolendar/repository"
	"github.com/alexedwards/argon2id"
	"github.com/google/uuid"
)

func createUser(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	var user struct {
		Username string
		Email    string
		Password string
	}

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if len(user.Password) < 8 {
		http.Error(w, "password must be at least 8 characters long", http.StatusBadRequest)
		return
	}

	hash, err := argon2id.CreateHash(user.Password, argon2id.DefaultParams)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	newUser, err := queries.CreateUser(ctx, repository.CreateUserParams{
		ID:       uuid.NewString(),
		Username: user.Username,
		Password: hash,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(newUser); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func loginUser(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	var creds struct {
		Username string
		Password string
	}

	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := queries.GetUserByUsername(ctx, creds.Username)
	if err != nil {
		http.Error(w, "invalid username or password", http.StatusUnauthorized)
		return
	}

	isPasswordCorrect, err := argon2id.ComparePasswordAndHash(creds.Password, user.Password)
	if err != nil || !isPasswordCorrect {
		http.Error(w, "invalid username or password", http.StatusUnauthorized)
		return
	}

	token, err := generateToken(user.ID)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := struct {
		Token string
	}{Token: token}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
