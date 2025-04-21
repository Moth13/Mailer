package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/moth13/mailer/mailer"
	"github.com/moth13/mailer/util"
)

var mailerInstance *mailer.Mailer

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		fmt.Println("Error loading config:", err)
		panic(err)
	}

	mailerInstance = mailer.NewMailer(config)
	if mailerInstance == nil {
		fmt.Println("Error creating mailer instance")
		panic("Mailer instance is nil")
	}

	router := gin.Default()

	router.POST("api/mailer/send", postMail)

	err = router.Run("0.0.0.0:8080")
	if err != nil {
		fmt.Println("Error starting server:", err)
		panic(err)
	}
}

func postMail(c *gin.Context) {
	var email mailer.Email
	if err := c.ShouldBindJSON(&email); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	err := mailerInstance.SendEmail(email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send email"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Email sent successfully"})
}
