package cmd

import (
	"combi/internal/cmd/run"
	"combi/internal/cmd/version"

	"github.com/spf13/cobra"
)

const (
	descriptionShort = `TODO`
	descriptionLong  = `
	TODO`
)

func NewRootCommand(name string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   name,
		Short: descriptionShort,
		Long:  descriptionLong,
	}

	cmd.AddCommand(
		version.NewCommand(),
		run.NewCommand(),
	)

	return cmd
}
