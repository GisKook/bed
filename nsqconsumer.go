package sha

import (
	"log"
	"sync"

	"github.com/bitly/go-nsq"
	"github.com/giskook/smarthome-access/pb"
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
		gatewayid, serialnum, command, err := CheckNsqProtocol(data)
		log.Println("recvnsq")
		log.Println("cmd %d\n", command.Type)
		if err == nil {
			switch command.Type {
			case Report.Command_CMT_REQLOGIN:
				packet := ParseNsqLogin(gatewayid, serialnum, command)
				if packet != nil {
					s.producer.Send(s.producer.GetTopic(), packet.Serialize())
				}
			case Report.Command_CMT_REQDEVICELIST:
				packet := ParseNsqDeviceList(gatewayid, serialnum, command)
				//NewConns().GetConn(gatewayid).SendToGateway(packet)
				if packet != nil {
					s.producer.Send(s.producer.GetTopic(), packet.Serialize())
				}
			case Report.Command_CMT_REQOP:
				packet := ParseNsqOp(gatewayid, serialnum, command)
				if packet != nil {
					NewConns().GetConn(gatewayid).SendToGateway(packet)
				}
			case Report.Command_CMT_REQONLINE:
				packet := ParseNsqCheckOnline(gatewayid, serialnum)
				if packet != nil {
					s.producer.Send(s.producer.GetTopic(), packet.Serialize())
				}
			case Report.Command_CMT_REQSETDEVICENAME:
				packet := ParseNsqSetDevicename(gatewayid, serialnum, command)
				if packet != nil {
					NewConns().GetConn(gatewayid).SendToGateway(packet)
				}
			case Report.Command_CMT_REQCHANGEPASSWD:
				packet := ParseNsqChangePasswd(gatewayid, serialnum, command)
				if packet != nil {
					s.producer.Send(s.producer.GetTopic(), packet.Serialize())
				}
			case Report.Command_CMT_REQDELDEVICE:
				packet := ParseNsqDelDevice(gatewayid, serialnum, command)
				if packet != nil {
					NewConns().GetConn(gatewayid).SendToGateway(packet)
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
