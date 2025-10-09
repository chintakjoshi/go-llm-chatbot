package models

type ChatRequest struct {
	Message   string `json:"message" binding:"required"`
	SessionID string `json:"session_id,omitempty"`
}

type AuthRequest struct {
	APIKey string `json:"api_key" binding:"required"`
}
