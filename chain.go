package itask

type Chain struct {
	terminate bool
}

func NewChain() *Chain {
	res := new(Chain)
	return res
}

func (c *Chain) MapIf(f interface{}, args ...interface{}) *Chain {
	if c.terminate {
		return c
	}
	h := NewHandler(f, args...)
	c.terminate = !h.BoolCall()
	return c
}

func (c *Chain) Call(f interface{}, args ...interface{}) *Chain {
	if c.terminate {
		return c
	}
	h := NewHandler(f, args...)
	h.Call()
	return c
}

func (c *Chain) Result() bool {
	return !c.terminate
}
