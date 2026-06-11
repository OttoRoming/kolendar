package server

import (
	"encoding/json"
	"net/http"

	"github.com/OttoRoming/kolendar/db"
)

type LibraryRequest struct {
	Name string `json:"name"`
}

func (s *Server) createLibrary(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := s.authenticateRequest(r)
	if err != nil {
		s.jsonError(w, http.StatusUnauthorized, "Unauthorized")
	}

	var req LibraryRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		s.jsonError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	library, err := s.queries.CreateLibrary(ctx, db.CreateLibraryParams{
		OwnerID: user.ID,
		Name:    req.Name,
	})

	s.jsonResponse(w, http.StatusCreated, library)
}

func (s *Server) deleteLibrary(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := s.authenticateRequest(r)
	if err != nil {
		s.jsonError(w, http.StatusUnauthorized, "Unauthorized")
	}

	id, err := s.pathValueUUID(r, "id")
	if err != nil {
		s.jsonError(w, http.StatusBadRequest, "Invalid library ID")
	}

	rowsAffected, err := s.queries.DeleteLibrarySecure(ctx, db.DeleteLibrarySecureParams{
		ID:      id,
		OwnerID: user.ID,
	})
	if err != nil {
		s.jsonError(w, http.StatusInternalServerError, "Failed to delete library")
		return
	}
	if rowsAffected == 0 {
		s.jsonError(w, http.StatusNotFound, "Library not found")
		return
	}

	s.jsonResponse(w, http.StatusOK, DeleteResponse{
		RowsAffected: rowsAffected,
	})
}

func (s *Server) updateLibrary(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := s.authenticateRequest(r)
	if err != nil {
		s.jsonError(w, http.StatusUnauthorized, "Unauthorized")
	}

	id, err := s.pathValueUUID(r, "id")
	if err != nil {
		s.jsonError(w, http.StatusBadRequest, "Invalid library ID")
	}

	var req LibraryRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		s.jsonError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	library, err := s.queries.UpdateLibrarySecure(ctx, db.UpdateLibrarySecureParams{
		ID:      id,
		OwnerID: user.ID,
		Name:    req.Name,
	})
	if err != nil {
		s.jsonError(w, http.StatusInternalServerError, "Failed to update library")
		return
	}

	s.jsonResponse(w, http.StatusOK, library)
}

func (s *Server) getLibraries(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := s.authenticateRequest(r)
	if err != nil {
		s.jsonError(w, http.StatusUnauthorized, "Unauthorized")
	}

	libraries, err := s.queries.GetLibrariesByOwnerID(ctx, user.ID)
	if err != nil {
		s.jsonError(w, http.StatusInternalServerError, "Failed to retrieve libraries")
		return
	}

	s.jsonResponse(w, http.StatusOK, libraries)
}
