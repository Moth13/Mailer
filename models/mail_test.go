package models

import (
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestEmail(t *testing.T) {
	to := createRandomEmailAdress()
	subject := createRandomText(20)
	body := createRandomText(100)

	email := Email{
		To:      to,
		Subject: subject,
		Body:    body,
	}

	assert.Equal(t, email.To, to)
	assert.Equal(t, email.Subject, subject)
	assert.Equal(t, email.Body, body)
}

func TestEmailIsScheduledNow(t *testing.T) {
	now := time.Now()

	scheduleAt := []struct {
		name       string
		scheduleAt time.Time
		threshold  time.Duration
		excepted   bool
	}{
		{"Zero", time.Time{}, 5 * time.Minute, false},
		{"InsideMinus", now.Add(-2 * time.Minute), 5 * time.Minute, true},
		{"InsidePlus", now.Add(2 * time.Minute), 5 * time.Minute, true},
		{"EqualNow", now, 5 * time.Minute, true},
		{"OutsideMinus", now.Add(-6 * time.Minute), 5 * time.Minute, false},
		{"OutsidePlus", now.Add(6 * time.Minute), 5 * time.Minute, false},
		{"EqualNoThreshold", now, 0 * time.Minute, true},
	}

	for _, test := range scheduleAt {
		t.Run(test.name, func(t *testing.T) {
			email := Email{ScheduledAt: test.scheduleAt}
			assert.Equal(t, email.IsScheduledNow(now, test.threshold), test.excepted)
		})
	}
}

func TestEmailIsScheduledInPast(t *testing.T) {
	now := time.Now()

	scheduleAt := []struct {
		name       string
		scheduleAt time.Time
		threshold  time.Duration
		excepted   bool
	}{
		{"Zero", time.Time{}, 5 * time.Minute, false},
		{"InsideMinus", now.Add(-2 * time.Minute), 5 * time.Minute, false},
		{"InsidePlus", now.Add(2 * time.Minute), 5 * time.Minute, false},
		{"EqualNow", now, 5 * time.Minute, false},
		{"OutsideMinus", now.Add(-6 * time.Minute), 5 * time.Minute, true},
		{"OutsidePlus", now.Add(6 * time.Minute), 5 * time.Minute, false},
		{"EqualNoThreshold", now, 0 * time.Minute, false},
	}

	for _, test := range scheduleAt {
		t.Run(test.name, func(t *testing.T) {
			email := Email{ScheduledAt: test.scheduleAt}
			assert.Equal(t, email.IsScheduledInPast(now, test.threshold), test.excepted)
		})
	}
}

func TestEmailIsScheduled(t *testing.T) {
	now := time.Now()

	scheduleAt := []struct {
		name       string
		scheduleAt time.Time
		threshold  time.Duration
		excepted   bool
	}{
		{"Zero", time.Time{}, 5 * time.Minute, false},
		{"InsideMinus", now.Add(-2 * time.Minute), 5 * time.Minute, false},
		{"InsidePlus", now.Add(2 * time.Minute), 5 * time.Minute, false},
		{"EqualNow", now, 5 * time.Minute, false},
		{"OutsideMinus", now.Add(-6 * time.Minute), 5 * time.Minute, false},
		{"OutsidePlus", now.Add(6 * time.Minute), 5 * time.Minute, true},
		{"EqualNoThreshold", now, 0 * time.Minute, false},
	}

	for _, test := range scheduleAt {
		t.Run(test.name, func(t *testing.T) {
			email := Email{ScheduledAt: test.scheduleAt}
			assert.Equal(t, email.IsScheduled(now, test.threshold), test.excepted)
		})
	}
}

func createRandomText(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func createRandomEmailAdress() string {
	return createRandomText(10) + "@example.com"
}
