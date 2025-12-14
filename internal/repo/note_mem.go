package repo

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/ybotet/pz12-notes-api/internal/core"
)

// NoteRepository определяет интерфейс для доступа к данным заметок
type NoteRepository interface {
	Create(ctx context.Context, note core.Note) (int64, error)
	GetByID(ctx context.Context, id int64) (*core.Note, error)
	GetAll(ctx context.Context) ([]core.Note, error)
	Update(ctx context.Context, id int64, note core.Note) error
	Delete(ctx context.Context, id int64) error
}

// NoteRepoMem реализует NoteRepository
type NoteRepoMem struct {
	mu    sync.RWMutex
	notes map[int64]*core.Note
	next  int64
}

func NewNoteRepoMem() *NoteRepoMem {
	return &NoteRepoMem{
		notes: make(map[int64]*core.Note),
		next:  1,
	}
}

func (r *NoteRepoMem) Create(ctx context.Context, n core.Note) (int64, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	n.ID = r.next
	n.CreatedAt = time.Now()
	r.notes[n.ID] = &n
	r.next++

	return n.ID, nil
}

func (r *NoteRepoMem) GetByID(ctx context.Context, id int64) (*core.Note, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	note, exists := r.notes[id]
	if !exists {
		return nil, errors.New("заметка не найдена")
	}

	// Вернуть копию
	noteCopy := *note
	return &noteCopy, nil
}

func (r *NoteRepoMem) GetAll(ctx context.Context) ([]core.Note, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	notes := make([]core.Note, 0, len(r.notes))
	for _, note := range r.notes {
		notes = append(notes, *note)
	}

	return notes, nil
}

func (r *NoteRepoMem) Update(ctx context.Context, id int64, updatedNote core.Note) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, exists := r.notes[id]
	if !exists {
		return errors.New("заметка не найдена")
	}

	updatedNote.ID = id
	now := time.Now()
	updatedNote.UpdatedAt = &now
	r.notes[id] = &updatedNote

	return nil
}

func (r *NoteRepoMem) Delete(ctx context.Context, id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, exists := r.notes[id]
	if !exists {
		return errors.New("заметка не найдена")
	}

	delete(r.notes, id)
	return nil
}
