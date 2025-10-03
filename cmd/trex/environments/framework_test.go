package environments

import (
	"os/exec"
	"reflect"
	"testing"

	"github.com/spf13/pflag"
)

func BenchmarkGetDynos(b *testing.B) {
	b.ReportAllocs()
	fn := func(b *testing.B) {
		cmd := exec.Command("ocm", "get", "/api/rh-trex/v1/dinosaurs", "params='size=2'")
		_, err := cmd.CombinedOutput()
		if err != nil {
			b.Errorf("ERROR %+v", err)
		}
	}
	for n := 0; n < b.N; n++ {
		fn(b)
	}
}

func TestLoadServices(t *testing.T) {
	env := Environment()
	// Override environment name
	env.Name = "testing"
	err := env.AddFlags(pflag.CommandLine)
	if err != nil {
		t.Errorf("Unable to add flags for testing environment: %s", err.Error())
		return
	}
	pflag.Parse()
	err = env.Initialize()
	if err != nil {
		t.Errorf("Unable to load testing environment: %s", err.Error())
		return
	}

	s := reflect.ValueOf(env.Services)

	for i := 0; i < s.NumField(); i++ {
		field := s.Field(i)
		// Skip non-pointer fields (like mutex)
		if field.Kind() == reflect.Ptr || field.Kind() == reflect.Interface || 
		   field.Kind() == reflect.Map || field.Kind() == reflect.Slice || 
		   field.Kind() == reflect.Chan || field.Kind() == reflect.Func {
			if field.IsNil() {
				t.Errorf("Service field %d is nil", i)
			}
		}
	}
}
