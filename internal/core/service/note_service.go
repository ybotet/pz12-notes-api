package service

import (
	"context"
	"errors"
	"strings"

	"github.com/ybotet/notes-api/internal/core"
	"github.com/ybotet/notes-api/internal/repo"
)

// NoteService define la interfaz del servicio
type NoteService interface {
	CreateNote(ctx context.Context, note core.Note) (int64, error)
	GetNote(ctx context.Context, id int64) (*core.Note, error)
	GetAllNotes(ctx context.Context) ([]core.Note, error)
	UpdateNote(ctx context.Context, id int64, updates UpdateNoteRequest) error
	DeleteNote(ctx context.Context, id int64) error
}

// Struct para actualizaciones parciales
type UpdateNoteRequest struct {
	Title   *string `json:"title,omitempty"`
	Content *string `json:"content,omitempty"`
}

// Implementación concreta del servicio
type noteServiceImpl struct {
	repo repo.NoteRepository
}

// NewNoteService crea una nueva instancia del servicio
func NewNoteService(repo repo.NoteRepository) NoteService {
	return &noteServiceImpl{repo: repo}
}

func (s *noteServiceImpl) CreateNote(ctx context.Context, note core.Note) (int64, error) {
	// Validaciones de negocio
	if strings.TrimSpace(note.Title) == "" {
		return 0, errors.New("el título no puede estar vacío")
	}

	if len(note.Content) > 1000 {
		return 0, errors.New("el contenido no puede exceder 1000 caracteres")
	}

	// Procesamiento del contenido
	note.Title = strings.TrimSpace(note.Title)
	note.Content = strings.TrimSpace(note.Content)

	// Crear la nota
	return s.repo.Create(ctx, note)
}

func (s *noteServiceImpl) GetNote(ctx context.Context, id int64) (*core.Note, error) {
	if id <= 0 {
		return nil, errors.New("ID inválido")
	}

	return s.repo.GetByID(ctx, id)
}

func (s *noteServiceImpl) GetAllNotes(ctx context.Context) ([]core.Note, error) {
	return s.repo.GetAll(ctx)
}

func (s *noteServiceImpl) UpdateNote(ctx context.Context, id int64, updates UpdateNoteRequest) error {
	if id <= 0 {
		return errors.New("ID inválido")
	}

	// Obtener la nota existente
	existingNote, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Aplicar actualizaciones parciales
	if updates.Title != nil {
		title := strings.TrimSpace(*updates.Title)
		if title == "" {
			return errors.New("el título no puede estar vacío")
		}
		existingNote.Title = title
	}

	if updates.Content != nil {
		content := strings.TrimSpace(*updates.Content)
		if len(content) > 1000 {
			return errors.New("el contenido no puede exceder 1000 caracteres")
		}
		existingNote.Content = content
	}

	// Guardar cambios
	return s.repo.Update(ctx, id, *existingNote)
}

func (s *noteServiceImpl) DeleteNote(ctx context.Context, id int64) error {
	if id <= 0 {
		return errors.New("ID inválido")
	}

	return s.repo.Delete(ctx, id)
}
