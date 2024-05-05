package gocron

import (
	"time"

	"github.com/go-co-op/gocron/v2"
)

type CronItf interface {
	DeleteVerificationCode(job func() error) error
}

type Cron struct {
	gocron gocron.Scheduler
}

func NewCron(gocron gocron.Scheduler) CronItf {
	return &Cron{gocron: gocron}
}

func (c *Cron) DeleteVerificationCode(job func() error) error {
	_, err := c.gocron.NewJob(
		gocron.DurationJob(time.Minute * 5),
		gocron.NewTask(job),
		gocron.WithName("delete verification code"),
	)
	if err != nil {
		return err
	}
	return nil
}