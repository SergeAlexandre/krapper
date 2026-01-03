package wrap

import (
	"fmt"
	"krapper/internal/misc"
	"strings"
)

type Cel string
type UiComponent string
type Alignment string
type WrTemplate string
type MenuMode string

const (
	leftAlign   Alignment = "left"
	centerAlign Alignment = "center"
	rightAlign  Alignment = "right"
)

const (
	gridMode    MenuMode = "grid"    // Entities as grid in main pane
	subMenuMode MenuMode = "subMenu" // Entities as subMenu. Selected one in view mode
)

type Wrap struct {
	// required
	ApiVersion string `yaml:"apiVersion" json:"apiVersion"`
	// Always 'wrap' for now
	Kind string `yaml:"kind" json:"kind"`
	// Required
	Name string `yaml:"name" json:"name"`
	// Required
	Version string `yaml:"version" json:"version"`
	// Optional
	Label string `yaml:"label" json:"label"`
	// optional
	Description string `yaml:"description,omitempty" json:"description,omitempty"`
	// Required
	MenuMode MenuMode `yaml:"menuMode,omitempty" json:"menuMode,omitempty"`

	Source struct {
		// Required
		ApiVersion string `yaml:"apiVersion" json:"apiVersion"`
		// Required
		Kind          string            `yaml:"kind" json:"kind"`
		Namespace     string            `yaml:"namespace,omitempty" json:"namespace,omitempty"`
		Selector      map[string]string `yaml:"selector,omitempty" json:"selector,omitempty"`
		ClusterScoped bool              `yaml:"clusterScoped" json:"clusterScoped"`
	} `yaml:"source" json:"source"`

	Operations struct {
		View   bool `yaml:"view" json:"view"`
		Create bool `yaml:"create" json:"create"`
		Update bool `yaml:"update" json:"update"`
		Delete bool `yaml:"delete" json:"delete"`
	} `yaml:"operations" json:"operations"`

	Schema struct {
		Validation *Validation `yaml:"validation,omitempty" json:"validation,omitempty"`
		ValuePath  string      `yaml:"valuePath,omitempty" json:"valuePath,omitempty"`
		Fields     []Field     `yaml:"fields,omitempty" json:"fields,omitempty"`
	} `yaml:"schema" json:"schema"`

	Template WrTemplate `yaml:"template,omitempty" json:"template,omitempty"`
}

var _ valuePathProvider = &Wrap{}

func (w *Wrap) GetValuePath() string {
	return w.Schema.ValuePath
}

func (w *Wrap) Groom() error {
	if w.ApiVersion != "krapper.kubotal.io/v1alpha1" {
		return fmt.Errorf("invalid api version: %s", w.ApiVersion)
	}
	if w.Kind != "Wrap" {
		return fmt.Errorf("invalid Wrap type: %s", w.Kind)
	}
	if w.Name == "" {
		return fmt.Errorf("name is required")
	}
	if w.Version == "" {
		return fmt.Errorf("version is required")
	}
	if w.Label == "" {
		w.Label = misc.Labelize(w.Name)
	}
	if w.MenuMode == "" {
		return fmt.Errorf("menuMode is required")
	}
	if !validMenuModes[w.MenuMode] {
		return fmt.Errorf("invalid menuMode: %s", w.MenuMode)
	}

	if w.Source.ApiVersion == "" {
		return fmt.Errorf("no apiVersion defined for source")
	}
	if w.Source.Kind == "" {
		return fmt.Errorf("no kind defined for source")
	}

	if w.Schema.Validation != nil {
		err := w.Schema.Validation.groom()
		if err != nil {
			return fmt.Errorf("invalid global schema validation: %v", err)
		}
	}

	for idx := range w.Schema.Fields {
		err := w.Schema.Fields[idx].groom(w)
		if err != nil {
			return fmt.Errorf("field '%s': %v", w.Schema.Fields[idx].Name, err)
		}
	}
	return nil
}

var validMenuModes = map[MenuMode]bool{
	gridMode:    true,
	subMenuMode: true,
}

type Validation struct {
	Test    Cel    `yaml:"test,omitempty" json:"test,omitempty"`       // If false, then error
	Message string `yaml:"message,omitempty" json:"message,omitempty"` // Error message
}

func (v *Validation) groom() error {
	return validCel(v.Test)
}

func validCel(exp Cel) error {
	if exp == "" {
		return nil
	}
	// TODO
	//env, err := cel.NewEnv()
	//if err != nil {
	//	return err
	//}
	//ast, issues := env.Compile(string(exp))
	//if issues != nil && issues.Err() != nil {
	//	return issues.Err()
	//}
	//_, err = env.Program(ast)
	//if err != nil {
	//	return err
	//}
	return nil
}

var alignmentSet = map[Alignment]bool{leftAlign: true, centerAlign: true, rightAlign: true}

func validAlignment(a Alignment) error {
	if !alignmentSet[a] {
		return fmt.Errorf("invalid alignment: %s. Must be one of 'left', 'center' or 'right'", a)
	}
	return nil
}

func joinPath(elem1 string, elem2 string) string {
	s := strings.HasSuffix(elem1, ".")
	p := strings.HasPrefix(elem2, ".")
	if s && p {
		return elem1 + elem2[1:]
	}
	if !s && !p {
		return elem1 + "." + elem2
	}
	return elem1 + elem2
}

type valuePathProvider interface {
	GetValuePath() string
}
