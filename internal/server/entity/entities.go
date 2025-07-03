package entity

type User struct {
	ID       int64
	Login    string
	Password string // хэш пароля
}

type Record struct {
	ID     int64  `json:"id"`
	UserID int64  `json:"user_id"`
	Type   string `json:"type"` // password, card, binary, text
	Data   []byte `json:"data"` // base64-encoded blob
	Meta   string `json:"meta"` // произвольный текст
}
