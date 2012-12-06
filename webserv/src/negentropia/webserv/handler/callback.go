package handler

import (
	//"io"
	//"os"
	"fmt"
	"log"
	//"errors"
	//"time"
	"io/ioutil"
	"net/http"
	"encoding/json"
	
	"code.google.com/p/goauth2/oauth"
	"github.com/bradfitz/gomemcache/memcache"
)

type GoogleProfile struct {
	Id   string
	Name string
}

func GoogleCallback(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	code := r.FormValue("code");
	
	/*
	str := "handler.googleCallback url=" + path + " code=" + code
	log.Printf(str)
	fmt.Fprintf(w, str)
	*/

	// See previous step under loginAuth handler

	config := googleOauth2Config()

	// Set up a Transport with our config, define the cache
	transp := &oauth.Transport{Config: config}
	//tokenCache = oauth.CacheFile(*cachefile)
	
	// Step two, exchange the authorization code for an access token.
	
	tok, err := transp.Exchange(code)
	if err != nil {
		msg := fmt.Sprintf("handler.googleCallback url=%s Exchange: %s", path, err)
		log.Printf(msg)

		if err := sendLogin(w, Page{GoogleAuthMsg: msg}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}

	// FIXME: Load cached token, if available.
	
	transp.Token = &oauth.Token{AccessToken: tok.AccessToken}

	// FIXME: Tack on the extra parameters, if specified.
	apiRequest := "https://www.googleapis.com/oauth2/v1/userinfo"
	/*
	if *authparam != "" {
		*apiRequest += *authparam + ctoken.AccessToken
	}
	*/

	// Send sequest
	resp, err := transp.Client().Get(apiRequest)
	if err != nil {
		msg := fmt.Sprintf("handler.googleCallback url=%s Request: %s", path, err)
		log.Printf(msg)

		if err := sendLogin(w, Page{GoogleAuthMsg: msg}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
	
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		msg := fmt.Sprintf("handler.googleCallback url=%s ioutil.ReadAll: %s", path, err)
		log.Printf(msg)

		if err := sendLogin(w, Page{GoogleAuthMsg: msg}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
	
	var profile GoogleProfile
	
	err = json.Unmarshal(body, &profile)
	if err != nil {
		msg := fmt.Sprintf("handler.googleCallback url=%s Unmarshal: %s", path, err)
		log.Printf(msg)

		if err := sendLogin(w, Page{GoogleAuthMsg: msg}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
	
	log.Printf("handler.googleCallback url=%s name=%s id=%s", path, profile.Name, profile.Id)
	
    err = mc.Set(&memcache.Item{Key: "session", Value: []byte("123456"), Expiration: 24*3600})
	if err != nil {
		log.Printf("handler.googleCallback mc.Set url=%s err=%s", path, err)
	}
	var it *memcache.Item
    it, err = mc.Get("session")
	if err != nil {
		log.Printf("handler.googleCallback mc.Get url=%s err=%s", path, err)
	} else {
		log.Printf("handler.googleCallback url=%s session=%s", path, it.Value)
	}
	
	http.Redirect(w, r, "/n/", http.StatusFound)
}
