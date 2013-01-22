package store

import (
	"log"
	
	"github.com/vmihailenco/redis"
)

type KeyField struct {
	key   string
	field string
}

var (
	redisPassword  string        = ""
	redisDb        int64         = -1
	redisClient   *redis.Client
	redisExpire    int64         = 2 * 86400 // expire keys after 2 days
	req            chan KeyField = make(chan KeyField)
	rep            chan string   = make(chan string)
)

func serveField() {
	log.Printf("store.serveField: goroutine started")
	// receives values from the channel repeatedly until it is closed
	for r := range req {
		rep <- redisClient.HGet(r.key, r.field).Val()
	}
	log.Printf("store.serveField: PANIC: req channel closed")
}

func Init(serverAddr string) {
	log.Printf("store.Init(): redis client for: %s", serverAddr)
	redisClient = redis.NewTCPClient(serverAddr, redisPassword, redisDb)
	go serveField()
}

func QueryField(key, field string) string {
	req <- KeyField{key, field} // send key,field
	return <- rep // read reply and return it
}
