package mq

import (
	"fmt"
	"log"

	"encoding/json"

	"github.com/nsqio/go-nsq"
)

type (
	MQ interface {
		AddConsumer(Consumer)
		StartConsumers()
		Publish(string, string, interface{})
	}
	Consumer struct {
		Host        string
		Topic       string
		Channel     string
		MaxAttempt  uint16
		MaxInFlight int
		Handler     MQHandler
	}
	nsqMQ struct {
		consumers []Consumer
	}
	MQHandler func([]byte) error
)

func (q *nsqMQ) AddConsumer(consumer Consumer) {
	q.consumers = append(q.consumers, consumer)
}
func (q *nsqMQ) StartConsumers() {
	for _, c := range q.consumers {
		NSQConfig := nsq.NewConfig()
		NSQConfig.MaxAttempts = c.MaxAttempt
		NSQConfig.MaxInFlight = c.MaxInFlight

		NSQConsumer, err := nsq.NewConsumer(c.Topic, c.Channel, NSQConfig)
		if err != nil {
			log.Fatal(err)
		} else {
			var handler nsq.HandlerFunc = func(message *nsq.Message) (err error) {
				fmt.Printf("consumer %s is consuming \n", c.Channel)
				fmt.Print("Message : ", message)
				err = c.Handler(message.Body)
				return
			}
			NSQConsumer.AddHandler(handler)
			err = NSQConsumer.ConnectToNSQLookupd(c.Host)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

func (q *nsqMQ) Publish(host string, topic string, data interface{}) {
	p, err := nsq.NewProducer(host, nsq.NewConfig())
	if err != nil {
		log.Fatal(p)
	} else {
		json, jsonErr := json.Marshal(data)
		if jsonErr != nil {
			log.Fatal(jsonErr)
		} else {
			err = p.Publish(topic, json)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

func New() MQ {
	return &nsqMQ{}
}
