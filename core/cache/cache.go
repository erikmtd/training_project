package cache

import (
	"log"
	"sync"

	r "github.com/garyburd/redigo/redis"
)

type (
	Cache interface {
		SetWithExpiredTime(key string, data string, time int)
		Get(key string) string
	}
	redis struct {
		pool r.Pool
	}
)

var lock = &sync.Mutex{}
var instance *redis

func (r *redis) SetWithExpiredTime(key string, data string, time int) {
	conn := r.pool.Get()
	defer conn.Close()

	_, err := conn.Do("SET", key, data)
	if err != nil {
		log.Fatal(err)
	}
	_, err = conn.Do("EXPIRE", key, time)
	if err != nil {
		log.Fatal(err)
	}
	conn.Flush()

}
func (rd *redis) Get(key string) string {
	conn := rd.pool.Get()
	defer conn.Close()

	v, err := r.String(conn.Do("GET", key))

	if err != nil {
		return ""
	}
	return v
}

func New() Cache {
	lock.Lock()
	defer lock.Unlock()

	if instance == nil {

		instance = &redis{
			pool: r.Pool{
				IdleTimeout: 100,
				MaxActive:   30,
				MaxIdle:     50,
				Dial: func() (r.Conn, error) {
					return r.Dial("tcp", "devel-redis.tkpd:6379")
				},
			},
		}
	}
	return instance
}
