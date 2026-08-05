// Package apidoc содержит общие для всех сущностей вспомогательные типы,
// которые нужны только Tonic'у для генерации OpenAPI-схемы.
package apidoc

// Empty — пустое тело запроса/ответа. Используется как generic-параметр
// D или R в ftonic.For[D, R], когда эндпоинт не принимает JSON-тело
// (например, читает данные из cookie) или не возвращает JSON (просто 200).
type Empty struct{}

// Pagination — общие query-параметры "page"/"limit" для списковых эндпоинтов.
// Теги нужны ТОЛЬКО Tonic'у для генерации OpenAPI-параметров (query:"...") —
// сам биндинг в хендлере по-прежнему делается вручную через c.Query("page")/
// c.Query("limit"), Tonic ничего сам не парсит и не валидирует за нас.
type Pagination struct {
	Page  int `query:"page" validate:"omitempty,min=1"`
	Limit int `query:"limit" validate:"omitempty,min=1,max=100"`
}
