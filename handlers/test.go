package handlers

import (
	"errors"
	"io"
	"net/http"
	"os/exec"

	"github.com/sevoma/SeriousApiarist/models"
	util "github.com/sevoma/goutil"
)

// Test endpoint enables running service tests
func Test(task *models.Task, fw models.FlushWriter, w http.ResponseWriter,
	r *http.Request) *models.AppTrace {
	handler := util.FuncName()

	if len(task.Test) == 0 {
		return &models.AppTrace{Handler: handler, Task: task,
			Error:   errors.New("Empty test param"),
			Message: "No test provided", Code: 400}
	}

	io.WriteString(fw.W, "\n\nAttempting test\n\n")
	imageName := task.ImageName + "-test"
	cmd := exec.Command("docker", "run", "--rm", "-t", imageName, task.Test)
	cmd.Stdout = &fw
	cmd.Stderr = &fw
	err := cmd.Run()
	if err != nil {
		return &models.AppTrace{Handler: handler, Task: task, Error: err,
			Message: "Error occured during test", Code: 500}
	}

	return &models.AppTrace{Handler: handler, Task: task, Error: nil,
		Message: "Success", Code: 200}
}
