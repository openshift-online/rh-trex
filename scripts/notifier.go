package main

import (
	"context"
	"flag"
	"fmt"
	"time"

	"github.com/openshift-online/rh-trex/cmd/ocm-example-service/environments"
	"github.com/spf13/pflag"
)

func init() {
	_ = flag.Set("logtostderr", "true")
	flags := pflag.CommandLine
	flags.AddGoFlagSet(flag.CommandLine)
}

func main() {
	// Parse flags
	pflag.Parse()

	err := environments.Environment().Initialize()
	if err != nil {
		fmt.Errorf("%s", err)
		return
	}

	env := environments.Environment()
	gorm := env.Database.SessionFactory.New(context.Background())

	for i := 0; i < 10; i++ {
		sql := fmt.Sprintf("select pg_notify('events','%s')", time.Now().String())

		fmt.Printf("attempting: %s\n", sql)

		err = gorm.Exec(sql).Error
		if err != nil {
			fmt.Errorf("%s", err)
			return
		}

		time.Sleep(500 * time.Millisecond)

	}
}
