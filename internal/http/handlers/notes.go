package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/ybotet/notes-api/internal/core"
	"github.com/ybotet/notes-api/internal/core/service"
)

type Handler struct {
	NoteService service.NoteService
}

func NewHandler(noteService service.NoteService) *Handler {
	return &Handler{NoteService: noteService}
}

func (h *Handler) CreateNote(w http.ResponseWriter, r *http.Request) {
	var note core.Note
	if err := json.NewDecoder(r.Body).Decode(&note); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	id, err := h.NoteService.CreateNote(r.Context(), note)
	if err != nil {
		if strings.Contains(err.Error(), "no puede estar vacío") {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Obtener la nota creada
	createdNote, err := h.NoteService.GetNote(r.Context(), id)
	if err != nil {
		http.Error(w, "Error retrieving created note", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdNote)
}

func (h *Handler) GetNote(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	note, err := h.NoteService.GetNote(r.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "no encontrada") {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(note)
}

func (h *Handler) GetAllNotes(w http.ResponseWriter, r *http.Request) {
	notes, err := h.NoteService.GetAllNotes(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notes)
}

func (h *Handler) UpdateNote(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var updates service.UpdateNoteRequest
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	err = h.NoteService.UpdateNote(r.Context(), id, updates)
	if err != nil {
		if strings.Contains(err.Error(), "no encontrada") {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else if strings.Contains(err.Error(), "no puede estar vacío") {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) DeleteNote(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	err = h.NoteService.DeleteNote(r.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "no encontrada") {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
