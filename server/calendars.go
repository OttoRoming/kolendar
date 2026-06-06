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

	calendar, err := s.queries.GetCalendarByID(ctx, calendarID)
	if err != nil {
		s.jsonError(w, http.StatusNotFound, "Calendar not found")
		return
	}

	if calendar.UserID != user.ID {
		s.jsonError(w, http.StatusForbidden, "Forbidden")
		return
	}

	err = s.queries.DeleteCalendarByID(ctx, calendarID)
	if err != nil {
		s.jsonError(w, http.StatusInternalServerError, "Failed to delete calendar")
		return
	}

	s.jsonResponse(w, http.StatusNoContent, nil)
}

func (s *Server) updateCalendar(w http.ResponseWriter, r *http.Request) {
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

	calendar, err := s.queries.GetCalendarByID(ctx, calendarID)
	if err != nil {
		s.jsonError(w, http.StatusNotFound, "Calendar not found")
		return
	}

	if calendar.UserID != user.ID {
		s.jsonError(w, http.StatusForbidden, "Forbidden")
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

	updatedCalendar, err := s.queries.UpdateCalendarByID(ctx, db.UpdateCalendarByIDParams{
		ID:    calendarID,
		Name:  payload.Name,
		Color: payload.Color,
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
	}

	calendarID, err := s.pathValueUUID(r, "id")
	if err != nil {
		s.jsonError(w, http.StatusBadRequest, "Invalid calendar ID")
		return
	}

	calendar, err := s.queries.GetCalendarByID(ctx, calendarID)
	if err != nil {
		s.jsonError(w, http.StatusNotFound, "Calendar not found")
		return
	}

	if calendar.UserID != user.ID {
		s.jsonError(w, http.StatusForbidden, "Forbidden")
		return
	}

	events, err := s.queries.GetEventsByCalendarID(ctx, calendarID)
	if err != nil {
		s.jsonError(w, http.StatusInternalServerError, "Failed to retrieve events")
		return
	}

	s.jsonResponse(w, http.StatusOK, events)
}
