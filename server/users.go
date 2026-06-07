package server

import (
	"encoding/json"
	"net/http"

	"fmt"
	"github.com/OttoRoming/kolendar/db"
	"github.com/alexedwards/argon2id"
	"github.com/jackc/pgx/v5/pgtype"
)

const (
	MaxPasswordLength = 1024
)

type UserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserResponse struct {
	ID       pgtype.UUID `json:"id"`
	Username string      `json:"username"`
	Token    string      `json:"token"`
}

func (s *Server) createUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	req := &UserRequest{}
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		s.jsonError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if len(req.Password) > MaxPasswordLength {
		s.jsonError(w, http.StatusBadRequest, fmt.Sprintf("password is too long, max length is %d", MaxPasswordLength))
		return
	}

	hashedPassword, err := argon2id.CreateHash(req.Password, argon2id.DefaultParams)
	if err != nil {
		s.jsonError(w, http.StatusInternalServerError, "failed to hash password")
		return
	}

	user, err := s.queries.CreateUser(ctx, db.CreateUserParams{
		Username: req.Username,
		Password: hashedPassword,
	})
	if err != nil {
		s.jsonError(w, http.StatusInternalServerError, "failed to create user")
		return
	}

	session, err := s.queries.CreateSession(ctx, user.ID)

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    session.Token,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Expires:  session.ExpiresAt.Time.UTC(),
	})

	s.jsonResponse(w, http.StatusCreated, &UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Token:    session.Token,
	})
}

func (s *Server) loginUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	req := &UserRequest{}
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		s.jsonError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	user, err := s.queries.GetUserByUsername(ctx, req.Username)
	if err != nil {
		s.jsonError(w, http.StatusUnauthorized, "invalid username or password")
		return
	}

	match, err := argon2id.ComparePasswordAndHash(req.Password, user.Password)
	if err != nil || !match {
		s.jsonError(w, http.StatusUnauthorized, "invalid username or password")
		return
	}

	session, err := s.queries.CreateSession(ctx, user.ID)
	if err != nil {
		s.jsonError(w, http.StatusInternalServerError, "failed to create session")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    session.Token,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Expires:  session.ExpiresAt.Time.UTC(),
	})

	s.jsonResponse(w, http.StatusOK, &UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Token:    session.Token,
	})
}
