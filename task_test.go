package itask

import (
	"fmt"
	"testing"
	"time"
)

func TestTask(t *testing.T) {
	task := NewTask()
	task.PreProcess(test2, 1)
	task.Process(test2, 2)
	task.Process(testPanic)
	task.PostProcess(test2, 3)
	task.SetRecover(testRecover)
	task.Run()      // sync, execute one by one
	task.Parallel() // run parallel, wait for all go routine finish
	task.Async()    // async, return immediately without waiting for response
	time.Sleep(100 * time.Millisecond)
}

func test2(a int) {
	fmt.Println(a)
}

func testPanic() {
	panic("aaa")
}

func testRecover(r *RecoverMsg) {
	fmt.Println("recover: ", r.FuncName, r.StartTime, r.Err)
}
