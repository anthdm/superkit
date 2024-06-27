package scheduler

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

const cronRegex = `^([0-5][0-9]|\*) ([0-1][0-9]|2[0-4]|\*) ([0-2][0-9]|3[0-1]|\*) ([0-9]|1[0-2]|\*) ([0-6]|\*)$`

type schedule struct {
	min     int
	hour    int
	day     int
	month   int
	weekday int
}

type job struct {
	cmd      func()
	name     string
	schedule schedule
}

type Scheduler struct {
	cv       *sync.Cond
	jobs     []job
	logger   *log.Logger
	mutex    *sync.Mutex
	numReady int
	ready    sync.Map
}

func New(logger *log.Logger) *Scheduler {
	mutex := &sync.Mutex{}
	cv := sync.NewCond(mutex)

	return &Scheduler{
		cv:       cv,
		jobs:     nil,
		logger:   logger,
		mutex:    mutex,
		numReady: 0,
		ready:    sync.Map{},
	}
}

func (s *Scheduler) Add(schedule string, cmd func()) (bool, error) {
	return s.addJob(schedule, cmd, fmt.Sprintf("job-%d", time.Now().UTC().Unix()))
}

func (s *Scheduler) AddWithName(schedule string, cmd func(), name string) (bool, error) {
	return s.addJob(schedule, cmd, name)
}

func (s *Scheduler) addJob(schedule string, cmd func(), name string) (bool, error) {
	ok, err := validateSchedule(schedule)
	if err != nil {
		return false, err
	}

	if !ok {
		return false, nil
	}

	job := job{
		cmd:      cmd,
		schedule: parseSchedule(schedule),
		name:     name,
	}

	s.mutex.Lock()
	s.ready.Store(s.numReady, true)
	s.jobs = append(s.jobs, job)
	s.numReady += 1
	s.mutex.Unlock()

	return true, nil
}

func (s *Scheduler) Start(ctx context.Context) {
	s.logger.Printf("Starting CRON scheduling with [%d] jobs.", len(s.jobs))
	go func() {
		for {
			select {
			case <-ctx.Done():
			default:
				s.mutex.Lock()
				defer s.mutex.Unlock()

				for s.numReady == 0 {
					s.cv.Wait()
				}

				s.ready.Range(func(key, _ any) bool {
					go func() {
						job := s.jobs[key.(int)]

						duration := getDurationTillNextProc(job.schedule)
						s.logger.Printf("Scheduling job [%s] | Time till next proc [%s]", job.name, duration)

						time.Sleep(duration)
						job.cmd()

						s.mutex.Lock()
						defer s.mutex.Unlock()
						s.ready.Store(key, true)
						s.numReady += 1
						s.cv.Signal()
					}()

					s.ready.Delete(key)
					s.numReady -= 1

					return true
				})
			}
		}
	}()
}

func getDurationTillNextProc(s schedule) time.Duration {
	currentDate := time.Now()

	var nextMin int
	if s.min == -1 {
		nextMin = currentDate.Minute() + 1
	} else {
		nextMin = s.min
	}

	var nextHour int
	if s.hour == -1 {
		nextHour = currentDate.Hour()
		if nextMin < currentDate.Minute() {
			nextHour += 1
		}
	} else {
		nextHour = s.hour
	}

	var nextDay int
	if s.day == -1 {
		nextDay = currentDate.Day()
		if nextHour < currentDate.Hour() {
			nextDay += 1
		}
	} else {
		nextDay = s.day
	}

	var nextMonth int
	if s.month == -1 {
		nextMonth = int(currentDate.Month())
		if nextDay < currentDate.Day() {
			nextMonth += 1
		}
	} else {
		nextMonth = s.month
	}

	var nextYear int = currentDate.Year()
	if nextMonth < int(currentDate.Month()) {
		nextYear += 1
	}

	nextDate := time.Date(nextYear, time.Month(nextMonth), nextDay, nextHour, nextMin, 0, 0, currentDate.Location())
	return nextDate.Sub(currentDate)
}

func validateSchedule(schedule string) (bool, error) {
	// For now just match basic numerics and wildcards
	ok, err := regexp.MatchString(cronRegex, schedule)
	if err != nil {
		return false, err
	}

	return ok, nil
}

func parseSchedule(s string) schedule {
	timings := strings.Split(s, " ")

	min := convCronTiming(timings[0], -1)
	hour := convCronTiming(timings[1], -1)
	day := convCronTiming(timings[2], -1)
	month := convCronTiming(timings[3], -1)
	weekday := convCronTiming(timings[4], -1)

	return schedule{
		min:     min,
		hour:    hour,
		day:     day,
		month:   month,
		weekday: weekday,
	}
}

func convCronTiming(timing string, defaultVal int) int {
	if timing != "*" {
		val, err := strconv.Atoi(timing)

		// Conversion should not fail because of regex matching string prior
		if err != nil {
			panic("Conversion of cron timing should not have failed")
		}
		return val
	}

	return defaultVal
}
