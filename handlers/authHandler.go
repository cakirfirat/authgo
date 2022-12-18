package handlers

import (
	. "authgo/helpers"
	. "authgo/models"
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

var userStore = make(map[int]User)

var i int = 0

func RegisterHandler(w http.ResponseWriter, r *http.Request) {

	i++
	locale := r.Header.Get("Accept-Language")
	username := r.FormValue("username")
	phone_number := r.FormValue("phone")
	email := r.FormValue("email")
	password := r.FormValue("password")

	if username == "" || phone_number == "" || email == "" || password == "" {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(ConvertJson(Localizate(locale, "Missing data!")))

		return
	}

	for _, v := range userStore {
		if v.Phone == phone_number || v.Id == i || v.Username == username || v.Email == email {
			w.WriteHeader(http.StatusConflict)
			w.Write(ConvertJson(Localizate(locale, "Duplicate registration error!")))
			return
		}
	}
	token := GenerateToken(phone_number)
	otp := CreateOtp()
	created_time := time.Now()

	expire := created_time.AddDate(0, 1, 0)

	//SendSms(phone_number, Localizate(locale, "OTP Message")+otp)
	fmt.Println(Localizate(locale, "OTP Message"))
	var user = User{
		Id:          i,
		Token:       token,
		Status:      0,
		Username:    username,
		Phone:       phone_number,
		Email:       email,
		Password:    password,
		Otp:         otp,
		CreatedDate: created_time,
		ExpireDate:  expire,
	}
	userStore[i] = user

	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(userStore)
	CheckError(err)

	w.WriteHeader(http.StatusCreated)
	CheckError(err)
	w.Write(ConvertJson(Localizate(locale, "Registeration successfull.")))

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
				//data, err := json.Marshal(userStore[key])
				//CheckError(err)
				//w.Write(data)
				response_user := map[string]interface{}{
					"id":     "delicious",
					"token":  value.Token,
					"status": status,
				}
				data, err := json.Marshal(response_user)
				CheckError(err)
				w.Write(data)
			}
		}
	}

}
