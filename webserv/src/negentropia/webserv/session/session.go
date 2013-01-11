package session

import (
	"log"
	"time"
	//"errors"
	"strings"
	"net/http"
	
	"github.com/bradfitz/gomemcache/memcache"
)

var (
	mcServerList     []string         = []string{"127.0.0.1:11211", "127.0.0.1:12000"}
	mc               *memcache.Client
)

type Session struct {
	Id string
}

func init() {
	log.Printf("session.init(): memcache client for: " + strings.Join(mcServerList, ","))
	mc = memcache.New(mcServerList...)
}

func sessionLookup() string {
	//var it *memcache.Item
    it, err := mc.Get("session")
	if err != nil {
		log.Printf("handler.googleCallback mc.Get err=%s", err)
		return ""
	}
	return string(it.Value)
}

func sessionSave(value string) error {
    err := mc.Set(&memcache.Item{Key: "session", Value: []byte(value), Expiration: 24*3600})
	if err != nil {
		log.Printf("handler.googleCallback mc.Set err=%s", err)
	}
	return err
}

func newCookie(name, value string) *http.Cookie {
	var maxAge int = 86400
	var expires time.Time

	if maxAge > 0 {
		expires = time.Now().Add(time.Duration(maxAge) * time.Second)
	} else if maxAge < 0 {
		expires = time.Unix(1, 0)
	}

	cookie := &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		Domain:   "",
		MaxAge:   maxAge,
		Expires:  expires,
		Secure:   false,
		HttpOnly: true,
	}
	
	return cookie
}

func newSession(id string) *Session {
	return &Session{id}
}

func Get(w http.ResponseWriter, r *http.Request) *Session {

	cook, err := r.Cookie("session")
	if err == nil {
		log.Printf("session.Get FOUND cookie:session=%s", cook.Value)
		
		sId := sessionLookup();

		log.Printf("session.Get FOUND DB:session=%s", sId)
		
		return newSession(cook.Value)
	}
	
	// Create cookie
	
	sessionId := "test123456"

	ck := newCookie("session", sessionId)
	
	http.SetCookie(w, ck)
	
	sessionSave(sessionId)
	
	return newSession(ck.Value)
}
