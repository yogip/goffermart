package model

type UserRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type User struct {
	ID           int64   `json:"id"`
	Login        string  `json:"login"`
	PasswordHash *[]byte `json:"password,omitempty"`
}

type Order struct {
	ID        int64  `json:"number"`
	Status    string `json:"status"`
	Accrual   int32  `json:"accrual"`
	CreatedAt int32  `json:"uploaded_at"`
}
