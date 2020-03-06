package fixss

import (
	fix44mkd "github.com/quickfixgo/fix44/marketdatarequest"
	"github.com/quickfixgo/quickfix"
	"github.com/quickfixgo/tag"
)

type Executor struct {
	*quickfix.MessageRouter
}

func (e *Executor) OnCreate(sessionID quickfix.SessionID) {
	return
}
func (e *Executor) OnLogon(sessionID quickfix.SessionID) {
	CreateQuoteSession(sessionID)
}
func (e *Executor) OnLogout(sessionID quickfix.SessionID) {
	RemoveQuoteSession(sessionID)
}
func (e *Executor) ToAdmin(msg *quickfix.Message, sessionID quickfix.SessionID) {
	return
}
func (e *Executor) ToApp(msg *quickfix.Message, sessionID quickfix.SessionID) error {
	return nil
}
func (e *Executor) FromApp(message *quickfix.Message, sessionID quickfix.SessionID) quickfix.MessageRejectError {
	return e.Route(message, sessionID)
}
func (e *Executor) FromAdmin(msg *quickfix.Message, sessionID quickfix.SessionID) quickfix.MessageRejectError {
	return nil
}

func CreateExecutor() *Executor {
	e := &Executor{MessageRouter: quickfix.NewMessageRouter()}
	e.AddRoute(fix44mkd.Route(e.OnFIX44MarketDataRequest))

	return e
}

func (e *Executor) OnFIX44MarketDataRequest(msg fix44mkd.MarketDataRequest, sessionID quickfix.SessionID) quickfix.MessageRejectError {
	symbol, err := msg.GetString(tag.Symbol)
	if err == nil {
		SubscribeToSymbol(symbol, sessionID)
	}
	return nil
}
