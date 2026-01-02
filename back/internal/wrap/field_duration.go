package wrap

import (
	"fmt"
	"time"
)

type FieldDuration struct {
	Default     time.Duration `yaml:"default,omitempty" json:"default,omitempty"`
	Value       Cel           `yaml:"value,omitempty" json:"value,omitempty"`
	UiComponent UiComponent   `yaml:"uiComponent,omitempty" json:"uiComponent,omitempty"`
	InList      *struct {
		Hidden      bool        `yaml:"hidden,omitempty" json:"hidden,omitempty"`
		Header      string      `yaml:"header,omitempty" json:"header,omitempty"` // Default to label
		UiComponent UiComponent `yaml:"uiComponent,omitempty" json:"uiComponent,omitempty"`
		Alignment   Alignment   `yaml:"alignment,omitempty" json:"alignment,omitempty"`
	} `yaml:"inList,omitempty" json:"inList,omitempty"`
}

func (f *FieldDuration) groom(defaultValueCel Cel, label string) error {
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
	if !validDurationUiComponents[f.UiComponent] {
		return fmt.Errorf("invalid UiComponent: %s", f.UiComponent)
	}
	if f.InList != nil {
		if f.InList.Header == "" {
			f.InList.Header = label
		}
		if f.InList.UiComponent == "" {
			f.InList.UiComponent = f.UiComponent
		}
		if !validDurationUiComponents[f.InList.UiComponent] {
			return fmt.Errorf("invalid UiComponent: %s", f.InList.UiComponent)
		}
		if f.InList.Alignment == "" {
			f.InList.Alignment = "left"
		}
		err = validAlignment(f.InList.Alignment)
		if err != nil {
			return err
		}
	}
	return nil
}

var validDurationUiComponents = map[UiComponent]bool{
	"raw": true,
}
