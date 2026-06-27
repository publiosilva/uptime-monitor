package monitor

import (
	"encoding/json"
	"errors"
	"net/http"

	"uptime-monitor-backend/internal/auth"

	"github.com/go-chi/chi"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req CreateMonitorRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	monitor, err := h.service.Create(userID, req)
	if err != nil {
		switch {
		case errors.As(err, &ErrRequiredFieldMissing{}),
			errors.Is(err, ErrInvalidName),
			errors.Is(err, ErrInvalidURL),
			errors.Is(err, ErrInvalidMethod),
			errors.Is(err, ErrInvalidFrequency):
			writeError(w, http.StatusBadRequest, err.Error())
		default:
			writeError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	writeJSON(w, http.StatusCreated, monitor.ToResponse())
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	monitors, err := h.service.List(userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	responses := make([]MonitorResponse, len(monitors))
	for i, monitor := range monitors {
		responses[i] = monitor.ToResponse()
	}

	writeJSON(w, http.StatusOK, responses)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	monitorID := chi.URLParam(r, "id")
	if monitorID == "" {
		writeError(w, http.StatusBadRequest, "monitor id is required")
		return
	}

	err := h.service.Delete(userID, monitorID)
	if err != nil {
		if errors.Is(err, ErrMonitorNotFound) {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
