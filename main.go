package main

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/moth13/mailer/mailer"
	"github.com/moth13/mailer/util"
	"github.com/moth13/mailer/worker"
)

var mailerInstance *mailer.Mailer
var workerPool worker.WorkerPool
var maxWorkers = 5

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

	wg := &sync.WaitGroup{}
	workerPool = worker.NewWorkerPool(maxWorkers, wg)
	if workerPool == nil {
		fmt.Println("Error creating worker pool")
		panic("Worker pool is nil")
	}

	workerPool.Run()

	router := gin.Default()

	router.POST("api/mailer/send", postMail)

	err = router.Run("0.0.0.0:8080")
	if err != nil {
		fmt.Println("Error starting server:", err)
		panic(err)
	}
	workerPool.Stop()
}

func postMail(c *gin.Context) {
	var email mailer.Email
	if err := c.ShouldBindJSON(&email); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	workerPool.AddTask(func() {
		if err := mailerInstance.SendEmail(email); err != nil {
			fmt.Println("Error sending email:", err)
		}
	})

	c.JSON(http.StatusOK, gin.H{"message": "Email sent successfully"})
}
