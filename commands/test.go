package main

import (
	"log"
	"github.com/Lavos/casket/contentstorers"
)

func main () {
	r := contentstorers.NewRedis("localhost", "casket", 6379)

	log.Printf("r: %#v", r)

	s, err := r.Put([]byte("blah"))

	log.Printf("%#v %#v", s, err)

	c, err := r.Get(s)

	log.Printf("%s %#v", c, err)
}
