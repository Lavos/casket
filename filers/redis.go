package filers

import (
	"errors"
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

func (r *Redis) Put(f *casket.File) (error) {
	conn := r.RedisStorer.Pool.Get()
	defer conn.Close()

	conn.Do("HSET", r.RedisStorer.GetFullKey(f.Name), "content-type", f.ContentType)
	conn.Do("DEL", r.RedisStorer.GetFullKey(f.Name, "revisions"))

	args := redis.Args{}.Add(r.RedisStorer.GetFullKey(f.Name, "revisions"))
	shas := make([]string, len(f.Revisions))

	for i, rev := range(f.Revisions) {
		shas[i] = rev.String()
	}

	args.AddFlat(shas)
	conn.Do("LPUSH", args...)

	return nil
}

func (r *Redis) Get(name string) (*casket.File, error) {
	conn := r.RedisStorer.Pool.Get()
	defer conn.Close()

	contentType, content_err := redis.String(conn.Do("HGET", r.RedisStorer.GetFullKey(name), "content-type"))

	if content_err != nil {
		return nil, errors.New("File not found.")
	}

	revisions, _ := redis.Strings(conn.Do("LRANGE", r.RedisStorer.GetFullKey(name, "revisions"), 0, -1))

	shas := make([]casket.SHA1Sum, len(revisions))
	for i, s := range revisions {
		shas[i] = casket.NewSHA1SumFromString(s)
	}

	return &casket.File{
		Filer: r,

		Name: name,
		ContentType: contentType,
		Revisions: shas,
	}, nil
}

func (r *Redis) AddRevision(f *casket.File, s casket.SHA1Sum) error {
	conn := r.RedisStorer.Pool.Get()
	defer conn.Close()

	conn.Do("RPUSH", r.RedisStorer.GetFullKey(f.Name, "revisions"), s)
	f.Revisions = append(f.Revisions, s)
	return nil
}

func (r *Redis) NewFile(name, contentType string) (*casket.File, error) {
	exists, err := r.Exists(name)

	if exists || err != nil {
		return nil, errors.New("Exists error.")
	}

	file := &casket.File{
		Filer: r,

		Name: name,
		ContentType: contentType,
	}

	err = r.Put(file)
	return file, err
}

func (r *Redis) Exists(name string) (bool, error) {
	conn := r.RedisStorer.Pool.Get()
	defer conn.Close()

	return redis.Bool(conn.Do("EXISTS", r.RedisStorer.GetFullKey(name)))
}
