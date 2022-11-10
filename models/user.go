package models

import "time"

type User struct {
	Id          int
	Token       string
	Status      int
	Username    string
	Phone       string
	Email       string
	Password    string
	Otp         string
	CreatedDate time.Time
	ExpireDate  time.Time
}
