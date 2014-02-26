package handler

import (
	"log"
	"path"
	"strings"
	"time"

	"negentropia/webserv/session"
)

var (
	templateRootPath string
	GoogleId         *string
	GoogleSecret     *string
	FacebookId       *string
	FacebookSecret   *string
)

func Init(templateRoot string) {

	log.Printf("handler.Init: unconfirmedExpire = %d seconds", unconfirmedExpire/time.Second)

	setTemplateRoot(templateRoot)
}

func setTemplateRoot(p string) {
	templateRootPath = p
}

func TemplateRoot() string {
	return templateRootPath
}

func TemplatePath(p string) string {
	return path.Join(TemplateRoot(), p)
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
