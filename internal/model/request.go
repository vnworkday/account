package model

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
	Field           string    `json:"field"`
	Order           SortOrder `json:"order"`
	IsCaseSensitive bool      `json:"is_case_sensitive"`
}

type FilterOperator int

const (
	_ FilterOperator = iota
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
	Field           string         `json:"field"`
	Value           any            `json:"value"`
	Operator        FilterOperator `json:"operator"`
	IsCaseSensitive bool           `json:"is_case_sensitive"`
}
