package wrap

import "fmt"

type FieldObject struct {
	BasePath    string      `yaml:"basePath,omitempty" json:"basePath,omitempty"` // From root if defined
	Fields      []Field     `yaml:"fields" json:"fields"`
	UiComponent UiComponent `yaml:"uiComponent,omitempty" json:"uiComponent,omitempty"`
	InList      *struct {
		Hidden      bool        `yaml:"hidden,omitempty" json:"hidden,omitempty"`
		Header      string      `yaml:"header,omitempty" json:"header,omitempty"` // Default to label
		UiComponent UiComponent `yaml:"uiComponent,omitempty" json:"uiComponent,omitempty"`
		Alignment   Alignment   `yaml:"alignment,omitempty" json:"alignment,omitempty"`
		Width       int         `yaml:"width,omitempty" json:"width,omitempty"`
		Height      int         `yaml:"height,omitempty" json:"height,omitempty"`
		Value       Cel         `yaml:"value,omitempty" json:"value,omitempty"`
	} `yaml:"inList,omitempty" json:"inList,omitempty"`
}

var _ valuePathProvider = &FieldObject{}

func (f *FieldObject) GetValuePath() string {
	return f.BasePath
}

func (f *FieldObject) groom(defaultValueCel Cel, label string) error {
	if f.BasePath == "" {
		f.BasePath = string(defaultValueCel)
	}
	for idx := range f.Fields {
		err := f.Fields[idx].groom(f)
		if err != nil {
			return fmt.Errorf("field '%s': %v", f.Fields[idx].Name, err)
		}
	}
	if f.UiComponent == "" {
		f.UiComponent = "fieldSet"
	}
	if !validObjectCardUiComponents[f.UiComponent] {
		return fmt.Errorf("invalid UiComponent: %s", f.UiComponent)
	}
	if f.InList != nil {
		if f.InList.Header == "" {
			f.InList.Header = label
		}
		if f.InList.UiComponent == "" {
			f.InList.UiComponent = "raw"
		}
		if !validObjectListUiComponents[f.InList.UiComponent] {
			return fmt.Errorf("invalid UiComponent: %s", f.InList.UiComponent)
		}
		if f.InList.Alignment == "" {
			f.InList.Alignment = "left"
		}
		err := validAlignment(f.InList.Alignment)
		if err != nil {
			return err
		}
		if f.InList.Value == "" {
			f.InList.Value = "'{...}'"
		}
	}
	return nil
}

var validObjectCardUiComponents = map[UiComponent]bool{
	"fieldSet": true,
}
var validObjectListUiComponents = map[UiComponent]bool{
	"raw": true,
}
