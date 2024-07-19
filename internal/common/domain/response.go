package domain

type Page struct {
	Next       int `json:"next"`
	Prev       int `json:"prev"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

type ListResponse[T any] struct {
	Items []*T `json:"items"`
	Count int  `json:"count"`
}
