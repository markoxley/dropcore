package dropcore

type RingBuffer struct {
	rp   uint64
	wp   uint64
	sz   uint64
	mx   uint64
	sp   bool
	ao   bool
	data []interface{}
}

func NewRingBuffer(size int, allowOverwrite bool) *RingBuffer {
	return &RingBuffer{
		rp:   0,
		wp:   0,
		sz:   0,
		mx:   uint64(size),
		ao:   allowOverwrite,
		data: make([]interface{}, size),
	}
}

func (r *RingBuffer) Push(i interface{}) {
	if r.sz > 0 && r.rp == r.wp && !r.ao {
		return
	}
	r.data[r.wp] = i
	if r.sz > 0 && r.rp == r.wp {
		r.rp = (r.rp + 1) % r.mx
	}
	r.wp = (r.wp + 1) % r.mx
	r.sz++
	if r.sz >= r.mx {
		r.sz = r.mx
	}
}

func (r *RingBuffer) Pop() (interface{}, bool) {
	if r.sz == 0 {
		return nil, false
	}
	res := r.data[r.rp]
	r.rp = (r.rp + 1) % r.mx
	r.sz--
	return res, true
}

func (r *RingBuffer) Peek() (interface{}, bool) {
	if r.sz == 0 {
		return nil, false
	}
	return r.data[r.rp], true

}

func (r *RingBuffer) Count() uint64 {
	return r.sz
}

func (r *RingBuffer) Clear() {
	r.wp = 0
	r.rp = 0
	r.sz = 0
}
