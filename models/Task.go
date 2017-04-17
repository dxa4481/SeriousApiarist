package models

import (
	"errors"
	"fmt"
	"hash/fnv"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"runtime"
	"time"

	"golang.org/x/crypto/ssh"

	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	gitssh "gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"

	"github.com/inconshreveable/log15"
	util "github.com/sevoma/goutil"
	"github.com/spf13/viper"

	"goji.io/pat"
)

// Task - struct for task information
type Task struct {
	Group          string
	Repo           string
	Service        string
	Test           string
	Ref            string
	Commit         string
	Pipeline       string
	CommitTime     time.Time
	Committer      string
	CommitterEmail string
	ProjectPath    string
	StartTime      time.Time
	ImageName      string
	ImageTag       string
}

// NewTask - Checkout on the commit specified and returns task info object
func NewTask(r *http.Request, fw FlushWriter) (Task, error) {
	t0 := time.Now().UTC()

	group := pat.Param(r, "group")
	repo := pat.Param(r, "repo")

	service := r.PostFormValue("service")
	imageTag := r.PostFormValue("imageTag")
	test := r.PostFormValue("test")
	ref := r.PostFormValue("ref")
	commit := r.PostFormValue("commit")
	pipeline := r.PostFormValue("pipeline")

	finalTag := pipeline
	if len(imageTag) > 0 {
		finalTag = imageTag + "-" + pipeline
	}

	// Using path.Join handles the service case nicely, which could be blank
	// if a service is not provided. If service is blank, it's not included in
	// the image name/'path'
	imageName := path.Join(
		viper.GetString("registry"),
		group,
		repo,
		service,
	) + ":" + finalTag

	task := Task{
		Group:    group,
		Repo:     repo,
		Service:  service,
		Test:     test,
		Ref:      ref,
		Commit:   commit,
		Pipeline: pipeline,
		ImageTag: imageTag,
		// These are generated or generated from validated params
		StartTime: t0,
		ImageName: imageName,
	}

	if !util.ValidName(task.Group) {
		return task, errors.New("Invalid group param")
	}
	if viper.GetBool("whitelistGroups") {
		if !util.StringInSlice(task.Group, viper.GetStringSlice("groups")) {
			return task, errors.New("Provided group not whitelisted")
		}
	}

	if !util.ValidName(task.Repo) {
		return task, errors.New("Invalid repo param")
	}

	if !util.ValidName(task.Service) {
		return task, errors.New("Invalid service param")
	}

	if !util.ValidName(task.Test) {
		return task, errors.New("Invalid test param")
	}

	if !util.ValidRef(task.Ref) {
		return task, errors.New("Invalid ref param")
	}

	if !util.ValidName(task.Commit) {
		return task, errors.New("Invalid commit param")
	}

	if !util.ValidInt(task.Pipeline) {
		return task, errors.New("Invalid pipeline param")
	}

	if !util.ValidName(task.ImageTag) {
		return task, errors.New("Invalid imageTag param")
	}

	projectPath, err := checkout(fw, &task)
	if err != nil {
		return task, err
	}
	io.WriteString(fw.W, "Successfully checked out commit\n\n")
	task.ProjectPath = projectPath

	if viper.GetBool("whitelistCommitters") &&
		!util.StringInSlice(task.CommitterEmail,
			viper.GetStringSlice("committers")) {
		return task, errors.New("You are not a whitelisted committer")
	}

	return task, nil
}

// Trace task execution
func Trace(task *Task, handler string, err error) error {

	endingTime := time.Now().UTC()
	duration := endingTime.Sub(task.StartTime)
	deltaMillis := duration.Nanoseconds() / 1e6

	srvLog := log15.New(
		"group", task.Group,
		"repo", task.Repo,
		"service", task.Service,
		"pipeline", task.Pipeline,
		"ref", task.Ref,
		"commit", task.Commit,
		"commitTime", task.CommitTime,
		"committer", task.Committer,
		"committerEmail", task.CommitterEmail,
		"handler", handler,
		"startTime", task.StartTime,
		"numGoroutines", fmt.Sprint(runtime.NumGoroutine()),
		"durationMillis", deltaMillis,
		"env", viper.GetString("env"),
	)

	srvLog.SetHandler(log15.MultiHandler(log15.StreamHandler(os.Stderr,
		log15.JsonFormat())))

	if err != nil {
		srvLog.Error(
			fmt.Sprint(err),
		)
	} else {
		srvLog.Info(
			fmt.Sprint(err),
		)
	}
	return nil
}

func checkout(fw FlushWriter, task *Task) (string, error) {
	projectPath, pathErr := getPath(task)

	buffer, err := ioutil.ReadFile(viper.GetString("gitPrivateKey"))
	if err != nil {
		return projectPath, err
	}
	signer, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		return projectPath, err
	}

	if pathErr != nil {
		io.WriteString(fw.W, "Cloning project.\n")
		err = os.MkdirAll(projectPath, 0700)
		if err != nil {
			return projectPath, err
		}

		gitURL := viper.GetString("gitServer") + ":" +
			task.Group + "/" + task.Repo + ".git"

		_, err = git.PlainClone(projectPath, false, &git.CloneOptions{
			URL:      gitURL,
			Progress: &fw,
			Auth: &gitssh.PublicKeys{
				User:   "git",
				Signer: signer,
			},
		})
		if err != nil {
			return projectPath, err
		}
	}
	r, err := git.PlainOpen(projectPath)
	if err != nil {
		return projectPath, err
	}
	r.Pull(&git.PullOptions{
		RemoteName: task.Ref,
		Progress:   &fw,
		Auth: &gitssh.PublicKeys{
			User:   "git",
			Signer: signer,
		},
	})
	if err != nil {
		return projectPath, err
	}

	wt, err := r.Worktree()
	if err != nil {
		return projectPath, err
	}

	err = wt.Checkout(plumbing.NewHash(task.Commit))
	if err != nil {
		return projectPath, err
	}

	commit, err := r.CommitObject(plumbing.NewHash(task.Commit))
	if err != nil {
		return projectPath, err
	}
	task.Committer = commit.Author.Name
	task.CommitterEmail = commit.Author.Email
	task.CommitTime = commit.Author.When

	return projectPath, err
}

// GetPath returns project path and existence (err)
// Pipeline gets cached if the steps run on the same node
func getPath(task *Task) (string, error) {
	s := []byte(task.Group + task.Repo + task.Ref + task.Commit)
	h := fnv.New32a()
	h.Write(s)

	projectPath := "/" + fmt.Sprint(
		path.Join(
			"tmp",
			"SeriousApiarist",
			fmt.Sprint(h.Sum32())))

	_, err := os.Stat(projectPath)

	return projectPath, err
}
