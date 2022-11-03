package handlers

import (
	. "authgo/helpers"
	. "authgo/models"
	"encoding/json"
	"fmt"
	"net/http"
)

var userStore = make(map[int]User)

var i int = 0

func RegisterHandler(w http.ResponseWriter, r *http.Request) {

	i++
	phone_number := r.FormValue("phone")

	for _, v := range userStore {
		if v.Phone == phone_number {
			fmt.Println("Mükerrer kayıt.")
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
		Username: r.FormValue("username"),
		Phone:    phone_number,
		Password: r.FormValue("password"),
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
