package main

import (
	"log"
	"github.com/Lavos/casket/storers"
)

func main () {
	// r := contentstorers.NewRedis("localhost", "casket", 6379)
	// f := filers.NewRedis("localhost", "casket", 6379)

	b, _ := storers.NewBolt("one.bolt")

	// file, err := b.NewFile("abc/xyz/index.html", "text/html")
	file, err := b.GetFile("abc/xyz/index.html")

	log.Printf("file: %#v, %#v", file, err)

	if err != nil {
		log.Fatal(err)
	}

	sha, err := b.PutContent([]byte("abc12345"))

	err = file.AddRevision(sha)

	log.Printf("add revision error: %#v", err)

	log.Printf("file: %#v, err: %#v", file, err)
	log.Printf("number of revisions: %d", len(file.Revisions))

	file2, err := b.GetFile("abc/xyz/index.html")
	log.Printf("File2: %#v %#v", file2, err)
	log.Printf("number of revisions: %d", len(file2.Revisions))
}
