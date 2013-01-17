package main

import (
	"os"
	"io"
	"bufio"
	//"fmt"
	"log"
	"flag"
	"errors"
	"strings"
	//"time"
	//"io/ioutil"
	"net/http"	

	"negentropia/webserv/handler"
	"negentropia/webserv/session"
)

type portList []string

var (
	staticPath   string         = "/tmp/devel/negentropia/wwwroot"
	templatePath string         = "/tmp/devel/negentropia/template"
	configFile   string	
	listenOn     portList       = []string{":8000", ":8080"}
	configFlags  *flag.FlagSet  = flag.NewFlagSet("config flags", flag.ExitOnError)
)

// Initialize package main
func init() {
	/*
	handler.GoogleId = flag.String("gId", "", "google client id")
	handler.GoogleSecret = flag.String("gSecret", "", "google client secret")
	flag.StringVar(&configFile, "config", "", "load config flags from this file")
	flag.Var(&listenOn, "listenOn", "comma-separated list of [addr]:port pairs")
	*/

	handler.GoogleId = configFlags.String("gId", "", "google client id")
	handler.GoogleSecret = configFlags.String("gSecret", "", "google client secret")
	configFlags.StringVar(&configFile, "config", "", "load config flags from this file")
	configFlags.Var(&listenOn, "listenOn", "comma-separated list of [addr]:port pairs")
	configFlags.Parse(os.Args[1:])
	
	handler.SetTemplateRoot(templatePath)
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
	path := r.URL.Path
	log.Printf("StaticHandler.ServeHTTP url=%s", path)
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
	if err != nil {
		log.Panicf("ListenAndServe: %s: %s", addr, err)
	}
}

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
			break;
		}
		num++
		if pref {
			log.Printf("very long flags config line at %d", num)
			return nil, errors.New("very long flags config line")
		}
		if len(line) == 0 || line[0] == '#' {
			continue
		}
		//log.Printf("line [%d]: %s", num, line)
		flags = append(flags, string(line))
	}
	
	return flags, nil
}

func loadConfig() error {
	f, err := loadFlagsFromFile(configFile)
	if err != nil {
		log.Printf("failure reading config flags: %s", err);		
		return err
	}
	
	err = configFlags.Parse(f)
	if err != nil {
		log.Printf("failure parsing config flags: %s", err);	
		return err
	}
	
	//log.Printf("loaded %d flags", len(f))
	
	return nil
}

func main() {
	//flag.Parse()
	
	err := loadConfig()
	if err != nil {
		log.Printf("failure loading config flags: %s", err);
		return
	}
	
	if *handler.GoogleId == "" {
		log.Printf("warning: google client id is UNDEFINED: google login won't be available")
	}
	
	if *handler.GoogleSecret == "" {
		log.Printf("warning: google client secret is UNDEFINED: google login won't be available")
	}
	
	http.Handle("/", StaticHandler{http.FileServer(http.Dir(staticPath))})
	http.HandleFunc("/n/",               func (w http.ResponseWriter, r *http.Request) { trapHandle(w, r, handler.Home) } )
	http.HandleFunc("/n/logout",         func (w http.ResponseWriter, r *http.Request) { trapHandle(w, r, handler.Logout) } )
	http.HandleFunc("/n/login",          func (w http.ResponseWriter, r *http.Request) { trapHandle(w, r, handler.Login) } )
	http.HandleFunc("/n/loginAuth",      func (w http.ResponseWriter, r *http.Request) { trapHandle(w, r, handler.LoginAuth) } )
	http.HandleFunc("/n/googleCallback", func (w http.ResponseWriter, r *http.Request) { trapHandle(w, r, handler.GoogleCallback) } )	

	last := len(listenOn) - 1
	// serve ports except the last one
	for _, port := range listenOn[:last] {
		go serve(port)
	}
	serve(listenOn[last]) // serve last port
}
