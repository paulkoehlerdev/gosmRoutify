package datadecoder

import (
	"fmt"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/osmpbfreader/getdata"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/osmpbfreader/osmpbfreaderdata"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/osmpbfreader/osmproto"
	"google.golang.org/protobuf/proto"
	"time"
)

func DecodeHeaderBlock(blob *osmproto.Blob) (*osmpbfreaderdata.Header, error) {
	data, err := getdata.GetData(blob)
	if err != nil {
		return nil, err
	}

	headerBlock := new(osmproto.HeaderBlock)
	if err := proto.Unmarshal(data, headerBlock); err != nil {
		return nil, err
	}

	parseCapabilities := map[string]bool{
		"OsmSchema-V0.6": true,
		"DenseNodes":     true,
	}

	requiredFeatures := headerBlock.GetRequiredFeatures()
	for _, feature := range requiredFeatures {
		if !parseCapabilities[feature] {
			return nil, fmt.Errorf("parser does not have %s capability", feature)
		}
	}

	header := &osmpbfreaderdata.Header{
		RequiredFeatures:                 headerBlock.GetRequiredFeatures(),
		OptionalFeatures:                 headerBlock.GetOptionalFeatures(),
		WritingProgram:                   headerBlock.GetWritingprogram(),
		Source:                           headerBlock.GetSource(),
		OsmosisReplicationBaseUrl:        headerBlock.GetOsmosisReplicationBaseUrl(),
		OsmosisReplicationSequenceNumber: headerBlock.GetOsmosisReplicationSequenceNumber(),
	}

	if headerBlock.OsmosisReplicationTimestamp != nil {
		header.OsmosisReplicationTimestamp = time.Unix(*headerBlock.OsmosisReplicationTimestamp, 0)
	}
	if headerBlock.Bbox != nil {
		header.BoundingBox = &osmpbfreaderdata.BoundingBox{
			Left:   degMultiplier * float64(*headerBlock.Bbox.Left),
			Right:  degMultiplier * float64(*headerBlock.Bbox.Right),
			Bottom: degMultiplier * float64(*headerBlock.Bbox.Bottom),
			Top:    degMultiplier * float64(*headerBlock.Bbox.Top),
		}
	}

	return header, nil
}
