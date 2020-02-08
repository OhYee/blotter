package api

// SimpleResponse a standard response type
type SimpleResponse struct {
	Success bool   `json:"success"`
	Title   string `json:"title"`
	Content string `json:"content"`
}
