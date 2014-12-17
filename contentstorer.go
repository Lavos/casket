package casket

type ContentStorer interface {
	Put([]byte) (SHA1Sum, error)
	Get(SHA1Sum) ([]byte, error)
	Exists(SHA1Sum) (bool, error)
}
