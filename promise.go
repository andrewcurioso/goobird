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
}

type SPromise struct {
	err   interface{}
	value interface{}
	wg    *sync.WaitGroup
	m1    sync.Mutex // Ensures that the chain happens in order
	m2    sync.Mutex // Ensures that only one link is running at a time
}

var pType = reflect.TypeOf((*IPromise)(nil)).Elem()

func run(p *SPromise, f func() ) {
	p.m1.Unlock()
	p.m2.Lock()

	f()

	if p.err == nil {

		vType := reflect.TypeOf(p.value)

		if vType.Kind() == reflect.Ptr && vType.Implements(pType) {
			p.value, p.err = p.value.(*SPromise).Wait()
		}
	}

	p.wg.Done()
	p.m2.Unlock()
}

func (p *SPromise) Wait() (interface{}, interface{}) {
	p.wg.Done()
	p.wg.Wait()

	return p.value, p.err
}

func (p *SPromise) Then(ret promiseThenCallback) *SPromise {
	p.wg.Add(1)

	p.m1.Lock()
	go run(p, func() {
		if p.err == nil {
			p.value, p.err = ret(p.value)
		}
	})

	return p
}

func Create(cb promiseCallback) *SPromise {
	var p SPromise

	p.wg = &sync.WaitGroup{}
	p.wg.Add(2)

	p.m1.Lock()
	go run(&p, func() {
		p.value, p.err = cb()
	})

	return &p
}
