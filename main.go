package main

import (
	"flag"
	"github.com/dencat/fixss/fixss"
	log "github.com/jeanphorn/log4go"
	"os"
	"os/signal"
)

func main() {
	log.LoadConfiguration("./config/log.json")

	log.Info("Start application")

	port := 8080
	webServerPortPtr := flag.Int("port", port, "control port")
	flag.Parse()

	if webServerPortPtr != nil {
		port = *webServerPortPtr
	}

	fixss.LoadDefaultQuoteConfig()

	fixss.StartWebServer(port)

	err := fixss.StartAcceptor()
	if err != nil {
		log.Error("Can't start acceptor ", err)
		return
	}
	interrupt := make(chan os.Signal)
	signal.Notify(interrupt, os.Interrupt, os.Kill)
	<-interrupt

	fixss.StopAcceptor()

	log.Close()
}
