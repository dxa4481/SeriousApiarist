package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path"

	"github.com/sevoma/SeriousApiarist/models"
	"github.com/sevoma/SeriousApiarist/util"
	"github.com/sevoma/goutil"
	"github.com/spf13/viper"
)

// Deploy endpoint enables stack deploys
func Deploy(task *models.Task, fw models.FlushWriter, w http.ResponseWriter,
	r *http.Request) *models.AppTrace {
	handler := goutil.FuncName()

	stackFilePath := path.Join(task.ProjectPath, "docker-compose.yml")
	_, err := os.Stat(stackFilePath)
	if err != nil {
		return &models.AppTrace{Handler: handler, Task: task, Error: err,
			Message: "docker-compose.yml not found", Code: 500}
	}

	//  docker login -u gitlab-ci-token -p $CI_BUILD_TOKEN $REGISTRY
	io.WriteString(fw.W, "\n\nLogging into the registry\n\n")
	registryUser, err := goutil.GetSecret("registryUser")
	if err != nil {
		return &models.AppTrace{Handler: handler, Task: task, Error: err,
			Message: "Error occured while fetching registry login secret", Code: 500}
	}
	registryPassword, err := goutil.GetSecret("registryPassword")
	if err != nil {
		return &models.AppTrace{Handler: handler, Task: task, Error: err,
			Message: "Error occured while fetching registry password secret", Code: 500}
	}
	cmd := exec.Command("docker", "login",
		"-u", registryUser,
		"-p", registryPassword,
		viper.GetString("registry"))
	cmd.Stdout = &fw
	cmd.Stderr = &fw
	err = cmd.Run()
	if err != nil {
		return &models.AppTrace{Handler: handler, Task: task, Error: err,
			Message: "Error occured during deploy", Code: 500}
	}

	io.WriteString(fw.W, "\n\nSending Duo 2FA push to your device\n\n")
	err = util.DuoPush(task.Committer, task.CommitterEmail,
		task.Group, task.Repo)
	if err != nil {
		return &models.AppTrace{Handler: handler, Task: task, Error: err,
			Message: "Duo 2FA push failed", Code: 401}
	}

	io.WriteString(fw.W, "\n\nAttempting stack deploy\n\n")
	cmd = exec.Command("docker", "stack", "deploy",
		"--with-registry-auth",
		"--compose-file",
		stackFilePath, task.Repo)
	cmd.Stdout = &fw
	cmd.Stderr = &fw
	err = cmd.Run()
	if err != nil {
		return &models.AppTrace{Handler: handler, Task: task, Error: err,
			Message: "Error occured during deploy", Code: 500}
	}

	util.Alert(fmt.Sprintf("%s <%s> successfully deployed %s/%s",
		task.Committer, task.CommitterEmail, task.Group, task.Repo))

	return &models.AppTrace{Handler: handler, Task: task, Error: nil,
		Message: "Success", Code: 200}
}
