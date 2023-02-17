package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/quokkamail/quokka/smtp"
)

func main() {
	smtpSrv := &smtp.Server{}

	// Start the "smtp" server on standard port 25
	go func() {
		if err := smtpSrv.ListenAndServe(); err != nil {
			log.Fatalf("smtp server: %v", err)
		}
	}()

	submissionSrv := &smtp.SubmissionServer{}

	// Start the "submission" server on standard port 587
	go func() {
		if err := submissionSrv.ListenAndServe(); err != nil {
			log.Fatalf("submission server: %v", err)
		}
	}()

	submissionsSrv := &smtp.SubmissionsServer{}

	// Start the "submissions" server on standard port 465
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
