package contentstorers

import (
	"github.com/Lavos/casket"
	"github.com/Lavos/casket/storers"
	"github.com/garyburd/redigo/redis"
)

type Redis struct {
	storers.RedisStorer
}

func NewRedis (host, namespace string, port int64) *Redis {
	r := &Redis{}
	r.Init(host, namespace, port)
	return r
}

func (r *Redis) Put(content []byte) (casket.SHA1Sum, error) {
	conn := r.RedisStorer.Pool.Get()
	defer conn.Close()

	s := casket.NewSHA1Sum(content)

	_, err := conn.Do("SET", r.RedisStorer.GetFullKey(s.String()), content)
	return s, err
}

func (r *Redis) Get(s casket.SHA1Sum) ([]byte, error) {
	conn := r.RedisStorer.Pool.Get()
	defer conn.Close()

	return redis.Bytes(conn.Do("GET", r.RedisStorer.GetFullKey(s.String())))
}

func (r *Redis) Exists(s casket.SHA1Sum) (bool, error) {
	conn := r.RedisStorer.Pool.Get()
	defer conn.Close()

	return redis.Bool(conn.Do("EXISTS", r.RedisStorer.GetFullKey(s.String())))
}
