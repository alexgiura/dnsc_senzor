package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	apperrors "senzor/internal/errors"
	"senzor/internal/models"
	"senzor/internal/services"

	"github.com/gorilla/mux"
)

const maxNetworkAlertBodyBytes = 1 << 20 // 1 MiB

// RegisterNetworkAlertRoutes registers POST /api/v1/network-alerts for agent JSON ingestion.
func RegisterNetworkAlertRoutes(router *mux.Router, svc services.NetworkAlertService) {
	router.HandleFunc("/api/v1/network-alerts", postNetworkAlertHandler(svc)).Methods(http.MethodPost)
}

func postNetworkAlertHandler(svc services.NetworkAlertService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if svc == nil {
			http.Error(w, "service unavailable", http.StatusServiceUnavailable)
			return
		}

		r.Body = http.MaxBytesReader(w, r.Body, maxNetworkAlertBodyBytes)
		defer r.Body.Close()

		var payload models.NetworkAlert
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			writeJSONError(w, http.StatusBadRequest, "invalid_json", "request body must be valid JSON")
			return
		}

		if err := svc.Ingest(r.Context(), &payload); err != nil {
			if errors.Is(err, apperrors.ErrValidation) {
				writeJSONError(w, http.StatusBadRequest, "validation_failed", err.Error())
				return
			}
			writeJSONError(w, http.StatusInternalServerError, "internal_error", "failed to store alert")
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"status":"created"}`))
	}
}

func writeJSONError(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	resp := models.ErrorResponse{Code: code, Message: message}
	_ = json.NewEncoder(w).Encode(resp)
}
