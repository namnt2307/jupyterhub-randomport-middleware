package main

import "net/http"

func (App *Application) Routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/getSpawnNode", App.GetPort)
	return mux
}
