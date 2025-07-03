package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	fmt.Println("Hello world!")
	serverPort := "8080"
	router := http.NewServeMux()
	router.HandleFunc("/", handleHome)
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
}

func checkOtp(w http.ResponseWriter, r *http.Request) {
	// compare otp sent to one in database with corresponding id
}
