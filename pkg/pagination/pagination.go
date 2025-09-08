package pagination

type Meta struct {
	Page  int   `json:"page"`
	Limit int   `json:"limit"`
	Total int64 `json:"total"`
}

type Response[T any] struct {
	Data []T  `json:"data"`
	Meta Meta `json:"meta"`
}
