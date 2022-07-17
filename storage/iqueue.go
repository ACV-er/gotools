package storage

type Queue interface {
	Push(interface{})
	Pop() interface{}
	Front() interface{}
	Back() interface{}
	Size() int
	Clean()
}
