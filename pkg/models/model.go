package models

//User struct
type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Level    string `json:"level"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

//Response struct
type Response struct {
	Message string `json:"message"`
}

// RespToken struct
type RespToken struct {
	Message string `json:"message"`
	Token   string `json:"token"`
}
