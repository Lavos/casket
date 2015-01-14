package filers

import (
	"fmt"
	"bytes"

	"github.com/lavos/casket"
	"github.com/lavos/casket/storers"
	"github.com/boltdb/bolt"
)

var (
	BucketFiles = []byte("files")
)

type Bolt struct {
	storers.Bolt
}

func NewBolt(path string) (*Bolt, error) {
	b := &Bolt{}
	err := b.Init(path)
	return b, err
}

func (b *Bolt) CreateKey(filename, property string) []byte {
	return fmt.Sprintf("%s:%s", filename, property)
}

func (b *Bolt) Get(filename string) (*casket.File, error) {
	var file Casket.File
	file.Filer = b
	file.Name = filename

	bolt_err := b.DB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(BucketFiles)

		if bucket == nil {
			return fmt.Errorf("Could not find bucket %s", BucketFiles)
		}

		file.ContentType = bucket.Get(b.CreateKey(filename, "content_type"))

		buf := bytes.NewBuffer(bucket.Get(b.CreateKey(filename, "revisions")))

		count := len(buf) / 20
		file.Revisions = make([]SHA1Sum, count)

		for x := 0; x < count; x++ {
			file.Revisions[x] = casket.NewSHA1SumFromBytes(buf[(x * 20):((x + 1) * 20)])
		}

		return nil
	})

	if bolt_err == nil {
		return nil, bolt_err
	}

	return file, nil
}

func (b *Bolt) Put(file *casket.File) error {
	return b.DB.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(BucketFiles)

		if bucket == nil {
			return fmt.Errorf("Could not find bucket %s", BucketFiles)
		}

		err := bucket.Put(b.CreateKey(file.Name, "content_type"), []byte(file.ContentType))

		if err != nil {
			return err
		}

		var buf bytes.Buffer

		for i, rev := range file.Revisions {
			buf.Write(rev[:])
		}

		err = bucket.Put(b.CreateKey(file.Name, "revisions"), buf.Bytes())

		return err
	})
}

func (b *Bolt) NewFile(filename string, content_type string) (*casket.File, error) {

}

func (b *Bolt) AddRevision(file *casket.File, casket.SHA1Sum) error {

}

func (b *Bolt) Exists(filename string) (bool, error) {


}
