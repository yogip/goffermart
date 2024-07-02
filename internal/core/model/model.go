package model

type UserRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type User struct {
	ID           int64   `json:"id"`
	Login        string  `json:"login"`
	PasswordHash *[]byte `json:"password"`
}

// type User struct {
// 	ID           int64   `json:"id"`
// 	Login        string  `json:"login"`
// 	PasswordHash *string `json:"password,omitempty"`
// }
