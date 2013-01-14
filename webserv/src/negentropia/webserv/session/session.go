package session

import (
	"log"
	"time"
	//"errors"
	//"strings"
	"strconv"
	"net/http"
	
	//"github.com/bradfitz/gomemcache/memcache"
	"github.com/vmihailenco/redis"
)

const (
	AUTH_PROV_CRED     = 0
	AUTH_PROV_GOOGLE   = 1
	AUTH_PROV_FACEBOOK = 2
)

var (
	/*
	mcServerList     []string         = []string{"127.0.0.1:11211", "127.0.0.1:12000"}
	mc               *memcache.Client
	*/
	redisAddr     string  = "localhost:6379"
	redisPassword string  = ""
	redisDb       int64   = -1
	redisClient   *redis.Client
	redisExpire   int64   = 2 * 86400 // expire keys after 2 days
)

type Session struct {
	SessionId        string
	AuthProvider     int    // 1 = Google
	AuthProviderId   string // "102990441336549717697" (Google Profile)
	AuthProviderName string // "Everton Marques"
}

func init() {
	/*
	log.Printf("session.init(): memcache client for: " + strings.Join(mcServerList, ","))
	mc = memcache.New(mcServerList...)
	*/
	log.Printf("session.init(): redis client for: %s", redisAddr)
	redisClient = redis.NewTCPClient(redisAddr, redisPassword, redisDb)
}

/*
func sessionGet() string {
	//var it *memcache.Item
    it, err := mc.Get("session")
	if err != nil {
		log.Printf("handler.googleCallback mc.Get err=%s", err)
		return ""
	}
	return string(it.Value)
}

func sessionSet(value string) error {
    err := mc.Set(&memcache.Item{Key: "session", Value: []byte(value), Expiration: 24*3600})
	if err != nil {
		log.Printf("handler.googleCallback mc.Set err=%s", err)
	}
	return err
}
*/

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
		HttpOnly: true,
	}
	
	return cookie
}

func newSession(sid string, provider int, acctId string, acctName string) *Session {
	return &Session{sid, provider, acctId, acctName}
}

func sessionLoad(sessionId string) *Session {

	if !redisClient.Exists(sessionId).Val() {
		return nil
	}
	
	var (
		//err          error
		provider     int 
		profileId    string
		profileName  string
	)
	
	provider, _   = strconv.Atoi(redisClient.HGet(sessionId, "AuthProvider").Val())
	profileId     = redisClient.HGet(sessionId, "AuthProviderId").Val()
	profileName   = redisClient.HGet(sessionId, "AuthProviderName").Val()
	
	return newSession(sessionId, provider, profileId, profileName)
}

func sessionSave(session *Session) error {
	redisClient.HSet(session.SessionId, "AuthProvider",     strconv.Itoa(session.AuthProvider))
	redisClient.HSet(session.SessionId, "AuthProviderId",   session.AuthProviderId)
	redisClient.HSet(session.SessionId, "AuthProviderName", session.AuthProviderName)

	redisClient.Expire(session.SessionId, redisExpire)
	
	return nil
}

func newSessionId() string {
	return strconv.FormatInt(redisClient.Incr("sessionIdGenerator").Val(), 10)
}

func Get(r *http.Request) *Session {

	cook, err := r.Cookie("session")
	if err != nil {
		log.Printf("session.Get cookie NOT FOUND: err=%s", err)
		return nil
	}

	log.Printf("session.Get FOUND cookie session=%s", cook.Value)
		
	session := sessionLoad(cook.Value);
	if session == nil {
		log.Printf("session.Get: failure loading session id=%s", cook.Value)	
		return nil
	}

	log.Printf("session.Get LOADED session=%s", cook.Value)
		
	return session
}

func Set(w http.ResponseWriter, provider int, acctId string, acctName string) *Session {
		
	sessionId := newSessionId()

	session := newSession(sessionId, provider, acctId, acctName)
	
	err := sessionSave(session)
	if (err != nil) {
		log.Printf("session.Set: failure saving session id=%s error=[%s]", sessionId, err)
		return nil
	}

	// MaxAge=0 means no 'Max-Age' attribute specified.
	cook := newCookie("session", sessionId, 0)

	http.SetCookie(w, cook)
	
	return session
}

func Delete(w http.ResponseWriter, session *Session) {

	redisClient.Del(session.SessionId)
	
	// MaxAge<0 means delete cookie now
	cook := newCookie("session", "", -1)

	http.SetCookie(w, cook)
}
