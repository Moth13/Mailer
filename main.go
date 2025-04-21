package main

import (
	"fmt"

	"github.com/moth13/mailer/mailer"
	"github.com/moth13/mailer/util"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		fmt.Println("Error loading config:", err)
		return
	}

	mailerInstance := mailer.NewMailer(config)

	mail := mailer.Email{
		To:      "jeremie.guerinel@gmail.com",
		Subject: "Test mail",
		Body:    "Email body",
	}

	err = mailerInstance.SendEmail(mail)
	if err != nil {
		fmt.Println("Error sending email:", err)
		return
	}
	fmt.Println("Email sent successfully!")
}
