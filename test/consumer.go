package main

import (
	"log"
	"sync"

	"github.com/bitly/go-nsq"
)

func main() {
	wg := &sync.WaitGroup{}
	wg.Add(1)

	config := nsq.NewConfig()
	var w *nsq.Producer
	w, _ = nsq.NewProducer("127.0.0.1:4150", config)

	log.Printf("%v", w)
	//w.Publish("write_test", []byte("test"))
	q, _ := nsq.NewConsumer("write_test", "ch", config)
	q.AddHandler(nsq.HandlerFunc(func(message *nsq.Message) error {
		log.Printf("Got a message: %v", message)
		//wg.Done()
		return nil
	}))
	log.Printf("before conn")
	err := q.ConnectToNSQD("127.0.0.1:4150")
	if err != nil {
		log.Panic("Could not connect")
	}

	wg.Wait()

}
