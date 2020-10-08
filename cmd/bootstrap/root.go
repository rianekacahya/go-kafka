package bootstrap

import (
	"github.com/spf13/cobra"
	"log"
	"os"
)
func Execute() {
	var command = new(cobra.Command)

	command.AddCommand(
		migrationCommand,
		orderConsumerCommand,
		httpCommand,
	)

	if err := command.Execute(); err != nil {
		log.Printf("got an error while initialize, error: %s", err)
		os.Exit(1)
	}
}