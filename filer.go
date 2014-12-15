package casket

type Filer interface {
	Get(string) (*File, error)
	Put(*File) error
	Exists(*File) (bool, error)
}
