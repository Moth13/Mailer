package main

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"
	"github.com/moth13/mailer/mailer"
	"github.com/moth13/mailer/models"
	"github.com/moth13/mailer/scheduler"
	"github.com/moth13/mailer/util"
	"github.com/moth13/mailer/views"
	"github.com/moth13/mailer/worker"
)

var mailerInstance *mailer.Mailer
var workerPool worker.WorkerPool
var schedulerInstance *scheduler.Scheduler
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

	schedulerInstance = scheduler.NewScheduler(time.Second*5, &workerPool, mailerInstance)
	if schedulerInstance == nil {
		fmt.Println("Error creating scheduler instance")
		panic("Scheduler instance is nil")
	}

	schedulerInstance.Run()

	router := gin.Default()

	router.GET("/", indexPageHandler)
	router.GET("/mails", mailsPageHandler)
	router.POST("/mails", postMailsPageHandler)

	router.POST("api/mailer/send", postMail)

	err = router.Run("0.0.0.0:8080")
	if err != nil {
		fmt.Println("Error starting server:", err)
		panic(err)
	}
	workerPool.Stop()
}

func render(ctx *gin.Context, status int, template templ.Component) error {
	ctx.Status(status)
	return template.Render(ctx.Request.Context(), ctx.Writer)
}

func indexPageHandler(c *gin.Context) {
	_, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	err := render(c, http.StatusOK, views.Index())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}
}

func mailsPageHandler(c *gin.Context) {
	_, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	err := render(c, http.StatusOK, views.Mails())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}
}

func postMailsPageHandler(c *gin.Context) {
	var email models.Email
	if err := c.ShouldBind(&email); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid request %s", err.Error())})
		return
	}
	if email.ScheduledAt.IsZero() {
		workerPool.AddTask(func() error { return mailerInstance.SendEmail(email) })
	} else {
		schedulerInstance.AddMail(email)
	}

	_, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	err := render(c, http.StatusOK, views.Mails())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}
}

func postMail(c *gin.Context) {
	var email models.Email
	if err := c.ShouldBindJSON(&email); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid request %s", err.Error())})
		return
	}

	if email.ScheduledAt.IsZero() {
		workerPool.AddTask(func() error { return mailerInstance.SendEmail(email) })
	} else {
		schedulerInstance.AddMail(email)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Email sent successfully"})
}
