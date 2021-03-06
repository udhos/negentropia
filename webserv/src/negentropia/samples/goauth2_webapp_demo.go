/*
	goauth2_webapp_demo.go

	Copyright (c) 2014 Everton da Silva Marques

	Redistribution and use in source and binary forms, with or without
	modification, are permitted provided that the following conditions are
	met:

	   * Redistributions of source code must retain the above copyright
	notice, this list of conditions and the following disclaimer.
	   * Redistributions in binary form must reproduce the above
	copyright notice, this list of conditions and the following disclaimer
	in the documentation and/or other materials provided with the
	distribution.
	   * Neither the name of its authors nor the names of any of its
	contributors may be used to endorse or promote products derived from
	this software without specific prior written permission.

	THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
	"AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
	LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
	A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
	OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
	SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
	LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
	DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
	THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
	(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
	OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

	HOW TO RUN:

	1. Install GoLang -- http://golang.org/doc/install
	2. Install Mercurial -- http://mercurial.selenic.com
	3. Set the GOPATH environment variable -- http://golang.org/doc/code.html#GOPATH
	4. Install goauth2 -- go get code.google.com/p/goauth2/oauth
	5. Register the callback URL http://localhost:8080/callback in the oauth2 provider's authorized list
	   For Google: https://console.developers.google.com
	   For Facebook: https://developers.facebook.com/apps
	6. Run the demo -- go run goauth2_webapp_demo.go
	7. Point the browser to: http://localhost:8080
	    FOR GOOGLE:
	7.1 In the browser, click the "Preload Google Data" button
	7.2 In the browser, manually enter both Google ClientId and Google ClientSecret
	7.3 In the browser, click the "Run Oauth2" button
	    FOR FACEBOOK:
	7.1 In the browser, click the "Preload Facebook Data" button
	7.2 In the browser, manually enter both Facebook ClientId and Facebook ClientSecret
	7.3 In the browser, click the "Run Oauth2" button
*/

package main

import (
	"fmt"
	"html/template"
	//"io"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	//"github.com/HairyMezican/goauth2/oauth"
	"code.google.com/p/goauth2/oauth"
)

const baseTemplate = `
<!DOCTYPE html>
<html lang="en">
    <head>
        <title>{{ template "title" . }}</title>
    </head>
    <body>
		{{if .ShowNavHome}}<a href="{{.HomePath}}">home</a>{{end}}
		
        {{ template "content" . }}
				
		{{ template "script" . }}
    </body>
</html>
`

const homeTemplate = `
{{ define "title" }}goauth2_webapp_demo.go{{ end }}
{{ define "script" }}
<script>
	function clear_form() {
		document.getElementById("ClientId").value = "";
		document.getElementById("ClientSecret").value = "";
		document.getElementById("Scope").value = "";
		document.getElementById("AuthURL").value = "";
		document.getElementById("TokenURL").value = "";
		document.getElementById("RedirectURL").value = "{{.CallbackURL}}";
		document.getElementById("ApiRequest").value = "";
	}
	
	function preload_google() {
		document.getElementById("Scope").value = "https://www.googleapis.com/auth/userinfo.profile https://www.googleapis.com/auth/userinfo.email";
		document.getElementById("AuthURL").value = "https://accounts.google.com/o/oauth2/auth";
		document.getElementById("TokenURL").value = "https://accounts.google.com/o/oauth2/token";
		document.getElementById("ApiRequest").value = "https://www.googleapis.com/oauth2/v1/userinfo"; 
	}

	function preload_facebook() {
		document.getElementById("Scope").value = "email";
		document.getElementById("AuthURL").value = "https://www.facebook.com/dialog/oauth";
		document.getElementById("TokenURL").value = "https://graph.facebook.com/oauth/access_token";
		document.getElementById("ApiRequest").value = "https://graph.facebook.com/me";
	}
	
</script>
{{ end }}
{{ define "content" }}

<h1>{{.Header}}</h1>

<form name="login" action="/login" method="POST">

<div>ClientId: <input type="text" id="ClientId" name="ClientId" size="80" value="{{.ClientIdValue}}">{{.ClientIdMsg}}</div>
<div>ClientSecret: <input type="text" id="ClientSecret" name="ClientSecret" size="80" value="{{.ClientSecretValue}}">{{.ClientSecretMsg}}</div>
<div>Scope: <input type="text" id="Scope" name="Scope" size="80" value="{{.ScopeValue}}">{{.ScopeMsg}}</div>
<div>AuthURL: <input type="text" id="AuthURL" name="AuthURL" size="80" value="{{.AuthURLValue}}">{{.AuthURLMsg}}</div>
<div>TokenURL: <input type="text" id="TokenURL" name="TokenURL" size="80" value="{{.TokenURLValue}}">{{.TokenURLMsg}}</div>
<div>RedirectURL: <input type="text" id="RedirectURL" name="RedirectURL" size="80" value="{{.RedirectURLValue}}">{{.RedirectURLMsg}}</div>
<div>ApiRequest: <input type="text" id="ApiRequest" name="ApiRequest" size="80" value="{{.ApiRequestValue}}">{{.ApiRequestMsg}}</div>

<div><input type="submit" name="Oauth2Login" value="Run Oauth2"><span class="failmsg">{{.AuthMsg}}</span></div>
<div><input type="reset" name="ResetButton" value="Reset Form"></div>

</form>

<div><button onclick="clear_form()">Clear Form</button></div>
<div><button onclick="preload_google()">Preload Google Data</button></div>
<div><button onclick="preload_facebook()">Preload Facebook Data</button></div>

{{ end }}
`

