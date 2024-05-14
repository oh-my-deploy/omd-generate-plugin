package cmd

import "github.com/spf13/cobra"

func CreateRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "omd-generate-plugin",
		Short: "omd-generate-plugin is a generating k8s manifest for using program resource",
		Long:  "omd-generate-plugin is a generating k8s manifest for using program resource",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
	return cmd
}
