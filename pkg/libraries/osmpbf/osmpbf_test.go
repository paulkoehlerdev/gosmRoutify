package osmpbf_test

import (
	"errors"
	"fmt"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/osmpbf"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/osmpbf/osmpbfData"
	"io"
	"os"
	"runtime"
	"testing"
)

func TestImpl_Decode(t *testing.T) {
	filepath := "../../../resources/sample/sample.pbf"
	file, err := os.OpenFile(filepath, os.O_RDONLY, 0)
	if err != nil {
		t.Fatalf("error opening file: %s", err.Error())
	}

	decoder := osmpbf.New(file)
	err = decoder.Start(osmpbf.ProcsCount(runtime.GOMAXPROCS(-1)))
	if err != nil {
		t.Fatalf("error starting decoder: %s", err.Error())
	}

	nodeCount := 0
	wayCount := 0
	relationCount := 0

	for {
		data, err := decoder.Decode()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			t.Fatalf("error decoding: %s", err.Error())
		}

		switch d := data.(type) {
		case *osmpbfData.Node:
			nodeCount++
		case *osmpbfData.Way:
			wayCount++
		case *osmpbfData.Relation:
			relationCount++
		default:
			t.Fatalf("unexpected data type: %T", d)
		}
	}

	fmt.Printf("node count: %d\n", nodeCount)
	fmt.Printf("way count: %d\n", wayCount)
	fmt.Printf("relation count: %d\n", relationCount)
}
