package main

import (
	"bufio"
	"github.com/giskook/smarthome-access/client"
	"io"
	"log"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"sync"
	"syscall"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	file, _ := os.OpenFile("./smarthomebox.txt", os.O_RDONLY, 0666)
	reader := bufio.NewReader(file)
	wg := &sync.WaitGroup{}
	boxcount := 1
	counter := 0
	for {
		buf, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		values := strings.Split(string(buf), "-")
		log.Println(values)
		box := shb.NewSmarthomebox(values[0], values[1])
		for i := 0; i < 4; i++ {
			buf, _, err = reader.ReadLine()
			values = strings.Split(string(buf), "-")
			box.Add(values[0], 1, 1, values[1], 1)
		}

		go box.Do("192.168.2.67:8989", wg)
		counter++
		if counter == boxcount {
			break
		}
	}
	wg.Wait()

	chSig := make(chan os.Signal)
	signal.Notify(chSig, syscall.SIGINT, syscall.SIGTERM)
	log.Println("Signal: ", <-chSig)
}
