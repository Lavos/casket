package casket

type ContentStorer interface {
	PutContent([]byte) (SHA1Sum, error)
	GetContent(SHA1Sum) ([]byte, error)
	ContentExists(SHA1Sum) (bool, error)
}

type Filer interface {
	GetFile(string) (*File, error)
	PutFile(*File) error
	NewFile(string, string) (*File, error)
	AddRevision(*File, SHA1Sum) error
	FileExists(string) (bool, error)
}
