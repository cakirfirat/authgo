package handlers

import (
	. "authgo/helpers"
	. "authgo/models"
	"bytes"
	"encoding/gob"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

var userStore = make(map[int]User)

var i int = 0

func RegisterHandler(w http.ResponseWriter, r *http.Request) {

	i++
	locale := r.Header.Get("Accept-Language")
	username := r.FormValue("username")
	phone_number := r.FormValue("phone")
	email := r.FormValue("email")
	password := Md5Hash(r.FormValue("password"))

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
	token, _ := CreateJwt()
	otp := CreateOtp()
	created_time := time.Now()

	expire := created_time.AddDate(0, 1, 0)

	SendSms(phone_number, Localizate(locale, "OTP Message")+otp)
	//fmt.Println(Localizate(locale, "OTP Message"))
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
		UpdatedDate: created_time,
	}
	userStore[i] = user

	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(userStore)
	CheckError(err)

	userdata, _ := json.Marshal(userStore[i])

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(userdata)
	CheckError(err)
	w.Write(ConvertJson(Localizate(locale, "Registeration successfull.")))

}

func LoginHandler(w http.ResponseWriter, r *http.Request) {

	var userUpdate User

	locale := r.Header.Get("Accept-Language")

	if len(userStore) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(ConvertJson(Localizate(locale, "There is no data")))
		return
	}

	userName := r.FormValue("username")
	userMail := r.FormValue("email")
	userPassword := Md5Hash(r.FormValue("password"))
	updDate := time.Now()
	expDate := updDate.AddDate(0, 1, 0)
	timeForToken := strconv.FormatInt(time.Now().Unix(), 10)

	for k, v := range userStore {
		if (userName == v.Username) || (userMail == v.Email) {

			if userPassword == v.Password {

				token := GenerateToken(v.Phone + timeForToken)
				userUpdate.Status = v.Status
				userUpdate.Id = v.Id
				userUpdate.Token = token
				userUpdate.Username = v.Username
				userUpdate.Phone = v.Phone
				userUpdate.Email = v.Email
				userUpdate.Password = v.Password
				userUpdate.Otp = v.Otp
				userUpdate.CreatedDate = v.CreatedDate
				userUpdate.UpdatedDate = updDate
				userUpdate.ExpireDate = expDate
				delete(userStore, v.Id)
				userStore[k] = userUpdate

				response_user := map[string]interface{}{
					"token": token,
				}

				data, err := json.Marshal(response_user)

				CheckError(err)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write(data)

			} else {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				w.Write(ConvertJson(Localizate(locale, "Login failed")))
			}

		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(ConvertJson(Localizate(locale, "Login failed")))

		}

	}

}

func ForgotPasswordHandler(w http.ResponseWriter, r *http.Request) {

	userName := r.FormValue("username")
	userMail := r.FormValue("email")
	//locale := r.Header.Get("Accept-Language")

	for _, v := range userStore {
		if (userName == v.Username) || (userMail == v.Email) {

			fEmail := v.Email

			id := strconv.Itoa(v.Id)

			SendEmail("Please click the link below to reset your password. </br> <a href='http://localhost:8090/api/v1/resetpassword/"+id+"'>Click</a>", "Reset password", fEmail)

		}
	}

}

func ResetPassword(w http.ResponseWriter, r *http.Request) {
	locale := r.Header.Get("Accept-Language")
	var userUpdate User

	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	for k, v := range userStore {

		if id == v.Id {

			newPassword := String(8)

			token := v.Token
			userUpdate.Status = v.Status
			userUpdate.Id = v.Id
			userUpdate.Token = token
			userUpdate.Username = v.Username
			userUpdate.Phone = v.Phone
			userUpdate.Email = v.Email
			userUpdate.Password = Md5Hash(newPassword)
			userUpdate.Otp = v.Otp
			userUpdate.CreatedDate = v.CreatedDate
			userUpdate.UpdatedDate = v.UpdatedDate
			userUpdate.ExpireDate = v.ExpireDate
			delete(userStore, v.Id)
			userStore[k] = userUpdate

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(ConvertJson(Localizate(locale, "New password created:") + newPassword))

		}

	}

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

func ValidMethod(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Giriş başarılı"))
}
