package fixss

import (
	log "github.com/jeanphorn/log4go"
	"github.com/quickfixgo/enum"
	"github.com/quickfixgo/field"
	er "github.com/quickfixgo/fix44/executionreport"
	fix44mkd "github.com/quickfixgo/fix44/marketdatarequest"
	nos "github.com/quickfixgo/fix44/newordersingle"
	"github.com/quickfixgo/quickfix"
	"github.com/quickfixgo/tag"
	"github.com/shopspring/decimal"
	"runtime"
	"strconv"
	"sync/atomic"
	"time"
)

type execIdCounter struct {
	val int64
}

func (c *execIdCounter) getNext() string {
	atomic.AddInt64(&c.val, 1)
	runtime.Gosched()
	return strconv.FormatInt(atomic.LoadInt64(&c.val), 10)
}

var nextExecId = execIdCounter{}

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
	e.AddRoute(nos.Route(e.OnFIX44NewOrderSingle))

	return e
}

func (e *Executor) OnFIX44MarketDataRequest(msg fix44mkd.MarketDataRequest, sessionID quickfix.SessionID) quickfix.MessageRejectError {
	symbol, err := msg.GetString(tag.Symbol)
	if err == nil {
		SubscribeToSymbol(symbol, sessionID)
	}
	return nil
}

func (e *Executor) OnFIX44NewOrderSingle(msg nos.NewOrderSingle, id quickfix.SessionID) quickfix.MessageRejectError {
	processOrder(msg, id)
	return nil
}

func processOrder(msg nos.NewOrderSingle, id quickfix.SessionID) {
	orderId, _ := msg.GetClOrdID()
	symbol, _ := msg.GetSymbol()
	side, _ := msg.GetSide()
	price, _ := msg.GetPrice()
	orderQty, _ := msg.GetOrderQty()

	orderConfig := GetOrderConfig(symbol)

	log.Info("New order %s %s %s %s %s, strategy: ", orderId, symbol, side, price, orderQty, orderConfig.Strategy)

	switch orderConfig.Strategy {
	case Accept:
		tryToExecuteOrder(orderId, symbol, side, price, orderQty, id)
	case Reject:
		fallthrough
	default:
		sendRejectEr(orderId, symbol, side, id, "reject all orders")
	}

}

func tryToExecuteOrder(orderId string, symbol string, side enum.Side, price decimal.Decimal, qty decimal.Decimal, session quickfix.SessionID) {
	size, _ := qty.Float64()
	marketPrice := GetMarketPrice(symbol, side, size)
	if marketPrice == nil {
		log.Info("Reject order %s because there is no market price", orderId)
		sendRejectEr(orderId, symbol, side, session, "no market price")
		return
	}

	sendRejectEr(orderId, symbol, side, session, "not implement yet")
}

func sendRejectEr(orderId string, symbol string, side enum.Side, session quickfix.SessionID, reason string) {
	executionReport := er.New(
		field.NewOrderID(orderId),
		field.NewExecID(nextExecId.getNext()),
		field.NewExecType(enum.ExecType_REJECTED),
		field.NewOrdStatus(enum.OrdStatus_REJECTED),
		field.NewSide(side),
		field.NewLeavesQty(decimal.NewFromInt(0), 4),
		field.NewCumQty(decimal.NewFromInt(0), 4),
		field.NewAvgPx(decimal.NewFromInt(0), 4),
	)
	executionReport.SetClOrdID(orderId)
	executionReport.SetTransactTime(time.Now())
	executionReport.SetSymbol(symbol)
	executionReport.SetText(reason)

	err := quickfix.SendToTarget(executionReport, session)
	if err != nil {
		log.Error("Can't send Execution report %s", err.Error())
	}
}
