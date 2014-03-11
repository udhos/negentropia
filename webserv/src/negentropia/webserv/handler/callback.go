package handler

import (
	//"io"
	//"os"
	"fmt"
	"log"
	//"errors"
	//"time"
	//"strings"
	"encoding/json"
	"io/ioutil"
	"net/http"

	//"code.google.com/p/goauth2/oauth" // facebook broken
	//"github.com/robfig/goauth2/oauth" // google broken
	//"code.google.com/r/jasonmcvetta-goauth2" // go get broken
	"github.com/HairyMezican/goauth2/oauth"

	"negentropia/webserv/cfg"
	"negentropia/webserv/session"
)

type GoogleProfile struct {
	Id    string
	Name  string
	Email string
}

type FacebookProfile struct {
	Id    string
	Name  string
	Email string
}

func GoogleCallback(w http.ResponseWriter, r *http.Request, s *session.Session) {
	path := r.URL.Path

	account := accountLabel(s)

	code := r.FormValue("code")

	/*
		str := "handler.googleCallback url=" + path + " code=" + code
		log.Printf(str)
		fmt.Fprintf(w, str)
	*/

	// See previous step under loginAuth handler

	config := googleOauth2Config()

	// Set up a Transport with our config, define the cache
	transp := &oauth.Transport{Config: config}

	var tok *oauth.Token
	if *GoogleTokenCacheFile != "" {
		// Try to pull the token from the cache; if this fails, we need to get one.
		var err error
		tok, err = config.TokenCache.Token()
		if err == nil {
			log.Printf("handler.googleCallback: found token from cache")
		}
	}

	if tok == nil {
		// Exchange the authorization code for an access token.
		// ("Here's the code you gave the user, now give me a token!")
		log.Printf("handler.googleCallback: requesting token")
		var err error
		tok, err = transp.Exchange(code)
		if err != nil {
			msg := fmt.Sprintf("handler.googleCallback url=%s Exchange: %s", path, err)
			log.Printf(msg)

			if err := sendLogin(w, Page{Account: account, ShowNavAccount: true, ShowNavHome: true, GoogleAuthMsg: msg}); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			return
		}
		if *GoogleTokenCacheFile != "" {
			// (The Exchange method will automatically cache the token.)
			log.Printf("Google token is cached in %v\n", config.TokenCache)
		}
	}

	// Make the actual request using the cached token to authenticate.
	// ("Here's the token, let me in!")
	transp.Token = tok

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

		if err := sendLogin(w, Page{Account: account, ShowNavAccount: true, ShowNavHome: true, GoogleAuthMsg: msg}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		msg := fmt.Sprintf("handler.googleCallback url=%s ioutil.ReadAll: %s", path, err)
		log.Printf(msg)

		if err := sendLogin(w, Page{Account: account, ShowNavAccount: true, ShowNavHome: true, GoogleAuthMsg: msg}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		return
	}

	//log.Printf("handler.googleCallback url=%s body=%s", path, body)

	var profile GoogleProfile

	err = json.Unmarshal(body, &profile)
	if err != nil {
		msg := fmt.Sprintf("handler.googleCallback url=%s Unmarshal: %s", path, err)
		log.Printf(msg)

		if err := sendLogin(w, Page{Account: account, ShowNavAccount: true, ShowNavHome: true, GoogleAuthMsg: msg}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		return
	}

	profile.Email = formatEmail(profile.Email)

	log.Printf("handler.googleCallback url=%s name=%s id=%s email=%s", path, profile.Name, profile.Id, profile.Email)

	// required non-empty email
	if profile.Email == "" {
		if err := sendLogin(w, Page{Account: account, ShowNavAccount: true, ShowNavHome: true, GoogleAuthMsg: "Google email is required"}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if s != nil {
		session.Delete(w, s)
		s = nil
	}
	s = session.Set(w, session.AUTH_PROV_GOOGLE, profile.Id, profile.Name, profile.Email)
	if s == nil {
		log.Printf("handler.googleCallback url=%s could not establish session", path)
		http.Error(w, "handler.googleCallback could not establish session", http.StatusInternalServerError)
		return
	}

	log.Printf("handler.googleCallback url=%s session=%s DONE", path, s.SessionId)

	http.Redirect(w, r, cfg.HomePath(), http.StatusFound)
}

func FacebookCallback(w http.ResponseWriter, r *http.Request, s *session.Session) {
	path := r.URL.Path

	account := accountLabel(s)

	code := r.FormValue("code")

	// See previous step under loginAuth handler

	config := facebookOauth2Config()

	// Set up a Transport with our config, define the cache
	transp := &oauth.Transport{Config: config}

	var tok *oauth.Token
	if *FacebookTokenCacheFile != "" {
		// Try to pull the token from the cache; if this fails, we need to get one.
		var err error
		tok, err = config.TokenCache.Token()
		if err == nil {
			log.Printf("handler.facebookCallback: found token from cache")
		}
	}

	if tok == nil {
		// Exchange the authorization code for an access token.
		// ("Here's the code you gave the user, now give me a token!")
		log.Printf("handler.facebookCallback: requesting token")
		var err error
		tok, err = transp.Exchange(code)
		if err != nil {
			msg := fmt.Sprintf("handler.facebookCallback url=%s Exchange: %s", path, err)
			log.Printf(msg)

			if err := sendLogin(w, Page{Account: account, ShowNavAccount: true, ShowNavHome: true, FacebookAuthMsg: msg}); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			return
		}
		if *FacebookTokenCacheFile != "" {
			// (The Exchange method will automatically cache the token.)
			log.Printf("Facebook token is cached in %v\n", config.TokenCache)
		}
	}

	// Make the actual request using the cached token to authenticate.
	// ("Here's the token, let me in!")
	transp.Token = tok

	// FIXME: Tack on the extra parameters, if specified.
	//apiRequest := "https://graph.facebook.com/me?fields=name,email"
	apiRequest := "https://graph.facebook.com/me"
	/*
		if *authparam != "" {
			*apiRequest += *authparam + ctoken.AccessToken
		}
	*/

	// Send request
	resp, err := transp.Client().Get(apiRequest)
	if err != nil {
		msg := fmt.Sprintf("handler.facebookCallback url=%s Request: %s", path, err)
		log.Printf(msg)

		if err := sendLogin(w, Page{Account: account, ShowNavAccount: true, ShowNavHome: true, FacebookAuthMsg: msg}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		msg := fmt.Sprintf("handler.facebookCallback url=%s ioutil.ReadAll: %s", path, err)
		log.Printf(msg)

		if err := sendLogin(w, Page{Account: account, ShowNavAccount: true, ShowNavHome: true, FacebookAuthMsg: msg}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		return
	}

	log.Printf("handler.facebookCallback url=%s body=%s", path, body)

	var profile FacebookProfile

	err = json.Unmarshal(body, &profile)
	if err != nil {
		msg := fmt.Sprintf("handler.facebookCallback url=%s Unmarshal: %s", path, err)
		log.Printf(msg)

		if err := sendLogin(w, Page{Account: account, ShowNavAccount: true, ShowNavHome: true, FacebookAuthMsg: msg}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		return
	}

	profile.Email = formatEmail(profile.Email)

	log.Printf("handler.facebookCallback url=%s name=%s id=%s email=%s", path, profile.Name, profile.Id, profile.Email)

	// required non-empty email
	if profile.Email == "" {
		if err := sendLogin(w, Page{Account: account, ShowNavAccount: true, ShowNavHome: true, FacebookAuthMsg: "Facebook email is required"}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if s != nil {
		session.Delete(w, s)
		s = nil
	}
	s = session.Set(w, session.AUTH_PROV_FACEBOOK, profile.Id, profile.Name, profile.Email)
	if s == nil {
		log.Printf("handler.facebookCallback url=%s could not establish session", path)
		http.Error(w, "handler.facebookCallback could not establish session", http.StatusInternalServerError)
		return
	}

	log.Printf("handler.facebookCallback url=%s session=%s DONE", path, s.SessionId)

	http.Redirect(w, r, cfg.HomePath(), http.StatusFound)
}
