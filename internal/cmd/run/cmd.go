package run

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"combi/internal/combi"

	"github.com/spf13/cobra"
)

const (
	configFlagName = "config"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "run",
		DisableFlagsInUseLine: true,
		Short:                 `TODO: short description`,
		Long:                  `TODO: long description`,

		Run: RunCommand,
	}

	cmd.Flags().String(configFlagName, "combi.yaml", "combi configuration file")

	return cmd
}

// RunCommand TODO
// Ref: https://pkg.go.dev/github.com/spf13/pflag#StringSlice
func RunCommand(cmd *cobra.Command, args []string) {
	configFilePath, err := cmd.Flags().GetString(configFlagName)
	if err != nil {
		log.Fatalf("unable to get flag --config: %s", err.Error())
	}

	/////////////////////////////
	// EXECUTION FLOW RELATED
	/////////////////////////////

	c, err := combi.NewCombi(configFilePath)
	if err != nil {
		log.Fatalf("unable to init combi instance: %s", err.Error())
	}

	go c.Run()
	defer c.Stop()

	// Wait for the process to be shutdown.
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
}
