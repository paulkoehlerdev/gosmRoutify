package osmpbfreader_test

import (
	"errors"
	"fmt"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/osmpbfreader"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/osmpbfreader/osmpbfreaderdata"
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

	decoder := osmpbfreader.New(file)
	err = decoder.Start(osmpbfreader.ProcsCount(runtime.GOMAXPROCS(-1)))
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
		case *osmpbfreaderdata.Node:
			nodeCount++
		case *osmpbfreaderdata.Way:
			wayCount++
		case *osmpbfreaderdata.Relation:
			relationCount++
		default:
			t.Fatalf("unexpected data type: %T", d)
		}
	}

	fmt.Printf("node count: %d\n", nodeCount)
	fmt.Printf("way count: %d\n", wayCount)
	fmt.Printf("relation count: %d\n", relationCount)
}
