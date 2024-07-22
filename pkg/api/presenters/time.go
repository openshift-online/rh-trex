package presenters

import (
	"time"

	"github.com/openshift-online/rh-trex/pkg/util"
)

func PresentTime(t time.Time) *time.Time {
	if t.IsZero() {
		return util.ToPtr(time.Time{})
	}
	return util.ToPtr(t.Round(time.Microsecond))
}
