package main

import (
	"errors"
	"flag"
	"github.com/dencat/fixss/fixss"
	log "github.com/jeanphorn/log4go"
	"os"
	"os/signal"
)

func main() {
	cfgPath, err := ParseFlags()
	if err != nil {
		println(err)
		return
	}
	cfg, err := fixss.NewConfig(cfgPath)
	if err != nil {
		println(err)
		return
	}

	log.LoadConfiguration(cfg.Logging.Config)
	defer log.Close()
	log.LOGGER("app").Info("Start application")

	err = fixss.LoadDefaultQuoteConfig(cfg)
	if err != nil {
		log.LOGGER("app").Error("Can't load default quote config: %s", err)
		return
	}

	fixss.StartWebServer(cfg)

	err = fixss.StartAcceptor(cfg)
	if err != nil {
		log.LOGGER("app").Error("Can't start acceptor ", err)
		return
	}
	interrupt := make(chan os.Signal)
	signal.Notify(interrupt, os.Interrupt, os.Kill)
	<-interrupt

	fixss.StopAcceptor()

}

func ParseFlags() (string, error) {
	var configPath string

	flag.StringVar(&configPath, "config", "./config/config.yml", "path to config file")
	flag.Parse()

	if !fixss.FileExists(configPath) {
		return "", errors.New(configPath + " not found")
	}

	return configPath, nil
}
