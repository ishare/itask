package chain

import (
	"fmt"
	"testing"
)

func TestTask(t *testing.T) {
	task := NewTask()
	task.PreProcess(test2, 1)
	task.Process(test2, 2)
	task.PostProcess(test2, 3)
	task.Run()      // sync, execute one by one
	task.Parallel() // run parallel, wait for all go routine finish
	task.Async()    // async, return immediately without waiting for response
}

func test2(a int) {
	fmt.Println(a)
}
