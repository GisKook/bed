package sha

import (
	"github.com/giskook/gotcp"
	"log"
	"net"
	"time"
)

type ServerConfig struct {
	Listener      *net.TCPListener
	AcceptTimeout time.Duration
	Uptopic       string
}

type Server struct {
	config           *ServerConfig
	srv              *gotcp.Server
	nsqproducer      *NsqProducer
	nsqconsumer      *NsqConsumer
	database         *ExecDatabase
	checkconnsticker *time.Ticker
}

var Gserver *Server

func SetServer(server *Server) {
	Gserver = server
}

func GetServer() *Server {
	return Gserver
}

func NewServer(srv *gotcp.Server, nsqproducer *NsqProducer, nsqconsumer *NsqConsumer, config *ServerConfig, db *ExecDatabase) *Server {
	serverstatistics := GetConfiguration().GetServerStatistics()
	return &Server{
		config:           config,
		srv:              srv,
		nsqproducer:      nsqproducer,
		nsqconsumer:      nsqconsumer,
		database:         db,
		checkconnsticker: time.NewTicker(time.Duration(serverstatistics) * time.Second),
	}
}

func (s *Server) GetProducer() *NsqProducer {
	return s.nsqproducer
}

func (s *Server) GetConsumer() *NsqConsumer {
	return s.nsqconsumer
}

func (s *Server) GetDatabase() *ExecDatabase {
	return s.database
}

func (s *Server) GetTopic() string {
	return s.config.Uptopic
}

func (s *Server) Start() {
	go s.nsqproducer.Start()
	go s.nsqconsumer.Start()

	go s.srv.Start(s.config.Listener, s.config.AcceptTimeout)
	go s.checkStatistics()
}

func (s *Server) Stop() {
	s.srv.Stop()
	s.nsqproducer.Stop()
	s.nsqconsumer.Stop()
	s.checkconnsticker.Stop()
}

func (s *Server) checkStatistics() {
	for {
		<-s.checkconnsticker.C
		log.Printf("---------------------Totol Connections : %d---------------------\n", NewConns().GetCount())
	}
}
