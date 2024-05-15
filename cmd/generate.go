package cmd

import (
	"fmt"

	"github.com/oh-my-deploy/omd-generate-plugin/utils"
	"github.com/oh-my-deploy/omd-operator/api/v1alpha1"

	"io"
	"reflect"
	"strings"

	"github.com/spf13/cobra"
)

func InitGenerateCmd() *cobra.Command {
	return newGenerateCommand()
}

func newGenerateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate <path>",
		Short: "Generate manifests from Program resource",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return fmt.Errorf("<path> argument required to generate manifests")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			var manifests []v1alpha1.Program
			var err error
			path := args[0]
			files, err := utils.ListFiles(path)
			if len(files) < 1 {
				return fmt.Errorf("no YAML or JSON files were found in %s", path)
			}
			if err != nil {
				return err
			}
			var errs []error
			manifests, errs = utils.ReadFilesAsManifests(files)
			if len(errs) != 0 {
				errMessages := make([]string, len(errs))
				for idx, err := range errs {
					errMessages[idx] = err.Error()
				}
				return fmt.Errorf("could not read YAML/JSON files:\n%s", strings.Join(errMessages, "\n"))
			}
			for _, manifest := range manifests {
				if reflect.ValueOf(manifest.Spec).IsZero() {
					continue
				}

				if err = generateManifestYaml(cmd.OutOrStdout(), &manifest, true, "default"); err != nil {
					errs = append(errs, err)
				}
				isDeploymentEnabled := manifest.Spec.App.AppType == "server" || manifest.Spec.App.AppType == "ssr"
				if err = generateManifestYaml(cmd.OutOrStdout(), &manifest, isDeploymentEnabled, "deployment"); err != nil {
					errs = append(errs, err)
				}
				if err = generateManifestYaml(cmd.OutOrStdout(), &manifest, manifest.Spec.ServiceAccount.Create, "sa"); err != nil {
					errs = append(errs, err)
				}
				if err = generateManifestYaml(cmd.OutOrStdout(), &manifest, manifest.Spec.Service.Enabled, "service"); err != nil {
					errs = append(errs, err)
				}
				if err = generateManifestYaml(cmd.OutOrStdout(), &manifest, manifest.Spec.Ingress.Enabled, "ingress"); err != nil {
					errs = append(errs, err)
				}
				if err = generateManifestYaml(cmd.OutOrStdout(), &manifest, manifest.Spec.Scheduler.HorizontalPodAutoScaler.Enabled, "hpa"); err != nil {
					errs = append(errs, err)
				}
				if err = generateManifestYaml(cmd.OutOrStdout(), &manifest, manifest.Spec.Scheduler.PodDisruptionBudget.Enabled, "pdb"); err != nil {
					errs = append(errs, err)
				}
			}
			if len(errs) != 0 {
				return fmt.Errorf("could not generate manifests: %v", errs)
			}
			return nil
		},
	}

	return cmd
}

func generateManifestYaml(w io.Writer, manifest *v1alpha1.Program, isEnabled bool, resourceType string) error {
	if !isEnabled {
		return nil
	}
	resource := createResource(manifest, resourceType)
	return printYaml(w, resource)
}

func printYaml(w io.Writer, obj any) error {
	output, err := utils.ConvertToYaml(obj)
	if err != nil {
		return err
	}
	output = strings.ReplaceAll(output, "status: {}", "")
	fmt.Fprintf(w, "%s---\n", output)
	return nil
}

func createResource(manifest *v1alpha1.Program, resourceType string) any {
	var resource any
	switch resourceType {
	case "deployment":
		resource = manifest.ConvertToDeployment()
	case "service":
		resource = manifest.ConvertToService()
	case "sa":
		resource = manifest.ConvertToServiceAccount()
	case "ingress":
		resource = manifest.ConvertToIngress()
	case "hpa":
		resource = manifest.ConvertToHPA()
	case "pdb":
		resource = manifest.ConvertToPdb()
	default:
		resource = manifest
	}
	return resource
}
