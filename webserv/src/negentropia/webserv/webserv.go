package main

import (
	"os"
	//"fmt"
	"flag"
	"log"
	//"time"
	"strings"
	//"io/ioutil"
	"net/http"

	"negentropia/webserv/cfg"
	"negentropia/webserv/handler"
	"negentropia/webserv/session"
	"negentropia/webserv/store"
	"negentropia/webserv/configflag"
	"negentropia/webserv/util"
)

//type portList []string

var (
	configFlags  *flag.FlagSet = flag.NewFlagSet("config flags", flag.ExitOnError)
	configList   configflag.FileList
	staticMap    string
	templatePath string
	listenAddr   string
	redisAddr    string
	basePath     string
)

// Initialize package main
func init() {
	handler.GoogleId = configFlags.String("gId", "", "google client id")
	handler.GoogleSecret = configFlags.String("gSecret", "", "google client secret")
	handler.FacebookId = configFlags.String("fId", "", "facebook client id")
	handler.FacebookSecret = configFlags.String("fSecret", "", "facebook client secret")
	configFlags.Var(&configList, "config", "load config flags from this file")
	configFlags.StringVar(&listenAddr, "listenOn", ":8080", "http listen address [addr]:port")
	configFlags.StringVar(&cfg.RedirectHost, "redirectHost", "localhost", "host part of redirect in proto://host:port/path")
	configFlags.StringVar(&redisAddr, "redisAddr", "localhost:6379", "redis server address")
	configFlags.StringVar(&basePath, "path", "/ne", "www base path")
	configFlags.StringVar(&templatePath, "template", "", "template root path")
	configFlags.StringVar(&staticMap, "static", "", "www static mapping")
	configFlags.StringVar(&cfg.Protocol, "proto", "http", "protocol")

	configFlags.StringVar(&cfg.SmtpAuthUser, "smtpAuthUser", "user@domain.com", "smtp auth user")
	configFlags.StringVar(&cfg.SmtpAuthPass, "smtpAuthPass", "", "smtp auth password")
	configFlags.StringVar(&cfg.SmtpAuthServer, "smtpAuthServer", "smtp.gmail.com", "smtp server")
	configFlags.StringVar(&cfg.SmtpHostPort, "smtpHostPort", "smtp.gmail.com:587", "smtp server host:port")
}

// Wrapper type for Handler
type StaticHandler struct {
	innerHandler http.Handler // save trapped/wrapped Handler
}

func (handler StaticHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("StaticHandler.ServeHTTP url=%s", r.URL.Path)
	handler.innerHandler.ServeHTTP(w, r) // call trapped/wrapped Handler
}

func serve(addr string) {
	if addr == "" {
		log.Printf("server starting on :http (empty address)")
	} else {
		log.Printf("server starting on " + addr)
	}

	err := http.ListenAndServe(addr, nil)
	/*
		s := &http.Server{
			Addr:           addr,
			Handler:        nil,
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Second,
			MaxHeaderBytes: 1 << 20,
		}
		err := s.ListenAndServe()
	*/
	if err != nil {
		log.Panicf("ListenAndServe: %s: %s", addr, err)
	}
}

// Add session parameter to handle
func trapHandle(w http.ResponseWriter, r *http.Request, handler func(http.ResponseWriter, *http.Request, *session.Session)) {
	s := session.Get(r)
	handler(w, r, s)
}

