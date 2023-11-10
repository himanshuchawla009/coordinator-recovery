package service

import (
	"github.com/dgraph-io/badger/v3"
)

type Data struct {
	Key   string
	Value []byte
}

type DatabaseService interface {
	Get(key string) ([]byte, error)
	Set(key string, value []byte) error
	GetFromPrefix(prefix string) ([]Data, error)
}

type database struct {
	db *badger.DB
}

func NewDatabaseService(db *badger.DB) DatabaseService {
	return &database{db: db}
}

func (d *database) Get(key string) ([]byte, error) {
	var value []byte
	err := d.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}
		err = item.Value(func(val []byte) error {
			value = val
			return nil
		})
		return nil
	})
	return value, err
}

func (d *database) Set(key string, value []byte) error {
	return d.db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(key), value)
	})
}

func (d *database) GetFromPrefix(prefix string) ([]Data, error) {
	var res []Data
	pre := []byte(prefix)
	err := d.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		for it.Seek(pre); it.ValidForPrefix(pre); it.Next() {
			item := it.Item()
			k := item.Key()
			err := item.Value(func(v []byte) error {
				value := make([]byte, len(v))
				copy(value, v)
				data := Data{
					Key:   string(k),
					Value: value,
				}
				res = append(res, data)
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}
