package wrap

import "fmt"

type FieldNumber struct {
	Default     *float64    `yaml:"default,omitempty" json:"default,omitempty"`
	Enum        []float64   `yaml:"enum,omitempty" json:"enum,omitempty"`
	Value       Cel         `yaml:"value,omitempty" json:"value,omitempty"`
	UiComponent UiComponent `yaml:"uiComponent,omitempty" json:"uiComponent,omitempty"`
	Format      string      `yaml:"format,omitempty" json:"format,omitempty"` // fmt.Sprintf format expression
	Inlist      *struct {
		Hidden      bool        `yaml:"hidden,omitempty" json:"hidden,omitempty"`
		Header      string      `yaml:"header,omitempty" json:"header,omitempty"` // Default to label
		UiComponent UiComponent `yaml:"uiComponent,omitempty" json:"uiComponent,omitempty"`
		Alignment   Alignment   `yaml:"alignment,omitempty" json:"alignment,omitempty"`
		Format      string      `yaml:"format,omitempty" json:"format,omitempty"` // fmt.Sprintf format expression
	} `yaml:"inlist,omitempty" json:"inlist,omitempty"`
}

func (f *FieldNumber) groom(defaultValueCel Cel, label string) error {
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
	if !validNumberUiComponents[f.UiComponent] {
		return fmt.Errorf("invalid UiComponent: %s", f.UiComponent)
	}
	if f.Format == "" {
		f.Format = "%f"
	}
	if f.Inlist != nil {
		if f.Inlist.Header == "" {
			f.Inlist.Header = label
		}
		if f.Inlist.UiComponent == "" {
			f.Inlist.UiComponent = f.UiComponent
		}
		if !validNumberUiComponents[f.Inlist.UiComponent] {
			return fmt.Errorf("invalid UiComponent: %s", f.Inlist.UiComponent)
		}
		if f.Inlist.Alignment == "" {
			f.Inlist.Alignment = "left"
		}
		err = validAlignment(f.Inlist.Alignment)
		if err != nil {
			return err
		}
		if f.Inlist.Format == "" {
			f.Inlist.Format = f.Format
		}
	}
	return nil
}

var validNumberUiComponents = map[UiComponent]bool{
	"raw": true,
}
