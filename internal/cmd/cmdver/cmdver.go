package cmdver

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	version = "dev"
	golang  = "none"
	commit  = "none"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "version",
		DisableFlagsInUseLine: true,
		Short:                 `Print the current version`,
		Long:                  `version command show the current client version and commit hash.`,

		Run: RunCommand,
	}

	return cmd
}

func RunCommand(cmd *cobra.Command, args []string) {
	fmt.Printf("%-8s %s\n%-8s %s\n%-8s %s\n", "version:", version, "golang:", golang, "commit:", commit)
}
