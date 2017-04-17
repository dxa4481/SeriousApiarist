package models

// AppTrace - Request error handling wrapper on the handler
type AppTrace struct {
	Handler string
	Task    *Task
	Error   error
	Message string
	Code    int
}
