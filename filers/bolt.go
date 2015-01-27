package filers

import (
	"fmt"
	"bytes"
	"log"

	"github.com/Lavos/casket"
	"github.com/Lavos/casket/storers"
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
	return []byte(fmt.Sprintf("%s:%s", filename, property))
}

func (b *Bolt) Get(filename string) (*casket.File, error) {
	var file casket.File
	file.Filer = b
	file.Name = filename

	bolt_err := b.DB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(BucketFiles)

		if bucket == nil {
			return fmt.Errorf("Could not find bucket %s", BucketFiles)
		}

		file.ContentType = string(bucket.Get(b.CreateKey(filename, "content_type")))

		rev := bucket.Get(b.CreateKey(filename, "revisions"))

		if len(rev) >= 20 {
			count := len(rev) / 20
			file.Revisions = make([]casket.SHA1Sum, count)

			for x := 0; x < count; x++ {
				file.Revisions[x] = casket.NewSHA1SumFromBytes(rev[(x * 20):((x + 1) * 20)])
			}
		}

		return nil
	})

	if bolt_err != nil {
		return nil, bolt_err
	}

	return &file, nil
}

func (b *Bolt) Put(file *casket.File) error {
	return b.DB.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(BucketFiles)

		if err != nil {
			return err
		}

		err = bucket.Put(b.CreateKey(file.Name, "content_type"), []byte(file.ContentType))

		if err != nil {
			return err
		}

		var buf bytes.Buffer

		for _, rev := range file.Revisions {
			buf.Write(rev[:])
		}

		err = bucket.Put(b.CreateKey(file.Name, "revisions"), buf.Bytes())

		return err
	})
}

func (b *Bolt) NewFile(filename string, content_type string) (*casket.File, error) {
	var file casket.File
	file.Name = filename
	file.ContentType = content_type
	file.Filer = b
	file.Revisions = make([]casket.SHA1Sum, 0)

	err := b.Put(&file)

	if err != nil {
		return nil, err
	}

	return &file, nil
}

func (b *Bolt) AddRevision(file *casket.File, sha1sum casket.SHA1Sum) error {
	log.Printf("ADD REVISION %#v %#v", file, sha1sum)

	file.Revisions = append(file.Revisions, sha1sum)

	log.Printf("ADDED %#v", file)

	return b.DB.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(BucketFiles)

		if bucket == nil {
			return fmt.Errorf("Could not find bucket %s", BucketFiles)
		}

		buf := bucket.Get(b.CreateKey(file.Name, "revisions"))
		size := len(buf) + len(sha1sum)
		newbuf := make([]byte, len(buf), size)
		copy(newbuf, buf)
		newbuf = append(newbuf, sha1sum[:]...)

		return bucket.Put(b.CreateKey(file.Name, "revisions"), newbuf)
	})
}

func (b *Bolt) Exists(filename string) (bool, error) {
	var exists bool

	err := b.DB.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(BucketFiles)

		if bucket == nil {
			return fmt.Errorf("Could not find bucket %s", BucketFiles)
		}

		p := bucket.Get(b.CreateKey(filename, "content_type"))

		if len(p) > 0 {
			exists = true
		} else {
			exists = false
		}

		return nil
	})

	return exists, err
}
