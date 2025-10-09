package models

type ChatResponse struct {
	Response  string `json:"response"`
	SessionID string `json:"session_id,omitempty"`
	Timestamp int64  `json:"timestamp"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

type AuthResponse struct {
	Token     string `json:"token"`
	ExpiresIn int64  `json:"expires_in"`
}
