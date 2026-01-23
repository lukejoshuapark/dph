package schema

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Definition struct {
	Flags map[string]*FlagDefinition `yaml:"flags"`
}

type FlagDefinition struct {
	Description      string   `yaml:"description"`
	DocumentationURL string   `yaml:"documentationUrl,omitempty"`
	Exclude          []string `yaml:"exclude,omitempty"`
}

func LoadFromFile(file string) (*Definition, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	definition := &Definition{}
	if err := decoder.Decode(definition); err != nil {
		return nil, err
	}

	return definition, nil
}

func PersistToFile(file string, definition *Definition) error {
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()

	encoder := yaml.NewEncoder(f)
	if err := encoder.Encode(definition); err != nil {
		return err
	}

	return nil
}
