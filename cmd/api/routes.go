package main

import (
	"encoding/json"
	"net/http"

	k8s "github.com/namnt2307/jupyterhub-freeport/pkg/kubernetes"
)

type PostDataFormat struct {
	Namespace     string `json:"namespace"`
	PodName       string `json:"podName"`
	NodeSelector  string `json:"nodeSelector"`
	CpuLimit      string `json:"cpuLimit"`
	CpuRequest    string `json:"cpuRequest"`
	MemoryLimit   string `json:"memoryLimit"`
	MemoryRequest string `json:"memoryRequest"`
}
type ReturnDataFormat struct {
	HostIP string `json:"hostIP"`
}

func (app *application) getport(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/api/getSpawnNode" {
		http.NotFound(w, r)
	}
	switch r.Method {
	case http.MethodGet:
		w.Write([]byte("hello"))
	case http.MethodPost:
		var myData PostDataFormat
		err := json.NewDecoder(r.Body).Decode(&myData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		//get ip
		hostIP := k8s.MakePod(app.clientSet, myData.Namespace, myData.PodName, myData.NodeSelector, myData.CpuLimit, myData.CpuRequest, myData.MemoryLimit, myData.MemoryRequest)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(&ReturnDataFormat{HostIP: hostIP})
	default:
		app.clientError(w, http.StatusMethodNotAllowed)
	}

}
