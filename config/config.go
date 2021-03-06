package config

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v3"
)

// Config is a parsed alpaca.json file.
type Config struct {
	Author      string
	BundleID    string            `yaml:"bundle-id"`
	Description string            `yaml:"description"`
	Icon        string            `yaml:"icon"`
	Name        string            `yaml:"name"`
	Objects     ObjectMap         `yaml:"objects"`
	Readme      string            `yaml:"readme"`
	URL         string            `yaml:"url"`
	Variables   map[string]string `yaml:"variables"`
	Version     string            `yaml:"version"`
}

// ObjectMap is a mapping of object names to objects
type ObjectMap map[string]Object

// ThenList is a list of Then structs
type ThenList []Then

func (l *ThenList) UnmarshalYAML(node *yaml.Node) error {
	var s string
	if err := node.Decode(&s); err == nil {
		*l = ThenList{Then{Object: s}}
		return nil
	}

	type alias ThenList
	var as alias
	if err := node.Decode(&as); err != nil {
		return err
	}

	*l = ThenList(as)

	return nil
}

// Then is an object following another object.
type Then struct {
	Object string `yaml:"object"`
}

func (t *Then) UnmarshalYAML(node *yaml.Node) error {
	var s string
	if err := node.Decode(&s); err == nil {
		t.Object = s
		return nil
	}

	type alias Then
	var as alias
	if err := node.Decode(&as); err != nil {
		return err
	}

	*t = Then(as)

	return nil
}

// UnmarshalYAML unmarshals an object.
func (o *ObjectMap) UnmarshalYAML(node *yaml.Node) error {
	var m map[string]Object
	if err := node.Decode(&m); err != nil {
		return err
	}

	*o = make(ObjectMap)

	for name, obj := range m {
		obj.Name = name
		(*o)[obj.Name] = obj
	}

	return nil
}

// Read parses an alpaca.json file.
func Read(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(bytes, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
