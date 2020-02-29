package main

import (
	"github.com/go-cmd/cmd"
	"github.com/kataras/iris/v12"
)

var (
	tasksMap    = map[uint64]*cmd.Cmd{}    // all tasks
	runningChan = make(chan *cmd.Cmd, 10)  // 10 tasks can be running simultaneously
	tasksChan   = make(chan *cmd.Cmd, 100) // total queued task is 100
)

// fail before start condition:
//		status.Error != nil
// pending condition:
//		status.StartTs == 0
// running condition:
//		status.StopTs == 0
// finished with error condition:
//		status.Complete && status.Exit != 0
func taskStatusHandler(context iris.Context) {
	userID := getUserID(context.GetHeader("Authorization"))
	c := tasksMap[userID]
	if c == nil {
		context.StatusCode(iris.StatusBadRequest)
		context.JSON(iris.Map{
			"error": "no previous task exists",
		})
		return
	}
	status := c.Status()
	if status.Error != nil {
		context.StatusCode(iris.StatusInternalServerError)
		context.JSON(iris.Map{
			"cmdStatus": "failed",
			"command":   append([]string{c.Name}, c.Args...),
			"error":     status.Error,
		})
		return
	}
	if status.StartTs == 0 && !status.Complete {
		context.StatusCode(iris.StatusOK)
		context.JSON(iris.Map{
			"cmdStatus": "pending",
			"command":   append([]string{c.Name}, c.Args...),
		})
		return
	}
	if status.StopTs == 0 && !status.Complete {
		context.StatusCode(iris.StatusOK)
		context.JSON(iris.Map{
			"cmdStatus": "running",
			"command":   append([]string{c.Name}, c.Args...),
		})
		return
	}
	if status.Complete && status.Exit != 0 {
		context.StatusCode(iris.StatusInternalServerError)
		context.JSON(iris.Map{
			"cmdStatus": "failed",
			"command":   append([]string{c.Name}, c.Args...),
		})
		return
	}
	if status.Complete && status.Exit == 0 {
		context.StatusCode(iris.StatusOK)
		context.JSON(iris.Map{
			"cmdStatus": "finished",
			"command":   append([]string{c.Name}, c.Args...),
		})
		return
	}
}

func taskResultHandler(context iris.Context) {
	userID := getUserID(context.GetHeader("Authorization"))
	c := tasksMap[userID]
	if c == nil {
		context.StatusCode(iris.StatusBadRequest)
		context.JSON(iris.Map{
			"error": "no previous task exists",
		})
		return
	}
	status := c.Status()
	context.StatusCode(iris.StatusOK)
	context.JSON(iris.Map{
		"result": status,
	})
}

// taskRunner limits Running tasks to 10 at a time
func taskRunner() {
	for t := range tasksChan {
		runningChan <- t
		go func() {
			<-t.Done()
			<-runningChan
		}()
		t.Start()
	}
}
