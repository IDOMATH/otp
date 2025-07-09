package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"

	"github.com/IDOMATH/CheetahMath/base_change"
)

func main() {
	fmt.Println("Hello world!")
	serverPort := "8080"
	router := http.NewServeMux()
	router.HandleFunc("GET /", handleHome)
	router.HandleFunc("POST /otp", sendOtp)
	router.HandleFunc("POST /otp/{id}", checkOtp)
	server := http.Server{
		Addr:    fmt.Sprint(":", serverPort),
		Handler: router,
	}

	log.Fatal(server.ListenAndServe())
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello world!"))
}

func sendOtp(w http.ResponseWriter, r *http.Request) {
	// find a way to make a unique id.  Something with the user email
	// add otp and id to database
	// send email
	w.Write([]byte(MakeOtp(8)))
}

func MakeOtp(len int) string {
	otp := strings.Builder{}
	r := rand.New(rand.NewSource(99))
	for i := 0; i < len; i++ {
		fmt.Println("Random: ", r.Int63n(36))
		otp.Write([]byte(base_change.FromTen(36, r.Int63n(36))))
	}
	return otp.String()
}

func checkOtp(w http.ResponseWriter, r *http.Request) {
	// compare otp sent to one in database with corresponding id
}
