package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

func (App *Application) ServerError(w http.ResponseWriter, err error) {

	//debug.stack to get a stack trace for current goroutine and append to log msg
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	App.errorLog.Output(2, trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (App *Application) ClientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (App *Application) NotFound(w http.ResponseWriter) {
	App.ClientError(w, http.StatusNotFound)
}
