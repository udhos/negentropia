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
	
	"negentropia/webserv/store"
)

const (
	AUTH_PROV_PASSWORD = 0 // local password
	AUTH_PROV_GOOGLE   = 1
	AUTH_PROV_FACEBOOK = 2
)

var (
	/*
	mcServerList     []string         = []string{"127.0.0.1:11211", "127.0.0.1:12000"}
	mc               *memcache.Client
	*/
	//RedisAddr     string
	redisPassword string  = ""
	redisDb       int64   = -1
	redisClient   *redis.Client
	redisExpire   int64   = 2 * 86400 // expire keys after 2 days
)

type Session struct {
	SessionId    string
	AuthProvider int    // 1 = Google
	ProfileId    string // "102990441336549717697" (Google Profile)
	ProfileName  string // "Everton Marques"
	ProfileEmail string
}

/*
func init() {
	log.Printf("session.init(): memcache client for: " + strings.Join(mcServerList, ","))
	mc = memcache.New(mcServerList...)
}
*/

func Init(serverAddr string) {
	log.Printf("session.Init(): redis client for: %s", serverAddr)
	redisClient = redis.NewTCPClient(serverAddr, redisPassword, redisDb)
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

func newSession(sid string, provider int, profId, profName, profEmail string) *Session {
	s := &Session{sid, provider, profId, profName, profEmail}
	log.Printf("newSession sessionId=%s email=%s", s.SessionId, s.ProfileEmail)
	return s
}

func RedisQueryField(key, field string) string {
	//return redisClient.HGet(key, field).Val()
	return store.QueryField(key, field)
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
		profileEmail string
	)
	
	provider, _  = strconv.Atoi(redisClient.HGet(sessionId, "AuthProvider").Val())
	profileId    = redisClient.HGet(sessionId, "ProfileId").Val()
	profileName  = redisClient.HGet(sessionId, "ProfileName").Val()
	profileEmail = redisClient.HGet(sessionId, "ProfileEmail").Val()
	
	return newSession(sessionId, provider, profileId, profileName, profileEmail)
}

func sessionSave(session *Session) error {
	redisClient.HSet(session.SessionId, "AuthProvider", strconv.Itoa(session.AuthProvider))
	redisClient.HSet(session.SessionId, "ProfileId",    session.ProfileId)
	redisClient.HSet(session.SessionId, "ProfileName",  session.ProfileName)
	redisClient.HSet(session.SessionId, "ProfileEmail", session.ProfileEmail)

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

func Set(w http.ResponseWriter, provider int, profId, profName, profEmail string) *Session {
		
	sessionId := newSessionId()

	session := newSession(sessionId, provider, profId, profName, profEmail)
	
	err := sessionSave(session)
	if err != nil {
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
