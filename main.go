package main

import (
	"fmt"
	"io"
	"net/http"

	goji "goji.io"

	"github.com/sevoma/SeriousApiarist/handlers"
	"github.com/sevoma/SeriousApiarist/models"
	"github.com/spf13/viper"

	"goji.io/pat"
)

type appHandler func(*models.Task, models.FlushWriter,
	http.ResponseWriter, *http.Request) *models.AppTrace

// Middleware providing task information to handlers, setting content type,
// and logging detailed result traces for each request
func (fn appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fw := models.NewFlushWriter(w)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	// Validate and produce Task object from params
	// It's provided to every handler
	task, err := models.NewTask(r, fw)
	if err != nil {
		models.Trace(&task, "Task", err)
		io.WriteString(fw.W, err.Error())
		io.WriteString(fw.W, "\n\nSTAGE FAILED")
		return
	}

	// Collect result traces from every request to the API
	e := fn(&task, fw, w, r)
	if e != nil {
		if e.Error != nil { // e is *appError, not os.Error.
			models.Trace(e.Task, e.Handler, e.Error)
			// Because it's a streaming request, status code 200 was sent
			// at the beginning. Nothing we can do about that now :[
			// http.Error(w, e.Message, e.Code)
			// Instead let's include it in the final line so the client can parse it.
			io.WriteString(fw.W, e.Message)
			io.WriteString(fw.W, "\n\nSTAGE FAILED")
		} else {
			models.Trace(e.Task, e.Handler, nil)
			io.WriteString(fw.W, "\n\nSTAGE SUCCESSFUL")
		}
	} else {
		models.Trace(&task, "nil returned instead of *models.AppTrace", e.Error)
	}
}

func main() {
	// Get config
	viper.SetConfigName("config") // name of config file (without extension)
	viper.AddConfigPath("/")
	viper.AddConfigPath("$HOME")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig() // yaml, toml, json, ini, whatever
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s", err))
	}

	mux := goji.NewMux()
	mux.Handle(pat.Post("/build/:group/:repo"), appHandler(handlers.Build))
	mux.Handle(pat.Post("/test/:group/:repo"), appHandler(handlers.Test))
	mux.Handle(pat.Post("/scan/:group/:repo"), appHandler(handlers.Scan))
	mux.Handle(pat.Post("/release/:group/:repo"), appHandler(handlers.Release))
	mux.Handle(pat.Post("/deploy/:group/:repo"), appHandler(handlers.Deploy))
	http.ListenAndServe(":8080", mux)
}
