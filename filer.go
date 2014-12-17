package casket

type Filer interface {
	Get(string) (*File, error)
	Put(*File) error
	NewFile(string, string) (*File, error)
	AddRevision(*File, SHA1Sum) error
	Exists(string) (bool, error)
}
