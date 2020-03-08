package tests

import (
	"fmt"
	"github.com/quickfixgo/quickfix"
	"os"
)

type TradeClient struct {
	loginDone chan bool
}

func (e TradeClient) OnCreate(sessionID quickfix.SessionID) {
	return
}

func (e TradeClient) OnLogon(sessionID quickfix.SessionID) {
	e.loginDone <- true
}

func (e TradeClient) OnLogout(sessionID quickfix.SessionID) {
	return
}

func (e TradeClient) FromAdmin(msg *quickfix.Message, sessionID quickfix.SessionID) (reject quickfix.MessageRejectError) {
	return
}

func (e TradeClient) ToAdmin(msg *quickfix.Message, sessionID quickfix.SessionID) {
	return
}

func (e TradeClient) ToApp(msg *quickfix.Message, sessionID quickfix.SessionID) (err error) {
	fmt.Printf("Sending %s\n", msg)
	return
}

func (e TradeClient) FromApp(msg *quickfix.Message, sessionID quickfix.SessionID) (reject quickfix.MessageRejectError) {
	fmt.Printf("FromApp: %s\n", msg.String())
	return
}

func CreateInitiator(loginDone chan bool) (*quickfix.Initiator, *TradeClient, error) {
	cfg, err := os.Open("config/client.cfg")
	if err != nil {
		return nil, nil, err
	}

	appSettings, err := quickfix.ParseSettings(cfg)
	if err != nil {
		return nil, nil, err
	}

	app := TradeClient{
		loginDone: loginDone,
	}

	fileLogFactory, err := quickfix.NewFileLogFactory(appSettings)

	if err != nil {
		return nil, nil, err
	}

	initiator, err := quickfix.NewInitiator(app, quickfix.NewMemoryStoreFactory(), appSettings, fileLogFactory)
	return initiator, &app, err

}
