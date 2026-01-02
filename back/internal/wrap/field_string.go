package wrap

import "fmt"

type FieldString struct {
	Default     string      `yaml:"default,omitempty" json:"default,omitempty"`
	Enum        []string    `yaml:"enum,omitempty" json:"enum,omitempty"`
	Value       Cel         `yaml:"value,omitempty" json:"value,omitempty"`
	UiComponent UiComponent `yaml:"uiComponent,omitempty" json:"uiComponent,omitempty"`
	Width       int         `yaml:"width,omitempty" json:"width,omitempty"`
	Height      int         `yaml:"height,omitempty" json:"height,omitempty"`
	Inlist      *struct {
		Hidden      bool        `yaml:"hidden,omitempty" json:"hidden,omitempty"`
		Header      string      `yaml:"header,omitempty" json:"header,omitempty"` // Default to label
		UiComponent UiComponent `yaml:"uiComponent,omitempty" json:"uiComponent,omitempty"`
		Alignment   Alignment   `yaml:"alignment,omitempty" json:"alignment,omitempty"`
		Value       Cel         `yaml:"value,omitempty" json:"value,omitempty"`
		Width       int         `yaml:"width,omitempty" json:"width,omitempty"`
		Height      int         `yaml:"height,omitempty" json:"height,omitempty"`
	} `yaml:"inlist,omitempty" json:"inlist,omitempty"`
}

func (f *FieldString) groom(defaultValueCel Cel, label string) error {
	if f.Value == "" {
		f.Value = defaultValueCel
	}
	err := validCel(f.Value)
	if err != nil {
		return err
	}
	if f.Width == 0 {
		f.Width = 30
	}
	if f.Height == 0 {
		f.Height = 1
	}
	if f.UiComponent == "" {
		if f.Height == 1 {
			f.UiComponent = "input"
		} else {
			f.UiComponent = "textarea"
		}
	}
	if !validStringUiComponents[f.UiComponent] {
		return fmt.Errorf("invalid UiComponent: %s", f.UiComponent)
	}
	if f.Inlist != nil {
		if f.Inlist.Header == "" {
			f.Inlist.Header = label
		}
		if f.Inlist.UiComponent == "" {
			f.Inlist.UiComponent = f.UiComponent
		}
		if !validStringUiComponents[f.Inlist.UiComponent] {
			return fmt.Errorf("invalid UiComponent: %s", f.Inlist.UiComponent)
		}
		if f.Inlist.Alignment == "" {
			f.Inlist.Alignment = "left"
		}
		err = validAlignment(f.Inlist.Alignment)
		if err != nil {
			return err
		}
		if f.Inlist.Value == "" {
			f.Inlist.Value = f.Value
		}
		err := validCel(f.Value)
		if err != nil {
			return err
		}
		if f.Inlist.Width == 0 {
			f.Inlist.Width = f.Width
		}
		if f.Inlist.Height == 0 {
			f.Inlist.Height = 1
		}
	}
	return nil
}

var validStringUiComponents = map[UiComponent]bool{
	"input":    true,
	"textarea": true,
}
