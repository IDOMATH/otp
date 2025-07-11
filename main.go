package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/IDOMATH/CheetahMath/base_change"
	"github.com/IDOMATH/session/memorystore"
	"github.com/idomath/CheetahFarm/otp/mailer"
)

type Repository struct {
	mStore *memorystore.MemoryStore
	mail   *mailer.Mailer
}

func main() {
	fmt.Println("Hello world!")
	serverPort := "8080"
	router := http.NewServeMux()

	memstore := memorystore.New()
	repo := Repository{mStore: memstore}
	repo.mail = setUpMailer()

	router.HandleFunc("GET /", handleHome)
	router.HandleFunc("POST /otp", repo.sendOtp)
	router.HandleFunc("POST /otp/{id}", repo.checkOtp)
	server := http.Server{
		Addr:    fmt.Sprint(":", serverPort),
		Handler: router,
	}

	log.Fatal(server.ListenAndServe())
}

func setUpMailer() *mailer.Mailer {

	from := os.Getenv("EMAIL")
	if from == "" {
		log.Fatal("error getting EMAIL from env")
	}

	password := os.Getenv("PASSWORD")
	if password == "" {
		log.Fatal("error geting PASSWORD from env")
	}

	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	return mailer.NewMailer(smtpHost, smtpPort, from, password)
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello world!"))
}

func (repo *Repository) sendOtp(w http.ResponseWriter, r *http.Request) {
	// find a way to make a unique id.  Something with the user email
	// add otp and id to database
	// send email
	pass := MakeOtp(8)
	w.Write([]byte(pass))
	err := repo.mail.SendEmail(r.PostFormValue("email"), fmt.Sprintf("Your OTP is: %s", pass))
	if err != nil {
		fmt.Println("error sending email: ", err)
		return
	}
	repo.mStore.Insert("001", []byte(pass), time.Now().Add(time.Hour))
}

func MakeOtp(len int) string {
	otp := strings.Builder{}
	r := rand.New(rand.NewSource(time.Now().UnixMicro()))
	for i := 0; i < len; i++ {
		otp.Write([]byte(base_change.FromTen(36, r.Int63n(36))))
	}
	return otp.String()
}

func (repo *Repository) checkOtp(w http.ResponseWriter, r *http.Request) {
	// compare otp sent to one in database with corresponding id
	passIn := r.PostFormValue("pass")
	storedPass, _ := repo.mStore.Get(r.PathValue("id"))
	if passIn == string(storedPass) {
		w.Write([]byte("OTP good"))
		return
	}
	w.Write([]byte("OTP not match"))
}
