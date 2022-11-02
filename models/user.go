package models

type User struct {
	Id       int
	Token    string
	Status   int
	Username string
	Phone    string
	Password string
	Otp      string
}
