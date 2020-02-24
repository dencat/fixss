package main

import (
	"os"
	"os/signal"

	//"github.com/quickfixgo/tag"
	//"github.com/shopspring/decimal"
	"github.com/juju/loggo"
)

const CONFIG_PATH = "config/server.cfg"

var Log = loggo.GetLogger("")

func main() {
	loggo.GetLogger("").SetLogLevel(loggo.INFO)
	loggo.GetLogger("fix").SetLogLevel(loggo.INFO)

	Log.Infof("Start application")

	LoadDefaultQuoteConfig()

	err := StartAcceptor()
	if err != nil {
		Log.Errorf("Can't start acceptor ", err)
		return
	}
	interrupt := make(chan os.Signal)
	signal.Notify(interrupt, os.Interrupt, os.Kill)
	<-interrupt

	StopAcceptor()
}
