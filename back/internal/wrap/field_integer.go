package wrap

import "fmt"

type FieldInteger struct {
	Default     *int        `yaml:"default,omitempty" json:"default,omitempty"`
	Enum        []int       `yaml:"enum,omitempty" json:"enum,omitempty"`
	Value       Cel         `yaml:"value,omitempty" json:"value,omitempty"`
	UiComponent UiComponent `yaml:"uiComponent,omitempty" json:"uiComponent,omitempty"`
	Inlist      *struct {
		Hidden      bool        `yaml:"hidden,omitempty" json:"hidden,omitempty"`
		Header      string      `yaml:"header,omitempty" json:"header,omitempty"` // Default to label
		UiComponent UiComponent `yaml:"uiComponent,omitempty" json:"uiComponent,omitempty"`
		Alignment   Alignment   `yaml:"alignment,omitempty" json:"alignment,omitempty"`
	} `yaml:"inlist,omitempty" json:"inlist,omitempty"`
}

func (f *FieldInteger) groom(defaultValueCel Cel, label string) error {
	if f.Value == "" {
		f.Value = defaultValueCel
	}
	err := validCel(f.Value)
	if err != nil {
		return err
	}
	if f.UiComponent == "" {
		f.UiComponent = "raw"
	}
	if !validIntegerUiComponents[f.UiComponent] {
		return fmt.Errorf("invalid UiComponent: %s", f.UiComponent)
	}
	if f.Inlist != nil {
		if f.Inlist.Header == "" {
			f.Inlist.Header = label
		}
		if f.Inlist.UiComponent == "" {
			f.Inlist.UiComponent = f.UiComponent
		}
		if !validIntegerUiComponents[f.Inlist.UiComponent] {
			return fmt.Errorf("invalid UiComponent: %s", f.Inlist.UiComponent)
		}
		if f.Inlist.Alignment == "" {
			f.Inlist.Alignment = "left"
		}
		err = validAlignment(f.Inlist.Alignment)
		if err != nil {
			return err
		}
	}
	return nil
}

var validIntegerUiComponents = map[UiComponent]bool{
	"raw": true,
}
