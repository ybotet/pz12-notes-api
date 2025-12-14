package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/ybotet/pz12-notes-api/internal/http/handlers"
)

func NewRouter(h *handlers.Handler) *chi.Mux {
	r := chi.NewRouter()

	// Middlewares
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Rutas de la API
	r.Route("/api/v1/notes", func(r chi.Router) {
		r.Get("/", h.GetAllNotes)
		r.Post("/", h.CreateNote)
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", h.GetNote)
			r.Put("/", h.UpdateNote)
			r.Delete("/", h.DeleteNote)
		})
	})

	// Ruta de salud
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	return r
}
