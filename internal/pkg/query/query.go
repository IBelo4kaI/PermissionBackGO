package query

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v3"
)

func QueryInt(c fiber.Ctx, key string, fallback int) int {
	raw := c.Query(key)
	if raw == "" {
		return fallback
	}
	v, err := strconv.Atoi(raw)
	if err != nil {
		return fallback
	}
	return v
}

// QuerySearch читает query-параметр "search" (поиск по текстовым полям).
// Пустая строка (не передан или сплошные пробелы) трактуется как "без поиска" —
// nil, а не пустой указатель, чтобы дальше по коду можно было отличить
// "искать пустую строку" от "поиск не задан" и передать nil прямо в
// nullable.StringPtr для sqlc.narg.
func QuerySearch(c fiber.Ctx) *string {
	raw := strings.TrimSpace(c.Query("search"))
	if raw == "" {
		return nil
	}
	return &raw
}

// QueryBoolPtr читает необязательный булев query-параметр ("true"/"1" —
// true, "false"/"0" — false, всё остальное, включая отсутствие параметра —
// nil, то есть "фильтр не задан"). Используется, когда важно отличить
// "показать только false" от "параметр не передан" (например is_global).
func QueryBoolPtr(c fiber.Ctx, key string) *bool {
	switch strings.ToLower(strings.TrimSpace(c.Query(key))) {
	case "true", "1":
		v := true
		return &v
	case "false", "0":
		v := false
		return &v
	default:
		return nil
	}
}

// QuerySort парсит "sort_by"/"sort_dir" из query-параметров. sort_by
// сверяется со списком allowed (белый список колонок конкретного маршрута,
// который умеет сортировать соответствующий SQL-запрос) — если значение не
// входит в список (или не передано), используется defaultSortBy. sort_dir
// принимает только "asc"/"desc", любое другое значение (включая пустое)
// нормализуется в "desc".
func QuerySort(c fiber.Ctx, allowed []string, defaultSortBy string) (sortBy, sortDir string) {
	sortBy = c.Query("sort_by")
	valid := false
	for _, a := range allowed {
		if sortBy == a {
			valid = true
			break
		}
	}
	if !valid {
		sortBy = defaultSortBy
	}

	sortDir = strings.ToLower(c.Query("sort_dir"))
	if sortDir != "asc" && sortDir != "desc" {
		sortDir = "desc"
	}

	return sortBy, sortDir
}
