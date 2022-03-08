package main

import "net/http"

func (App *Application) Routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/getSpawnNode", App.GetPort)
	return mux
}
