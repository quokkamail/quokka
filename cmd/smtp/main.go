package main

import (
	"crypto/tls"
	"log"
	"os"
	"os/signal"

	"github.com/quokkamail/quokka/smtp"
)

func main() {
	// private key: openssl ecparam -name prime256v1 -genkey -noout -out key.pem
	// cert: openssl req -new -x509 -key key.pem -out cert.pem -days 365
	cert, err := tls.LoadX509KeyPair("cert.pem", "key.pem")
	if err != nil {
		log.Fatalf("tls: %v", err)
	}

	tlsConfig := &tls.Config{
		MinVersion:   tls.VersionTLS12,
		Certificates: []tls.Certificate{cert},
	}

	// Start the "smtp" server on standard port 25
	smtpSrv := &smtp.Server{
		Addr:                    ":smtp",
		Domain:                  "quokka.local",
		TLSConfig:               tlsConfig,
		AuthenticationEncrypted: true,
	}

	go func() {
		if err := smtpSrv.ListenAndServe(); err != nil {
			log.Fatalf("smtp server: %v", err)
		}
	}()

	// Start the "submission" server on standard port 587
	submissionSrv := &smtp.Server{
		Addr:                    ":submission",
		Domain:                  "quokka.local",
		TLSConfig:               tlsConfig,
		AuthenticationEncrypted: true,
		AuthenticationMandatory: true,
	}

	go func() {
		if err := submissionSrv.ListenAndServe(); err != nil {
			log.Fatalf("submission server: %v", err)
		}
	}()

	// Start the "submissions" server on standard port 465
	submissionsSrv := &smtp.Server{
		Addr:                    ":465",
		Domain:                  "quokka.local",
		TLSConfig:               tlsConfig,
		AuthenticationEncrypted: true,
		AuthenticationMandatory: true,
	}

	go func() {
		if err := submissionsSrv.ListenAndServeTLS(); err != nil {
			log.Fatalf("submissions server: %v", err)
		}
	}()

	interruptSignal := make(chan os.Signal, 1)
	signal.Notify(interruptSignal, os.Interrupt)

	<-interruptSignal
	log.Println("got an interrupt signal")

	if err := smtpSrv.Close(); err != nil {
		log.Fatal(err)
	}

	if err := submissionSrv.Close(); err != nil {
		log.Fatal(err)
	}

	if err := submissionsSrv.Close(); err != nil {
		log.Fatal(err)
	}
}
