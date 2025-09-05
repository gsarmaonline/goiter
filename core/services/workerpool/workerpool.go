package workerpool

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/gsarmaonline/goiter/core/services/cache"
	"gorm.io/gorm"
)

type (
	EventTypeT string

	WorkerJobHandler func(input *Event) error

	WorkerPool struct {
		Namespace string

		db *gorm.DB

		redisPool *redis.Pool

		jobs         map[EventTypeT][]WorkerJobHandler
		asyncEventCh chan *Event
		eventCh      chan *Event

		workerCount int
		workers     []*Worker
	}

	Worker struct {
		ID     string
		wp     *WorkerPool
		exitCh chan bool
	}

	Event struct {
		EventType EventTypeT

		SourceType string
		SourceID   string

		Data      interface{}
		EmittedAt time.Time

		UserID uint
	}
)

func NewWorker(id string, wp *WorkerPool, exitCh chan bool) (w *Worker, err error) {
	w = &Worker{
		ID:     id,
		wp:     wp,
		exitCh: exitCh,
	}

	return
}

func (w *Worker) Start() (err error) {
	for {
		select {
		case event := <-w.wp.asyncEventCh:
			w.handleEvent(event)
		case event := <-w.wp.eventCh:
			w.handleEvent(event)
		case <-w.exitCh:
			return
		}
	}
	return
}

func (w *Worker) handleEvent(event *Event) (err error) {
	handlers, ok := w.wp.jobs[event.EventType]
	if !ok {
		return
	}
	for _, handler := range handlers {
		err = w.runHandler(event, handler)
		if err != nil {
			return
		}
	}
	return
}

func (w *Worker) runHandler(event *Event, handler WorkerJobHandler) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic recovered in handler: %v", r)
		}
	}()
	err = handler(event)
	if err != nil {
		return
	}
	return
}

func NewWorkerPool(namespace string, db *gorm.DB) (wp *WorkerPool, err error) {
	wp = &WorkerPool{
		Namespace:    namespace,
		db:           db,
		redisPool:    cache.NewCache().Pool,
		workerCount:  4,
		asyncEventCh: make(chan *Event, 1000),
		eventCh:      make(chan *Event),
		jobs:         make(map[EventTypeT][]WorkerJobHandler),
	}
	for i := 0; i < wp.workerCount; i++ {
		var (
			worker *Worker
		)
		exitCh := make(chan bool)
		worker, err = NewWorker(fmt.Sprintf("worker-%d", i), wp, exitCh)
		if err != nil {
			return
		}
		wp.workers = append(wp.workers, worker)
	}

	return
}

func (wp *WorkerPool) Start() {

	for _, worker := range wp.workers {
		go worker.Start()
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, os.Kill)
	<-signalChan

}

func (wp *WorkerPool) RegisterJob(eventName EventTypeT, handler WorkerJobHandler) (err error) {
	var (
		existing  []WorkerJobHandler
		isPresent bool
	)
	if existing, isPresent = wp.jobs[eventName]; !isPresent {
		existing = []WorkerJobHandler{handler}
	}
	existing = append(existing, handler)
	wp.jobs[eventName] = existing

	return
}

func (wp *WorkerPool) EmitAsync(event *Event) (err error) {
	wp.asyncEventCh <- event
	return
}

func (wp *WorkerPool) EmitSync(event *Event) (err error) {
	wp.eventCh <- event
	return
}

func (wp *WorkerPool) Shutdown() (err error) {
	for _, worker := range wp.workers {
		worker.exitCh <- true
	}
	close(wp.asyncEventCh)
	close(wp.eventCh)
	return
}
