package handlers

import (
	. "authgo/helpers"
	. "authgo/models"
	"net/http"
)

var userStore = make(map[int]User)

var i int = 0

func RegisterHandler(w http.ResponseWriter, r *http.Request) {

	i++

	//token := GenerateToken("**********")
	otp := CreateOtp()
	SendSms("**********", otp)

}

func LoginHandler(w http.ResponseWriter, r *http.Request) {

}

func ForgotPasswordHandler(w http.ResponseWriter, r *http.Request) {

}
