package models

type LoginDetails struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

type User struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}
