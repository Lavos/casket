package main

import (
	"github.com/Lavos/casket/examples"
	"github.com/Lavos/casket/contentstorers"
	"github.com/Lavos/casket/filers"
)

func main () {
	r := contentstorers.NewRedis("localhost", "casket", 6379)
	f := filers.NewRedis("localhost", "casket", 6379)
	s := examples.NewServer(r, f, ":8035")
	s.Run()
}
