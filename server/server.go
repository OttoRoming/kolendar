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

	slog := slog.New(slog.NewTextHandler(os.Stdout, nil))

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

	slog.Info("Database connection established")

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

func (s *Server) authenticateRequest(r *http.Request) (db.User, error) {
	ctx := r.Context()
	var user db.User

	cookie, err := r.Cookie("token")
	if err != nil {
		return user, err
	}
	user, err = s.queries.GetUserBySessionToken(ctx, cookie.Value)
	if err != nil {
		return user, err
	}

	return user, nil
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

	router.HandleFunc("POST   /api/users/", s.createUser)
	router.HandleFunc("POST   /api/users/login", s.loginUser)

	router.HandleFunc("POST   /api/libraries/", s.createLibrary)
	router.HandleFunc("DELETE /api/libraries/{id}/", s.deleteLibrary)
	router.HandleFunc("UPDATE /api/libraries/{id}/", s.updateLibrary)
	router.HandleFunc("GET    /api/libraries/", s.getLibraries)

	router.HandleFunc("POST   /api/libraries/", s.createLibrary)

	address := os.Getenv("ADDRESS")
	if address == "" {
		address = "localhost:8080"
	}

	s.slog.Info("Starting server", "address", address)
	err := http.ListenAndServe(address, s.slogMiddleware(router))
	if err != nil {
		return err
	}

	return nil
}

func (s *Server) Close() {
	s.pool.Close()
}
