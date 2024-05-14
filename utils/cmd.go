package utils

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/oh-my-deploy/omd-operator/api/v1alpha1"
	"github.com/spf13/cobra"
	k8yaml "k8s.io/apimachinery/pkg/util/yaml"
)

type GeneratorCmdFn func() *cobra.Command

func RegisterSubCommands(parent *cobra.Command, generateCommands ...GeneratorCmdFn) {
	for _, gc := range generateCommands {
		parent.AddCommand(gc())
	}
}

func GenerateGroups(parent *cobra.Command, groups ...*cobra.Group) {
	for _, g := range groups {
		parent.AddGroup(g)
	}
}

func ListFiles(root string) ([]string, error) {
	var files []string

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) == ".yaml" || filepath.Ext(path) == ".yml" || filepath.Ext(path) == ".json" {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return files, err
	}

	return files, nil
}

func ReadFilesAsManifests(paths []string) (result []v1alpha1.Program, errs []error) {
	for _, path := range paths {
		rawdata, err := os.ReadFile(path)
		if err != nil {
			errs = append(errs, fmt.Errorf("could not read file: %s from disk: %s", path, err))
		}
		manifest, err := readManifestData(bytes.NewReader(rawdata))
		if err != nil {
			errs = append(errs, fmt.Errorf("could not read file: %s from disk: %s", path, err))
		}
		result = append(result, manifest...)
	}

	return result, errs
}

func readManifestData(yamlData io.Reader) ([]v1alpha1.Program, error) {
	decoder := k8yaml.NewYAMLOrJSONDecoder(yamlData, 1)

	var manifests []v1alpha1.Program
	for {
		nxtManifest := v1alpha1.Program{}
		err := decoder.Decode(&nxtManifest)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		// Skip empty manifests
		manifests = append(manifests, nxtManifest)
	}

	return manifests, nil
}
