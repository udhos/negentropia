package cfg

// Exposes constant variables to multiple readers (read-only goroutines)

import "log"

var (
	homePath				string
	logoutPath				string
	loginPath				string
	loginAuthPath			string
	googleCallbackPath		string
	facebookCallbackPath	string
)

func SetBasePath(basePath string) {
	log.Printf("cfg.SetBasePath: www base path: %s", basePath)
	
	homePath             = basePath + "/"
	logoutPath           = basePath + "/logout"
	loginPath            = basePath + "/login"
	loginAuthPath        = basePath + "/loginAuth"
	googleCallbackPath   = basePath + "/googleCallback"
	facebookCallbackPath = basePath + "/facebookCallback"
}

func HomePath()             string { return homePath }
func LogoutPath()           string { return logoutPath }
func LoginPath()            string { return loginPath }
func LoginAuthPath()        string { return loginAuthPath }
func GoogleCallbackPath()   string { return googleCallbackPath }
func FacebookCallbackPath() string { return facebookCallbackPath }
