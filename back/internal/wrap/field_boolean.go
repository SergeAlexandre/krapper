package wrap

import "fmt"

type FieldBoolean struct {
	Default     bool        `yaml:"default,omitempty" json:"default,omitempty"`
	Value       Cel         `yaml:"value,omitempty" json:"value,omitempty"`
	UiComponent UiComponent `yaml:"uiComponent,omitempty" json:"uiComponent,omitempty"`
	Inlist      *struct {
		Hidden      bool        `yaml:"hidden,omitempty" json:"hidden,omitempty"`
		Header      string      `yaml:"header,omitempty" json:"header,omitempty"` // Default to label
		UiComponent UiComponent `yaml:"uiComponent,omitempty" json:"uiComponent,omitempty"`
	} `yaml:"inlist,omitempty" json:"inlist,omitempty"`
}

func (f *FieldBoolean) groom(defaultValueCel Cel, label string) error {
	if f.Value == "" {
		f.Value = defaultValueCel
	}
	err := validCel(f.Value)
	if err != nil {
		return err
	}
	if f.UiComponent == "" {
		f.UiComponent = "checkbox"
	}
	if !validBooleanUiComponents[f.UiComponent] {
		return fmt.Errorf("invalid UiComponent: %s", f.UiComponent)
	}
	if f.Inlist != nil {
		if f.Inlist.Header == "" {
			f.Inlist.Header = label
		}
		if f.Inlist.UiComponent == "" {
			f.Inlist.UiComponent = f.UiComponent
		}
		if !validBooleanUiComponents[f.Inlist.UiComponent] {
			return fmt.Errorf("invalid UiComponent: %s", f.Inlist.UiComponent)
		}
	}
	return nil
}

var validBooleanUiComponents = map[UiComponent]bool{
	"checkbox": true,
}
