package helpers

import (
	"bytes"
	"crypto/md5"
	"crypto/tls"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/text/language"
	gomail "gopkg.in/mail.v2"
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
	default:
		bundle.LoadMessageFile("../helpers/lang/en-EN.json")
	}

	localizer := i18n.NewLocalizer(bundle, lang)

	return localizer.MustLocalize(&i18n.LocalizeConfig{DefaultMessage: &i18n.Message{ID: text}})
}

func Md5Hash(str string) string {
	hasher := md5.New()
	hasher.Write([]byte(str))
	return hex.EncodeToString(hasher.Sum(nil))
}

func SendEmail(msg, sbj, to string) string {
	_ = godotenv.Load("../.env")

	from := os.Getenv("EMAIL_ADRESS")
	password := os.Getenv("EMAIL_PASSWORD")

	m := gomail.NewMessage()

	// Set E-Mail sender
	m.SetHeader("From", from)

	// Set E-Mail receivers
	m.SetHeader("To", to)

	// Set E-Mail subject
	m.SetHeader("Subject", sbj)

	// Set E-Mail body. You can set plain text or html with text/html
	m.SetBody("text/html", msg)

	// t := template.Must(template.New("../templates/reset-password.html").Parse("Hello {{.}}!"))
	// m.SetBodyWriter("text/plain", func(w io.Writer) error {
	// 	return t.Execute(w, "Bob")
	// })

	// Settings for SMTP server
	d := gomail.NewDialer("smtp.gmail.com", 587, from, password)

	// This is only needed when SSL/TLS certificate is not valid on server.
	// In production this should be set to false.
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Now send E-Mail
	if err := d.DialAndSend(m); err != nil {
		fmt.Println(err)
		panic(err)
	}

	return "true"
}

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func StringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func String(length int) string {
	return StringWithCharset(length, charset)
}

var SECRET = []byte("super-secret-auth")

func CreateJwt() (string, error) {

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["exp"] = time.Now().Add(time.Hour * 730).Unix()

	tokenStr, err := token.SignedString(SECRET)

	CheckError(err)

	return tokenStr, nil

}

func ValidateJwt(next func(w http.ResponseWriter, r *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header["Token"] != nil {
			token, err := jwt.Parse(r.Header["Token"][0], func(t *jwt.Token) (interface{}, error) {
				_, ok := t.Method.(*jwt.SigningMethodHMAC)
				if !ok {
					w.WriteHeader(http.StatusUnauthorized)
					w.Write([]byte("Not authorized"))
				}
				return SECRET, nil
			})
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Not authorized"))
			}

			if token.Valid {
				next(w, r)
			}
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Not authorized"))
		}
	})
}
