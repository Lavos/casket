package main

import (
	"log"
	// "github.com/Lavos/casket"
	"github.com/Lavos/casket/contentstorers"
	"github.com/Lavos/casket/filers"
)

func main () {
	// r := contentstorers.NewRedis("localhost", "casket", 6379)
	// f := filers.NewRedis("localhost", "casket", 6379)

	r, _ := contentstorers.NewBolt("content.bolt")
	f, _ := filers.NewBolt("files.bolt")

	file, err := f.Get("abc/xyz/index.html")

	log.Printf("file: %#v, %#v", file, err)

	if err != nil {
		log.Fatal(err)
	}

	sha, err := r.Put([]byte("abc12345"))

	err = file.AddRevision(sha)

	log.Printf("add revision error: %#v", err)

	log.Printf("file: %#v, err: %#v", file, err)
	log.Printf("number of revisions: %d", len(file.Revisions))

	file2, err := f.Get("abc/xyz/index.html")
	log.Printf("File2: %#v %#v", file2, err)
	log.Printf("number of revisions: %d", len(file2.Revisions))
}
