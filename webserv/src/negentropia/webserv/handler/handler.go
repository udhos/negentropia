package handler

import (
	//"log"
)

var (
	templateRootPath string
	GoogleId         *string
	GoogleSecret     *string
)

func SetTemplateRoot(path string) {
	templateRootPath = path
}

func TemplateRoot() string {
	return templateRootPath
}

func TemplatePath(path string) string {
	return TemplateRoot() + "/" + path
}
