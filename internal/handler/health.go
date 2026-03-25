package handler

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/savrin-sharif/go-rest-template/pkg/httputil"
)

// HealthHandler exposes service readiness endpoints.
type HealthHandler struct {
	logger  *slog.Logger
	appName string
}

func NewHealthHandler(appName string, logger *slog.Logger) *HealthHandler {
	return &HealthHandler{appName: appName, logger: logger}
}

// Welcome returns a simple welcome payload.
func (h *HealthHandler) Welcome(w http.ResponseWriter, r *http.Request) {
	payload := map[string]string{
		"message": fmt.Sprintf("Welcome to %s", h.appName),
	}
	httputil.WriteJSON(w, http.StatusOK, payload)
}

// Health reports basic liveness information.
func (h *HealthHandler) Health(w http.ResponseWriter, r *http.Request) {
	payload := map[string]any{
		"status":    "ok",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	}

	httputil.WriteJSON(w, http.StatusOK, payload)
}
