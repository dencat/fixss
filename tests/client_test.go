package tests

import (
	"fmt"
	"github.com/quickfixgo/enum"
	"github.com/quickfixgo/field"
	er "github.com/quickfixgo/fix44/executionreport"
	fix44mkdir "github.com/quickfixgo/fix44/marketdataincrementalrefresh"
	fix44mdr "github.com/quickfixgo/fix44/marketdatarequest"
	fix44full "github.com/quickfixgo/fix44/marketdatasnapshotfullrefresh"
	nos "github.com/quickfixgo/fix44/newordersingle"
	"github.com/quickfixgo/quickfix"
	"github.com/shopspring/decimal"
	"os"
	"sync"
	"time"
)

var quotes = map[string]string{}
var orderStatuses = map[string]enum.OrdStatus{}
var someMapMutex = sync.RWMutex{}
var mux = sync.Mutex{}

type TradeClient struct {
	*quickfix.MessageRouter
	loginDone chan bool
	sessionID *quickfix.SessionID
}

func (e *TradeClient) OnCreate(sessionID quickfix.SessionID) {
	return
}

func (e *TradeClient) OnLogon(sessionID quickfix.SessionID) {
	e.sessionID = &sessionID
	e.loginDone <- true
}

func (e *TradeClient) OnLogout(sessionID quickfix.SessionID) {
	return
}

func (e *TradeClient) FromAdmin(msg *quickfix.Message, sessionID quickfix.SessionID) (reject quickfix.MessageRejectError) {
	return
}

func (e *TradeClient) ToAdmin(msg *quickfix.Message, sessionID quickfix.SessionID) {
	return
}

func (e *TradeClient) ToApp(msg *quickfix.Message, sessionID quickfix.SessionID) (err error) {
	fmt.Printf("Sending %s\n", msg)
	return
}

func (e *TradeClient) FromApp(msg *quickfix.Message, sessionID quickfix.SessionID) (reject quickfix.MessageRejectError) {
	return e.Route(msg, sessionID)
}

func (e *TradeClient) SendMarketDataRequest(symbol string) {
	quickfix.SendToTarget(queryMarketDataRequest(symbol), *e.sessionID)
}

func (e *TradeClient) SendOrder(orderId string, symbol string, orderQty decimal.Decimal, price decimal.Decimal, side enum.Side) {
	order := nos.New(
		field.NewClOrdID(orderId),
		field.NewSide(side),
		field.NewTransactTime(time.Now()),
		field.NewOrdType(enum.OrdType_LIMIT),
	)
	order.SetOrderQty(orderQty, 4)
	order.SetPrice(price, 4)
	order.SetSymbol(symbol)
	order.SetTimeInForce(enum.TimeInForce_IMMEDIATE_OR_CANCEL)
	quickfix.SendToTarget(order, *e.sessionID)
}

func (e *TradeClient) OnMarketDataIncrementalRefresh(msg fix44mkdir.MarketDataIncrementalRefresh, id quickfix.SessionID) quickfix.MessageRejectError {
	mux.Lock()
	entries, _ := msg.GetNoMDEntries()
	for i := 0; i < entries.Len(); i++ {
		px, _ := entries.Get(i).GetMDEntryPx()
		qty, _ := entries.Get(i).GetMDEntrySize()
		quoteType, _ := entries.Get(i).GetMDEntryType()
		symbol, _ := entries.Get(i).GetSymbol()
		key := symbol + "_" + string(quoteType) + "_" + qty.String()
		quotes[key] = px.String()
	}
	defer mux.Unlock()
	return nil
}

func (e *TradeClient) OnMarketDataSnapshotFullRefresh(msg fix44full.MarketDataSnapshotFullRefresh, id quickfix.SessionID) quickfix.MessageRejectError {
	mux.Lock()
	defer mux.Unlock()

	symbol, _ := msg.GetSymbol()
	groups, _ := msg.GetNoMDEntries()
	if groups.Len() == 0 {
		quotes = map[string]string{}
		quotes[symbol] = "drop"
	}

	return nil
}

func (e *TradeClient) GetLastQuote(key string) string {
	mux.Lock()
	defer mux.Unlock()
	return quotes[key]
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

	app := &TradeClient{
		MessageRouter: quickfix.NewMessageRouter(),
		loginDone:     loginDone,
	}
	app.AddRoute(fix44mkdir.Route(app.OnMarketDataIncrementalRefresh))
	app.AddRoute(fix44full.Route(app.OnMarketDataSnapshotFullRefresh))
	app.AddRoute(er.Route(app.OnExecutionReport))

	fileLogFactory, err := quickfix.NewFileLogFactory(appSettings)

	if err != nil {
		return nil, nil, err
	}

	initiator, err := quickfix.NewInitiator(app, quickfix.NewMemoryStoreFactory(), appSettings, fileLogFactory)
	return initiator, app, err

}

func queryMarketDataRequest(symbol string) fix44mdr.MarketDataRequest {
	request := fix44mdr.New(field.NewMDReqID("MARKETDATAID"),
		field.NewSubscriptionRequestType(enum.SubscriptionRequestType_SNAPSHOT),
		field.NewMarketDepth(0),
	)

	entryTypes := fix44mdr.NewNoMDEntryTypesRepeatingGroup()
	entryTypes.Add().SetMDEntryType(enum.MDEntryType_TRADE)
	request.SetNoMDEntryTypes(entryTypes)

	relatedSym := fix44mdr.NewNoRelatedSymRepeatingGroup()
	relatedSym.Add().SetSymbol(symbol)
	request.SetNoRelatedSym(relatedSym)

	return request
}

func (e *TradeClient) OnExecutionReport(msg er.ExecutionReport, id quickfix.SessionID) quickfix.MessageRejectError {
	orderId, _ := msg.GetOrderID()
	status, _ := msg.GetOrdStatus()
	someMapMutex.Lock()
	defer someMapMutex.Unlock()
	orderStatuses[orderId] = status

	return nil
}

func (e *TradeClient) GetOrderStatus(key string) enum.OrdStatus {
	someMapMutex.RLock()
	defer someMapMutex.RUnlock()
	if res, ok := orderStatuses[key]; ok {
		return res
	}
	return enum.OrdStatus_NEW
}