const callbackTemplate = `
{{ define "title" }}goauth2_webapp_demo.go{{ end }}
{{ define "script" }}
{{ end }}
{{ define "content" }}

<h1>{{.Header}}</h1>

<div>{{.AuthMsg}}</div>

{{ end }}
`

type Page struct {
	Header      string
	HomePath    string
	ShowNavHome bool
	CallbackURL string
	AuthMsg     string

	ClientIdValue     string
	ClientSecretValue string
	ScopeValue        string
	AuthURLValue      string
	TokenURLValue     string
	RedirectURLValue  string
	ApiRequestValue   string

	ClientIdMsg     string
	ClientSecretMsg string
	ScopeMsg        string
	AuthURLMsg      string
	TokenURLMsg     string
	RedirectURLMsg  string
	ApiRequestMsg   string
}

func sendHome(w http.ResponseWriter, p Page) error {
	var err error
	t := template.New("home")
	if t, err = t.Parse(baseTemplate); err != nil {
		return err
	}
	if t, err = t.Parse(homeTemplate); err != nil {
		return err
	}

	p.Header = "Home"
	p.ShowNavHome = false
	p.CallbackURL = callbackURL()

	if strings.TrimSpace(p.RedirectURLValue) == "" {
		p.RedirectURLValue = p.CallbackURL
	}

	if err = t.Execute(w, p); err != nil {
		return err
	}

	return nil
}

func sendLogin(w http.ResponseWriter, p Page) error {
	var err error
	t := template.New("login")
	if t, err = t.Parse(baseTemplate); err != nil {
		return err
	}
	if t, err = t.Parse(homeTemplate); err != nil {
		return err
	}

	p.Header = "Login"
	p.HomePath = "/"
	p.ShowNavHome = true
	p.CallbackURL = callbackURL()

	if strings.TrimSpace(p.RedirectURLValue) == "" {
		p.RedirectURLValue = p.CallbackURL
	}

	if err = t.Execute(w, p); err != nil {
		return err
	}

	return nil
}

func sendCallback(w http.ResponseWriter, p Page) error {
	var err error
	t := template.New("callback")
	if t, err = t.Parse(baseTemplate); err != nil {
		return err
	}
	if t, err = t.Parse(callbackTemplate); err != nil {
		return err
	}

	p.Header = "Callback"
	p.HomePath = "/"
	p.ShowNavHome = true

	if err = t.Execute(w, p); err != nil {
		return err
	}

	return nil
}

