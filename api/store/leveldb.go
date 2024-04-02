package store

import (
	"encoding/json"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
)

type LevelDB struct {
	driver *leveldb.DB
}

func NewLevelDB() (*LevelDB, error) {
	db, err := leveldb.OpenFile("data", nil)
	if err != nil {
		return nil, err
	}
	return &LevelDB{
		driver: db,
	}, nil
}

func (db *LevelDB) Put(key string, value interface{}) error {
	var byteData []byte
	if v, ok := value.(string); ok {
		byteData = []byte(v)
	} else {
		b, err := json.Marshal(value)
		if err != nil {
			return err
		}
		byteData = b
	}
	return db.driver.Put([]byte(key), byteData, nil)
}

func (db *LevelDB) Get(key string) ([]byte, error) {
	bytes, err := db.driver.Get([]byte(key), nil)
	if err != nil {
		return nil, err
	}

	return bytes, nil
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

type PageVo struct {
	Items     []string
	Page      int
	PageSize  int
	Total     int
	TotalPage int
}

func (db *LevelDB) SearchPage(prefix string, page int, pageSize int) *PageVo {
	var items = make([]string, 0)
	iter := db.driver.NewIterator(util.BytesPrefix([]byte(prefix)), nil)
	defer iter.Release()

	res := &PageVo{Page: page, PageSize: pageSize}
	// 计算数据总数和总页数
	total := 0
	for iter.Next() {
		total++
	}
	res.TotalPage = (total + pageSize - 1) / pageSize
	res.Total = total

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
