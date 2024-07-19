package domain

type ListRequest struct {
	Pagination Pagination
	Filters    []Filter
	Sorts      []Sort
}

type Pagination struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

type SortOrder string

const (
	Asc  SortOrder = "asc"
	Desc SortOrder = "desc"
)

type Sort struct {
	Field         string    `json:"field"`
	Order         SortOrder `json:"order"`
	CaseSensitive bool      `json:"case_sensitive"`
}

type Op int

const (
	_ Op = iota
	Eq
	Ne
	Gt
	Lt
	Ge
	Le
	In
	NotIn
	Contains
	NotContains
	StartsWith
	EndsWith
	Null
	NotNull
	Between
)

type FilterValueType int

const (
	_ FilterValueType = iota
	String
	Integer
	Float
	Boolean
	Date
	Time
	DateTime
)

type Filter struct {
	Field         string `json:"field"`
	Value         any    `json:"value"`
	Op            Op     `json:"op"`
	CaseSensitive bool   `json:"is_case_sensitive"`
}
