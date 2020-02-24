package main

import (
	"github.com/quickfixgo/quickfix"
	"os"
)

var acceptor *quickfix.Acceptor

func StartAcceptor() error {
	config, err := os.Open(CONFIG_PATH)
	if err != nil {
		return err
	}

	settings, err := quickfix.ParseSettings(config)
	if err != nil {
		return err
	}

	logFactory := NewFixLogFactory()
	executor := CreateExecutor()

	acceptor, err = quickfix.NewAcceptor(executor, quickfix.NewMemoryStoreFactory(), settings, logFactory)
	if err != nil {
		return err
	}

	err = acceptor.Start()
	if err != nil {
		return err
	}

	return nil
}

func StopAcceptor() {
	if acceptor != nil {
		acceptor.Stop()
	}
}