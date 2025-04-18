package model

type Role string

const (
	Owner Role = "owner"
	Admin Role = "admin"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Role     Role   `json:"role"`
	Password string `json:"password"`
}
