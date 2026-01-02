package wrapstore

import "krapper/internal/wrap"

type Catalog struct {
	Wraps []struct {
		Name  string `yaml:"name" json:"name"`
		Label string `yaml:"label" json:"label"`
	} `yaml:"wraps" json:"wraps"`
}

type WrapStore interface {
	GetCatalog() *Catalog
	GetWrap(name string) wrap.Wrap
}

func New(path string) (WrapStore, error) {
	return nil, nil
}
