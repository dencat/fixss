package fixss

import (
	"fmt"
	"github.com/juju/loggo"
	"github.com/quickfixgo/quickfix"
)

var Log = loggo.GetLogger("")

var fixLogger = loggo.GetLogger("fix")

type fixLog struct {
}

func (l fixLog) OnIncoming(s []byte) {
	fixLogger.Debugf("incoming ", string(s))
}

func (l fixLog) OnOutgoing(s []byte) {
	fixLogger.Debugf("outgoing ", string(s))
}

func (l fixLog) OnEvent(s string) {
	fixLogger.Infof("event ", s)
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
