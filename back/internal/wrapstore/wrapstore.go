package wrapstore

import (
	"krapper/internal/wrap"
	"log/slog"
)

type Catalog struct {
	Wraps []struct {
		Name     string        `yaml:"name" json:"name"`
		Label    string        `yaml:"label" json:"label"`
		MenuMode wrap.MenuMode `yaml:"menuMode" json:"menuMode"`
	} `yaml:"wraps" json:"wraps"`
}

type WrapStore interface {
	GetCatalog() *Catalog
	GetWrap(name string) wrap.Wrap
}

func New(path string, logger *slog.Logger) (WrapStore, error) {
	return nil, nil
}
