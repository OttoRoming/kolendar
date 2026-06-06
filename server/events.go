package server

import (
	"encoding/json"
	"github.com/OttoRoming/kolendar/db"
	"github.com/jackc/pgx/v5/pgtype"
	"net/http"
)

type EventRequest struct {
	CalendarID  pgtype.UUID        `json:"calendar_id"`
	Title       string             `json:"title"`
	Description string             `json:"description"`
	AllDay      bool               `json:"all_day"`
	StartTime   pgtype.Timestamptz `json:"start_time"`
	EndTime     pgtype.Timestamptz `json:"end_time"`
	Location    string             `json:"location"`
}

func (s *Server) createEvent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := s.authenticateRequest(r)
	if err != nil {
		s.jsonError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	payload := &EventRequest{}
	err = json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		s.jsonError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	event, err := s.queries.CreateEventSecure(ctx, db.CreateEventSecureParams{
		UserID:      user.ID,
		CalendarID:  payload.CalendarID,
		Title:       payload.Title,
		Description: payload.Description,
		AllDay:      payload.AllDay,
		StartTime:   payload.StartTime,
		EndTime:     payload.EndTime,
		Location:    payload.Location,
	})
	if err != nil {
		s.jsonError(w, http.StatusInternalServerError, "Failed to create event")
		return
	}

	s.jsonResponse(w, http.StatusCreated, event)
}

func (s *Server) deleteEvent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := s.authenticateRequest(r)
	if err != nil {
		s.jsonError(w, http.StatusUnauthorized, "Unauthorized")
	}

	eventID, err := s.pathValueUUID(r, "id")
	if err != nil {
		s.jsonError(w, http.StatusBadRequest, "Invalid event ID")
		return
	}

	rowsAffected, err := s.queries.DeleteEventByIDAndUserID(ctx, db.DeleteEventByIDAndUserIDParams{
		ID:     eventID,
		UserID: user.ID,
	})
	if err != nil {
		s.jsonError(w, http.StatusInternalServerError, "Failed to delete event")
		return
	}

	s.jsonResponse(w, http.StatusOK, DeleteResponse{
		RowsAffected: rowsAffected,
	})
}

func (s *Server) updateEvent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := s.authenticateRequest(r)
	if err != nil {
		s.jsonError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	eventID, err := s.pathValueUUID(r, "id")
	if err != nil {
		s.jsonError(w, http.StatusBadRequest, "Invalid event ID")
		return
	}

	payload := &EventRequest{}
	err = json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		s.jsonError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	updatedEvent, err := s.queries.UpdateEventByIDAndUserID(ctx, db.UpdateEventByIDAndUserIDParams{
		ID:          eventID,
		UserID:      user.ID,
		CalendarID:  payload.CalendarID,
		Title:       payload.Title,
		Description: payload.Description,
		AllDay:      payload.AllDay,
		StartTime:   payload.StartTime,
		EndTime:     payload.EndTime,
		Location:    payload.Location,
	})

	s.jsonResponse(w, http.StatusOK, updatedEvent)
}
