package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"

	auth "github.com/abbot/go-http-auth"
	"github.com/golang/glog"
	"github.com/gorilla/mux"
)

type HttpServer struct {
	graph    *Graph
	listener net.Listener
	iface    string
	port     string
	htpasswd string

	listenerMtx sync.Mutex
}

func NewHttpServer(iface, port, htpasswd string, graph *Graph) *HttpServer {
	h := &HttpServer{
		graph:    graph,
		iface:    iface,
		port:     port,
		htpasswd: htpasswd,
	}

	return h
}

func DefaultSecret(user, realm string) string {
	if user == "admin" {
		return "$1$dlPL2MqE$oQmn16q49SqdmhenQuNgs1"
	}
	return ""
}

func (h *HttpServer) Start(start chan tankSettings, stop chan bool) {
	var (
		htpasswd      auth.SecretProvider
		authenticator *auth.BasicAuth
	)
	if h.htpasswd != "" {
		htpasswd = auth.HtpasswdFileProvider(h.htpasswd)
		authenticator = auth.NewBasicAuthenticator("Tank auth", htpasswd)
	} else {
		authenticator = auth.NewBasicAuthenticator("Tank default auth", DefaultSecret)
		glog.Infoln("Started with default user 'admin' and password 'hello'")
	}
	mx := mux.NewRouter()

	mx.HandleFunc("/", authenticator.Wrap(func(w http.ResponseWriter, req *auth.AuthenticatedRequest) {
		glog.V(1).Infoln("Settings:", h.graph.Settings)
		h.graph.Write(w)
	}))

	mx.HandleFunc("/graph.json", authenticator.Wrap(func(w http.ResponseWriter, req *auth.AuthenticatedRequest) {
		w.Header().Set("Content-Type", "application/json")
		encoder := json.NewEncoder(w)
		if err := encoder.Encode(h.graph); err != nil {
			glog.Error("An error occurred while serving JSON endpoint: %v", err)
		}
	}))

	mx.HandleFunc("/settings", authenticator.Wrap(func(w http.ResponseWriter, req *auth.AuthenticatedRequest) {
		if req.Method != "POST" {
			http.Error(w, "Methos isn't support", http.StatusInternalServerError)
			return
		}
		req.ParseForm()
		count, err := strconv.ParseInt(req.PostFormValue("parallel"), 10, 0)
		if err != nil {
			http.Error(w, "Parallel value isn't valid: "+err.Error(), http.StatusInternalServerError)
			return
		}
		time, err := strconv.ParseInt(req.PostFormValue("time"), 10, 0)
		if err != nil {
			http.Error(w, "Time value isn't valid: "+err.Error(), http.StatusInternalServerError)
			return
		}
		timeout, err := strconv.ParseInt(req.PostFormValue("timeout"), 10, 0)
		if err != nil {
			http.Error(w, "TimeOut value isn't valid: "+err.Error(), http.StatusInternalServerError)
			return
		}
		set := tankSettings{
			req.PostFormValue("url"),
			timeout, count, time,
			req.PostFormValue("username"),
			req.PostFormValue("password"),
			req.PostFormValue("useragent"),
			req.PostFormValue("cookie"),
		}
		h.graph.setSettings(set)
		stop <- true
		start <- set
	}))

	server := http.Server{
		Handler:      mx,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	server.Serve(h.Listener())
}

func (h *HttpServer) Close() {
	h.Listener().Close()
}

func (h *HttpServer) Url() string {
	return fmt.Sprintf("http://%s/", h.Listener().Addr())
}

func (h *HttpServer) Listener() net.Listener {
	h.listenerMtx.Lock()
	defer h.listenerMtx.Unlock()

	if h.listener != nil {
		return h.listener
	}

	ifaceAndPort := fmt.Sprintf("%v:%v", h.iface, h.port)
	listener, err := net.Listen("tcp4", ifaceAndPort)
	if err != nil {
		glog.Fatalln("Error listen socket:", err)
	}

	h.listener = listener
	return h.listener
}
