package model

type Page struct {
	Next       int `json:"next"`
	Prev       int `json:"prev"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}
