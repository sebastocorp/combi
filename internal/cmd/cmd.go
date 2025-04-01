package cmd

import (
	"combi/internal/cmd/cmdrun"
	"combi/internal/cmd/cmdver"

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
		cmdver.NewCommand(),
		cmdrun.NewCommand(),
	)

	return cmd
}
