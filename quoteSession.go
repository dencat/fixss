package main

import (
	"github.com/quickfixgo/enum"
	fix44mkdir "github.com/quickfixgo/fix44/marketdataincrementalrefresh"
	"github.com/quickfixgo/quickfix"
	"github.com/shopspring/decimal"
	"time"
)

type quoteSession struct {
	symbols map[string]bool
}

var sessions = map[quickfix.SessionID]quoteSession{}

func CreateQuoteSession(sessionID quickfix.SessionID) {
	Log.Infof("Create session %s", sessionID)
	sessions[sessionID] = quoteSession{
		map[string]bool{},
	}
}

func RemoveQuoteSession(sessionID quickfix.SessionID) {
	Log.Infof("Remove session %s", sessionID)
	delete(sessions, sessionID)
}

func SubscribeToSymbol(symbol string, sessionID quickfix.SessionID) {
	if session, ok := sessions[sessionID]; ok {
		if _, exists := session.symbols[symbol]; !exists {
			Log.Infof("Subscribe to %s", symbol)
			session.symbols[symbol] = true
			go startSendingMarketData(symbol, sessionID)
		}
	}
}

func startSendingMarketData(symbol string, id quickfix.SessionID) {
	Log.Infof("Start sending quotes %s to %s", symbol, id)
	for {
		session, ok := sessions[id]
		if !ok {
			break
		}
		_, ok = session.symbols[symbol]
		if !ok {
			break
		}
		quoteConfig := GetQuoteConfig(symbol)
		if quoteConfig != nil {
			sendMarketData(symbol, *quoteConfig, id)
			time.Sleep(time.Duration(quoteConfig.interval * 1000000))
		} else {
			time.Sleep(10000 * time.Millisecond)
		}

	}
	Log.Infof("Finish sending quotes %s to %s", symbol, id)
}

func sendMarketData(symbol string, config quoteConfig, id quickfix.SessionID) {
	res := fix44mkdir.New()

	res.SetMDReqID(symbol + "0000")
	group := fix44mkdir.NewNoMDEntriesRepeatingGroup()

	for _, entity := range config.entities {
		entry := group.Add()
		entry.SetMDUpdateAction(enum.MDUpdateAction_CHANGE)
		entry.SetSymbol(symbol)
		entry.SetMDEntrySize(decimal.NewFromFloat(entity.size), 4)
		entry.SetMDEntryPx(decimal.NewFromFloat(entity.price), 4)
		if entity.direction == BID {
			entry.SetMDEntryType(enum.MDEntryType_BID)
		} else {
			entry.SetMDEntryType(enum.MDEntryType_OFFER)
		}
	}

	res.SetNoMDEntries(group)
	quickfix.SendToTarget(res.ToMessage(), id)
}
