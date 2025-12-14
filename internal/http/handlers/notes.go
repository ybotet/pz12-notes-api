package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/ybotet/pz12-notes-api/internal/core"
	"github.com/ybotet/pz12-notes-api/internal/core/service"
)

type Handler struct {
	NoteService service.NoteService
}

func NewHandler(noteService service.NoteService) *Handler {
	return &Handler{NoteService: noteService}
}

// GetAllNotes godoc
// @Summary Получить все заметки
// @Description Возвращает список всех заметок в системе
// @Tags notes
// @Accept json
// @Produce json
// @Success 200 {array} core.Note
// @Failure 500 {object} core.ErrorResponse
// @Router /api/v1/notes [get]
func (h *Handler) GetAllNotes(w http.ResponseWriter, r *http.Request) {
	notes, err := h.NoteService.GetAllNotes(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notes)
}

// CreateNote godoc
// @Summary Создать новую заметку
// @Description Создает новую заметку с предоставленными данными
// @Tags notes
// @Accept json
// @Produce json
// @Param input body core.NoteCreateRequest true "Данные новой заметки"
// @Success 201 {object} core.Note
// @Failure 400 {object} core.ErrorResponse
// @Failure 500 {object} core.ErrorResponse
// @Router /api/v1/notes [post]
func (h *Handler) CreateNote(w http.ResponseWriter, r *http.Request) {
	var noteReq core.NoteCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&noteReq); err != nil {
		http.Error(w, "Неверный ввод", http.StatusBadRequest)
		return
	}

	// Преобразовать DTO в сущность
	note := core.Note{
		Title:   noteReq.Title,
		Content: noteReq.Content,
	}

	id, err := h.NoteService.CreateNote(r.Context(), note)
	if err != nil {
		if strings.Contains(err.Error(), "не может быть пустым") {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Получить созданную заметку
	createdNote, err := h.NoteService.GetNote(r.Context(), id)
	if err != nil {
		http.Error(w, "Ошибка при получении созданной заметки", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdNote)
}

// GetNote godoc
// @Summary Получить заметку по ID
// @Description Возвращает конкретную заметку по её ID
// @Tags notes
// @Accept json
// @Produce json
// @Param id path int true "ID заметки"
// @Success 200 {object} core.Note
// @Failure 400 {object} core.ErrorResponse
// @Failure 404 {object} core.ErrorResponse
// @Router /api/v1/notes/{id} [get]
func (h *Handler) GetNote(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}

	note, err := h.NoteService.GetNote(r.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "не найдена") {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(note)
}

// UpdateNote godoc
// @Summary Обновить существующую заметку
// @Description Обновляет существующую заметку предоставленными данными (частичное обновление)
// @Tags notes
// @Accept json
// @Produce json
// @Param id path int true "ID заметки"
// @Param input body core.NoteUpdateRequest true "Поля для обновления"
// @Success 200 {object} core.Note
// @Failure 400 {object} core.ErrorResponse
// @Failure 404 {object} core.ErrorResponse
// @Failure 500 {object} core.ErrorResponse
// @Router /api/v1/notes/{id} [put]
func (h *Handler) UpdateNote(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}

	var updates core.NoteUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		http.Error(w, "Неверный ввод", http.StatusBadRequest)
		return
	}

	// Преобразовать DTO в структуру сервиса
	updateReq := service.UpdateNoteRequest{
		Title:   updates.Title,
		Content: updates.Content,
	}

	err = h.NoteService.UpdateNote(r.Context(), id, updateReq)
	if err != nil {
		if strings.Contains(err.Error(), "не найдена") {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else if strings.Contains(err.Error(), "не может быть пустым") {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Получить обновленную заметку
	updatedNote, err := h.NoteService.GetNote(r.Context(), id)
	if err != nil {
		http.Error(w, "Ошибка при получении обновленной заметки", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedNote)
}

// DeleteNote godoc
// @Summary Удалить заметку
// @Description Удаляет конкретную заметку по её ID
// @Tags notes
// @Accept json
// @Produce json
// @Param id path int true "ID заметки"
// @Success 204 "No Content"
// @Failure 400 {object} core.ErrorResponse
// @Failure 404 {object} core.ErrorResponse
// @Failure 500 {object} core.ErrorResponse
// @Router /api/v1/notes/{id} [delete]
func (h *Handler) DeleteNote(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}

	err = h.NoteService.DeleteNote(r.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "не найдена") {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
