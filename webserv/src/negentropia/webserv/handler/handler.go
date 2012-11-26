package handler

import (
	//"log"
)

var (
	templateRootPath string
)

func SetRootPath(path string) {
	templateRootPath = path
}

func RootPath() string {
	return templateRootPath
}

func TemplatePath(path string) string {
	return RootPath() + "\\" + path
}
