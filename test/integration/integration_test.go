package integration

import (
	"flag"
	"os"
	"runtime"
	"testing"

	"github.com/golang/glog"

	"gitlab.cee.redhat.com/service/sdb-ocm-example-service/test"
)

func TestMain(m *testing.M) {
	flag.Parse()
	glog.Infof("Starting integration test using go version %s", runtime.Version())
	helper := test.NewHelper(&testing.T{})
	exitCode := m.Run()
	helper.Teardown()
	os.Exit(exitCode)
}
