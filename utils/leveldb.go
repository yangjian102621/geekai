package utils

import (
	"encoding/json"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
)

type LevelDB struct {
	driver *leveldb.DB
}

func NewLevelDB(path string) (*LevelDB, error) {
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return nil, err
	}
	return &LevelDB{
		driver: db,
	}, nil
}

func (db *LevelDB) Put(key string, value interface{}) error {
	bytes, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return db.driver.Put([]byte(key), bytes, nil)
}

func (db *LevelDB) Get(key string) (interface{}, error) {
	bytes, err := db.driver.Get([]byte(key), nil)
	if err != nil {
		return nil, err
	}

	var value interface{}
	err = json.Unmarshal(bytes, &value)
	if err != nil {
		return nil, err
	}

	return value, nil
}

func (db *LevelDB) Search(prefix string) []string {
	var items = make([]string, 0)
	iter := db.driver.NewIterator(util.BytesPrefix([]byte(prefix)), nil)
	for iter.Next() {
		items = append(items, string(iter.Value()))
	}
	iter.Release()
	return items
}

func (db *LevelDB) Delete(key string) error {
	return db.driver.Delete([]byte(key), nil)
}

// Close release resources
func (db *LevelDB) Close() error {
	return db.driver.Close()
}
