package bed

import (
	"log"
	"sync"

	"github.com/bitly/go-nsq"
	"github.com/giskook/bed/pb"
)

type NsqConsumerConfig struct {
	Addr    string
	Topic   string
	Channel string
}

type NsqConsumer struct {
	config    *NsqConsumerConfig
	waitGroup *sync.WaitGroup

	consumer *nsq.Consumer
	producer *NsqProducer
}

func NewNsqConsumer(config *NsqConsumerConfig, producer *NsqProducer) *NsqConsumer {
	return &NsqConsumer{
		config:    config,
		waitGroup: &sync.WaitGroup{},
		producer:  producer,
	}
}

func (s *NsqConsumer) recvNsq() {
	s.consumer.AddHandler(nsq.HandlerFunc(func(message *nsq.Message) error {
		data := message.Body
		bedid, serialnum, command, err := CheckNsqProtocol(data)
		log.Printf("recvnsq bedid %x ", bedid)
		log.Println("cmd %d\n", command.Type)
		if err == nil {
			switch command.Type {
			case Report.Command_CMT_REQBEDRUN:
				packet := ParseNsqBedControl(serialnum, command)
				if packet != nil {
					if NewConns().Check(bedid) {
						NewConns().GetConn(bedid).SendToBed(packet)
					}
				}
			case Report.Command_CMT_REQTOILET:
				packet := ParseNsqPotty(serialnum)
				if packet != nil {
					if NewConns().Check(bedid) {
						NewConns().GetConn(bedid).SendToBed(packet)
					}
				}
			case Report.Command_CMT_REQBEDRESET:
				packetahead := ParseNsqBedResetAhead(serialnum)
				if packetahead != nil {
					s.producer.Send(s.producer.GetTopic(), packetahead.Serialize())
				}
				packet := ParseNsqBedReset(serialnum)
				if packet != nil {
					s.producer.Send(s.producer.GetTopic(), packet.Serialize())
				}
			}
		}

		return nil
	}))
}

func (s *NsqConsumer) Start() {
	s.waitGroup.Add(1)
	defer func() {
		s.waitGroup.Done()
		err := recover()
		if err != nil {
			log.Println("err found")
			s.Stop()
		}

	}()

	config := nsq.NewConfig()

	var errmsg error
	s.consumer, errmsg = nsq.NewConsumer(s.config.Topic, s.config.Channel, config)

	if errmsg != nil {
		//	panic("create consumer error -> " + errmsg.Error())
		log.Println("create consumer error -> " + errmsg.Error())
	}
	s.recvNsq()

	err := s.consumer.ConnectToNSQD(s.config.Addr)
	if err != nil {
		panic("Counld not connect to nsq -> " + err.Error())
	}
}

func (s *NsqConsumer) Stop() {
	s.waitGroup.Wait()

	errmsg := s.consumer.DisconnectFromNSQD(s.config.Addr)

	if errmsg != nil {
		log.Printf("stop consumer error ", errmsg.Error())
	}

	s.consumer.Stop()
}
