package storers

import (
	"time"
	"fmt"
	"strings"
	"github.com/garyburd/redigo/redis"
	"errors"
	"github.com/Lavos/casket"
)

type Redis struct {
	Pool      *redis.Pool
	Namespace string
}

func NewRedis(host, namespace string, port int64) *Redis {
	pool := &redis.Pool{
		MaxIdle:     5,
		MaxActive:   1000,
		IdleTimeout: 5 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", fmt.Sprintf("%s:%d", host, port))

			if err != nil {
				return nil, err
			}

			return c, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}

	return &Redis{pool, namespace}
}

func (r *Redis) GetFullKey(args ...string) string {
	return fmt.Sprintf("%s:%s", r.Namespace, strings.Join(args, ":"))
}

func (r *Redis) GetSubKey(key string) string {
	return strings.TrimPrefix(key, fmt.Sprintf("%s:", r.Namespace))
}

func (r *Redis) PutContent(content []byte) (casket.SHA1Sum, error) {
	conn := r.Pool.Get()
	defer conn.Close()

	s := casket.NewSHA1Sum(content)

	_, err := conn.Do("SET", r.GetFullKey(s.String()), content)
	return s, err
}

func (r *Redis) GetContent(s casket.SHA1Sum) ([]byte, error) {
	conn := r.Pool.Get()
	defer conn.Close()

	return redis.Bytes(conn.Do("GET", r.GetFullKey(s.String())))
}

func (r *Redis) ContentExists(s casket.SHA1Sum) (bool, error) {
	conn := r.Pool.Get()
	defer conn.Close()

	return redis.Bool(conn.Do("EXISTS", r.GetFullKey(s.String())))
}

func (r *Redis) PutFile(f *casket.File) (error) {
	conn := r.Pool.Get()
	defer conn.Close()

	conn.Do("HSET", r.GetFullKey(f.Name), "content-type", f.ContentType)
	conn.Do("DEL", r.GetFullKey(f.Name, "revisions"))

	args := redis.Args{}.Add(r.GetFullKey(f.Name, "revisions"))
	shas := make([]string, len(f.Revisions))

	for i, rev := range(f.Revisions) {
		shas[i] = rev.String()
	}

	args.AddFlat(shas)
	conn.Do("LPUSH", args...)

	return nil
}

func (r *Redis) GetFile(name string) (*casket.File, error) {
	conn := r.Pool.Get()
	defer conn.Close()

	contentType, content_err := redis.String(conn.Do("HGET", r.GetFullKey(name), "content-type"))

	if content_err != nil {
		return nil, errors.New("File not found.")
	}

	revisions, _ := redis.Strings(conn.Do("LRANGE", r.GetFullKey(name, "revisions"), 0, -1))

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
	conn := r.Pool.Get()
	defer conn.Close()

	conn.Do("RPUSH", r.GetFullKey(f.Name, "revisions"), s)
	f.Revisions = append(f.Revisions, s)
	return nil
}

func (r *Redis) NewFile(name, contentType string) (*casket.File, error) {
	exists, err := r.FileExists(name)

	if exists || err != nil {
		return nil, errors.New("Exists error.")
	}

	file := &casket.File{
		Filer: r,

		Name: name,
		ContentType: contentType,
	}

	err = r.PutFile(file)
	return file, err
}

func (r *Redis) FileExists(name string) (bool, error) {
	conn := r.Pool.Get()
	defer conn.Close()

	return redis.Bool(conn.Do("EXISTS", r.GetFullKey(name)))
}

