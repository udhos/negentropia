package session

import (
	"log"
	"time"
	//"errors"
	//"strings"
	"net/http"
	"strconv"

	"negentropia/webserv/store"
	"negentropia/webserv/util"
)

const (
	AUTH_PROV_PASSWORD = 0 // local password
	AUTH_PROV_GOOGLE   = 1
	AUTH_PROV_FACEBOOK = 2

	SESSION_ID = "sid" // session id
)

var (
	sessionKeyExpire int64 = 2 * 86400 // expire keys after 2 days
)

type Session struct {
	SessionId    string
	AuthProvider int    // 1 = Google
	ProfileId    string // "102990441336549717697" (Google Profile)
	ProfileName  string // "Everton Marques"
	ProfileEmail string
}

func newCookie(name, value string, maxAge int) *http.Cookie {
	var expires time.Time

	// MaxAge=0 means no 'Max-Age' attribute specified.
	// MaxAge<0 means delete cookie now, equivalently 'Max-Age: 0'
	// MaxAge>0 means Max-Age attribute present and given in seconds
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
		HttpOnly: false, /* http-only cookie can't be read by javascript */
	}

	return cookie
}

func newSession(sid string, provider int, profId, profName, profEmail string) *Session {
	s := &Session{sid, provider, profId, profName, profEmail}
	//log.Printf("newSession sessionId=%s email=%s", s.SessionId, s.ProfileEmail)
	return s
}

func RedisQueryField(key, field string) string {
	return store.QueryField(key, field)
}

func Load(sessionId string) *Session {

	if !store.Exists(sessionId) {
		return nil
	}

	var (
		provider     int
		profileId    string
		profileName  string
		profileEmail string
	)

	provider, _ = strconv.Atoi(store.QueryField(sessionId, "AuthProvider"))
	profileId = store.QueryField(sessionId, "ProfileId")
	profileName = store.QueryField(sessionId, "ProfileName")
	profileEmail = store.QueryField(sessionId, "ProfileEmail")

	return newSession(sessionId, provider, profileId, profileName, profileEmail)
}

func sessionSave(session *Session) error {
	store.SetField(session.SessionId, "AuthProvider", strconv.Itoa(session.AuthProvider))
	store.SetField(session.SessionId, "ProfileId", session.ProfileId)
	store.SetField(session.SessionId, "ProfileName", session.ProfileName)
	store.SetField(session.SessionId, "ProfileEmail", session.ProfileEmail)

	store.Expire(session.SessionId, sessionKeyExpire)

	return nil
}

func newSessionId() string {
	return "s:" + strconv.FormatInt(store.Incr("i:sessionIdGenerator"), 10) + util.RandomSuffix()
}

func Get(r *http.Request) *Session {

	cook, err := r.Cookie(SESSION_ID)
	if err != nil {
		//log.Printf("session.Get cookie NOT FOUND: err=%s", err)
		return nil
	}

	//log.Printf("session.Get FOUND cookie session=%s", cook.Value)

	session := Load(cook.Value)
	if session == nil {
		log.Printf("session.Get: failure loading session id=%s", cook.Value)
		return nil
	}

	//log.Printf("session.Get LOADED session=%s", cook.Value)

	return session
}

func Set(w http.ResponseWriter, provider int, profId, profName, profEmail string) *Session {

	sessionId := newSessionId()

	session := newSession(sessionId, provider, profId, profName, profEmail)

	err := sessionSave(session)
	if err != nil {
		log.Printf("session.Set: failure saving session id=%s error=[%s]", sessionId, err)
		return nil
	}

	// MaxAge=0 means no 'Max-Age' attribute specified.
	cook := newCookie(SESSION_ID, sessionId, 0)

	http.SetCookie(w, cook)

	return session
}

func Delete(w http.ResponseWriter, session *Session) {

	store.Del(session.SessionId)

	// MaxAge<0 means delete cookie now
	cook := newCookie(SESSION_ID, "", -1)

	http.SetCookie(w, cook)
}
