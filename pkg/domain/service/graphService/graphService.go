package graphService

import (
	"errors"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/repository/graphRepository"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/repository/temporaryOSMDataRepository"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/domain/value/graph"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/logging"
	"io"
)

type GraphService interface {
	AddOSMData(point any) error
	BuildGraph()
}

type impl struct {
	repository         graphRepository.GraphRepository
	tempDataRepository temporaryOSMDataRepository.TemporaryOSMDataRepository
	logger             logging.Logger
}

func New(repository graphRepository.GraphRepository, tempDataRepository temporaryOSMDataRepository.TemporaryOSMDataRepository, logger logging.Logger) GraphService {
	return &impl{
		repository:         repository,
		tempDataRepository: tempDataRepository,
		logger:             logger,
	}
}

func (i *impl) AddOSMData(data any) error {
	return i.tempDataRepository.AddOSMData(data)
}

func (i *impl) BuildGraph() {
	nodeToWayMap := make(map[int64][]graph.GraphID)

	i.logger.Info().Msgf("building graph")
	i.logger.Info().Msgf("adding ways")

	counter := 0

	for {
		counter++
		if counter%10000 == 0 {
			i.logger.Info().Msgf("imported %d K ways", counter/10000)
		}

		way, err := i.tempDataRepository.Next()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			} else {
				i.logger.Error().Msgf("error while reading way: %s", err.Error())
				return
			}
		}

		highwayClass, ok := way.Tags["highway"]
		if !ok {
			highwayClass = "unclassified"
		}

		startingNode := i.tempDataRepository.FindNode(way.NodeIDs[0])
		if startingNode == nil {
			continue
		}

		lID := graph.LevelIDFromHighwayClass(highwayClass)
		gID := i.repository.AddWay(way, graph.TileIDFromNode(startingNode, lID))

		for _, nodeID := range way.NodeIDs {
			nodeToWayMap[nodeID] = append(nodeToWayMap[nodeID], gID)
		}
	}

	i.logger.Info().Msgf("adding intersections")
	counter = 0

	for nodeID, wayIDs := range nodeToWayMap {
		if len(wayIDs) < 2 {
			continue
		}

		counter++
		if counter%1000 == 0 {
			i.logger.Info().Msgf("imported %d K intersections", counter/1000)
		}

		node := i.tempDataRepository.FindNode(nodeID)
		if node == nil {
			continue
		}

		for iter := 0; iter < len(wayIDs); iter++ {
			for j := iter + 1; j < len(wayIDs); j++ {
				i.repository.AddIntersection(node, wayIDs[iter], wayIDs[j])
			}
		}
	}

	i.logger.Info().Msgf("building graph finished")

	i.logger.Info().Msgf("cleaning up")
	i.tempDataRepository.Cleanup()
	i.logger.Info().Msgf("cleaning up finished")
}