func handlerHome(w http.ResponseWriter, r *http.Request) {
	msg := fmt.Sprintf("home: URL=%s", r.URL.Path)
	log.Printf(msg)

	if r.URL.Path != "/" {
		log.Printf("home: URL=%s refusing to serve", r.URL.Path)
		http.NotFound(w, r)
		return
	}

	if err := sendHome(w, Page{AuthMsg: ""}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

type StateMessage struct {
	Config     oauth.Config
	ApiRequest string
}

func handlerLogin(w http.ResponseWriter, r *http.Request) {
	msg := fmt.Sprintf("login: URL=%s", r.URL.Path)
	log.Printf(msg)

	var (
		authMsg         string
		clientIdMsg     string
		clientSecretMsg string
		scopeMsg        string
		authURLMsg      string
		tokenURLMsg     string
		redirectURLMsg  string
		apiRequestMsg   string
	)

	clientId := strings.TrimSpace(r.FormValue("ClientId"))
	clientSecret := strings.TrimSpace(r.FormValue("ClientSecret"))
	scope := strings.TrimSpace(r.FormValue("Scope"))
	authURL := strings.TrimSpace(r.FormValue("AuthURL"))
	tokenURL := strings.TrimSpace(r.FormValue("TokenURL"))
	redirectURL := strings.TrimSpace(r.FormValue("RedirectURL"))
	apiRequest := strings.TrimSpace(r.FormValue("ApiRequest"))

	if clientId == "" {
		authMsg = "missing required field"
		clientIdMsg = "missing ClientId"
	}
	if clientSecret == "" {
		authMsg = "missing required field"
		clientSecretMsg = "missing ClientSecret"
	}
	if scope == "" {
		authMsg = "missing required field"
		scopeMsg = "missing Scope"
	}
	if authURL == "" {
		authMsg = "missing required field"
		authURLMsg = "missing AuthURL"
	}
	if tokenURL == "" {
		authMsg = "missing required field"
		tokenURLMsg = "missing TokenURL"
	}
	if redirectURL == "" {
		redirectURL = callbackURL()
	}
	if apiRequest == "" {
		authMsg = "missing required field"
		apiRequestMsg = "missing ApiRequest"
	}

	if authMsg != "" {
		if err := sendLogin(w,
			Page{AuthMsg: authMsg,

				ClientIdValue:     clientId,
				ClientSecretValue: clientSecret,
				ScopeValue:        scope,
				AuthURLValue:      authURL,
				TokenURLValue:     tokenURL,
				RedirectURLValue:  redirectURL,
				ApiRequestValue:   apiRequest,

				ClientIdMsg:     clientIdMsg,
				ClientSecretMsg: clientSecretMsg,
				ScopeMsg:        scopeMsg,
				AuthURLMsg:      authURLMsg,
				TokenURLMsg:     tokenURLMsg,
				RedirectURLMsg:  redirectURLMsg,
				ApiRequestMsg:   apiRequestMsg}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	stateMessage := StateMessage{
		Config: oauth.Config{
			ClientId:     clientId,
			ClientSecret: clientSecret,
			Scope:        scope,
			AuthURL:      authURL,
			TokenURL:     tokenURL,
			RedirectURL:  redirectURL,
			//TokenCache: oauth.CacheFile(TokenCacheFile),
		},
		ApiRequest: apiRequest,
	}

	// encode config in state (will be sent back to callback handler)
	state, err := json.Marshal(stateMessage)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	url := stateMessage.Config.AuthCodeURL(string(state))

	http.Redirect(w, r, url, http.StatusFound)

	// See next steps under callback handler
}

type Profile struct {
	Id    string
	Name  string
	Email string
}

func handlerCallback(w http.ResponseWriter, r *http.Request) {

	code := r.FormValue("code")
	state := r.FormValue("state")

	msg := fmt.Sprintf("callback: URL=%s code=%s state=%s", r.URL.Path, code, state)
	log.Printf(msg)

	// Load config from "state" parameter
	var stateMessage StateMessage
	if err := json.Unmarshal([]byte(state), &stateMessage); err != nil {
		e := fmt.Sprintf("callback: json.Unmarshal(state): %s", err.Error())
		http.Error(w, e, http.StatusInternalServerError)
		return
	}

	// Set up a Transport with our config, define the cache
	transp := &oauth.Transport{Config: &stateMessage.Config}

	// Exchange the authorization code for an access token.
	// ("Here's the code you gave the user, now give me a token!")
	var tok *oauth.Token
	var err error
	if tok, err = transp.Exchange(code); err != nil {
		e := fmt.Sprintf("callback: transp.Exchange(code): %s", err.Error())
		http.Error(w, e, http.StatusInternalServerError)
		return
	}

	// Make the actual request using the cached token to authenticate.
	// ("Here's the token, let me in!")
	transp.Token = tok

	// FIXME: Tack on the extra parameters, if specified.
	//apiRequest := "https://www.googleapis.com/oauth2/v1/userinfo"
	apiRequest := stateMessage.ApiRequest

	// Send request
	var resp *http.Response
	if resp, err = transp.Client().Get(apiRequest); err != nil {
		e := fmt.Sprintf("callback: transp.Client().Get(apiRequest): %s", err.Error())
		http.Error(w, e, http.StatusInternalServerError)
		return
	}

	// Read response
	var body []byte
	if body, err = ioutil.ReadAll(resp.Body); err != nil {
		e := fmt.Sprintf("callback: ioutil.ReadAll(resp.Body): %s", err.Error())
		http.Error(w, e, http.StatusInternalServerError)
		return
	}

	// Parse JSON
	var profile Profile
	if err = json.Unmarshal(body, &profile); err != nil {
		e := fmt.Sprintf("callback: json.Unmarshal(body): %s", err.Error())
		http.Error(w, e, http.StatusInternalServerError)
		return
	}

	msg = fmt.Sprintf("callback: URL=%s code=%s state=%s name=%s id=%s email=%s",
		r.URL.Path, code, state, profile.Name, profile.Id, profile.Email)
	log.Printf(msg)

	if err = sendCallback(w, Page{AuthMsg: msg}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

var addr string = "localhost:8080"
var callbackPath string = "/callback"

func homeURL() string {
	return "http://" + addr
}

func callbackURL() string {
	return "http://" + addr + callbackPath
}

func main() {

	http.HandleFunc("/", handlerHome)
	http.HandleFunc("/login", handlerLogin)
	http.HandleFunc(callbackPath, handlerCallback)

	log.Printf("serving at %s", homeURL())
	log.Printf("callback is %s", callbackURL())

	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Panicf("ListenAndServe: %s: %s", addr, err)
	}
}
