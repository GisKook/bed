package main

import (
	"fmt"
	"github.com/giskook/gotcp"
	"github.com/giskook/smarthome-access"
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
	configuration, err := sha.ReadConfig("./conf.json")
	sha.SetConfiguration(configuration)

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
	srv := gotcp.NewServer(config, &sha.Callback{}, &sha.ShaProtocol{})

	// creates a nsqproducer server
	nsqpconfig := &sha.NsqProducerConfig{
		Addr:  configuration.NsqConfig.Addr,
		Topic: configuration.NsqConfig.UpTopic,
	}
	nsqpserver := sha.NewNsqProducer(nsqpconfig)

	// creates a nsqconsumer server
	nsqcconfig := &sha.NsqConsumerConfig{
		Addr:    configuration.NsqConfig.Addr,
		Topic:   configuration.NsqConfig.DownTopic,
		Channel: configuration.NsqConfig.Downchannel,
	}
	nsqcserver := sha.NewNsqConsumer(nsqcconfig, nsqpserver)

	// database passwd_monitor
	dbconfig := &sha.DBConfig{
		Host:   configuration.DbConfig.Host,
		Port:   configuration.DbConfig.Port,
		User:   configuration.DbConfig.User,
		Passwd: configuration.DbConfig.Passwd,
		Dbname: configuration.DbConfig.Dbname,
	}

	database, dberr := sha.NewExecDatabase(dbconfig)
	if dberr != nil {
		log.Println(dberr.Error())
	} else {
		log.Println("conn to database success")
	}

	// create sha server
	shaserverconfig := &sha.ServerConfig{
		Listener:      listener,
		AcceptTimeout: time.Duration(configuration.ServerConfig.ConnTimeout) * time.Second,
		Uptopic:       configuration.NsqConfig.UpTopic,
	}
	shaserver := sha.NewServer(srv, nsqpserver, nsqcserver, shaserverconfig, database)
	sha.SetServer(shaserver)
	shaserver.Start()

	userhub, err := sha.NewUserPasswdHub(dbconfig)
	sha.SetUserPasswdHub(userhub)
	err = userhub.LoadAll()
	if err != nil {
		fmt.Println("connect to db error " + err.Error())
		return
	}
	fmt.Println("user info has been loaded")
	err = userhub.Listen(configuration.DbConfig.Monitortable)
	if err != nil {
		panic(err)
	}
	go userhub.WaitForNotification()

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
