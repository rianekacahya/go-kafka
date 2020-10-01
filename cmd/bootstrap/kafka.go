package bootstrap

import (
	"context"
	"github.com/rianekacahya/go-kafka/infrastructure/invoker"
	"github.com/rianekacahya/go-kafka/infrastructure/persistence"
	"github.com/rianekacahya/go-kafka/interface/consumers/orders"
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
			kafkago := gokafka.New(
				"v4",
				goconf.Config().GetString("kafka.address"),
				goconf.Config().GetInt("kafka.max_retry"),
			)

			// init event invoker
			eventInvoker := invoker.EventInitialize(kafkago)

			// init repository
			ordersRepository := persistence.NewOrdersRepository(database)

			// init usecase
			ordersUsecase := usecase.NewOrdersUsecase(ordersRepository, eventInvoker)

			// starting order consumers
			if err := orders.NewHandler(ctx, kafkago, ordersUsecase); err != nil {
				log.Fatalf("got an error while bootstraping order consumers, error: %s", err)
			}

			quit := make(chan os.Signal)
			signal.Notify(quit, os.Interrupt)
			<-quit


			log.Println("kafka subcriber stopped.")
		},
	}
)

