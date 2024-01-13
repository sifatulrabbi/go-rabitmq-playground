package scheduler

import (
	"time"
)

type Task struct {
	ID            string            `json:"id"`             // auto generated UUID
	Name          string            `json:"name"`           // name for the task
	ExecutionerID string            `json:"executioner_id"` // the function name/id
	Body          map[string]string `json:"body"`           // required info for the task
	Completed     bool              `json:"completed"`      // if the task is completed or not
	ScheduledAt   time.Time         `json:"scheduled_at"`   // when the task should run
	CompletedAt   time.Time         `json:"completed_at"`   // if the task is completed then this is the completion timestamp
}

func NewSchedulerTask(name string, executionerId string, scheduledAt time.Time, body map[string]string) Task {
	st := Task{
		ID:            "",
		Name:          name,
		ExecutionerID: executionerId,
		Body:          body,
		ScheduledAt:   scheduledAt,
		Completed:     false,
	}
	return st
}
