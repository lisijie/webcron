package cron

import (
	"bytes"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/lisijie/webcron/app/models"
	"os/exec"
	"runtime/debug"
	"sync"
	"time"
)

type Job struct {
	id         int
	name       string
	runFunc    func() ([]byte, []byte, error)
	running    sync.Mutex
	Concurrent bool
}

func NewJobFromTask(task *models.Task) (*Job, error) {
	if task.Id < 1 {
		return nil, fmt.Errorf("ToJob: 缺少id")
	}
	job := NewCommandJob(task.Id, task.TaskName, task.Command)
	job.Concurrent = task.Concurrent == 0
	return job, nil
}

func NewCommandJob(id int, name string, command string) *Job {
	job := &Job{
		id:   id,
		name: name,
	}
	job.runFunc = func() ([]byte, []byte, error) {
		bufOut := new(bytes.Buffer)
		bufErr := new(bytes.Buffer)
		cmd := exec.Command("/bin/bash", "-c", command)
		cmd.Stdout = bufOut
		cmd.Stderr = bufErr
		err := cmd.Run()

		return bufOut.Bytes(), bufErr.Bytes(), err
	}
	return job
}

func (j *Job) Run() {
	defer func() {
		if err := recover(); err != nil {
			beego.Error(err, "\n", string(debug.Stack()))
		}
	}()

	t := time.Now()

	if j.Concurrent {
		j.running.Lock()
		defer j.running.Unlock()
	}

	if workPool != nil {
		workPool <- true
		defer func() {
			<-workPool
		}()
	}

	bout, berr, err := j.runFunc()

	ut := time.Now().Sub(t) / time.Millisecond

	log := new(models.TaskLog)
	log.TaskId = j.id
	log.Output = string(bout)
	log.Error = string(berr)
	log.ProcessTime = int(ut)
	log.CreateTime = t.Unix()
	if err != nil {
		log.Status = -1
		log.Error = err.Error() + ":" + string(berr)
	}
	models.TaskLogAdd(log)
}
