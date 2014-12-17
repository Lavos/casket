package main

import (
	"log"
	// "github.com/Lavos/casket"
	"github.com/Lavos/casket/contentstorers"
	"github.com/Lavos/casket/filers"
)

func main () {
	r := contentstorers.NewRedis("localhost", "casket", 6379)
	f := filers.NewRedis("localhost", "casket", 6379)

	// file, err := f.NewFile("abc/xyz/index.html", "text/html")
	file, err := f.Get("abc/xyz/index.html")

	if err != nil {
		log.Fatal(err)
	}

	sha, err := r.Put([]byte("12345"))

	file.AddRevision(sha)
	log.Printf("file: %#v, err: %#v", file, err)
}
