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
	signupPath				string
	signupProcessPath		string
	confirmPath				string
	confirmProcessPath		string	
)

func SetBasePath(basePath string) {
	log.Printf("cfg.SetBasePath: www base path: %s", basePath)
	
	homePath             = basePath + "/"
	logoutPath           = basePath + "/logout"
	loginPath            = basePath + "/login"
	loginAuthPath        = basePath + "/loginAuth"
	googleCallbackPath   = basePath + "/googleCallback"
	facebookCallbackPath = basePath + "/facebookCallback"
	signupPath           = basePath + "/signup"
	signupProcessPath    = basePath + "/signupProcess"
	confirmPath          = basePath + "/confirm"
	confirmProcessPath   = basePath + "/confirmProcess"
}

func HomePath()             string { return homePath }
func LogoutPath()           string { return logoutPath }
func LoginPath()            string { return loginPath }
func LoginAuthPath()        string { return loginAuthPath }
func GoogleCallbackPath()   string { return googleCallbackPath }
func FacebookCallbackPath() string { return facebookCallbackPath }
func SignupPath()           string { return signupPath }
func SignupProcessPath()    string { return signupProcessPath }
func ConfirmPath()          string { return confirmPath }
func ConfirmProcessPath()   string { return confirmProcessPath }