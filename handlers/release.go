package handlers

import (
	"io"
	"net/http"
	"os/exec"

	"github.com/sevoma/SeriousApiarist/models"
	util "github.com/sevoma/goutil"
)

// Release endpoint enables releasing images for deploy
func Release(task *models.Task, fw models.FlushWriter, w http.ResponseWriter,
	r *http.Request) *models.AppTrace {
	handler := util.FuncName()

	io.WriteString(fw.W, "\n\nPulling test image\n\n")
	imageName := task.ImageName + "-test"
	cmd := exec.Command("docker", "pull", imageName)
	cmd.Stdout = &fw
	cmd.Stderr = &fw
	err := cmd.Run()
	if err != nil {
		return &models.AppTrace{Handler: handler, Task: task, Error: err,
			Message: "Error occured during pull", Code: 500}
	}

	cmd = exec.Command("docker", "tag", imageName, task.ImageName)
	cmd.Stdout = &fw
	cmd.Stderr = &fw
	err = cmd.Run()
	if err != nil {
		return &models.AppTrace{Handler: handler, Task: task, Error: err,
			Message: "Error occured during tag", Code: 500}
	}

	io.WriteString(fw.W, "\n\nPushing release image\n\n")
	cmd = exec.Command("docker", "push", task.ImageName)
	cmd.Stdout = &fw
	cmd.Stderr = &fw
	err = cmd.Run()
	if err != nil {
		return &models.AppTrace{Handler: handler, Task: task, Error: err,
			Message: "Error occured during push", Code: 500}
	}

	return &models.AppTrace{Handler: handler, Task: task, Error: nil,
		Message: "Success", Code: 200}
}
