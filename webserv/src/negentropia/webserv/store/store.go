package store

import (
	"log"
	
	"github.com/vmihailenco/redis"
)

type KeyField struct {
	key   string
	field string
}

type KeyFieldValue struct {
	key   string
	field string
	value string
}

type KeyExpire struct {
	key    string
	expire int64
}

var (
	redisPassword  string        = ""
	redisDb        int64         = -1
	redisClient   *redis.Client
	redisExpire    int64         = 2 * 86400 // expire keys after 2 days
	queryReq       chan KeyField      = make(chan KeyField)
	queryRep       chan string        = make(chan string)
	setReq         chan KeyFieldValue = make(chan KeyFieldValue)
	expReq         chan KeyExpire     = make(chan KeyExpire)
	existsReq      chan string        = make(chan string)
	existsRep      chan bool          = make(chan bool)
	incrReq        chan string        = make(chan string)
	incrRep        chan int64         = make(chan int64)	
	delReq         chan string        = make(chan string)
	fieldExistsReq chan KeyField      = make(chan KeyField)
	fieldExistsRep chan bool          = make(chan bool)
)

func serve() {
	log.Printf("store.serve: goroutine started")
	for {
		select {
			case r := <- queryReq:
				queryRep <- redisClient.HGet(r.key, r.field).Val()
			case r := <- setReq:
				redisClient.HSet(r.key, r.field, r.value)
			case r := <- expReq:
				redisClient.Expire(r.key, r.expire)
			case key := <- existsReq:
				existsRep <- redisClient.Exists(key).Val()
			case key := <- incrReq:
				incrRep <- redisClient.Incr(key).Val()
			case key := <- delReq:
				redisClient.Incr(key)
			case r := <- fieldExistsReq:
				fieldExistsRep <- redisClient.HExists(r.key, r.field).Val()
		}
	}
}

/*
func serveQueryField() {
	log.Printf("store.serveQueryField: goroutine started")
	// receives values from the channel repeatedly until it is closed
	for r := range queryReq {
		queryRep <- redisClient.HGet(r.key, r.field).Val()
	}
	log.Printf("store.serveQueryField: PANIC: queryReq channel closed")
}

func serveSetField() {
	log.Printf("store.serveSetField: goroutine started")
	// receives values from the channel repeatedly until it is closed
	for r := range setReq {
		redisClient.HSet(r.key, r.field, r.value)
	}
	log.Printf("store.serveSetField: PANIC: queryReq channel closed")
}

func serveExpire() {
	log.Printf("store.serveExpire: goroutine started")
	// receives values from the channel repeatedly until it is closed
	for r := range expReq {
		redisClient.Expire(r.key, r.expire)
	}
	log.Printf("store.serveExpire: PANIC: expReq channel closed")
}

func serveExists() {
	log.Printf("store.serveExists: goroutine started")
	// receives values from the channel repeatedly until it is closed
	for key := range existsReq {
		existsRep <- redisClient.Exists(key).Val()
	}
	log.Printf("store.serveExists: PANIC: existsReq channel closed")
}

func serveIncr() {
	log.Printf("store.serveIncr: goroutine started")
	// receives values from the channel repeatedly until it is closed
	for key := range incrReq {
		incrRep <- redisClient.Incr(key).Val()
	}
	log.Printf("store.serveIncr: PANIC: incrReq channel closed")
}

func serveDel() {
	log.Printf("store.serveDel: goroutine started")
	// receives values from the channel repeatedly until it is closed
	for key := range delReq {
		redisClient.Incr(key)
	}
	log.Printf("store.serveDel: PANIC: delReq channel closed")
}
*/

func Init(serverAddr string) {
	log.Printf("store.Init: redis client for: %s", serverAddr)
	redisClient = redis.NewTCPClient(serverAddr, redisPassword, redisDb)
	/*
	go serveQueryField()
	go serveSetField()	
	go serveExpire()	
	go serveExists()
	go serveIncr()
	go serveDel()
	*/
	go serve()
}

func QueryField(key, field string) string {
	queryReq <- KeyField{key, field} // send key,field
	return <- queryRep // read reply and return it
}

func SetField(key, field, value string) {
	setReq <- KeyFieldValue{key, field, value} // send key,field,value
}

func Expire(key string, expire int64) {
	expReq <- KeyExpire{key, expire} // send key,expire
}

func Exists(key string) bool {
	existsReq <- key // send key
	return <- existsRep // read reply and return it
}

func Incr(key string) int64 {
	incrReq <- key // send key
	return <- incrRep // read reply and return it
}

func Del(key string) {
	delReq <- key // send key
}

func FieldExists(key, field string) bool {
	fieldExistsReq <- KeyField{key, field} // send key,field
	return <- fieldExistsRep // read reply and return it
}
