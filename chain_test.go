package itask

import (
	"fmt"
	"testing"
)

func TestChain(t *testing.T) {
	c := NewChain()
	c.Call(test1, 1).MapIf(func(x int) bool { return x > 0 }, -1).Call(test1, 2)
}

func test1(a int) {
	fmt.Println(a)
}

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
