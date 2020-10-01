package bootstrap

import (
	"github.com/spf13/cobra"
	"log"
)

var (
	orderConsumerCommand = &cobra.Command{
		Use:   "order-consumer",
		Short: "Starting Order Consumer",
		Run: func(cmd *cobra.Command, args []string) {
			log.Println("hello world")
		},
	}
)

