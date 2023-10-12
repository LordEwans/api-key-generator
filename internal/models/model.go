package models

type User struct {
	Username string `json:"username" validate:"required"`
	APIKey   string `json:"api_key,omitempty"`
	API      string `json:"api,omitempty" validate:"required"`
}
