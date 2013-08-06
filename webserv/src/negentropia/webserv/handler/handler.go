package handler

import (
	//"log"
	"path"
	"strings"

	"negentropia/webserv/session"
)

var (
	templateRootPath string
	GoogleId         *string
	GoogleSecret     *string
	FacebookId       *string
	FacebookSecret   *string
)

func SetTemplateRoot(p string) {
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
