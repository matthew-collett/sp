package config

import (
	"path/filepath"
	"strings"

	"github.com/matthew-collett/sp/internal/osutil"
)

const shelfFile = "shelf"

type Item struct {
	Name string `json:"name"`
	URI  string `json:"uri"`
	Type string `json:"type"`
}

type Shelf struct {
	path  string
	Items map[string]Item `json:"items"`
}

func NewShelf() (*Shelf, error) {
	dir, err := Dir()
	if err != nil {
		return nil, err
	}
	path := filepath.Join(dir, shelfFile)
	s := &Shelf{
		path:  path,
		Items: make(map[string]Item),
	}
	osutil.ReadFileSilent(path, s)
	if s.Items == nil {
		s.Items = make(map[string]Item)
	}
	return s, nil
}

func (s *Shelf) Write() error {
	return osutil.WriteFile(s, s.path, 0600)
}

func (s *Shelf) Has(name string) bool {
	_, ok := s.Items[name]
	return ok
}

func (s *Shelf) Get(name string) (Item, bool) {
	item, ok := s.Items[name]
	return item, ok
}

func (s *Shelf) Add(name, uri string) bool {
	if s.Has(name) {
		return false
	}
	s.Items[name] = Item{
		Name: name,
		URI:  uri,
		Type: uriType(uri),
	}
	return true
}

func (s *Shelf) Rename(old, new string) bool {
	item, ok := s.Items[old]
	if !ok {
		return false
	}
	item.Name = new
	s.Items[new] = item
	delete(s.Items, old)
	return true
}

func (s *Shelf) Drop(name string) bool {
	if !s.Has(name) {
		return false
	}
	delete(s.Items, name)
	return true
}

func (s *Shelf) IsEmpty() bool {
	return len(s.Items) == 0
}

func uriType(uri string) string {
	parts := strings.SplitN(uri, ":", 3)
	if len(parts) == 3 {
		return parts[1]
	}
	return "unknown"
}
