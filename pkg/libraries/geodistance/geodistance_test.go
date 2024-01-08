package geodistance_test

import (
	"fmt"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/geodistance"
	"testing"
)

func TestCalcDistanceInMeters(t *testing.T) {
	dist := geodistance.CalcDistanceInMeters(
		geodistance.NewPoint(48+07/60+22.5664/(60*60), 11+33/60+22.1335/(60*60)),
		geodistance.NewPoint(48+07/60+22.5664/(60*60), 11+33/60+22.1335/(60*60)),
	)

	fmt.Printf("dist = %f\n", dist)
}
