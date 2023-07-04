package store

import (
	"chatplus/store/vo"
	"encoding/json"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
)

type LevelDB struct {
	driver *leveldb.DB
}

func NewLevelDB() (*LevelDB, error) {
	db, err := leveldb.OpenFile("data/leveldb", nil)
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

func (db *LevelDB) Get(key string, value interface{}) error {
	bytes, err := db.driver.Get([]byte(key), nil)
	if err != nil {
		return err
	}

	return json.Unmarshal(bytes, &value)
}

func (db *LevelDB) Search(prefix string) []string {
	var items = make([]string, 0)
	iter := db.driver.NewIterator(util.BytesPrefix([]byte(prefix)), nil)
	defer iter.Release()

	for iter.Next() {
		items = append(items, string(iter.Value()))
	}
	return items
}

func (db *LevelDB) SearchPage(prefix string, page int, pageSize int) *vo.Page {
	var items = make([]string, 0)
	iter := db.driver.NewIterator(util.BytesPrefix([]byte(prefix)), nil)
	defer iter.Release()

	res := &vo.Page{Page: page, PageSize: pageSize}
	// 计算数据总数和总页数
	total := 0
	for iter.Next() {
		total++
	}
	res.TotalPage = (total + pageSize - 1) / pageSize
	res.Total = int64(total)

	// 计算目标页码的起始和结束位置
	start := (page - 1) * pageSize
	if start > total {
		return nil
	}
	end := start + pageSize
	if end > total {
		end = total
	}

	// 跳转到目标页码的起始位置
	count := 0
	for iter.Next() {
		if count >= start {
			items = append(items, string(iter.Value()))
		}
		count++
	}
	iter.Release()
	res.Items = items
	return res
}

func (db *LevelDB) Delete(key string) error {
	return db.driver.Delete([]byte(key), nil)
}

// Close release resources
func (db *LevelDB) Close() error {
	return db.driver.Close()
}
