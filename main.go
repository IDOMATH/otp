package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/IDOMATH/otp/db"
	"github.com/IDOMATH/otp/mailer"

	"github.com/IDOMATH/CheetahMath/base_change"
	"github.com/IDOMATH/CheetahUtil/env"
	"github.com/IDOMATH/session/memorystore"
)

type Repository struct {
	MemStore *memorystore.MemoryStore[string]
	OtpStore db.OtpStore
	Mail     *mailer.Mailer
}

func main() {
	fmt.Println("Hello world!")
	serverPort := env.GetEnvValueOrDefault("PORT", ":8080")
	router := http.NewServeMux()

	// Memstore := memorystore.New[string]()
	repo := Repository{}
	repo.Mail = setUpMailer()

	repo.OtpStore = db.NewOtpStore(setupDbConnection())

	router.HandleFunc("GET /", handleHome)
	router.HandleFunc("POST /otp", repo.sendOtp)
	router.HandleFunc("POST /otp/{id}", repo.checkOtp)
	server := http.Server{
		Addr:    fmt.Sprint(":", serverPort),
		Handler: router,
	}

	log.Fatal(server.ListenAndServe())
}

func setupDbConnection() *db.DB {
	dbHost := env.GetEnvValueOrDefault("DBHOST", "localhost")
	dbPort := env.GetEnvValueOrDefault("DBPORT", "5432")
	dbName := env.GetEnvValueOrDefault("DBNAME", "tournament-finder")
	dbUser := env.GetEnvValueOrDefault("DBUSER", "postgres")
	dbPass := env.GetEnvValueOrDefault("DBPASS", "postgres")
	dbSsl := env.GetEnvValueOrDefault("DBSSL", "disable")

	connectionString := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s", dbHost, dbPort, dbName, dbUser, dbPass, dbSsl)
	fmt.Println("Connecting to Postgres")
	postgresDb, err := db.ConnectSql(connectionString)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to Postgres")
	return postgresDb
}

func setUpMailer() *mailer.Mailer {

	from := env.GetEnvValue("EMAIL")
	if from == "" {
		log.Fatal("error getting EMAIL from env")
	}

	password := env.GetEnvValue("PASSWORD")
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
	pass := MakeOtp(8)
	id := repo.OtpStore.InsertOtp(pass)
	err := repo.Mail.SendEmail(r.PostFormValue("email"), fmt.Sprintf("Your OTP is: %s", pass))
	if err != nil {
		fmt.Println("error sending email: ", err)
		return
	}
	// Maybe a redirect to a page with a form and the ID.
	// Like a GET /otp/{id}
	w.Write([]byte(id))
}

func MakeOtp(len int) string {
	otp := strings.Builder{}
	r := rand.New(rand.NewSource(time.Now().UnixMicro()))
	for range len {
		otp.Write([]byte(base_change.FromTen(36, r.Int63n(36))))
	}
	return otp.String()
}

func (repo *Repository) checkOtp(w http.ResponseWriter, r *http.Request) {
	// compare otp sent to one in database with corresponding id
	passIn := r.PostFormValue("pass")
	storedPass := repo.OtpStore.Get(r.PathValue("id"))
	if passIn == string(storedPass) {
		w.WriteHeader(http.StatusOK)
		repo.OtpStore.DeleteOtp(r.PathValue("id"))
		return
	}
	w.WriteHeader(http.StatusForbidden)
	w.Write([]byte("OTP not match"))
}
