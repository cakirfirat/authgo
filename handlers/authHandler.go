package handlers

import (
	. "authgo/helpers"
	. "authgo/models"
	"encoding/json"
	"net/http"
)

var userStore = make(map[int]User)

var i int = 0

func RegisterHandler(w http.ResponseWriter, r *http.Request) {

	i++
	username := r.FormValue("username")
	phone_number := r.FormValue("phone")
	email := r.FormValue("email")
	password := r.FormValue("password")
	password_again := r.FormValue("password_again")

	if password != password_again {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(ConvertJson("Passwords not match!"))

		return
	}

	if username == "" || phone_number == "" || email == "" || password == "" || password_again == "" {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(ConvertJson("Missing data!"))
		return
	}

	for _, v := range userStore {
		if v.Phone == phone_number || v.Id == i || v.Username == username || v.Email == email {
			w.WriteHeader(http.StatusConflict)
			w.Write(ConvertJson("Duplicate registration error!"))
			return
		}
	}
	token := GenerateToken(phone_number)
	otp := CreateOtp()
	SendSms(phone_number, "Tek kullanımlık şifrenizi lütfen kimseyle paylaşmayınız. Şifreniz: "+otp)
	var user = User{
		Id:       i,
		Token:    token,
		Status:   0,
		Username: username,
		Phone:    phone_number,
		Email:    email,
		Password: password,
		Otp:      otp,
	}
	userStore[i] = user
	w.WriteHeader(http.StatusCreated)
	json_data, err := json.Marshal(user)
	CheckError(err)
	w.Write(json_data)

}

func LoginHandler(w http.ResponseWriter, r *http.Request) {

}

func ForgotPasswordHandler(w http.ResponseWriter, r *http.Request) {

}

func SendSmsAgainHandler(w http.ResponseWriter, r *http.Request) {

}
