package wrap

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type header struct {
	ApiVersion string `yaml:"apiVersion" json:"apiVersion"`
	Kind       string `yaml:"kind" json:"kind"`
}

// Load a wrap file. Return nil, nil if file s not a wrap one.
func Load(fName string) (*Wrap, error) {

	filename, err := filepath.Abs(fName)
	if err != nil {
		log.Fatalf("Error getting absolute path of yaml file: %v", err)
	}

	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening file %s: %v", filename, err)
	}
	defer func() { _ = file.Close() }()

	var h header
	err = yaml.NewDecoder(file).Decode(&h)
	if err != nil || h.ApiVersion != "krapper.kubotal.io/v1alpha1" || h.Kind != "Wrap" {
		return nil, nil // Non-wrap file
	}
	_, err = file.Seek(0, 0)
	if err != nil {
		return nil, fmt.Errorf("error seeking file %s: %v", filename, err)
	}
	var w Wrap
	decoder := yaml.NewDecoder(file)
	decoder.KnownFields(true)
	if err := decoder.Decode(&w); err != nil {
		return nil, fmt.Errorf("error decoding file %s: %v", filename, err)
	}

	err = w.Groom()
	if err != nil {
		return nil, err
	}
	return &w, nil
}
