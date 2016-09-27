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
	cb    func()
	next  *promiseLink
}

type SPromise struct {
	err   interface{}
	value interface{}
	wg    *sync.WaitGroup

	m1    sync.Mutex // Ensures that the chain happens in order
	m2    sync.Mutex // Ensures that only one link is running at a time

	mList *sync.Cond

	head  *promiseLink
	tail  *promiseLink
}

var pType = reflect.TypeOf((*IPromise)(nil)).Elem()

func (p *SPromise) push(f func()) {
	p.m2.Lock()

	p.wg.Add(1)

	var l promiseLink
	l.cb = f

	if p.tail != nil {
		p.tail.next = &l
	} else {
		p.head = &l
	}

	p.tail = &l

	p.mList.Signal()
	p.m2.Unlock()
}

func (p *SPromise) pop() *promiseLink {

	if p.head == nil {
		p.mList.L.Lock()
		p.mList.Wait()
		p.mList.L.Unlock()
	}

	p.m2.Lock()

	n := p.head

	if n != nil {
		p.head = n.next
	}

	p.m2.Unlock()

	return n 
}


func (p *SPromise) run() {
	p.m1.Unlock()

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
	p.mList.Signal()

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
	p.mList = sync.NewCond(&m)

	p.wg = &sync.WaitGroup{}
	p.wg.Add(1)

	p.m1.Lock()

	go p.run()

	p.push(func() {
		p.value, p.err = cb()
	})

	return &p
}
