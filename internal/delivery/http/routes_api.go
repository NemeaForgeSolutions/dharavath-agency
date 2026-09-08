package http

import (
	"net/http"
)

// registerAPIRoutes sets up JSON API endpoints (leads, inquiries, etc.).
func registerAPIRoutes(mux *http.ServeMux, cfg RouterConfig) {
	if cfg.LeadHandler != nil {
		mux.HandleFunc("POST /api/leads", cfg.LeadHandler.Submit)
	}
}
