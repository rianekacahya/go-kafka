package bootstrap

import (
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/rianekacahya/go-kafka/pkg/goconf"
	"github.com/spf13/cobra"
	"log"
)

var (
	migrationCommand = &cobra.Command{
		Use:   "migrations",
		Short: "this command uses to run database migration up",
		Run: func(cmd *cobra.Command, args []string) {
			dsn := fmt.Sprintf("mysql://%s", goconf.Config().GetString("mysql.dsn"))
			m, err := migrate.New("file://files/migrations", dsn)
			if err != nil {
				log.Fatalf("got an error while initialize database migrations, error: %s", err)
			}

			if err := m.Up(); err != nil {
				log.Fatalf("got an error while execute database migrations, error: %s", err)
			}
		},
	}
)
