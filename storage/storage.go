package storage

type Storage interface {
	Save(book *Book) error
	Read(name string) (*Book, error)
	IsExists(book *Book) (bool, error)
	Search(query string, count int) ([]Book, error)
}

type Book struct {
	Filename string
	FileID   string
	Tags     Tokens
}
