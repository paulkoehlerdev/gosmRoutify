package astar

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"testing"
)

func readMatrixFromFile(fileName string) [][]float64 {
	file, _ := os.Open(fileName)
	defer file.Close()

	var output [][]float64
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		numbersStr := strings.Split(scanner.Text(), ",")
		numbers := make([]float64, len(numbersStr))
		for i, str := range numbersStr {
			n, _ := strconv.ParseFloat(str, 32)
			numbers[i] = float64(n)
		}
		output = append(output, numbers)
	}

	return output
}

type grapNode struct {
	Id       int
	Northing float64
	Easting  float64
}

func generateTestHeuristic(end grapNode, b *testing.B) func(node grapNode) float64 {
	return func(node grapNode) float64 {
		if b != nil {
			b.StopTimer()
			defer b.StartTimer()
		}

		diffN := end.Northing - node.Northing
		diffE := end.Easting - node.Easting
		return math.Sqrt(diffN*diffN + diffE*diffE)
	}
}

func generateTestConnections(A [][]float64, L []grapNode, b *testing.B) func(node grapNode) map[grapNode]float64 {
	return func(node grapNode) map[grapNode]float64 {
		if b != nil {
			b.StopTimer()
			defer b.StartTimer()
		}

		out := make(map[grapNode]float64)
		for id, weight := range A[node.Id] {
			if weight <= 0 {
				continue
			}
			out[L[id]] = weight
		}
		return out
	}
}

func TestAStar(t *testing.T) {
	aMatrix := readMatrixFromFile("../../../data/test/A.txt")
	LMatrix := readMatrixFromFile("../../../data/test/L.txt")

	nodes := make([]grapNode, len(LMatrix))
	for i, mat := range LMatrix {
		nodes[i] = grapNode{
			Id:       i,
			Northing: mat[0],
			Easting:  mat[1],
		}
	}

	path, err := AStar(nodes[99], nodes[1999], generateTestConnections(aMatrix, nodes, nil), generateTestHeuristic(nodes[1999], nil))
	if err != nil {
		t.Fatal(err)
		return
	}

	fmt.Println(len(path))

	str := ""
	for _, node := range path {
		str += fmt.Sprintf("%d; ", node.Id)
	}
	fmt.Println(str)
}

func BenchmarkAStar(b *testing.B) {
	aMatrix := readMatrixFromFile("../../../data/test/A.txt")
	LMatrix := readMatrixFromFile("../../../data/test/L.txt")

	nodes := make([]grapNode, len(LMatrix))
	for i, mat := range LMatrix {
		nodes[i] = grapNode{
			Id:       i,
			Northing: mat[0],
			Easting:  mat[1],
		}
	}

	connections, heuristic := generateTestConnections(aMatrix, nodes, b), generateTestHeuristic(nodes[1999], b)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		path, err := AStar(nodes[99], nodes[1999], connections, heuristic)
		if err != nil {
			b.Fatal(err)
		}

		if len(path) != 91 {
			b.Fatal()
		}
	}
}
