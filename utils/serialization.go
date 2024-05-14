package utils

import (
	"encoding/json"
	"sigs.k8s.io/yaml"
)

func ConvertToYaml(obj any) (string, error) {
	jsonOutput, err := json.Marshal(obj)
	if err != nil {
		return "", err
	}

	yamlOutput, err := yaml.JSONToYAML(jsonOutput)
	if err != nil {
		return "", err
	}

	return string(yamlOutput), nil
}

func ConvertToObj(input string, out any) error {
	return yaml.Unmarshal([]byte(input), out)
}
