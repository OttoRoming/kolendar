package server

import (
	"encoding/json"
	"net/http"

	"github.com/OttoRoming/kolendar/db"
	"regexp"
)

var (
	colorRegex = regexp.MustCompile("^#[A-Fa-f0-9]{6}$")
)

type CalendarRequest struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

func validateColor(color string) bool {
	return colorRegex.MatchString(color)
}

func (s *Server) getCalendars(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := s.authenticateRequest(r)
	if err != nil {
		s.jsonError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	calendars, err := s.queries.GetCalendarsByUserID(ctx, user.ID)
	if err != nil {
		s.jsonError(w, http.StatusInternalServerError, "Failed to retrieve calendars")
		return
	}

	s.jsonResponse(w, http.StatusOK, calendars)
}

func (s *Server) createCalendar(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := s.authenticateRequest(r)
	if err != nil {
		s.jsonError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	payload := &CalendarRequest{}
	err = json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		s.jsonError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	ok := validateColor(payload.Color)
	if !ok {
		s.jsonError(w, http.StatusBadRequest, "Invalid color format")
		return
	}

	calendar, err := s.queries.CreateCalendar(ctx, db.CreateCalendarParams{
		UserID: user.ID,
		Name:   payload.Name,
		Color:  payload.Color,
	})
	if err != nil {
		s.jsonError(w, http.StatusInternalServerError, "Failed to create calendar")
		return
	}

	s.jsonResponse(w, http.StatusCreated, calendar)
}

func (s *Server) deleteCalendar(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := s.authenticateRequest(r)
	if err != nil {
		s.jsonError(w, http.StatusUnauthorized, "Unauthorized")
	}

	calendarID, err := s.pathValueUUID(r, "id")
	if err != nil {
		s.jsonError(w, http.StatusBadRequest, "Invalid calendar ID")
		return
	}

	rowsAffected, err := s.queries.DeleteCalendarByIDAndUserID(ctx, db.DeleteCalendarByIDAndUserIDParams{
		ID:     calendarID,
		UserID: user.ID,
	})
	if err != nil {
		s.jsonError(w, http.StatusInternalServerError, "Failed to delete calendar")
	}

	s.jsonResponse(w, http.StatusOK, DeleteResponse{
		RowsAffected: rowsAffected,
	})
}

func (s *Server) updateCalendar(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := s.authenticateRequest(r)
	if err != nil {
		s.jsonError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	calendarID, err := s.pathValueUUID(r, "id")
	if err != nil {
		s.jsonError(w, http.StatusBadRequest, "Invalid calendar ID")
		return
	}

	payload := &CalendarRequest{}
	err = json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		s.jsonError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	ok := validateColor(payload.Color)
	if !ok {
		s.jsonError(w, http.StatusBadRequest, "Invalid color format")
		return
	}

	updatedCalendar, err := s.queries.UpdateCalendarByIDAndUserID(ctx, db.UpdateCalendarByIDAndUserIDParams{
		ID:     calendarID,
		UserID: user.ID,
		Name:   payload.Name,
		Color:  payload.Color,
	})

	if err != nil {
		s.jsonError(w, http.StatusInternalServerError, "Failed to update calendar")
		return
	}

	s.jsonResponse(w, http.StatusOK, updatedCalendar)
}

func (s *Server) getCalendarEvents(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := s.authenticateRequest(r)
	if err != nil {
		s.jsonError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	calendarID, err := s.pathValueUUID(r, "id")
	if err != nil {
		s.jsonError(w, http.StatusBadRequest, "Invalid calendar ID")
		return
	}

	events, err := s.queries.GetEventsByCalendarIDAndUserID(ctx, db.GetEventsByCalendarIDAndUserIDParams{
		CalendarID: calendarID,
		UserID:     user.ID,
	})
	if err != nil {
		s.jsonError(w, http.StatusInternalServerError, "Failed to retrieve events")
		return
	}

	s.jsonResponse(w, http.StatusOK, events)
}
