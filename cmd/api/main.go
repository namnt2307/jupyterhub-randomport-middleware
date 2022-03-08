package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	k8s "github.com/namnt2307/jupyterhub-freeport/pkg/kubernetes"
	"k8s.io/client-go/kubernetes"
)

// logger struct
type Application struct {
	errorLog  *log.Logger
	infoLog   *log.Logger
	clientSet *kubernetes.Clientset
}

func main() {
	// default port
	addr := flag.String("addr", ":8000", "HTTP network address")
	flag.Parse()

	// init log
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// load kubernetes config
	clienSet := k8s.InitKubernetes()
	infoLog.Printf("Create clientset successfully")

	// init app
	App := &Application{
		errorLog:  errorLog,
		infoLog:   infoLog,
		clientSet: clienSet,
	}

	// init mux server
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  App.Routes(),
	}
	infoLog.Printf("starting server on %s", *addr)
	err := srv.ListenAndServe()
	errorLog.Fatal(err)
}
