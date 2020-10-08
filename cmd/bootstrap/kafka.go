package bootstrap

import (
	"context"
	"github.com/rianekacahya/go-kafka/infrastructure/invoker"
	"github.com/rianekacahya/go-kafka/infrastructure/persistence"
	"github.com/rianekacahya/go-kafka/interface/consumers/orders"
	"github.com/rianekacahya/go-kafka/middleware/event"
	"github.com/rianekacahya/go-kafka/pkg/goconf"
	"github.com/rianekacahya/go-kafka/pkg/gokafka"
	"github.com/rianekacahya/go-kafka/usecase"
	"github.com/spf13/cobra"
	"log"
	"os"
	"os/signal"
)

var (
	orderConsumerCommand = &cobra.Command{
		Use:   "order-consumer",
		Short: "Starting Order Consumer",
		Run: func(cmd *cobra.Command, args []string) {
			// init root context
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			// init infrastructure
			database := initMysql()
			redis := initRedis()
			telemetry := initTelemetry()
			kafkago := gokafka.New(
				goconf.Config().GetString("kafka.address"),
				event.Logger(),
				event.Telemetry(telemetry),
				event.RequestID(),
			)

			// init event invoker
			syncInvoker := invoker.SyncInitialize(redis)
			eventInvoker := invoker.EventInitialize(kafkago)

			// init repository
			ordersRepository := persistence.NewOrdersRepository(database)

			// init usecase
			ordersUsecase := usecase.NewOrdersUsecase(ordersRepository, eventInvoker, syncInvoker)

			// starting order consumers
			orders.NewHandler(ctx, kafkago, ordersUsecase)

			quit := make(chan os.Signal)
			signal.Notify(quit, os.Interrupt)
			<-quit


			log.Println("kafka subcriber stopped.")
		},
	}
)

