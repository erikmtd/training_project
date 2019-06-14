package mq

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/training_project/core/cache"
	"github.com/training_project/core/mq"
	"github.com/training_project/handler"
)

type mqHandler struct {
	mq mq.MQ
}

func (mq *mqHandler) Start() {
	fmt.Println("MQ starts")
	mq.mq.StartConsumers()
}

var lock = &sync.Mutex{}
var instance handler.Handler

func New() handler.Handler {
	lock.Lock()
	defer lock.Unlock()

	if instance == nil {
		q := mq.New()
		q.AddConsumer(mq.Consumer{
			Host:        "devel-go.tkpd:4161",
			Channel:     "erik_channel",
			Topic:       "tech_cur_nsq_0619_erik",
			MaxAttempt:  10,
			MaxInFlight: 100,
			Handler:     pageVisitorConsumerHandler,
		})

		instance = &mqHandler{
			mq: q,
		}
	}
	return instance
}

func pageVisitorConsumerHandler(message []byte) error {
	recordPageVisitorCount()
	return nil
}

var recLock = &sync.Mutex{}

func recordPageVisitorCount() {
	recLock.Lock()
	defer recLock.Unlock()

	chacheStore := cache.New()
	key := "training_project_0619_erik"

	var v string
	if v = chacheStore.Get(key); v == "" {
		v = "0"
	}

	count, err := strconv.Atoi(v)
	if err == nil {
		count += 1
		chacheStore.SetWithExpiredTime(key, strconv.Itoa(count), 3600)
	}
}
