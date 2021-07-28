package dropcore

import "sync"

type ThreadSafeRingBuffer struct {
	RingBuffer
	lock sync.Mutex
}

func NewTSRingBuffer(size int, allowOverwrite bool) *ThreadSafeRingBuffer {
	return &ThreadSafeRingBuffer{
		RingBuffer: RingBuffer{
			rp:   0,
			wp:   0,
			sz:   0,
			mx:   uint64(size),
			ao:   allowOverwrite,
			data: make([]interface{}, size),
		},
	}
}

func (r *ThreadSafeRingBuffer) Push(i interface{}) {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.RingBuffer.Push(i)
}

func (r *ThreadSafeRingBuffer) Pop() (interface{}, bool) {
	r.lock.Lock()
	defer r.lock.Unlock()
	return r.RingBuffer.Pop()
}

func (r *ThreadSafeRingBuffer) Peek() (interface{}, bool) {
	r.lock.Lock()
	defer r.lock.Unlock()
	return r.RingBuffer.Peek()
}

func (r *ThreadSafeRingBuffer) Count() uint64 {
	r.lock.Lock()
	defer r.lock.Unlock()
	return r.RingBuffer.Count()
}

func (r *ThreadSafeRingBuffer) Clear() {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.RingBuffer.Clear()
}
