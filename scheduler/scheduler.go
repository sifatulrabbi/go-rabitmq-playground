package scheduler

import (
	"fmt"
	"log"
	"time"
)

type Scheduler struct {
	Tasks         []Task
	Executioners  map[string]func(body map[string]string) error
	LastTaskIndex int
}

// this will keep looping forever make sure to run this in a separate go routine.
func (s *Scheduler) StartInBg(lastIndex int) {
	time.Sleep(time.Second * 5)
	if lastIndex >= len(s.Tasks) {
		s.StartInBg(0)
		return
	}

	task := &s.Tasks[lastIndex]
	if task.Completed {
		s.StartInBg(lastIndex + 1)
		return
	}

	execFn := s.Executioners[task.ExecutionerID]
	if execFn == nil {
		fmt.Println("no executioner with id:", task.ExecutionerID)
		s.StartInBg(lastIndex + 1)
		return
	}

	err := execFn(task.Body)
	if err != nil {
		log.Printf("Error in %s, err: %s", task.Name, err)
	}

	task.Completed = true
	task.CompletedAt = time.Now()

	s.StartInBg(lastIndex + 1)
}

func (s *Scheduler) AddNewExecutioner(id string, fn func(body map[string]string) error) {
	s.Executioners[id] = fn
}

func (s *Scheduler) RemoveExecutioner(id string) {
}

func (s *Scheduler) AddNewTask(task Task) {
	// verify the executioner id is available
	// verify the scheduled_at is not in the past
	s.Tasks = append(s.Tasks, task)
}

func (s *Scheduler) GetATaskById(id string) *Task {
	// find and return the task from the tasks repo
	return nil
}
