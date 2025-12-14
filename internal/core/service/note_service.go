package service

import (
	"context"
	"errors"
	"strings"

	"github.com/ybotet/pz12-notes-api/internal/core"
	"github.com/ybotet/pz12-notes-api/internal/repo"
)

// NoteService определяет интерфейс для бизнес-логики заметок
type NoteService interface {
	CreateNote(ctx context.Context, note core.Note) (int64, error)
	GetNote(ctx context.Context, id int64) (*core.Note, error)
	GetAllNotes(ctx context.Context) ([]core.Note, error)
	UpdateNote(ctx context.Context, id int64, updates UpdateNoteRequest) error
	DeleteNote(ctx context.Context, id int64) error
}

// UpdateNoteRequest представляет запрос на частичное обновление
type UpdateNoteRequest struct {
	Title   *string `json:"title,omitempty"`
	Content *string `json:"content,omitempty"`
}

// noteServiceImpl реализует NoteService
type noteServiceImpl struct {
	repo repo.NoteRepository
}

// NewNoteService создает новый экземпляр сервиса
func NewNoteService(repo repo.NoteRepository) NoteService {
	return &noteServiceImpl{repo: repo}
}

func (s *noteServiceImpl) CreateNote(ctx context.Context, note core.Note) (int64, error) {
	// Валидации бизнес-правил
	if strings.TrimSpace(note.Title) == "" {
		return 0, errors.New("заголовок не может быть пустым")
	}

	if len(note.Content) > 1000 {
		return 0, errors.New("содержание не может превышать 1000 символов")
	}

	// Обработка содержимого
	note.Title = strings.TrimSpace(note.Title)
	note.Content = strings.TrimSpace(note.Content)

	// Создать заметку
	return s.repo.Create(ctx, note)
}

func (s *noteServiceImpl) GetNote(ctx context.Context, id int64) (*core.Note, error) {
	if id <= 0 {
		return nil, errors.New("неверный ID")
	}

	return s.repo.GetByID(ctx, id)
}

func (s *noteServiceImpl) GetAllNotes(ctx context.Context) ([]core.Note, error) {
	return s.repo.GetAll(ctx)
}

func (s *noteServiceImpl) UpdateNote(ctx context.Context, id int64, updates UpdateNoteRequest) error {
	if id <= 0 {
		return errors.New("неверный ID")
	}

	// Получить существующую заметку
	existingNote, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Применить частичные обновления
	if updates.Title != nil {
		title := strings.TrimSpace(*updates.Title)
		if title == "" {
			return errors.New("заголовок не может быть пустым")
		}
		existingNote.Title = title
	}

	if updates.Content != nil {
		content := strings.TrimSpace(*updates.Content)
		if len(content) > 1000 {
			return errors.New("содержание не может превышать 1000 символов")
		}
		existingNote.Content = content
	}

	// Сохранить изменения
	return s.repo.Update(ctx, id, *existingNote)
}

func (s *noteServiceImpl) DeleteNote(ctx context.Context, id int64) error {
	if id <= 0 {
		return errors.New("неверный ID")
	}

	return s.repo.Delete(ctx, id)
}
