package domain

type Response struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

type ResponseSingleData[Data any] struct {
	Code    int    `json:"code"`    // number
	Status  string `json:"status"`  // string
	Data    Data   `json:"data"`    // of data
	Message string `json:"message"` // string
}

type ResponseMultipleData[Data any] struct {
	Code    int    `json:"code"`    // number
	Status  string `json:"status"`  // string
	Data    []Data `json:"data"`    // list of data
	Message string `json:"message"` // string
}

type Empty struct{}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code,omitempty"`
}

// PaginationInfo represents pagination metadata
type PaginationInfo struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

// PaginatedResponse represents a paginated response
type PaginatedResponse struct {
	Data       interface{}    `json:"data"`
	Pagination PaginationInfo `json:"pagination"`
}
