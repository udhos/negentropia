package main

import (
	"flag"
	"fmt"
	"log"
	"net/smtp"
)

var (
	from, to, pass string
)

func send(authUser, authPass, authServer, smtpHostPort, sender, recipient string) {

	log.Printf("auth=[%s] password=[%s] sender=[%s] recipient=[%s]", authUser, authPass, sender, recipient)

	// Set up authentication information.
	auth := smtp.PlainAuth(
		"",
		authUser,
		authPass,
		authServer,
	)

	from := fmt.Sprintf("From: <%s>\r\n", sender)
	to := fmt.Sprintf("To: <%s>\r\n", recipient)
	sub := "Subject: Hello\r\n\r\n"
	body := "This is the email body."

	// Connect to the server, authenticate, set the sender and recipient,
	// and send the email all in one step.
	err := smtp.SendMail(
		smtpHostPort,
		auth,
		sender,
		[]string{recipient},
		[]byte(from+to+sub+body),
	)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	flag.StringVar(&from, "from", "", "sender")
	flag.StringVar(&to, "to", "", "recipient")
	flag.StringVar(&pass, "pass", "", "password")
	flag.Parse()

	send(from, pass, "smtp.gmail.com", "smtp.gmail.com:587", from, to)
}
