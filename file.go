package casket

import (
	"errors"
)

type File struct {
	Name string
	ContentType string
	Revisions []SHA1Sum
}

func (f *File) Latest() (SHA1Sum, error) {
	if len(f.Revisions) == 0 {
		return [20]byte{}, errors.New("No revisions found.")
	}

	return f.Revisions[len(f.Revisions) - 1], nil
}