func main() {
	log.Printf("webserv booting")

	// Parse flags from command-line
	configFlags.Parse(os.Args[1:])

	// Parse flags from files
	log.Printf("config files: %d", len(configList))
	if err := configflag.Load(configFlags, configList); err != nil {
		log.Printf("failure loading config flags: %s", err)
		return
	}

	log.Printf("template root path: %s", templatePath)
	if templatePath == "" {
		log.Printf("template root path is required")
		return
	}
	handler.SetTemplateRoot(templatePath)

	cfg.RedirectPort = util.GetPort(listenAddr)

	if *handler.GoogleId == "" {
		log.Printf("warning: google client id is UNDEFINED: google login won't be available")
	}
	if *handler.GoogleSecret == "" {
		log.Printf("warning: google client secret is UNDEFINED: google login won't be available")
	}
	if *handler.FacebookId == "" {
		log.Printf("warning: facebook client id is UNDEFINED: facebook login won't be available")
	}
	if *handler.FacebookSecret == "" {
		log.Printf("warning: facebook client secret is UNDEFINED: facebook login won't be available")
	}
	if cfg.SmtpAuthPass == "" {
		log.Printf("warning: smtp auth password is UNDEFINED: automatic email validation for local accounts won't be available")
	}

	store.Init(redisAddr)

	cfg.SetBasePath(basePath)

	log.Printf("home URL: %s", cfg.HomeURL())

	staticMap := strings.TrimSpace(staticMap)
	if staticMap == "" {
		log.Printf("no static map defined")
	} else {
		for _, pair := range strings.Split(staticMap, ",") {
			pDir := strings.Split(pair, ":")
			if len(pDir) != 2 {
				log.Printf("bad static map pair: pair=[%s] map=[%s]", pair, staticMap)
				return
			}
			p := pDir[0]
			dir := pDir[1]
			log.Printf("installing static handler from path %s to directory %s", p, dir)
			http.Handle(p, StaticHandler{http.StripPrefix(p, http.FileServer(http.Dir(dir)))})
		}
	}

	http.HandleFunc(cfg.HomePath(), func(w http.ResponseWriter, r *http.Request) { trapHandle(w, r, handler.Home) })
	http.HandleFunc(cfg.HomeDartPath(), func(w http.ResponseWriter, r *http.Request) { trapHandle(w, r, handler.HomeDart) })	
	http.HandleFunc(cfg.LogoutPath(), func(w http.ResponseWriter, r *http.Request) { trapHandle(w, r, handler.Logout) })
	http.HandleFunc(cfg.LoginPath(), func(w http.ResponseWriter, r *http.Request) { trapHandle(w, r, handler.Login) })
	http.HandleFunc(cfg.LoginAuthPath(), func(w http.ResponseWriter, r *http.Request) { trapHandle(w, r, handler.LoginAuth) })
	http.HandleFunc(cfg.GoogleCallbackPath(), func(w http.ResponseWriter, r *http.Request) { trapHandle(w, r, handler.GoogleCallback) })
	http.HandleFunc(cfg.FacebookCallbackPath(), func(w http.ResponseWriter, r *http.Request) { trapHandle(w, r, handler.FacebookCallback) })
	http.HandleFunc(cfg.SignupPath(), func(w http.ResponseWriter, r *http.Request) { trapHandle(w, r, handler.Signup) })
	http.HandleFunc(cfg.SignupProcessPath(), func(w http.ResponseWriter, r *http.Request) { trapHandle(w, r, handler.SignupProcess) })
	http.HandleFunc(cfg.ConfirmPath(), func(w http.ResponseWriter, r *http.Request) { trapHandle(w, r, handler.Confirm) })
	http.HandleFunc(cfg.ConfirmProcessPath(), func(w http.ResponseWriter, r *http.Request) { trapHandle(w, r, handler.ConfirmProcess) })
	http.HandleFunc(cfg.ResetPassPath(), func(w http.ResponseWriter, r *http.Request) { trapHandle(w, r, handler.ResetPass) })
	http.HandleFunc(cfg.ResetPassProcessPath(), func(w http.ResponseWriter, r *http.Request) { trapHandle(w, r, handler.ResetPassProcess) })
	http.HandleFunc(cfg.ResetPassConfirmPath(), func(w http.ResponseWriter, r *http.Request) { trapHandle(w, r, handler.ResetPassConfirm) })
	http.HandleFunc(cfg.ResetPassConfirmProcessPath(), func(w http.ResponseWriter, r *http.Request) { trapHandle(w, r, handler.ResetPassConfirmProcess) })

	log.Printf("webserv boot complete")
	
	serve(listenAddr)
}
