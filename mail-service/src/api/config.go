package api

import (
	"log"
	"net/http"
	"os"
	"strconv"
)

type Config struct {
	WebPort string
	Mailer  Mail
}

func (app *Config) Serve() {
	// Serve the app
	srv := &http.Server{
		Addr:    ":" + app.WebPort,
		Handler: app.Routes(),
	}
	log.Println("Mail service litening on port", srv.Addr)
	if err := srv.ListenAndServe(); err != nil {
		panic(err)
	}
}

func CreateMailer() Mail {
	port, _ := strconv.Atoi(os.Getenv("MAIL_PORT"))
	return Mail{
		Domain:      os.Getenv("MAIL_DOMAIN"),
		Host:        os.Getenv("MAIL_HOST"),
		Port:        port,
		Username:    os.Getenv("MAIL_USERNAME"),
		Password:    os.Getenv("MAIL_PASSWORD"),
		Encryption:  os.Getenv("MAIL_ENCRYPTION"),
		FromName:    os.Getenv("MAIL_FROMNAME"),
		FromAddress: os.Getenv("MAIL_FROMADDRESS"),
	}
}
