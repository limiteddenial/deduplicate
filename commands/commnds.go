package commands

import "github.com/spf13/cobra"

func NewCmdRoot() *cobra.Command {
	var rootCmd = &cobra.Command{
		Use:   "deduplicate",
		Short: "Deduplicate files based on SHA256 hash",
	}
	rootCmd.AddCommand(
		NewCmdSearch(),
		NewCmdVersion(),
	)

	return rootCmd
}
