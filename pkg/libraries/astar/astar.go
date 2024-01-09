package astar

import (
	"fmt"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/arrayutil"
	"github.com/paulkoehlerdev/gosmRoutify/pkg/libraries/priorityQueue"
	"golang.org/x/exp/constraints"
)

type number interface {
	constraints.Float | constraints.Integer
}

func AStar[K comparable, N number](start K, end K, connections func(element K) map[K]N, heuristic func(K) N) ([]K, N, error) {
	open := priorityQueue.NewPriorityQueue[K, N]()
	open.Push(start, 0)

	parent := make(map[K]K)

	gScore := make(map[K]N)
	gScore[start] = 0

	count := 0
	for open.Len() > 0 {
		count++
		current := open.Pop()

		if current == end {
			return generatePath(parent, end), gScore[current], nil
		}

		neighbors := connections(current)
		for neighbor, weight := range neighbors {
			tentativeScore := gScore[current] + weight
			if score, ok := gScore[neighbor]; !ok || tentativeScore < score {
				parent[neighbor] = current
				gScore[neighbor] = tentativeScore
				open.Push(neighbor, -(tentativeScore + heuristic(neighbor)))
			}
		}
	}

	return nil, 0, fmt.Errorf("error: no route found, after %d iterations", count)
}

func generatePath[K comparable](parent map[K]K, end K) []K {
	var path []K
	for current, ok := end, true; ok; current, ok = parent[current] {
		path = append(path, current)
	}
	return arrayutil.Reverse(path)
}
