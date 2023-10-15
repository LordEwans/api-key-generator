package models

type User struct {
	Username string `json:"username" validate:"required" bson:"username"`
	APIKey   string `json:"api_key,omitempty" bson:"api_key"`
	API      string `json:"api,omitempty" validate:"required" bson:"api_key"`
}
