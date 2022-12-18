package helpers

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"log"
	"math/rand"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/text/language"
)

func CheckError(err error) {
	if err != nil {
		log.Println(err.Error())
	}
}

func SendSms(phoneno, message string) {
	_ = godotenv.Load("../.env")

	url := "https://api.netgsm.com.tr/sms/send/get"
	method := "POST"

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	_ = writer.WriteField("usercode", os.Getenv("USERNAME"))
	_ = writer.WriteField("password", os.Getenv("PASSWORD"))
	_ = writer.WriteField("gsmno", phoneno)
	_ = writer.WriteField("message", message)
	_ = writer.WriteField("msgheader", os.Getenv("USERNAME"))
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

	CheckError(err)
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

func ConvertJson(message string) []byte {
	msg, err := json.Marshal(message)
	CheckError(err)
	return msg
}

func Localizate(lang, text string) string {

	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

	switch lang {
	case "tr-tr":
		bundle.LoadMessageFile("../helpers/lang/tr-TR.json")
	case "en-en":
		bundle.LoadMessageFile("../helpers/lang/en-EN.json")
	}

	localizer := i18n.NewLocalizer(bundle, lang)

	return localizer.MustLocalize(&i18n.LocalizeConfig{DefaultMessage: &i18n.Message{ID: text}})
}
