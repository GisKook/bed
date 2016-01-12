package main

import (
	"fmt"
	"github.com/giskook/bed"
	"github.com/giskook/gotcp"
	"log"
	"net"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	// read configuration
	configuration, err := bed.ReadConfig("./conf.json")
	bed.SetConfiguration(configuration)

	checkError(err)
	// creates a tcp listener
	tcpAddr, err := net.ResolveTCPAddr("tcp4", ":8989")
	checkError(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)

	// creates a tcp server
	config := &gotcp.Config{
		PacketSendChanLimit:    20,
		PacketReceiveChanLimit: 20,
	}
	srv := gotcp.NewServer(config, &bed.Callback{}, &bed.BedProtocol{})

	// creates a nsqproducer server
	nsqpconfig := &bed.NsqProducerConfig{
		Addr:  configuration.NsqConfig.Addr,
		Topic: configuration.NsqConfig.UpTopic,
	}
	nsqpserver := bed.NewNsqProducer(nsqpconfig)

	// creates a nsqconsumer server
	nsqcconfig := &bed.NsqConsumerConfig{
		Addr:    configuration.NsqConfig.Addr,
		Topic:   configuration.NsqConfig.DownTopic,
		Channel: configuration.NsqConfig.Downchannel,
	}
	nsqcserver := bed.NewNsqConsumer(nsqcconfig, nsqpserver)

	// create bed server
	bedserverconfig := &bed.ServerConfig{
		Listener:      listener,
		AcceptTimeout: time.Duration(configuration.ServerConfig.ConnTimeout) * time.Second,
		Uptopic:       configuration.NsqConfig.UpTopic,
	}
	bedserver := bed.NewServer(srv, nsqpserver, nsqcserver, bedserverconfig)
	bed.SetServer(bedserver)
	bedserver.Start()

	// starts service
	fmt.Println("listening:", listener.Addr())

	// catchs system signal
	chSig := make(chan os.Signal)
	signal.Notify(chSig, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("Signal: ", <-chSig)

	// stops service
	srv.Stop()
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
