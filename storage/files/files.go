package files

import (
	"encoding/gob"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"

	"github.com/gleblug/library-bot/lib/e"
	"github.com/gleblug/library-bot/storage"
)

type Tags map[string]storage.Tokens

type Storage struct {
	basePath string
	tags     Tags
}

const (
	defaultPerm = 0774
)

func New(basePath string) (Storage, error) {
	tags, err := decodeTags(basePath)
	if err != nil {
		return Storage{}, e.Wrap("can't create files storage", err)
	}

	return Storage{
		basePath: basePath,
		tags:     tags,
	}, nil
}

func (s Storage) Save(book *storage.Book) (err error) {
	defer func() { err = e.WrapIfErr("can't save book", err) }()

	if err := os.MkdirAll(s.basePath, defaultPerm); err != nil {
		return err
	}

	fPath := filepath.Join(s.basePath, book.Filename)

	file, err := os.Create(fPath)
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()

	if err := gob.NewEncoder(file).Encode(book); err != nil {
		return err
	}

	s.tags[book.Filename] = book.Tags
	return nil
}

func (s Storage) Read(name string) (*storage.Book, error) {
	path := filepath.Join(s.basePath, name)

	book, err := decodeBook(path)
	if err != nil {
		return nil, e.Wrap("can't read book", err)
	}

	return book, nil
}

func (s Storage) IsExists(book *storage.Book) (bool, error) {
	path := filepath.Join(s.basePath, book.Filename)

	switch _, err := os.Stat(path); {
	case errors.Is(err, os.ErrNotExist):
		return false, nil
	case err != nil:
		msg := fmt.Sprintf("can't check '%s' file is exists", path)
		return false, e.Wrap(msg, err)
	}

	return true, nil
}

func (s Storage) Search(query string, count int) (books []storage.Book, err error) {
	defer func() { err = e.WrapIfErr("can't search book", err) }()

	queryTokens := storage.Analyze(query)
	relevance := map[string]int{}

	for name, tokens := range s.tags {
		relevance[name] = tokens.OverlapCoefficient(queryTokens)
	}

	weights := make([]int, 0)
	index := map[int][]string{}
	for name, w := range relevance {
		weights = append(weights, w)
		index[w] = append(index[w], name)
	}

	slices.Sort(weights)

	resNames := make([]string, 0)
	i := 0
	for i < len(weights) && len(resNames) < count {
		w := weights[i]
		resNames = append(resNames, index[w]...)
	}

	books = make([]storage.Book, 0)
	for _, name := range resNames {
		book, err := s.Read(name)
		if err != nil {
			return nil, err
		}

		books = append(books, *book)
		if len(books) == count {
			break
		}
	}

	return books, nil
}

func decodeBook(path string) (*storage.Book, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, e.Wrap("can't decode book", err)
	}
	defer func() { _ = f.Close() }()

	var b storage.Book

	if err := gob.NewDecoder(f).Decode(&b); err != nil {
		return nil, e.Wrap("can't decode book", err)
	}

	return &b, nil
}

func decodeTags(path string) (tags Tags, err error) {
	defer func() { err = e.WrapIfErr("can't decode tags", err) }()

	files, err := os.ReadDir(path)
	if errors.Is(err, os.ErrNotExist) {
		if err := os.MkdirAll(path, defaultPerm); err != nil {
			return nil, err
		}
		return Tags{}, nil
	}
	if err != nil {
		return nil, err
	}

	tags = Tags{}
	for _, file := range files {
		bookPath := filepath.Join(path, file.Name())

		book, err := decodeBook(bookPath)
		if err != nil {
			return nil, err
		}

		tags[book.Filename] = book.Tags
	}

	return tags, nil
}
