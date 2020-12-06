package fixss

import (
	"github.com/dencat/fixss/fixss/store"
	"github.com/quickfixgo/quickfix"
	"os"
)

const CONFIG_PATH = "config/server.cfg"

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

	acceptor, err = quickfix.NewAcceptor(executor, store.NewFixMemoryStoreFactory(), settings, logFactory)
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
