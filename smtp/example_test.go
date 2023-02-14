package smtp_test

import (
	"log"

	"github.com/quokkamail/quokka/smtp"
)

func ExampleServer_ListenAndServe() {
	var srv smtp.Server

	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("SMTP server ListenAndServe: %v", err)
	}
}
