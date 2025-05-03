package dto

// APIResponse represents the standard API response structure
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// PaginationRequest represents pagination parameters for list requests
type PaginationRequest struct {
	Page     int `query:"page" validate:"omitempty,min=1"`
	PageSize int `query:"pageSize" validate:"omitempty,min=1,max=100"`
}

// FilterRequest represents common filter parameters
type FilterRequest struct {
	StartDate string `query:"startDate" validate:"omitempty,datetime=2006-01-02"`
	EndDate   string `query:"endDate" validate:"omitempty,datetime=2006-01-02"`
	Status    string `query:"status" validate:"omitempty"`
	Type      string `query:"type" validate:"omitempty"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error"`
}

// SuccessResponse represents a simple success response
type SuccessResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}
