package main

import (
	"flag"
	"fmt"
	"log"
	"net/smtp"
	"strings"
)

var (
	from, to, pass string
)

func sendSmtp(authUser, authPass, authServer, smtpHostPort, sender, recipient string) {

	log.Printf("auth=[%s] password=[%s] sender=[%s] recipient=[%s]", authUser, authPass, sender, recipient)

	auth := smtp.PlainAuth(
		"",
		authUser,
		authPass,
		authServer,
	)

	from := fmt.Sprintf("From: <%s>\r\n", sender)
	to := fmt.Sprintf("To: <%s>\r\n", recipient)
	sub := "Subject: Hello\r\n\r\n"
	body := "This is the email body.\r\n"

	err := smtp.SendMail(
		smtpHostPort,
		auth,
		sender,
		[]string{recipient},
		[]byte(from+to+sub+body),
	)
	if err != nil {
		log.Printf("sendSmtp: failure: %q", strings.Split(err.Error(), "\n"))
		log.Fatal(err)
	}
}

func main() {
	flag.StringVar(&from, "from", "", "sender")
	flag.StringVar(&to, "to", "", "recipient")
	flag.StringVar(&pass, "pass", "", "password")
	flag.Parse()

	sendSmtp(from, pass, "smtp.gmail.com", "smtp.gmail.com:587", from, to)
}
