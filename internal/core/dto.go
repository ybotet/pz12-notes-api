package core

// NoteCreateRequest представляет данные для создания заметки
// @Description Структура для создания новой заметки
type NoteCreateRequest struct {
	Title   string `json:"title" example:"Моя первая заметка"`
	Content string `json:"content" example:"Текст заметки"`
}

// NoteUpdateRequest представляет данные для обновления заметки
// @Description Структура для обновления существующей заметки (частично)
type NoteUpdateRequest struct {
	Title   *string `json:"title,omitempty" example:"Обновленный заголовок"`
	Content *string `json:"content,omitempty" example:"Обновленный текст"`
}

// ErrorResponse представляет стандартный ответ об ошибке
// @Description Общий ответ об ошибке для API
type ErrorResponse struct {
	Error   string `json:"error" example:"сообщение об ошибке"`
	Message string `json:"message,omitempty" example:"Дополнительное описание"`
}
