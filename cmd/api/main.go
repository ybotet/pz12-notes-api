package main

import (
	"log"
	"net/http"

	"github.com/ybotet/notes-api/internal/core/service"
	httpx "github.com/ybotet/notes-api/internal/http"
	"github.com/ybotet/notes-api/internal/http/handlers"
	"github.com/ybotet/notes-api/internal/repo"
)

func main() {
	// Crear repositorio
	repo := repo.NewNoteRepoMem()

	// Crear servicio (con el repositorio como dependencia)
	noteService := service.NewNoteService(repo)

	// Crear handler (con el servicio como dependencia)
	h := &handlers.Handler{NoteService: noteService}

	r := httpx.NewRouter(h)

	log.Println("Server started at :8081")
	log.Fatal(http.ListenAndServe(":8081", r))
}
