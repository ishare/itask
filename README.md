# itask
Go task manager, including call chain, sync, async and parallel executor.

Usage:
-----------------
**basic chain call:**

If any `MapIf` func return false, the following call will not execute.
```go
func TestChain(t *testing.T) {
	NewChain().Call(test1, 1).MapIf(func(x int) bool { return x > 0 }, -1).Call(test1, 2)
}

func test1(a int) {
	fmt.Println(a)
}
```

**bool chain call:**

Every funcs must return `bool` value, and if previous func return false, the following funcs will not execute.
```go
func TestBoolChainCall(t *testing.T) {
	NewBoolChain().Call(less, 1, 2).Call(less, 5, 2).Call(2, 3)
}

func TestBoolChain(t *testing.T) {
	c := NewBoolChain()
	c.Append(less, 1, 2)
	c.Append(less, 5, 2)
	c.Prepend(less, 2, 3)
	fmt.Println(c.Run())
}

func less(a, b int) bool {
	fmt.Printf("compare %d, %d\n", a, b)
	return a < b
}
```

**task schedule:**

Register funcs, control process in pipeline, parallel, or async execute.
```go
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
```
