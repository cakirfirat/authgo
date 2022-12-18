package main

import (
	. "authgo/handlers"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {

	log.Println("Server starting...")

	r := mux.NewRouter()

	r.HandleFunc("/api/v1/register", RegisterHandler).Methods("POST")
	r.HandleFunc("/api/v1/login", LoginHandler).Methods("POST")
	r.HandleFunc("/api/v1/forgotpassword", ForgotPasswordHandler).Methods("POST")
	r.HandleFunc("/api/v1/sendotp", SendSmsAgainHandler).Methods("POST")
	r.HandleFunc("/api/v1/verifyotp", CheckOtpHandler).Methods("POST")

	server := &http.Server{
		Addr:    ":8090",
		Handler: r,
	}
	server.ListenAndServe()

}
