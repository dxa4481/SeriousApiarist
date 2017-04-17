package handlers

import (
	"io"
	"net/http"

	"github.com/sevoma/SeriousApiarist/models"
	util "github.com/sevoma/goutil"
)

// Scan endpoint enables running vuln scans on images
func Scan(task *models.Task, fw models.FlushWriter, w http.ResponseWriter,
	r *http.Request) *models.AppTrace {
	handler := util.FuncName()

	io.WriteString(fw.W, "\n\nVuln scan not implemented\n\n")

	return &models.AppTrace{Handler: handler, Task: task, Error: nil,
		Message: "Success", Code: 200}
}
