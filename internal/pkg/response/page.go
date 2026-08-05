package response

type Page[T any] struct {
	Items []T   `json:"items"`
	Total int64 `json:"total"`
	Page  int   `json:"page"`
	Limit int   `json:"limit"`
	Pages int   `json:"pages"`
}

func ToPageResponse[T any](items []T, total int64, page, limit, pages int) *Page[T] {
	return &Page[T]{
		Items: items,
		Total: total,
		Page:  page,
		Limit: limit,
		Pages: pages,
	}
}
