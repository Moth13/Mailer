package scheduler

import (
	"time"

	"github.com/moth13/mailer/mailer"
	"github.com/moth13/mailer/models"
	"github.com/moth13/mailer/worker"
)

type Scheduler struct {
	Interval       time.Duration
	Emails         chan models.Email
	WorkerPool     *worker.WorkerPool
	MailerInstance *mailer.Mailer
}

func NewScheduler(interval time.Duration, workerPool *worker.WorkerPool, mailerInstance *mailer.Mailer) *Scheduler {
	return &Scheduler{
		Interval:       interval,
		Emails:         make(chan models.Email),
		WorkerPool:     workerPool,
		MailerInstance: mailerInstance,
	}
}

func (s *Scheduler) Run() {
	go func() {
		ticker := time.NewTicker(s.Interval)
		defer ticker.Stop()

		for range ticker.C {
			go s.checkEmails()
		}
	}()
}

func (s *Scheduler) AddMail(email models.Email) {
	s.Emails <- email
}

func (s *Scheduler) checkEmails() {
	for {
		select {
		case email := <-s.Emails:
			if email.IsScheduledNow(time.Now(), s.Interval/2) {
				(*s.WorkerPool).AddTask(func() error {
					return s.MailerInstance.SendEmail(email)
				})
			} else if email.IsScheduledInPast(time.Now(), s.Interval/2) {
				email.Status = models.StatusErrorScheduledAt
			} else {
				email.Status = models.StatusPending
				s.Emails <- email
			}
		default:
			return
		}
	}
}
