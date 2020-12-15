package fixss

import (
	"fmt"
	log "github.com/jeanphorn/log4go"
	"github.com/quickfixgo/quickfix"
)

const CAT_APP = "app"

type fixLog struct {
	qualifier string
}

func (l fixLog) OnIncoming(s []byte) {
	log.LOGGER("incoming " + l.qualifier).Info(string(s))
}

func (l fixLog) OnOutgoing(s []byte) {
	log.LOGGER("outgoing " + l.qualifier).Info(string(s))
}

func (l fixLog) OnEvent(s string) {
	log.LOGGER("event " + l.qualifier).Info(s)
}

func (l fixLog) OnEventf(format string, a ...interface{}) {
	l.OnEvent(fmt.Sprintf(format, a...))
}

type fxLogFactory struct{}

func (fxLogFactory) Create() (quickfix.Log, error) {
	return fixLog{"default"}, nil
}

func (fxLogFactory) CreateSessionLog(sessionID quickfix.SessionID) (quickfix.Log, error) {
	return fixLog{sessionID.Qualifier}, nil
}

func NewFixLogFactory() quickfix.LogFactory {
	return fxLogFactory{}
}
