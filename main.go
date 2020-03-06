package main

import (
	"github.com/dencat/fixss/fixss"
	"github.com/juju/loggo"
	"os"
	"os/signal"
)

func main() {
	loggo.GetLogger("").SetLogLevel(loggo.INFO)
	loggo.GetLogger("fix").SetLogLevel(loggo.INFO)

	fixss.Log.Infof("Start application")

	fixss.LoadDefaultQuoteConfig()

	fixss.StartWebServer()

	err := fixss.StartAcceptor()
	if err != nil {
		fixss.Log.Errorf("Can't start acceptor ", err)
		return
	}
	interrupt := make(chan os.Signal)
	signal.Notify(interrupt, os.Interrupt, os.Kill)
	<-interrupt

	fixss.StopAcceptor()
}
