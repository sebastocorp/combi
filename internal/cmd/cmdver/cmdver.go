package cmdver

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	commit  = "none"
	golang  = "none"
	version = "dev"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "version",
		DisableFlagsInUseLine: true,
		Short:                 `print the current version`,
		Long:                  `print the current short commit, go version and client version.`,

		Run: RunCommand,
	}

	return cmd
}

func RunCommand(cmd *cobra.Command, args []string) {
	fmt.Printf("%-8s %s\n%-8s %s\n%-8s %s\n", "commit:", commit, "golang:", golang, "version:", version)
}
