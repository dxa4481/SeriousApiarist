package handlers

import (
	"io"
	"net/http"
	"os"
	"os/exec"
	"path"

	"github.com/sevoma/SeriousApiarist/models"
	util "github.com/sevoma/goutil"
)

// Build handler builds services and pushes them to the registry
func Build(task *models.Task, fw models.FlushWriter, w http.ResponseWriter,
	r *http.Request) *models.AppTrace {
	handler := util.FuncName()

	io.WriteString(fw.W, "\nAttempting to build\n\n")
	err := build(task, fw)
	if err != nil {
		return &models.AppTrace{Handler: handler, Task: task, Error: err,
			Message: "Error occured during build", Code: 500}
	}

	io.WriteString(fw.W, "\n\nAttempting to push\n\n")
	err = push(task, fw)
	if err != nil {
		return &models.AppTrace{Handler: handler, Task: task, Error: err,
			Message: "Error occured during push", Code: 500}
	}

	return &models.AppTrace{Handler: handler, Task: task, Error: nil,
		Message: "Success", Code: 200}
}

func build(task *models.Task, fw models.FlushWriter) error {
	// Using path.Join handles the service case nicely, which could be blank
	// if a service is not provided. If service is blank, it's not included in
	// the 'path'
	dockerfileFilePath := path.Join(task.DockerfileFolderPath, "Dockerfile")
	_, err := os.Stat(dockerfileFilePath)
	if err != nil {
		return err
	}

	imageName := task.ImageName + "-test"
	cmd := exec.Command("docker", "build", "-t", imageName, task.DockerfileFolderPath)
	cmd.Stdout = &fw
	cmd.Stderr = &fw
	err = cmd.Run()
	return err
}

func push(task *models.Task, fw models.FlushWriter) error {
	imageName := task.ImageName + "-test"
	cmd := exec.Command("docker", "push", imageName)
	cmd.Stdout = &fw
	cmd.Stderr = &fw
	err := cmd.Run()
	return err
}
