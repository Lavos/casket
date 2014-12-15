package casket

type Strorer interface {
	Put([]byte) (SHA1Sum, error)
	Get(SHA1Sum) ([]byte, error)
}
