package contentstorers

import (
	"github.com/lavos/casket"
	"github.com/lavos/casket/storers"
	"github.com/boltdb/bolt"
)

var (
	BucketRevisions = []byte("revisions")
}

type Bolt struct {
	storers.Bolt
}

func NewBolt(path string) (*Bolt, error) {
	b := &Bolt{}
	err := b.Init(path)
	return b, err
}

func (b *Bolt) Put(p []byte) (casket.SHA1Sum, error) {
	s := casket.NewSHA1Sum(content)

	return s, b.DB.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(BucketRevisions)

		if err != nil {
			return err
		}

		return bucket.Put(s[:], p)
	})
}

func (b *Bolt) Get(sha casket.SHA1Sum) ([]byte, error) {
	var p []byte
	bolt_err := b.DB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(BucketRevisions)

		if bucket == nil {
			return fmt.Errorf("Could not find bucket %s", BucketRevisions)
		}

		p = bucket.Get(sha[:])
		return nil
	})

	return p, bolt_err
}

func (b *Bolt) Exists(sha casket.SHA1Sum) (bool, error) {
	var exists bool

	bolt_err := b.DB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(BucketRevisions)

		if bucket == nil {
			return fmt.Errorf("Could not find bucket %s", BucketRevisions)
		}

		p := bucket.Get(sha[:])

		if p == nil {
			exists = false
		} else {
			exists = true
		}

		return nil
	})

	return exists, bolt_err
}
