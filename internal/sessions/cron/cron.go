package cron

import (
	"fmt"
	"sync"
)

type Job struct {
	ID       string
	Schedule string
	Task     string
}

type CronManager struct {
	mu   sync.RWMutex
	jobs map[string]*Job
}

var DefaultCronManager = &CronManager{
	jobs: make(map[string]*Job),
}

func (cm *CronManager) CreateSchedule(id, schedule, task string) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	cm.jobs[id] = &Job{
		ID:       id,
		Schedule: schedule,
		Task:     task,
	}
	return nil
}

func (cm *CronManager) ListSchedules() []Job {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	var list []Job
	for _, j := range cm.jobs {
		list = append(list, *j)
	}
	return list
}

func (cm *CronManager) DeleteSchedule(id string) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if _, ok := cm.jobs[id]; !ok {
		return fmt.Errorf("schedule not found: %s", id)
	}
	delete(cm.jobs, id)
	return nil
}
