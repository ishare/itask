package chain

import (
	"sync"
	"time"
)

type Task struct {
	preHandlers  []*Handler // sync execute
	handlers     []*Handler // can be sync or parallel
	postHandlers []*Handler // sync execute
}

func NewTask() *Task {
	res := new(Task)
	return res
}

func (c *Task) PreProcess(f interface{}, args ...interface{}) *Task {
	h := NewHandler(f, args...)
	c.preHandlers = append(c.preHandlers, h)
	return c
}

func (c *Task) Process(f interface{}, args ...interface{}) *Task {
	h := NewHandler(f, args...)
	c.handlers = append(c.handlers, h)
	return c
}

func (c *Task) PostProcess(f interface{}, args ...interface{}) *Task {
	h := NewHandler(f, args...)
	c.postHandlers = append(c.postHandlers, h)
	return c
}

// sync execution
func (c *Task) Run() {
	for _, f := range c.preHandlers {
		f.Call()
	}
	for _, f := range c.handlers {
		f.Call()
	}
	for _, f := range c.postHandlers {
		f.Call()
	}
}

// parallel execute, waiting for all go routine finished
func (c *Task) Parallel() {
	for _, f := range c.preHandlers {
		f.Call()
	}
	wg := &sync.WaitGroup{}
	for i := range c.handlers {
		wg.Add(1)
		go func(f *Handler) {
			f.Call()
			wg.Done()
		}(c.handlers[i])
	}
	wg.Wait()
	for _, f := range c.postHandlers {
		f.Call()
	}
}

// parallel execute with timeout
func (c *Task) ParallelWithTimeout(timeout time.Duration) error {
	for _, f := range c.preHandlers {
		f.Call()
	}
	wg := &sync.WaitGroup{}
	for i := range c.handlers {
		wg.Add(1)
		go func(f *Handler) {
			f.Call()
			wg.Done()
		}(c.handlers[i])
	}

	done := make(chan int)
	go func() {
		wg.Wait()
		close(done)
	}()
	select {
	case <-done: // all done
	case <-time.After(timeout): // timeout
		return ErrExecuteTimeout
	}
	for _, f := range c.postHandlers {
		f.Call()
	}
	return nil
}

// submit async job
func (c *Task) Async() {
	// TODO: write manager to schedule job
	go func() {
		for _, f := range c.preHandlers {
			f.Call()
		}
		for i := range c.handlers {
			go func(f *Handler) {
				f.Call()
			}(c.handlers[i])
		}
		for _, f := range c.postHandlers {
			f.Call()
		}
	}()
}
