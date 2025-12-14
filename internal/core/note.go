package core

import "time"

// Note представляет сущность заметки в системе
// @Description Основная структура заметки
type Note struct {
	ID        int64
	Title     string
	Content   string
	CreatedAt time.Time
	UpdatedAt *time.Time
}
