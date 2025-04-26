package dto

// @Description Bad request
type BadRequest struct {
	Error string `json:"error" example:"invalid request data"`
}

// @Description Internal server error
type InternalServerError struct {
	Error string `json:"error" example:"internal error"`
}

// @Description Service unavailable
type ServiceUnavailable struct {
	Error string `json:"error" example:"service unavailable"`
}

// @Description Resource not found
type NotFound struct {
	Error string `json:"error" example:"person not found"`
}
