package gatewayservice

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
)

/*
consul for service discovery
*/

func init() {

}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%v", vars["all"])
}

func _main() {
	r := mux.NewRouter()
	r.HandleFunc("/{all}/", HomeHandler)
	http.Handle("/", r)
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	go func() {
		if err := http.ListenAndServe("localhost:8080", nil); err != nil {
			log.Error(err.Error())
		}
	}()
	select {
	case <-ch:
		log.Println("Program exit")
	}
}

func Run() {
	r := mux.NewRouter()
	r.HandleFunc("/{all}/", HomeHandler)
	http.Handle("/", r)
	if err := http.ListenAndServe("localhost:8080", nil); err != nil {
		log.Error(err.Error())
	}
}
