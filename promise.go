package promise

import (
	"reflect"
	"sync"
)

type promiseCallback func() (response interface{}, err interface{})
type promiseThenCallback func(interface{}) (response interface{}, err interface{})

type IPromise interface {
	Then(ret promiseThenCallback) *SPromise
	Wait() (interface{}, interface{})
	run()
	push(f func())
	pop() *promiseLink
}

type promiseLink struct {
	cb   func()
	next *promiseLink
}

type SPromise struct {
	err   interface{}
	value interface{}
	wg    *sync.WaitGroup

	listMutex sync.RWMutex
	listCond  *sync.Cond

	head *promiseLink
	tail *promiseLink
}

var pType = reflect.TypeOf((*IPromise)(nil)).Elem()

func (p *SPromise) push(f func()) {
	p.listMutex.Lock()

	p.wg.Add(1)

	var l promiseLink
	l.cb = f

	if p.head != nil {
		p.tail.next = &l
	} else {
		p.head = &l
	}

	p.tail = &l

	p.listCond.Signal()
	p.listMutex.Unlock()
}

func (p *SPromise) pop() *promiseLink {
	p.listMutex.RLock()

	if p.head == nil {
		p.listMutex.RUnlock()
		p.listCond.L.Lock()
		p.listCond.Wait()
		p.listCond.L.Unlock()
	} else {
		p.listMutex.RUnlock()
	}

	p.listMutex.Lock()

	n := p.head

	if n != nil {
		p.head = n.next
	}

	p.listMutex.Unlock()

	return n
}

func (p *SPromise) run() {
	n := p.pop()

	for n != nil {
		n.cb()

		if p.err == nil {

			vType := reflect.TypeOf(p.value)

			if vType.Kind() == reflect.Ptr && vType.Implements(pType) {
				p.value, p.err = p.value.(*SPromise).Wait()
			}
		}

		p.wg.Done()

		n = p.pop()
	}
}

func (p *SPromise) Wait() (interface{}, interface{}) {
	p.wg.Done()
	p.wg.Wait()
	p.listCond.Signal()

	return p.value, p.err
}

func (p *SPromise) Then(ret promiseThenCallback) *SPromise {
	p.push(func() {
		if p.err == nil {
			p.value, p.err = ret(p.value)
		}
	})

	return p
}

func Create(cb promiseCallback) *SPromise {
	var p SPromise

	var m sync.Mutex
	p.listCond = sync.NewCond(&m)

	p.wg = &sync.WaitGroup{}
	p.wg.Add(1)

	go p.run()

	p.push(func() {
		p.value, p.err = cb()
	})

	return &p
}
