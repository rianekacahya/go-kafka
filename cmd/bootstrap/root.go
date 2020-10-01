package bootstrap

import (
	"github.com/rianekacahya/go-kafka/pkg/logger"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
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
		logger.Error("Bootstrap", zap.Any("error", err))
		os.Exit(1)
	}
}