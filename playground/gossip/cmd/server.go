package cmd

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/peizhong/letsgo/playground/gossip/app"
)

func whatip() error {
	interfaces, err := net.Interfaces()
	if err != nil {
		return err
	}
	for _, i := range interfaces {
		addrs, err := i.Addrs()
		if err != nil {
			return err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip != nil {
				println(ip.String(), ip.DefaultMask())
			}
		}
	}
	return nil
}

type SetRequest struct {
	Key, Value string
}

func setHandler(w http.ResponseWriter, r *http.Request) {
	var req SetRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Set error:%s", err.Error())
		return
	}
	err = ui.iCloud.Set(req.Key, req.Value)
	if !ui.CheckError("Set", err, w) {
		return
	}
	w.Write([]byte("OK"))
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Category: %v\n", vars["category"])
}

func infoHandler(w http.ResponseWriter, r *http.Request) {
	members, err := ui.iCloud.Info()
	if !ui.CheckError("Info", err, w) {
		return
	}
	for _, member := range members {
		fmt.Fprintln(w, member)
	}
}

type uiServer struct {
	iCloud app.ICloud
}

func (uiServer) CheckError(operation string, err error, w http.ResponseWriter) bool {
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintf(w, "%s: %s", operation, err.Error())
		return false
	}
	return true
}

var ui *uiServer

func serve(port int, icloud app.ICloud) error {
	if icloud == nil {
		panic("no iCloud")
	}
	ui = &uiServer{iCloud: icloud}
	println("start server", port)
	r := mux.NewRouter()
	r.HandleFunc("/set", setHandler).Methods("PUT", "POST", "DELETE")
	r.HandleFunc("/get", getHandler)
	r.HandleFunc("/info", infoHandler)
	srv := &http.Server{Handler: r, Addr: fmt.Sprintf(":%d", port)}
	if err := srv.ListenAndServe(); err != nil {
		return err
	}
	return nil
}
