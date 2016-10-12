package main

import (
	"encoding/json"
	"github.com/syndtr/goleveldb/leveldb"
	"time"
)

type CachedData struct {
	Path    string
	ModTime time.Time
	Owner   string
}

func (c *CachedData) Marshal() ([]byte, error) {
	return json.Marshal(c)
}

func (c *CachedData) Unmarshal(data []byte) error {
	return json.Unmarshal(data, &c)
}

func OpenDB(path string) (*leveldb.DB, error) {
	db, err := leveldb.OpenFile(DB_FILE, nil)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func GetDB(k string) *CachedData {
	data, err := db.Get([]byte(k), nil)
	if err != nil {
		return nil
	}

	c := &CachedData{}
	err = c.Unmarshal(data)
	if err != nil {
		return nil
	}

	return c
}

func PutDB(k string, c *CachedData) error {
	b, err := c.Marshal()
	if err != nil {
		return err
	}

	return db.Put([]byte(k), b, nil)
}
