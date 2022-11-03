package helpers

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"mime/multipart"
	"net/http"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

func CheckError(err error) {
	if err != nil {
		log.Println(err.Error())
	}
}

func SendSms(phoneno, message string) {

	url := "https://api.netgsm.com.tr/sms/send/get"
	method := "POST"

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	_ = writer.WriteField("usercode", "*********")
	_ = writer.WriteField("password", "*********")
	_ = writer.WriteField("gsmno", phoneno)
	_ = writer.WriteField("message", message)
	_ = writer.WriteField("msgheader", "*********")
	err := writer.Close()
	CheckError(err)

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	req, err := http.NewRequest(method, url, payload)

	CheckError(err)

	req.Header.Set("Content-Type", writer.FormDataContentType())

	res, err := client.Do(req)

	CheckError(err)

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	CheckError(err)
	fmt.Println(string(body))
}

func GenerateToken(phone string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(phone), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}
	return string(hash)

	hasher := md5.New()
	hasher.Write(hash)
	return hex.EncodeToString(hasher.Sum(nil))
}

func CreateOtp() string {
	randNumber := strconv.Itoa(rand.Intn(999999))
	return randNumber
}
