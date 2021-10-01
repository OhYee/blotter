package pool

// Pool is goroutine pool
type Pool struct {
	cap     int8
	running int8
	channel chan Func
}

// Func is gorouting function
type Func func()

// New returns a gorouting pool
func New(cap int) *Pool {
	p := &Pool{
		channel: make(chan Func, cap),
	}
	for i := 0; i < cap; i++ {
		go p.Thread()
	}
	return p
}

// Thread is the backend threading
func (p *Pool) Thread() {
	for {
		select {
		case f := <-p.channel:
			f()
		}
	}
}

// Do a threading
func (p *Pool) Do(f Func) {
	p.channel <- f
}

var pool = New(32)

// Do a threading
func Do(f Func) {
	pool.Do(f)
}
