package cfg

// Exposes constant variables to multiple readers (read-only goroutines)

import "log"

var (
	homePath				string
	homeJSPath	    		string
	logoutPath				string
	loginPath				string
	loginAuthPath			string
	googleCallbackPath		string
	facebookCallbackPath	string
	signupPath				string
	signupProcessPath		string
	confirmPath				string
	confirmProcessPath		string
	resetPassPath			string
	resetPassProcessPath	string
	resetPassConfirmPath	string	
	resetPassConfirmProcessPath string
	
	Protocol		string
	RedirectHost	string
	RedirectPort	string
		
	SmtpAuthUser	string
	SmtpAuthPass	string
	SmtpAuthServer	string
	SmtpHostPort	string
)

func SetBasePath(basePath string) {
	log.Printf("cfg.SetBasePath: www base path: %s", basePath)
	
	homePath             = basePath + "/"
	homeJSPath           = basePath + "/homeJS"
	logoutPath           = basePath + "/logout"
	loginPath            = basePath + "/login"
	loginAuthPath        = basePath + "/loginAuth"
	googleCallbackPath   = basePath + "/googleCallback"
	facebookCallbackPath = basePath + "/facebookCallback"
	signupPath           = basePath + "/signup"
	signupProcessPath    = basePath + "/signupProcess"
	confirmPath          = basePath + "/confirm"
	confirmProcessPath   = basePath + "/confirmProcess"
	resetPassPath        = basePath + "/resetPass"
	resetPassProcessPath = basePath + "/resetPassProcess"
	resetPassConfirmPath = basePath + "/resetPassConfirm"
	resetPassConfirmProcessPath = basePath + "/resetPassConfirmProcess"
}

func HomePath()             string { return homePath }
func HomeJSPath()           string { return homeJSPath }
func LogoutPath()           string { return logoutPath }
func LoginPath()            string { return loginPath }
func LoginAuthPath()        string { return loginAuthPath }
func GoogleCallbackPath()   string { return googleCallbackPath }
func FacebookCallbackPath() string { return facebookCallbackPath }
func SignupPath()           string { return signupPath }
func SignupProcessPath()    string { return signupProcessPath }
func ConfirmPath()          string { return confirmPath }
func ConfirmProcessPath()   string { return confirmProcessPath }
func ResetPassPath()        string { return resetPassPath }
func ResetPassProcessPath() string { return resetPassProcessPath }
func ResetPassConfirmPath() string { return resetPassConfirmPath }
func ResetPassConfirmProcessPath() string { return resetPassConfirmProcessPath }

// Full URLs to embed in email
func HomeURL()			   string { return Protocol + "://" + RedirectHost + RedirectPort + homePath }
func ConfirmURL()          string { return Protocol + "://" + RedirectHost + RedirectPort + confirmPath }
func ConfirmProcessURL()   string { return Protocol + "://" + RedirectHost + RedirectPort + confirmProcessPath }
func ResetPassConfirmURL() string { return Protocol + "://" + RedirectHost + RedirectPort + resetPassConfirmPath }
func GoogleCallbackURL()   string { return Protocol + "://" + RedirectHost + RedirectPort + googleCallbackPath }
func FacebookCallbackURL() string { return Protocol + "://" + RedirectHost + RedirectPort + facebookCallbackPath }
