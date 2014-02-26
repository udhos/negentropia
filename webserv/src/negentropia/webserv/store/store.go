package store

/*
	Some key conventions:
	s:xxx	session id
	c:xxx	signup confirmation id
	r:xxx	reset password confirmation id
	z:xxx	zone id
	p:xxx   shader program id
	o:xxx	object/model
	m:xxx	model instance
	l:xxx	instance list
	i:xxx	id generator		session.go		i:sessionIdGenerator
								signup.go		i:confirmationIdGenerator
								password.go		i:resetPassConfirmationIdGenerator
	xxx@yyy	email
*/

import (
	"log"
	"time"

	"github.com/vmihailenco/redis/v2"
)

type KeyField struct {
	key   string
	field string
}

type KeyValue struct {
	key   string
	value string
}

type KeyFieldValue struct {
	key   string
	field string
	value string
}

type KeyExpire struct {
	key    string
	expire time.Duration
}

var (
	redisPassword  string = ""
	redisClient    *redis.Client
	redisExpire    time.Duration      = 2 * 86400 * time.Second // expire keys after 2 days
	queryFieldReq  chan KeyField      = make(chan KeyField)
	queryFieldRep  chan string        = make(chan string)
	setFieldReq    chan KeyFieldValue = make(chan KeyFieldValue)
	expReq         chan KeyExpire     = make(chan KeyExpire)
	existsReq      chan string        = make(chan string)
	existsRep      chan bool          = make(chan bool)
	incrReq        chan string        = make(chan string)
	incrRep        chan int64         = make(chan int64)
	delReq         chan string        = make(chan string)
	fieldExistsReq chan KeyField      = make(chan KeyField)
	fieldExistsRep chan bool          = make(chan bool)
	setReq         chan KeyValue      = make(chan KeyValue)
	persistReq     chan string        = make(chan string)
	getReq         chan string        = make(chan string)
	getRep         chan string        = make(chan string)
	delFieldReq    chan KeyField      = make(chan KeyField)
	querySetReq    chan string        = make(chan string)
	querySetRep    chan []string      = make(chan []string)
)

func serve() {
	log.Printf("store.serve: goroutine started")
	for {
		select {
		case r := <-queryFieldReq:
			queryFieldRep <- redisClient.HGet(r.key, r.field).Val()
		case r := <-setFieldReq:
			redisClient.HSet(r.key, r.field, r.value)
		case r := <-expReq:
			redisClient.Expire(r.key, r.expire)
		case key := <-existsReq:
			existsRep <- redisClient.Exists(key).Val()
		case key := <-incrReq:
			incrRep <- redisClient.Incr(key).Val()
		case key := <-delReq:
			redisClient.Del(key)
		case r := <-fieldExistsReq:
			fieldExistsRep <- redisClient.HExists(r.key, r.field).Val()
		case r := <-setReq:
			redisClient.Set(r.key, r.value)
		case key := <-persistReq:
			redisClient.Persist(key)
		case key := <-getReq:
			getRep <- redisClient.Get(key).Val()
		case r := <-delFieldReq:
			redisClient.HDel(r.key, r.field)
		case key := <-querySetReq:
			querySetRep <- redisClient.SMembers(key).Val()
		}
	}
}

func Init(serverAddr string) {
	log.Printf("store.Init: redis client for: %s", serverAddr)

	log.Printf("store.Init: redisExpire = %d seconds", redisExpire/time.Second)

	//redisClient = redis.NewTCPClient(serverAddr, redisPassword, redisDb)
	redisClient := redis.NewTCPClient(&redis.Options{
		Addr:     serverAddr,
		Password: redisPassword,
		DB:       0, // use default DB
	})
	defer redisClient.Close()

	pong, err := redisClient.Ping().Result()
	log.Printf("store.Init: PING redis server: reply=%s err=%s", pong, err)

	go serve()
}

func QueryField(key, field string) string {
	queryFieldReq <- KeyField{key, field} // send key,field
	return <-queryFieldRep                // read reply and return it
}

func SetField(key, field, value string) {
	setFieldReq <- KeyFieldValue{key, field, value} // send key,field,value
}

func Expire(key string, expire time.Duration) {
	expReq <- KeyExpire{key, expire} // send key,expire
}

func Exists(key string) bool {
	existsReq <- key   // send key
	return <-existsRep // read reply and return it
}

func Incr(key string) int64 {
	incrReq <- key   // send key
	return <-incrRep // read reply and return it
}

func Del(key string) {
	delReq <- key // send key
}

func FieldExists(key, field string) bool {
	fieldExistsReq <- KeyField{key, field} // send key,field
	return <-fieldExistsRep                // read reply and return it
}

func Set(key, value string) {
	setReq <- KeyValue{key, value} // send key,value
}

func Persist(key string) {
	persistReq <- key // send key
}

func Get(key string) string {
	getReq <- key   // send key
	return <-getRep // read reply and return it
}

func DelField(key, field string) {
	delFieldReq <- KeyField{key, field} // send key,field
}

func QuerySet(key string) []string {
	querySetReq <- key   // send key,field
	return <-querySetRep // read reply and return it
}
