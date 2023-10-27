package integration

import (
	"flag"
	"os"
	"runtime"
	"testing"

	"github.com/golang/glog"

	"github.com/openshift-online/rh-trex/test"
)

func TestMain(m *testing.M) {
	flag.Parse()
	glog.Infof("Starting integration test using go version %s", runtime.Version())
	helper := test.NewHelper(&testing.T{})
	exitCode := m.Run()
	helper.Teardown()
	os.Exit(exitCode)
}
