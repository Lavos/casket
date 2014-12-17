package casket

import (
	"errors"
)

type File struct {
	Filer Filer `json:"-"`

	Name string `json:"name"`
	ContentType string `json:"content_type"`
	Revisions []SHA1Sum `json:"revisions"`
}

func (f *File) AddRevision(sha SHA1Sum) error {
	return f.Filer.AddRevision(f, sha)
}

func (f *File) Latest() (SHA1Sum, error) {
	if len(f.Revisions) == 0 {
		return [20]byte{}, errors.New("No revisions found.")
	}

	return f.Revisions[len(f.Revisions) - 1], nil
}
