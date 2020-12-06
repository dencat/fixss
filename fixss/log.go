package fixss

import (
	"fmt"
	log "github.com/jeanphorn/log4go"
	"github.com/quickfixgo/quickfix"
)

type fixLog struct {
}

func (l fixLog) OnIncoming(s []byte) {
	log.LOGGER("incoming").Info(string(s))
}

func (l fixLog) OnOutgoing(s []byte) {
	log.LOGGER("outgoing").Info(string(s))
}

func (l fixLog) OnEvent(s string) {
	log.LOGGER("event").Info(s)
}

func (l fixLog) OnEventf(format string, a ...interface{}) {
	l.OnEvent(fmt.Sprintf(format, a...))
}

type fxLogFactory struct{}

func (fxLogFactory) Create() (quickfix.Log, error) {
	return fixLog{}, nil
}

func (fxLogFactory) CreateSessionLog(sessionID quickfix.SessionID) (quickfix.Log, error) {
	return fixLog{}, nil
}

func NewFixLogFactory() quickfix.LogFactory {
	return fxLogFactory{}
}
