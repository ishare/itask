package itask

import (
	"log"
	"sync"
	"time"
)

type Task struct {
	preHandlers  []*Handler // sync execute
	handlers     []*Handler // can be sync or parallel
	postHandlers []*Handler // sync execute
	onRecover    func(*RecoverMsg)
}

type RecoverMsg struct {
	FuncName  string
	StartTime int64 // in nano seconds
	Err       interface{}
}

func NewTask() *Task {
	res := new(Task)
	return res
}

func (t *Task) PreProcess(f interface{}, args ...interface{}) *Task {
	h := NewHandler(f, args...)
	t.preHandlers = append(t.preHandlers, h)
	return t
}

func (t *Task) Process(f interface{}, args ...interface{}) *Task {
	h := NewHandler(f, args...)
	t.handlers = append(t.handlers, h)
	return t
}

func (t *Task) PostProcess(f interface{}, args ...interface{}) *Task {
	h := NewHandler(f, args...)
	t.postHandlers = append(t.postHandlers, h)
	return t
}

func (t *Task) SetRecover(f func(*RecoverMsg)) {
	t.onRecover = f
}

func (t *Task) deferFunc(wg *sync.WaitGroup, h *Handler, startTime int64) {
	if wg != nil {
		wg.Done()
	}
	if r := recover(); r != nil {
		if t.onRecover != nil {
			t.onRecover(&RecoverMsg{GetFuncName(h), startTime, r})
		} else {
			log.Println(GetFuncName(h), r)
		}
	}
}

// sync execution
func (t *Task) Run() {
	var h Handler
	defer t.deferFunc(nil, &h, time.Now().UnixNano())
	for i := range t.preHandlers {
		h = *t.preHandlers[i]
		h.Call()
	}
	for i := range t.handlers {
		h = *t.handlers[i]
		h.Call()
	}
	for i := range t.postHandlers {
		h = *t.postHandlers[i]
		h.Call()
	}
}

// parallel execute, waiting for all go routine finished
func (t *Task) Parallel() {
	var h Handler
	defer t.deferFunc(nil, &h, time.Now().UnixNano())
	for i := range t.preHandlers {
		h = *t.preHandlers[i]
		h.Call()
	}

	wg := &sync.WaitGroup{}
	for i := range t.handlers {
		wg.Add(1)
		go func(f *Handler) {
			defer t.deferFunc(wg, f, time.Now().UnixNano())
			f.Call()
		}(t.handlers[i])
	}
	wg.Wait()

	for i := range t.postHandlers {
		h = *t.postHandlers[i]
		h.Call()
	}
}

// parallel execute with timeout
func (t *Task) ParallelWithTimeout(timeout time.Duration) error {
	var h Handler
	defer t.deferFunc(nil, &h, time.Now().UnixNano())
	for i := range t.preHandlers {
		h = *t.preHandlers[i]
		h.Call()
	}

	wg := &sync.WaitGroup{}
	for i := range t.handlers {
		wg.Add(1)
		go func(f *Handler) {
			defer t.deferFunc(wg, f, time.Now().UnixNano())
			f.Call()
		}(t.handlers[i])
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

	for i := range t.postHandlers {
		h = *t.postHandlers[i]
		h.Call()
	}
	return nil
}

// submit async job
func (t *Task) Async() {
	// TODO: write manager to schedule job
	go func() {
		var h Handler
		defer t.deferFunc(nil, &h, time.Now().UnixNano())
		for i := range t.preHandlers {
			h = *t.preHandlers[i]
			h.Call()
		}

		for i := range t.handlers {
			go func(f *Handler) {
				defer t.deferFunc(nil, f, time.Now().UnixNano())
				f.Call()
			}(t.handlers[i])
		}

		for i := range t.postHandlers {
			h = *t.postHandlers[i]
			h.Call()
		}
	}()
}
