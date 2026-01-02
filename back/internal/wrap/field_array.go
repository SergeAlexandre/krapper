package wrap

import "fmt"

type FieldArray struct {
	Item struct {
		Validation *Validation `yaml:"validation,omitempty" json:"validation,omitempty"`
		Type       Type        `yaml:",inline" json:",inline"`
	} `yaml:"item" json:"item"`
	InList *struct {
		Hidden    bool        `yaml:"hidden,omitempty" json:"hidden,omitempty"`
		Header    string      `yaml:"header,omitempty" json:"header,omitempty"` // Default to label
		Display   UiComponent `yaml:"display,omitempty" json:"display,omitempty"`
		Width     int         `yaml:"width,omitempty" json:"width,omitempty"`
		Height    int         `yaml:"height,omitempty" json:"height,omitempty"`
		Value     Cel         `yaml:"value,omitempty" json:"value,omitempty"`
		Alignment Alignment   `yaml:"alignment,omitempty" json:"alignment,omitempty"`
	}
}

func (f *FieldArray) groom(defaultValueCel Cel, label string) error {
	if f.Item.Validation == nil {
		err := f.Item.Validation.groom()
		if err != nil {
			return fmt.Errorf("invalid validation: %w", err)
		}
	}
	err := f.Item.Type.groom(defaultValueCel, label)
	if err != nil {
		return err
	}
	if f.InList != nil {
		if f.InList.Header == "" {
			f.InList.Header = label
		}
		if f.InList.Display == "" {
			f.InList.Display = "raw"
		}
		if f.InList.Height == 0 {
			f.InList.Height = 1
		}
		if f.InList.Value == "" {
			f.InList.Value = "'[...]'"
		}
		err := validCel(f.InList.Value)
		if err != nil {
			return err
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
