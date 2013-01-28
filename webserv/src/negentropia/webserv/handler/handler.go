package handler

import (
	//"log"
	"strings"
	
	//"github.com/bradfitz/gomemcache/memcache"
	
	"negentropia/webserv/session"
)

var (
	templateRootPath string
	GoogleId       *string
	GoogleSecret   *string
	FacebookId     *string
	FacebookSecret *string
	RedirectHost    string
	RedirectPort    string
)

/*
func init() {
	log.Printf("handler.init(): memcache client for: " + strings.Join(mcServerList, ","))
	mc = memcache.New(mcServerList...)
}
*/

func SetTemplateRoot(path string) {
	templateRootPath = path
}

func TemplateRoot() string {
	return templateRootPath
}

func TemplatePath(path string) string {
	return TemplateRoot() + "/" + path
}

func accountLabel(s *session.Session) string {
	if s == nil {
		return ""
	}
	
	return s.ProfileName + " (" + s.ProfileEmail + ")"
}

func formatEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}


