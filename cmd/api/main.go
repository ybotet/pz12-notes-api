// Package main Notes API server.
//
// @title Notes API
// @version 1.0
// @description REST API для управления заметками (CRUD).
// @contact.name Институт перспективных технологий
// @contact.email example@university.ru
//
// @host localhost:8081
// @BasePath /api/v1
//
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Введите токен в формате: Bearer <token>
package main

import (
	"log"
	"net/http"
	"os"

	// _ "pz12-notes-api/docs"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/ybotet/pz12-notes-api/internal/core/service"
	"github.com/ybotet/pz12-notes-api/internal/http/handlers"
	"github.com/ybotet/pz12-notes-api/internal/repo"
)

// Intenta importar docs solo si existen
// Esto previene errores de compilación cuando docs/ no está generado
// var _ = func() error {
//     _, err := os.Stat("docs/docs.go")
//     if err == nil {
//         // Solo importa si el archivo existe
//         _ = "github.com/ybotet/pz12-notes-api/docs"
//     }
//     return nil
// }()

func main() {
	// Crear repositorio
	repo := repo.NewNoteRepoMem()

	// Crear servicio
	noteService := service.NewNoteService(repo)

	// Crear handlers
	handler := handlers.NewHandler(noteService)

	// Crear router
	r := chi.NewRouter()

	// Middlewares
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)

	// Rutas de la API
	r.Route("/api/v1/notes", func(r chi.Router) {
		r.Get("/", handler.GetAllNotes)
		r.Post("/", handler.CreateNote)
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", handler.GetNote)
			r.Put("/", handler.UpdateNote)
			r.Delete("/", handler.DeleteNote)
		})
	})

	// Ruta de salud
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "ok", "service": "notes-api"}`))
	})

	// Ruta de Swagger UI (condicional)
	if _, err := os.Stat("docs/swagger.json"); err == nil {
		// Configurar Swagger UI manualmente (sin importar http-swagger)
		r.Get("/docs", func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "/docs/index.html", http.StatusMovedPermanently)
		})

		r.Get("/docs/*", func(w http.ResponseWriter, r *http.Request) {
			http.StripPrefix("/docs", http.FileServer(http.Dir("docs"))).ServeHTTP(w, r)
		})

		log.Println("📚 Swagger UI disponible en http://localhost:8081/docs/index.html")
	} else {
		log.Println("⚠️  Swagger UI no disponible: ejecuta 'swag init' primero")
	}

	// Ruta de ReDoc (alternativa)
	r.Get("/redoc", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`
            <!DOCTYPE html>
            <html>
            <head>
                <title>Notes API - ReDoc</title>
                <meta charset="utf-8"/>
                <meta name="viewport" content="width=device-width, initial-scale=1">
            </head>
            <body>
                <redoc spec-url='/docs/swagger.json'></redoc>
                <script src="https://cdn.jsdelivr.net/npm/redoc@next/bundles/redoc.standalone.js"></script>
            </body>
            </html>
        `))
	})

	// Iniciar servidor
	log.Println("🚀 Сервер запущен на http://localhost:8081")
	log.Println("📝 API доступен по адресу http://localhost:8081/api/v1/notes")
	log.Fatal(http.ListenAndServe(":8081", r))
}
