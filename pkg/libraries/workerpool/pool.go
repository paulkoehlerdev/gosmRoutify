package workerpool

import (
	"io"
)

type ProcsCount int

type Pool[T any] interface {
	Start()
	Stop()
	Submit(func(id int) T)
	Result() (T, error)
}

type impl[T any] struct {
	jobqueue      chan func(id int) T
	workerqueues  []chan func(id int) T
	workerresults []chan T
	resultqueue   chan T
	workers       int
	done          chan struct{}
}

func New[T any](size ProcsCount) Pool[T] {
	return &impl[T]{
		jobqueue:    make(chan func(id int) T, size),
		resultqueue: make(chan T, size),
		workers:     int(size),
		done:        make(chan struct{}),
	}
}

func (i *impl[T]) Start() {
	i.workerqueues = make([]chan func(id int) T, i.workers)
	i.workerresults = make([]chan T, i.workers)

	for iter := 0; iter < i.workers; iter++ {
		i.workerqueues[iter] = make(chan func(id int) T, i.workers)
		i.workerresults[iter] = make(chan T, i.workers)

		go i.worker(iter)
	}

	go i.workerHandler()
	go i.workerOutputHandler()
}

func (i *impl[T]) Stop() {
	close(i.done)
}

func (i *impl[T]) workerHandler() {
	workerIndex := 0
	for {
		select {
		case job := <-i.jobqueue:
			i.workerqueues[workerIndex] <- job
		case <-i.done:
			close(i.jobqueue)
			return
		}
		workerIndex = (workerIndex + 1) % len(i.workerqueues)
	}
}

func (i *impl[T]) workerOutputHandler() {
	workerIndex := 0
	for {
		select {
		case result := <-i.workerresults[workerIndex]:
			i.resultqueue <- result
		case <-i.done:
			close(i.resultqueue)
			return
		}
		workerIndex = (workerIndex + 1) % len(i.workerqueues)
	}
}

func (i *impl[T]) worker(id int) {
	for {
		select {
		case job := <-i.workerqueues[id]:
			i.workerresults[id] <- job(id)
		case <-i.done:
			close(i.workerqueues[id])
			close(i.workerresults[id])
			return
		}
	}
}

func (i *impl[T]) Submit(job func(id int) T) {
	i.jobqueue <- job
}

func (i *impl[T]) Result() (T, error) {
	if result, ok := <-i.resultqueue; ok {
		return result, nil
	} else {
		return result, io.EOF
	}
}
