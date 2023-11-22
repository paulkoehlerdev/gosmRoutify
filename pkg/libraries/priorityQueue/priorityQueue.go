package priorityQueue

import (
	"container/heap"
	"golang.org/x/exp/constraints"
)

type number interface {
	constraints.Float | constraints.Integer
}

type PriorityQueue[K any, N number] interface {
	Push(item K, priority N)
	Pop() K
	Len() int
}

type impl[K any, N number] struct {
	pq priorityQueue[K, N]
}

func NewPriorityQueue[K any, N number]() PriorityQueue[K, N] {
	pq := priorityQueue[K, N]{}
	heap.Init(&pq)
	return &impl[K, N]{pq: pq}
}

func (i *impl[K, N]) Push(item K, priority N) {
	newQueueItem := &queueItem[K, N]{
		Value:    item,
		Priority: priority,
		Index:    -1,
	}

	heap.Push(&i.pq, newQueueItem)
}

func (i *impl[K, N]) Pop() K {
	removedQueueItem := heap.Pop(&i.pq).(*queueItem[K, N])
	return removedQueueItem.Value
}

func (i *impl[K, N]) Len() int {
	return i.pq.Len()
}

type queueItem[K any, N number] struct {
	Value    K
	Priority N
	Index    int
}

type priorityQueue[K any, N number] []*queueItem[K, N]

func (pq priorityQueue[K, N]) Less(i, j int) bool {
	return pq[i].Priority > pq[j].Priority
}

func (pq priorityQueue[K, N]) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].Index = i
	pq[j].Index = j
}

func (pq *priorityQueue[K, N]) Push(x any) {
	item := x.(*queueItem[K, N])
	item.Index = len(*pq)
	*pq = append(*pq, item)
}

func (pq *priorityQueue[K, N]) Pop() any {
	lastIndex := len(*pq) - 1
	item := (*pq)[lastIndex]
	(*pq)[lastIndex] = nil
	item.Index = -1
	*pq = (*pq)[0:lastIndex]
	return item
}

func (pq priorityQueue[K, N]) Len() int {
	return len(pq)
}
