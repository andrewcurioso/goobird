package promise

import (
	"testing"
	"time"
)

func TestBasicPromise(t *testing.T) {

	var p *SPromise

	p = Create(func() (res interface{}, err interface{}) {
		return 42, nil
	})

	x, err := p.Wait()

	if err != nil {
		t.Error("Expected err to be nil got", err)
	}

	if x != 42 {
		t.Error("Expected x to be 42 got", x)
	}

}

func TestThen(t *testing.T) {

	var p *SPromise

	p = Create(func() (res interface{}, err interface{}) {
		return 42, nil
	}).Then(func(ret interface{}) (interface{}, interface{}) {
		x := ret.(int)
		return x * 2, nil
	})

	x, err := p.Wait()

	if err != nil {
		t.Error("Expected err to be nil got", err)
	}

	if x != 84 {
		t.Error("Expected x to be 84 got", x)
	}

}

func TestNestedCreate(t *testing.T) {

	var p *SPromise

	p = Create(func() (res interface{}, err interface{}) {
		return Create(func() (res interface{}, err interface{}) { return 42, nil }), nil
	}).Then(func(ret interface{}) (interface{}, interface{}) {
		x := ret.(int)
		return x * 2, nil
	})

	x, err := p.Wait()

	if err != nil {
		t.Error("Expected err to be nil got", err)
	}

	if x != 84 {
		t.Error("Expected x to be 84 got", x)
	}

}

func TestNestedThen(t *testing.T) {

	var p *SPromise

	p = Create(func() (res interface{}, err interface{}) {
		return 42, nil
	}).Then(func(ret interface{}) (interface{}, interface{}) {
		x := ret.(int)
		return Create(func() (res interface{}, err interface{}) { return x * 2, nil }), nil
	})

	x, err := p.Wait()

	if err != nil {
		t.Error("Expected err to be nil got", err)
	}

	if x != 84 {
		t.Error("Expected x to be 84 got", x)
	}

}

func TestDelayedThen(t *testing.T) {

	var p *SPromise

	p = Create(func() (res interface{}, err interface{}) {
		return 42, nil
	})

	time.Sleep(100 * time.Millisecond)

	p.Then(func(ret interface{}) (interface{}, interface{}) {
		x := ret.(int)
		return x * 2, nil
	})

	x, err := p.Wait()

	if err != nil {
		t.Error("Expected err to be nil got", err)
	}

	if x != 84 {
		t.Error("Expected x to be 84 got", x)
	}

}
