package store

import (
	"github.com/quickfixgo/quickfix"
	"time"
)

type fixMemoryStore struct {
	senderMsgSeqNum, targetMsgSeqNum int
	creationTime                     time.Time
}

type MessageStoreFactory interface {
	Create(sessionID quickfix.SessionID) (quickfix.MessageStore, error)
}

func (store *fixMemoryStore) NextSenderMsgSeqNum() int {
	return store.senderMsgSeqNum + 1
}

func (store *fixMemoryStore) NextTargetMsgSeqNum() int {
	return store.targetMsgSeqNum + 1
}

func (store *fixMemoryStore) IncrNextSenderMsgSeqNum() error {
	store.senderMsgSeqNum++
	return nil
}

func (store *fixMemoryStore) IncrNextTargetMsgSeqNum() error {
	store.targetMsgSeqNum++
	return nil
}

func (store *fixMemoryStore) SetNextSenderMsgSeqNum(nextSeqNum int) error {
	store.senderMsgSeqNum = nextSeqNum - 1
	return nil
}
func (store *fixMemoryStore) SetNextTargetMsgSeqNum(nextSeqNum int) error {
	store.targetMsgSeqNum = nextSeqNum - 1
	return nil
}

func (store *fixMemoryStore) CreationTime() time.Time {
	return store.creationTime
}

func (store *fixMemoryStore) Reset() error {
	store.senderMsgSeqNum = 0
	store.targetMsgSeqNum = 0
	store.creationTime = time.Now()
	return nil
}

func (store *fixMemoryStore) Refresh() error {
	return nil
}

func (store *fixMemoryStore) Close() error {
	return nil
}

func (store *fixMemoryStore) SaveMessage(seqNum int, msg []byte) error {
	return nil
}

func (store *fixMemoryStore) GetMessages(beginSeqNum, endSeqNum int) ([][]byte, error) {
	return nil, nil
}

type fixMemoryStoreFactory struct{}

func (f fixMemoryStoreFactory) Create(sessionID quickfix.SessionID) (quickfix.MessageStore, error) {
	m := new(fixMemoryStore)
	m.Reset()
	return m, nil
}

func NewFixMemoryStoreFactory() MessageStoreFactory { return fixMemoryStoreFactory{} }
