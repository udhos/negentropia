package main

import (
	"bufio"
	"io"
	"os"
	//"fmt"
	"flag"
	"log"
	//"time"
	"errors"
	"strings"
	//"io/ioutil"
	"net/http"

	"negentropia/webserv/cfg"
	"negentropia/webserv/handler"
	"negentropia/webserv/session"
	"negentropia/webserv/store"
)

//type portList []string

var (
	staticMap    string
	templatePath string
	configFile   string
	listenAddr   string
	configFlags  *flag.FlagSet = flag.NewFlagSet("config flags", flag.ExitOnError)
	redisAddr    string
	basePath     string
)

// Initialize package main
func init() {
	handler.GoogleId = configFlags.String("gId", "", "google client id")
	handler.GoogleSecret = configFlags.String("gSecret", "", "google client secret")
	handler.FacebookId = configFlags.String("fId", "", "facebook client id")
	handler.FacebookSecret = configFlags.String("fSecret", "", "facebook client secret")
	configFlags.StringVar(&configFile, "config", "", "load config flags from this file")
	configFlags.StringVar(&listenAddr, "listenOn", ":8080", "listen address [addr]:port")
	configFlags.StringVar(&cfg.RedirectHost, "redirectHost", "localhost", "host part of redirect in proto://host:port/path")
	configFlags.StringVar(&redisAddr, "redisAddr", "localhost:6379", "redis server address")
	configFlags.StringVar(&basePath, "path", "/ne", "www base path")
	configFlags.StringVar(&templatePath, "template", "", "template root path")
	configFlags.StringVar(&staticMap, "static", "", "www static mapping")
	configFlags.StringVar(&cfg.Protocol, "proto", "http", "protocol")

	configFlags.StringVar(&cfg.SmtpAuthUser, "smtpAuthUser", "user@domain.com", "smtp auth user")
	configFlags.StringVar(&cfg.SmtpAuthPass, "smtpAuthPass", "putPasswordHere", "smtp auth password")
	configFlags.StringVar(&cfg.SmtpAuthServer, "smtpAuthServer", "smtp.gmail.com", "smtp server")
	configFlags.StringVar(&cfg.SmtpHostPort, "smtpHostPort", "smtp.gmail.com:587", "smtp server host:port")
}

/*
func flagSetInit(fs *flag.FlagSet) {
	handler.GoogleId = fs.String("gId", "", "google client id")
	handler.GoogleSecret = fs.String("gSecret", "", "google client secret")
	fs.StringVar(&configFile, "config", "", "load config flags from this file")
	fs.Var(&listenOn, "listenOn", "comma-separated list of [addr]:port pairs")
}
*/

// Wrapper type for Handler
type StaticHandler struct {
	innerHandler http.Handler // save trapped/wrapped Handler
}

func (handler StaticHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("StaticHandler.ServeHTTP url=%s", r.URL.Path)
	handler.innerHandler.ServeHTTP(w, r) // call trapped/wrapped Handler

	/*
		var delay time.Duration = 20 
		log.Printf("blocking for %d secs", delay)
		time.Sleep(delay * time.Second)
	*/
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

/*
// String is the method to get the flag value, part of the flag.Value interface.
func (pl *portList) String() string {
	return strings.Join(*pl, ",")
}

// Set is the method to set the flag value, part of the flag.Value interface.
// Set's argument is a string to be parsed to set the flag.
// It's a comma-separated list, so we split it.
func (pl *portList) Set(value string) error {
	*pl = strings.Split(value, ",") // redefine portList
	return nil
}
*/

// Add session parameter to handle
func trapHandle(w http.ResponseWriter, r *http.Request, handler func(http.ResponseWriter, *http.Request, *session.Session)) {
	s := session.Get(r)
	handler(w, r, s)
}

func loadFlagsFromFile(config string) ([]string, error) {
	log.Printf("loading config flags from file: %s", config)

	input, err := os.Open(config)
	if err != nil {
		log.Printf("failure opening flags config file: %s: %s", config, err)
		return nil, err
	}

	defer input.Close()

	var flags []string
	var num int = 0

	reader := bufio.NewReader(input)
	for line, pref, fail := reader.ReadLine(); fail != io.EOF; line, pref, fail = reader.ReadLine() {
		if fail != nil {
			log.Printf("failure reading line from flags config file: %s: %s", config, err)
			break
		}
		num++
		if pref {
			log.Printf("very long flags config line at %d", num)
			return nil, errors.New("very long flags config line")
		}
		f := strings.TrimSpace(string(line))
		if f == "" || f[:1] == "#" {
			continue
		}
		//log.Printf("flag config line [%d]: flag=[%s]", num, f)
		flags = append(flags, f)
	}

	return flags, nil
}

func loadConfig(config string) error {
	f, err := loadFlagsFromFile(config)
	if err != nil {
		log.Printf("failure reading config flags: %s", err)
		return err
	}

	err = configFlags.Parse(f)
	if err != nil {
		log.Printf("failure parsing config flags: %s", err)
		return err
	}

	//log.Printf("loaded %d flags", len(f))

	return nil
}

func getPort(hostPort string) string {
	pair := strings.Split(listenAddr, ":")
	if len(pair) == 1 {
		return ""
	}

	return ":" + pair[1]
}

func main() {
	// Parse flags from command-line
	//flag.Parse()
	configFlags.Parse(os.Args[1:])

	// Parse flags from file
	if configFile != "" {
		err := loadConfig(configFile)
		if err != nil {
			log.Printf("failure loading config flags: %s", err)
			return
		}
	}

	log.Printf("template root path: %s", templatePath)
	if templatePath == "" {
		log.Printf("template root path is required")
		return
	}
	handler.SetTemplateRoot(templatePath)

	cfg.RedirectPort = getPort(listenAddr)

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

	/*
		last := len(listenOn) - 1
		// serve ports except the last one
		for _, port := range listenOn[:last] {
			go serve(port)
		}
		serve(listenOn[last]) // serve last port
	*/
	serve(listenAddr)
}
