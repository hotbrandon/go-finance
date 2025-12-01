package models

// APIResponse provides a consistent structure for all API responses.
// It will contain either Data (for success) or Error (for failure), but not both.
type APIResponse struct {
	Data  interface{} `json:"data,omitempty"`
	Error string      `json:"error,omitempty"`
}

// NewSuccessResponse creates a new APIResponse for successful operations.
func NewSuccessResponse(data interface{}) APIResponse {
	return APIResponse{Data: data}
}

// NewErrorAPIResponse creates a new APIResponse for error conditions.
func NewErrorAPIResponse(err error) APIResponse {
	return APIResponse{Error: err.Error()}
}

// JobsAPIResponse is a specific APIResponse for Swagger documentation
// that explicitly defines the type of the 'data' field for the jobs endpoint.
type JobsAPIResponse struct {
	Data  []JobStatus `json:"data,omitempty"`
	Error string      `json:"error,omitempty"`
}
