package migrate

import (
	"context"
	"flag"

	"github.com/golang/glog"
	"github.com/spf13/cobra"
	"gitlab.cee.redhat.com/service/sdb-ocm-example-service/pkg/db/db_session"

	"gitlab.cee.redhat.com/service/sdb-ocm-example-service/pkg/config"
	"gitlab.cee.redhat.com/service/sdb-ocm-example-service/pkg/db"
)

var dbConfig = config.NewDatabaseConfig()

// migrate sub-command handles running migrations
func NewMigrateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migrate",
		Short: "Run OCM Example service data migrations",
		Long:  "Run OCM Example service data migrations",
		Run:   runMigrate,
	}

	dbConfig.AddFlags(cmd.PersistentFlags())
	cmd.PersistentFlags().AddGoFlagSet(flag.CommandLine)
	return cmd
}

func runMigrate(_ *cobra.Command, _ []string) {
	err := dbConfig.ReadFiles()
	if err != nil {
		glog.Fatal(err)
	}

	connection := db_session.NewProdFactory(dbConfig)
	if err := db.Migrate(connection.New(context.Background())); err != nil {
		glog.Fatal(err)
	}
}
