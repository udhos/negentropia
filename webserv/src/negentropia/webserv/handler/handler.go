package handler

import (
	"log"
	"strings"
	
	"github.com/bradfitz/gomemcache/memcache"
)

var (
	templateRootPath string
	GoogleId         *string
	GoogleSecret     *string
	mcServerList     []string         = []string{"127.0.0.1:11211", "127.0.0.1:12000"}
	mc               *memcache.Client
)

func init() {
	log.Printf("handler.init(): memcache client for: " + strings.Join(mcServerList, ","))
	mc = memcache.New(mcServerList...)
}

func SetTemplateRoot(path string) {
	templateRootPath = path
}

func TemplateRoot() string {
	return templateRootPath
}

func TemplatePath(path string) string {
	return TemplateRoot() + "/" + path
}
