package server

import (
	"context"
	_ "embed"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"

	"github.com/OttoRoming/kolendar/db"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	pool    *pgxpool.Pool
	queries *db.Queries
	slog    *slog.Logger
}

type DeleteResponse struct {
	RowsAffected int64 `json:"rows_affected"`
}

func NewServer(schema string) (*Server, error) {
	ctx := context.Background()

	connStr := "user=user password=password host=127.0.0.1 port=65432 dbname=app sslmode=disable"
	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}

	// create tables
	if _, err := pool.Exec(ctx, schema); err != nil {
		return nil, err
	}

	queries := db.New(pool)

	slog := slog.New(slog.NewTextHandler(os.Stdout, nil))

	return &Server{
		pool:    pool,
		queries: queries,
		slog:    slog,
	}, nil
}

func (s *Server) jsonResponse(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	err := json.NewEncoder(w).Encode(v)
	if err != nil {
		s.slog.Error("failed to write http response", "err", err)
	}
}

func (s *Server) jsonError(w http.ResponseWriter, status int, message string) {
	type ErrorResponse struct {
		Message string `json:"message"`
	}

	s.jsonResponse(w, status, ErrorResponse{
		Message: message,
	})
}

func (s *Server) authenticateRequest(r *http.Request) (*db.User, error) {
	cookie, err := r.Cookie("token")
	if err != nil {
		return nil, err
	}

	session, err := s.queries.GetSessionByToken(r.Context(), cookie.Value)
	if err != nil {
		return nil, err
	}

	user, err := s.queries.GetUserByID(r.Context(), session.UserID)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *Server) pathValueUUID(r *http.Request, key string) (pgtype.UUID, error) {
	value := r.PathValue(key)
	uuid := pgtype.UUID{}
	err := uuid.Scan(value)
	if err != nil {
		return pgtype.UUID{}, err
	}

	return uuid, nil
}

func (s *Server) slogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.slog.Info("incoming request",
			"ip", r.RemoteAddr,
			"method", r.Method,
			"route", r.URL.Path,
		)

		next.ServeHTTP(w, r)
	})
}

func (s *Server) Run() error {
	router := http.NewServeMux()

	router.HandleFunc("POST /users/", s.createUser)
	router.HandleFunc("POST /users/login", s.loginUser)

	router.HandleFunc("GET /calendars/", s.getCalendars)
	router.HandleFunc("POST /calendars/", s.createCalendar)
	router.HandleFunc("DELETE /calendars/{id}/", s.createCalendar)
	router.HandleFunc("UPDATE /calendars/{id}/", s.updateCalendar)
	router.HandleFunc("GET /calendars/{id}/events/", s.getCalendarEvents)

	router.HandleFunc("POST /events/", s.createEvent)
	router.HandleFunc("DELETE /events/{id}/", s.deleteEvent)
	router.HandleFunc("UPDATE /events/{id}/", s.updateEvent)

	err := http.ListenAndServe(":8000", s.slogMiddleware(router))
	if err != nil {
		return err
	}

	return nil
}

func (s *Server) Close() {
	s.pool.Close()
}
