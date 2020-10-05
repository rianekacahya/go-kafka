package bootstrap

import (
	"context"
	"fmt"
	"github.com/labstack/echo"
	"github.com/rianekacahya/go-kafka/infrastructure/invoker"
	"github.com/rianekacahya/go-kafka/infrastructure/persistence"
	"github.com/rianekacahya/go-kafka/interface/http/orders"
	"github.com/rianekacahya/go-kafka/middleware/event"
	"github.com/rianekacahya/go-kafka/middleware/rest"
	"github.com/rianekacahya/go-kafka/pkg/echoserver"
	"github.com/rianekacahya/go-kafka/pkg/goconf"
	"github.com/rianekacahya/go-kafka/pkg/gokafka"
	"github.com/rianekacahya/go-kafka/usecase"
	"github.com/spf13/cobra"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

var (
	httpCommand = &cobra.Command{
		Use:   "http",
		Short: "Starting HTTP Server",
		Run: func(cmd *cobra.Command, args []string) {
			// init root context
			ctx := context.Background()

			// init infrastructure
			database := initMysql()
			kafkago := gokafka.New(
				goconf.Config().GetString("kafka.address"),
				event.Logger(),
			)
			telemetry := initTelemetry()

			// init event invoker
			eventInvoker := invoker.EventInitialize(kafkago)

			// init repository
			ordersRepository := persistence.NewOrdersRepository(database)

			// init usecase
			ordersUsecase := usecase.NewOrdersUsecase(ordersRepository, eventInvoker)

			// init echo server
			server := echoserver.NewServer(goconf.Config().GetBool("debug"))

			// set default middleware
			server.Use(
				rest.CORS(),
				rest.Headers(),
				rest.Logger(),
				rest.RequestID(),
				rest.Telemetry(telemetry),
			)

			// versioning
			v1 := server.Group("/v1")

			// set healthcheck endpoint
			v1.GET("/healthcheck", func(c echo.Context) error {
				return c.JSON(http.StatusOK, "OK")
			})

			// init orders handlers
			orders.NewHandler(v1, ordersUsecase)

			// start server
			go func() {
				if err := server.StartServer(&http.Server{
					Addr:         fmt.Sprintf(":%v", goconf.Config().GetInt("rest.port")),
					WriteTimeout: time.Duration(goconf.Config().GetInt("rest.write_timeout")) * time.Second,
					ReadTimeout:  time.Duration(goconf.Config().GetInt("rest.read_timeout")) * time.Second,
					IdleTimeout:  time.Duration(goconf.Config().GetInt("rest.idle_timeout")) * time.Second,
				}); err != nil {
					log.Fatalf("got an error while starting http server, error: %s", err)
				}
			}()

			// shutdown server gracefully
			quit := make(chan os.Signal)
			signal.Notify(quit, os.Interrupt)
			<-quit

			ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
			defer cancel()

			if err := server.Shutdown(ctx); err != nil {
				log.Fatalf("got an error while shutdown http server, error: %s", err)
			}
		},
	}
)
