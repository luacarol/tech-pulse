package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/luacarol/tech-pulse/internal/model"
	"github.com/luacarol/tech-pulse/internal/service"
)

// Handler wires HTTP requests to the service layer.
type Handler struct {
	svc         *service.NewsService
	logger      *slog.Logger
	corsOrigins []string
}

func NewHandler(svc *service.NewsService, logger *slog.Logger, corsOrigins []string) *Handler {
	return &Handler{svc: svc, logger: logger, corsOrigins: corsOrigins}
}

func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", h.health)
	mux.HandleFunc("GET /api/v1/news", h.listNews)
	return h.withCORS(mux)
}

func (h *Handler) health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handler) listNews(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	filter := model.NewsFilter{
		Topics:  service.ParseTopics(q.Get("topics")),
		Sources: service.ParseSources(q.Get("sources")),
		Query:   q.Get("q"),
		Sort:    q.Get("sort"),
		Order:   q.Get("order"),
		Limit:   parseInt(q.Get("limit"), 20),
		Offset:  parseInt(q.Get("offset"), 0),
	}

	page, err := h.svc.List(r.Context(), filter)
	if err != nil {
		h.logger.Error("list news", "error", err)
		writeErr(w, http.StatusInternalServerError, "failed to list news")
		return
	}

	writeJSON(w, http.StatusOK, page)
}

func parseInt(s string, fallback int) int {
	if s == "" {
		return fallback
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return fallback
	}
	return n
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeErr(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

// withCORS only emits CORS headers when the request Origin is in the allowlist.
// Requests without an Origin (curl, server-to-server) and disallowed origins are
// served without CORS headers, so browsers block the response but non-browser
// clients are unaffected.
func (h *Handler) withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if origin := r.Header.Get("Origin"); originAllowed(origin, h.corsOrigins) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Vary", "Origin")
			w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		}
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func originAllowed(origin string, allowed []string) bool {
	if origin == "" {
		return false
	}
	for _, a := range allowed {
		if a == origin {
			return true
		}
	}
	return false
}
