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

	// Using path.Join handles the service case nicely, which could be blank
	// if a service is not provided. If service is blank, it's not included in
	// the 'path'
	dockerfileFolderPath := path.Join(task.ProjectPath, task.Service)
	dockerfileFilePath := path.Join(dockerfileFolderPath, "Dockerfile")
	_, err := os.Stat(dockerfileFilePath)
	if err != nil {
		return &models.AppTrace{Handler: handler, Task: task, Error: err,
			Message: "Could not find Dockerfile", Code: 500}
	}

	io.WriteString(fw.W, "\nAttempting to build\n\n")
	imageName := task.ImageName + "-test"
	cmd := exec.Command("docker", "build", "-t", imageName, dockerfileFolderPath)
	cmd.Stdout = &fw
	cmd.Stderr = &fw
	err = cmd.Run()
	if err != nil {
		return &models.AppTrace{Handler: handler, Task: task, Error: err,
			Message: "Error occured during build", Code: 500}
	}

	io.WriteString(fw.W, "\n\nAttempting to push\n\n")
	cmd = exec.Command("docker", "push", imageName)
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
