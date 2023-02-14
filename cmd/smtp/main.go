package main

import (
	"log"

	"github.com/quokkamail/quokka/smtp"
)

func main() {
	srv := &smtp.Server{}

	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("smtp server: %v", err)
	}
}
