package handlers

import (
	. "authgo/helpers"
	. "authgo/models"
	"bytes"
	"encoding/gob"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
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

	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(userStore)
	CheckError(err)

	ioutil.WriteFile("../database.txt", buf.Bytes(), 0644)

	w.WriteHeader(http.StatusCreated)
	//json_data, err := json.Marshal(user)
	CheckError(err)
	w.Write(ConvertJson("Registeration successfull."))

}

func LoginHandler(w http.ResponseWriter, r *http.Request) {

}

func ForgotPasswordHandler(w http.ResponseWriter, r *http.Request) {

}

func SendSmsAgainHandler(w http.ResponseWriter, r *http.Request) {
	if len(userStore) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(ConvertJson("There is no data"))
		return
	}
	id, _ := strconv.Atoi(r.FormValue("id"))
	for _, value := range userStore {
		if id == value.Id {
			phone := value.Phone
			otp := value.Otp
			SendSms(phone, "Tek kullanımlık şifrenizi lütfen kimseyle paylaşmayınız. Şifreniz: "+otp)
			w.WriteHeader(http.StatusOK)
			w.Write(ConvertJson("Sms sending operation is ok."))
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write(ConvertJson("Operation failed"))
	}
}

func CheckOtpHandler(w http.ResponseWriter, r *http.Request) {
	var userUpdate User

	if len(userStore) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(ConvertJson("There is no data"))
		return
	}
	id, err := strconv.Atoi(r.FormValue("id"))
	CheckError(err)
	status := 1
	otp := r.FormValue("otp")
	for key, value := range userStore {
		if id == value.Id {
			if otp == value.Otp {
				userUpdate.Status = status
				userUpdate.Id = value.Id
				userUpdate.Token = value.Token
				userUpdate.Username = value.Username
				userUpdate.Phone = value.Phone
				userUpdate.Email = value.Email
				userUpdate.Password = value.Password
				userUpdate.Otp = value.Otp
				delete(userStore, id)
				userStore[key] = userUpdate
				data, err := json.Marshal(userStore[key])
				CheckError(err)
				w.Write(data)
			}
		}
	}
}
