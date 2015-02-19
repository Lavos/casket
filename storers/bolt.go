package storers

import (
	"fmt"
	"github.com/Lavos/casket"
	"github.com/boltdb/bolt"
	"bytes"
	"log"
)

var (
	BucketRevisions = []byte("revisions")
	BucketFiles = []byte("files")
)

type Bolt struct {
	DB *bolt.DB
}

func NewBolt(path string) (*Bolt, error) {
	db, err := bolt.Open(path, 0666, nil)

	if err != nil {
		return nil, err
	}

	b := &Bolt{db}
	return b, err
}

func (b *Bolt) PutContent(p []byte) (casket.SHA1Sum, error) {
	s := casket.NewSHA1Sum(p)

	return s, b.DB.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(BucketRevisions)

		if err != nil {
			return err
		}

		return bucket.Put(s[:], p)
	})
}

func (b *Bolt) GetContent(sha casket.SHA1Sum) ([]byte, error) {
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

func (b *Bolt) ContentExists(sha casket.SHA1Sum) (bool, error) {
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

func (b *Bolt) CreateKey(filename, property string) []byte {
	return []byte(fmt.Sprintf("%s:%s", filename, property))
}

func (b *Bolt) GetFile(filename string) (*casket.File, error) {
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

func (b *Bolt) PutFile(file *casket.File) error {
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

	err := b.PutFile(&file)

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

func (b *Bolt) FileExists(filename string) (bool, error) {
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
