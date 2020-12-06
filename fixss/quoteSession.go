package fixss

import (
	log "github.com/jeanphorn/log4go"
	"github.com/quickfixgo/enum"
	fix44mkdir "github.com/quickfixgo/fix44/marketdataincrementalrefresh"
	fix44mkfull "github.com/quickfixgo/fix44/marketdatasnapshotfullrefresh"
	"github.com/quickfixgo/quickfix"
	"github.com/shopspring/decimal"
	"math/rand"
	"strconv"
	"time"
)

type quoteSession struct {
	symbols map[string]bool
}

var sessions = map[quickfix.SessionID]quoteSession{}

func CreateQuoteSession(sessionID quickfix.SessionID) {
	rand.Seed(time.Now().UnixNano())
	log.Info("Create session %s", sessionID)
	sessions[sessionID] = quoteSession{
		map[string]bool{},
	}
}

func RemoveQuoteSession(sessionID quickfix.SessionID) {
	log.Info("Remove session %s", sessionID)
	delete(sessions, sessionID)
}

func SubscribeToSymbol(symbol string, sessionID quickfix.SessionID) {
	if session, ok := sessions[sessionID]; ok {
		if _, exists := session.symbols[symbol]; !exists {
			log.Info("Subscribe to %s", symbol)
			session.symbols[symbol] = true
			go startSendingMarketData(symbol, sessionID)
		}
	}
}

func startSendingMarketData(symbol string, id quickfix.SessionID) {
	log.Info("Start sending quotes %s to %s", symbol, id)
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
			time.Sleep(time.Duration(quoteConfig.Interval * 1000000))
		} else {
			time.Sleep(10000 * time.Millisecond)
		}

	}
	log.Info("Finish sending quotes %s to %s", symbol, id)
}

func sendMarketData(symbol string, config QuoteConfig, id quickfix.SessionID) {
	if len(config.Entities) == 0 {
		sendDropQuote(symbol, config, id)
		return
	}
	res := fix44mkdir.New()

	res.SetMDReqID(symbol + "0000")
	group := fix44mkdir.NewNoMDEntriesRepeatingGroup()
	quoteId := 0
	for _, entity := range config.Entities {
		entry := group.Add()
		entry.SetMDUpdateAction(enum.MDUpdateAction_CHANGE)
		entry.SetSymbol(symbol)
		entry.SetMDEntrySize(decimal.NewFromFloat(entity.Size), 4)
		price := entity.MinPrice + rand.Float64()*(entity.MaxPrice-entity.MinPrice)
		entry.SetMDEntryPx(decimal.NewFromFloat(price), 4)
		entry.SetMDEntryRefID(strconv.Itoa(quoteId))
		quoteId = quoteId + 1
		if entity.Direction == BID {
			entry.SetMDEntryType(enum.MDEntryType_BID)
		} else {
			entry.SetMDEntryType(enum.MDEntryType_OFFER)
		}
	}

	res.SetNoMDEntries(group)
	quickfix.SendToTarget(res.ToMessage(), id)
}

func sendDropQuote(symbol string, config QuoteConfig, id quickfix.SessionID) {
	res := fix44mkfull.New()
	res.SetMDReqID(symbol + "0000")
	res.SetSymbol(symbol)
	group := fix44mkfull.NewNoMDEntriesRepeatingGroup()
	res.SetNoMDEntries(group)
	res.SetInt(11010, 2)
	quickfix.SendToTarget(res.ToMessage(), id)
}
