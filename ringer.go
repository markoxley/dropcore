package dropcore

type Ringer interface {
	Push(i interface{})
	Pop() (interface{}, bool)
	Peek() (interface{}, bool)
	Count() uint64
	Clear()
}
