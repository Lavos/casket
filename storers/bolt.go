package storer

import (
	"github.com/boltdb/bolt"
)

type Bolt struct {
	DB *bolt.DB
}

func (b *Bolt) Init(path string) error {
	db, err := bolt.Open(path, 0666, nil)

	if err != nil {
		return err
	}

	b.DB = db
	return nil
}
