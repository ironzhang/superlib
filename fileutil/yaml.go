package fileutil

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// ReadYAML read value from yaml file.
func ReadYAML(filename string, v interface{}) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	if err = yaml.Unmarshal(data, v); err != nil {
		return err
	}
	return nil
}

// WriteYAML write value to yaml file.
func WriteYAML(filename string, v interface{}) error {
	data, err := yaml.Marshal(v)
	if err != nil {
		return err
	}
	if err = ioutil.WriteFile(filename, data, 0666); err != nil {
		return err
	}
	return nil
}
