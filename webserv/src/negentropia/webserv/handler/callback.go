package handler

import (
	"io"
	"os"
	"fmt"
	"log"
	//"time"
	//"io/ioutil"
	"net/http"
	
	"code.google.com/p/goauth2/oauth"
)

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
	err = nil
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
	req, err := transp.Client().Get(apiRequest)
	if err != nil {
		msg := fmt.Sprintf("handler.googleCallback url=%s Request: %s", path, err)
		log.Printf(msg)

		if err := sendLogin(w, Page{GoogleAuthMsg: msg}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
	
	io.Copy(os.Stdout, req.Body)
	io.Copy(w, req.Body)
}
