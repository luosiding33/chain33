package db

import (
	"fmt"
)

type KVDB interface {
	Get(key []byte) (value []byte, err error)
	Set(key []byte, value []byte) (err error)
}

type DB interface {
	Get([]byte) []byte
	Set([]byte, []byte)
	SetSync([]byte, []byte)
	Delete([]byte)
	DeleteSync([]byte)
	Close()
	NewBatch(sync bool) Batch

	// For debugging
	Print()
	Iterator() Iterator
	Stats() map[string]string
	PrefixScan(key []byte) [][]byte
	IteratorScan(key []byte, count int32, direction int32) [][]byte
	IteratorScanFromLast(key []byte, count int32, direction int32) [][]byte
}

type Batch interface {
	Set(key, value []byte)
	Delete(key []byte)
	Write()
}

type Iterator interface {
	Next() bool

	Key() []byte
	Value() []byte
}

//-----------------------------------------------------------------------------

const (
	LevelDBBackendStr   = "leveldb" // legacy, defaults to goleveldb.
	CLevelDBBackendStr  = "cleveldb"
	GoLevelDBBackendStr = "goleveldb"
	MemDBBackendStr     = "memdb"
)

type dbCreator func(name string, dir string) (DB, error)

var backends = map[string]dbCreator{}

func registerDBCreator(backend string, creator dbCreator, force bool) {
	_, ok := backends[backend]
	if !force && ok {
		return
	}
	backends[backend] = creator
}

func NewDB(name string, backend string, dir string) DB {
	db, err := backends[backend](name, dir)
	if err != nil {
		fmt.Println("Error initializing DB: %v", err)
	}
	return db
}
