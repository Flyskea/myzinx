package skeapool

import (
	"errors"
	"sync"
	"sync/atomic"
)

const (
	// OPEN indicates that the pool is open to accept new tasks
	OPEN int32 = 1

	// CLOSED indicates that the pool is closed and won't accept new tasks
	CLOSED int32 = 0

	// DefaultPoolSize is a default size for number of workers in the pool
	DefaultPoolSize = 10
)

// ErrInvalidPoolSize indicates that the pool size is invalid
var ErrInvalidPoolSize = errors.New("invalid pool size: pool size must be a positive number")

// ErrInvalidPoolState indicates that the invalid pool state
var ErrInvalidPoolState = errors.New("pool is closed: cannot assign task to a closed pool")

// ErrNilFunction indicates that a nil function submitted
var ErrNilFunction = errors.New("cannot submit nil function()")

var ErrPoolOverload = errors.New("pool is full")

// pool represents a group of workers to whom tasks can be assigned.
type pool struct {
	// number of workers in the pool
	poolCapacity int32
	// queue to hold workers
	workerQueue *loopQueue
	// number of currently active workers
	activeWorkers int32
	// pool of available workers out of total poolCapacity
	availableWorkers sync.Pool
	// object which closes the pool and it can be called only once in the program scope
	closePool sync.Once
	// aovid data race
	lock sync.Mutex
	// wait for a available woker
	cond sync.Cond
	// represents the current state of the pool(OPEN/CLOSED)
	status int32
}

// NewPool returns an instance of pool with the size specified
func NewPool(newSize int) *pool {
	newPool := pool{
		poolCapacity: int32(newSize),
		status:       OPEN,
		workerQueue:  newWorkerLoopQueue(newSize),
	}
	newPool.cond = *sync.NewCond(&newPool.lock)
	newPool.availableWorkers.New = func() interface{} {
		return &worker{
			pool:     &newPool,
			taskChan: make(chan func(), 1),
		}
	}
	return &newPool
}

// retrieveWorker get a worker
func (p *pool) retrieveWorker() *worker {
	p.lock.Lock()
	defer p.lock.Unlock()
	w := p.workerQueue.detach()
	spawnWorker := func() {
		w = p.availableWorkers.Get().(*worker)
		w.run()
	}
	if w != nil {
		return w
	} else if capacity := p.PoolSize(); capacity > p.ActiveWorkers() {
		spawnWorker()
	} else {
		for w == nil {
			p.cond.Wait()
			var nw int
			if nw = p.ActiveWorkers(); nw == 0 {
				if !p.IsClosed() {
					spawnWorker()
				}
				return w
			}
			if w = p.workerQueue.detach(); w == nil {
				if nw < capacity {
					spawnWorker()
					return w
				}
			}
		}
	}
	return w
}

// revertWorker put the worker into workerQueue
func (p *pool) revertWorker(worker *worker) bool {
	p.lock.Lock()
	defer p.lock.Unlock()
	if capacity := p.PoolSize(); (capacity > 0 && p.ActiveWorkers() > capacity) || p.IsClosed() {
		return false
	}

	if err := p.workerQueue.insert(worker); err != nil {
		return false
	}
	// notify the waiters that there is a available woker
	p.cond.Signal()
	return true
}

// Done is called by a worker after completing its task
// when the workerQueue is full and notify there is a available worker
func (p *pool) done(w *worker) {
	p.availableWorkers.Put(w)
	p.decRunning()
	p.cond.Signal()
}

// incRunning increases the number of the currently running goroutines.
func (p *pool) incRunning() {
	atomic.AddInt32(&p.activeWorkers, 1)
}

// decRunning decreases the number of the currently running goroutines.
func (p *pool) decRunning() {
	atomic.AddInt32(&p.activeWorkers, -1)
}

// Close closes the pool and makes sure
// that no more tasks should be accepted
func (p *pool) Close() {
	p.closePool.Do(func() {
		atomic.StoreInt32(&p.status, CLOSED)
		p.workerQueue.reset()
		p.cond.Broadcast()
	})
}

// Submit submits a new task and assigns it to the pool
func (p *pool) Submit(task func()) error {
	if task == nil {
		return ErrNilFunction
	}

	if atomic.LoadInt32(&p.status) == CLOSED {
		return ErrInvalidPoolState
	}
	var w *worker
	if w = p.retrieveWorker(); w == nil {
		return ErrPoolOverload
	}
	w.taskChan <- task
	return nil
}

// AvailableWorkers returns available workers out of total workers
func (p *pool) AvailableWorkers() int {
	return p.PoolSize() - p.ActiveWorkers()
}

// ActiveWorkers returns number of active workers
func (p *pool) ActiveWorkers() int { return int(atomic.LoadInt32(&p.activeWorkers)) }

// PoolSize returns the pool size
func (p *pool) PoolSize() int { return int(atomic.LoadInt32(&p.poolCapacity)) }

// IsClosed return whether the pool is closed
func (p *pool) IsClosed() bool { return atomic.LoadInt32(&p.status) == CLOSED }
