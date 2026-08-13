// Package apidoc содержит общие для всех сущностей вспомогательные типы,
// которые нужны только Tonic'у для генерации OpenAPI-схемы.
package apidoc

// Empty — пустое тело запроса/ответа. Используется как generic-параметр
// D или R в ftonic.For[D, R], когда эндпоинт не принимает JSON-тело
// (например, читает данные из cookie) или не возвращает JSON (просто 200).
type Empty struct{}

// Pagination — общие query-параметры списковых эндпоинтов: "page"/"limit"
// для пагинации, "search" для поиска по текстовым полям сущности и
// "sort_by"/"sort_dir" для сортировки. Набор колонок, допустимых в sort_by,
// у каждого маршрута свой (белый список в самом хендлере) — здесь только
// документация для OpenAPI.
// Теги нужны ТОЛЬКО Tonic'у для генерации OpenAPI-параметров (query:"...") —
// сам биндинг в хендлере по-прежнему делается вручную через c.Query(...)/
// query.QuerySearch/query.QuerySort, Tonic ничего сам не парсит и не валидирует за нас.
type Pagination struct {
	Page    int    `query:"page" validate:"omitempty,min=1"`
	Limit   int    `query:"limit" validate:"omitempty,min=1,max=100"`
	Search  string `query:"search" validate:"omitempty,max=255" description:"Поиск по текстовым полям"`
	SortBy  string `query:"sort_by" description:"Колонка для сортировки (список зависит от маршрута)"`
	SortDir string `query:"sort_dir" validate:"omitempty,oneof=asc desc" description:"Направление сортировки: asc или desc (по умолчанию desc)"`
}
