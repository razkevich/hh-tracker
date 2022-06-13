package dto

// Error API format
type Error struct {
	Detail string `json:"detail" example:"there was a problem processing your request"`
	Status string `json:"status" example:"404"`
	Title  string `json:"title" example:"not found"`
} // @name Error
