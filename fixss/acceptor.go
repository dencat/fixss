package fixss

import (
	"github.com/dencat/fixss/fixss/store"
	"github.com/quickfixgo/quickfix"
	"os"
)

var acceptor *quickfix.Acceptor

func StartAcceptor(config *Config) error {
	configFile, err := os.Open(config.Fix.Config)
	if err != nil {
		return err
	}

	settings, err := quickfix.ParseSettings(configFile)
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
