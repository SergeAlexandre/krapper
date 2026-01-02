package wrap

import (
	"errors"
	"fmt"
	"krapper/internal/misc"
)

type Field struct {
	Name       string      `yaml:"name" json:"name"`
	Label      string      `yaml:"label,omitempty" json:"label,omitempty"`
	Tooltip    string      `yaml:"tooltip,omitempty" json:"tooltip,omitempty"`
	Validation *Validation `yaml:"validation,omitempty" json:"validation,omitempty"`
	Required   bool        `yaml:"required" json:"required"`

	Condition Cel `yaml:"condition,omitempty" json:"condition,omitempty"`
	ReadOnly  Cel `yaml:"readOnly,omitempty" json:"readOnly,omitempty"`

	Type Type `yaml:",inline" json:",inline"`

	pathProvider valuePathProvider
}

func (f *Field) GetValuePath() string {
	return f.pathProvider.GetValuePath()
}

var _ valuePathProvider = &Field{}

func (f *Field) groom(pathProvider valuePathProvider) error {
	f.pathProvider = pathProvider
	if f.Name == "" {
		return errors.New("field name is required")
	}
	if f.Label == "" {
		f.Label = misc.Labelize(f.Name)
	}
	if f.Validation != nil {
		err := f.Validation.groom()
		if err != nil {
			return fmt.Errorf("invalid validation: %w", err)
		}
	}
	err := validCel(f.Condition)
	if err != nil {
		return fmt.Errorf("invalid condition: %w", err)
	}
	err = validCel(f.ReadOnly)
	if err != nil {
		return fmt.Errorf("invalid readOnly expression: %w", err)
	}
	defaultValueCel := Cel(joinPath(pathProvider.GetValuePath(), f.Name))

	err = f.Type.groom(defaultValueCel, f.Label)
	if err != nil {
		return err
	}
	return nil
}

type Type struct {
	Array    *FieldArray    `yaml:"array,omitempty" json:"array,omitempty"`
	Boolean  *FieldBoolean  `yaml:"boolean,omitempty" json:"boolean,omitempty"`
	Duration *FieldDuration `yaml:"duration,omitempty" json:"duration,omitempty"`
	Integer  *FieldInteger  `yaml:"integer,omitempty" json:"integer,omitempty"`
	Number   *FieldNumber   `yaml:"number,omitempty" json:"number,omitempty"`
	Object   *FieldObject   `yaml:"object,omitempty" json:"object,omitempty"`
	String   *FieldString   `yaml:"string,omitempty" json:"string,omitempty"`
}

func (t *Type) groom(defaultValueCel Cel, label string) error {
	myType := "" // Just to detect duplication
	if t.Array != nil {
		//if myType != "" {
		//	return fmt.Errorf("can't be '%s' and 'array'", myType)
		//}
		myType = "array"
		err := t.Array.groom(defaultValueCel, label)
		if err != nil {
			return err
		}
	}
	if t.Boolean != nil {
		if myType != "" {
			return fmt.Errorf("can't be '%s' and 'boolean'", myType)
		}
		myType = "boolean"
		err := t.Boolean.groom(defaultValueCel, label)
		if err != nil {
			return err
		}
	}
	if t.Duration != nil {
		if myType != "" {
			return fmt.Errorf("can't be '%s' and 'duration'", myType)
		}
		myType = "duration"
		err := t.Duration.groom(defaultValueCel, label)
		if err != nil {
			return err
		}
	}
	if t.Integer != nil {
		if myType != "" {
			return fmt.Errorf("can't be '%s' and 'integer'", myType)
		}
		myType = "integer"
		err := t.Integer.groom(defaultValueCel, label)
		if err != nil {
			return err
		}
	}
	if t.Number != nil {
		if myType != "" {
			return fmt.Errorf("can't be '%s' and 'number'", myType)
		}
		myType = "number"
		err := t.Number.groom(defaultValueCel, label)
		if err != nil {
			return err
		}
	}
	if t.Object != nil {
		if myType != "" {
			return fmt.Errorf("can't be '%s' and 'object'", myType)
		}
		myType = "object"
		err := t.Object.groom(defaultValueCel, label)
		if err != nil {
			return err
		}
	}
	if myType == "" && t.String == nil {
		t.String = &FieldString{}
	}
	if t.String != nil {
		myType = "string"
		err := t.String.groom(defaultValueCel, label)
		if err != nil {
			return err
		}
	}
	return nil
}
