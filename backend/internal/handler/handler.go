package handler

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sysu-ecnc-dev/sso-system/backend/internal/repository"
)

type Handler struct {
	queries *repository.Queries

	Router *chi.Mux
}

func New(queries *repository.Queries) *Handler {
	return &Handler{
		queries: queries,
		Router:  chi.NewRouter(),
	}
}

func (h *Handler) RegisterRoutes() {
	h.Router.Use(middleware.Logger)

	h.Router.Route("/v1", func(r chi.Router) {
		r.Get("/health", healthCheck)
	})
}
