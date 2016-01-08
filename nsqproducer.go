package sha

import (
	"log"
	//"sync"

	"github.com/bitly/go-nsq"
)

type NsqProducerConfig struct {
	Addr  string
	Topic string
}

type NsqProducer struct {
	config *NsqProducerConfig
	//waitGroup *sync.WaitGroup

	producer *nsq.Producer
}

func NewNsqProducer(config *NsqProducerConfig) *NsqProducer {
	return &NsqProducer{
		config: config,
		//waitGroup: &sync.WaitGroup{},
	}
}

func (s *NsqProducer) Send(topic string, value []byte) error {
	err := s.producer.Publish(topic, value)

	return err
}

func (s *NsqProducer) Start() {
	//s.waitGroup.Add(1)
	defer func() {
		err := recover()
		//s.waitGroup.Done()
		if err != nil {
			log.Println("err found")
		}

	}()
	config := nsq.NewConfig()

	var errmsg error
	s.producer, errmsg = nsq.NewProducer(s.config.Addr, config)

	if errmsg != nil {
		//	log.Printf("create producer error" + errmsg.Error())
		panic("create producer error " + errmsg.Error())
	}
}

func (s *NsqProducer) GetTopic() string {
	return s.config.Topic
}

func (s *NsqProducer) Stop() {
	//	s.waitGroup.Done()
	//	s.waitGroup.Wait()

	s.producer.Stop()
}
