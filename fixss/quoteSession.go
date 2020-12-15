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
	log.LOGGER("app").Info("Create session %s", sessionID)
	sessions[sessionID] = quoteSession{
		map[string]bool{},
	}
}

func RemoveQuoteSession(sessionID quickfix.SessionID) {
	log.LOGGER(CAT_APP).Info("Remove session %s", sessionID)
	delete(sessions, sessionID)
}

func SubscribeToSymbol(symbol string, sessionID quickfix.SessionID) {
	if session, ok := sessions[sessionID]; ok {
		if _, exists := session.symbols[symbol]; !exists {
			log.LOGGER(CAT_APP).Info("Subscribe to %s, session: %s", symbol, sessionID)
			session.symbols[symbol] = true
			go startSendingMarketData(symbol, sessionID)
		}
	}
}

func startSendingMarketData(symbol string, id quickfix.SessionID) {
	log.LOGGER(CAT_APP).Info("Start sending quotes %s to %s", symbol, id)
	for {
		session, ok := sessions[id]
		if !ok {
			break
		}
		_, ok = session.symbols[symbol]
		if !ok {
			break
		}
		quotes, interval := GetNextQuotes(symbol)
		if quotes != nil {
			sendMarketData(symbol, quotes, id)
			time.Sleep(time.Duration(interval * 1000000))
		} else {
			time.Sleep(10000 * time.Millisecond)
		}

	}
	log.LOGGER(CAT_APP).Info("Finish sending quotes %s to %s", symbol, id)
}

func sendMarketData(symbol string, quotes []Quote, sessionId quickfix.SessionID) {
	if len(quotes) == 0 {
		sendDropQuote(symbol, sessionId)
		return
	}
	res := fix44mkdir.New()

	res.SetMDReqID(symbol + "0000")
	group := fix44mkdir.NewNoMDEntriesRepeatingGroup()
	quoteId := 0
	for _, quote := range quotes {
		entry := group.Add()
		entry.SetMDUpdateAction(enum.MDUpdateAction_CHANGE)
		entry.SetSymbol(symbol)
		entry.SetMDEntrySize(decimal.NewFromFloat(quote.Size), 4)
		entry.SetMDEntryPx(quote.price, 4)
		entry.SetMDEntryRefID(strconv.Itoa(quoteId))
		quoteId = quoteId + 1
		if quote.Direction == BID {
			entry.SetMDEntryType(enum.MDEntryType_BID)
		} else {
			entry.SetMDEntryType(enum.MDEntryType_OFFER)
		}
	}

	res.SetNoMDEntries(group)
	quickfix.SendToTarget(res.ToMessage(), sessionId)
}

func sendDropQuote(symbol string, id quickfix.SessionID) {
	res := fix44mkfull.New()
	res.SetMDReqID(symbol + "0000")
	res.SetSymbol(symbol)
	group := fix44mkfull.NewNoMDEntriesRepeatingGroup()
	res.SetNoMDEntries(group)
	res.SetInt(11010, 2)
	quickfix.SendToTarget(res.ToMessage(), id)
}
