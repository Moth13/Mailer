package models

import (
	"time"
)

type Status string

const (
	StatusPending          Status = "Pending"
	StatusSent             Status = "Sent"
	StatusErrorScheduledAt Status = "ErrorScheduledAt"
	StatusErrorAtSend      Status = "ErrorAtSend"
)

type Email struct {
	To          string    `json:"to" form:"to" binding:"required" format:"email"`
	Subject     string    `json:"subject" form:"subject" binding:"required"`
	Body        string    `json:"body" form:"body" binding:"required"`
	ScheduledAt time.Time `json:"scheduled_at" form:"scheduled_at" binding:"omitempty" format:"datetime" time_format:"2006-01-02T15:04"`
	Status      Status    `json:"status" form:"status" binding:"omitempty"`
}

func (e *Email) IsScheduledNow(now time.Time, threshold time.Duration) bool {
	// Has to be between now-thereshold and now+threshold
	// to be considered scheduled now
	return !e.ScheduledAt.IsZero() && !e.ScheduledAt.After(now.Add(threshold)) && !e.ScheduledAt.Before(now.Add(-threshold))
}

func (e *Email) IsScheduledInPast(now time.Time, threshold time.Duration) bool {
	return !e.ScheduledAt.IsZero() && e.ScheduledAt.Before(now.Add(-threshold))
}

func (e *Email) IsScheduled(now time.Time, threshold time.Duration) bool {
	return !e.ScheduledAt.IsZero() && e.ScheduledAt.After(now.Add(threshold))
}
