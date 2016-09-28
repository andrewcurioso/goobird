package promise

import (
	"testing"
	"time"
)

func TestBasicPromise(t *testing.T) {

	x, err := Create(func() (interface{}, interface{}) {
		return 42, nil
	}).Wait()

	if err != nil {
		t.Error("Expected err to be nil got", err)
	}

	if x != 42 {
		t.Error("Expected x to be 42 got", x)
	}

}

func TestSleep(t *testing.T) {

	x, err := Create(func() (interface{}, interface{}) {
		time.Sleep(100 * time.Millisecond)
		return 42, nil
	}).Wait()

	if err != nil {
		t.Error("Expected err to be nil got", err)
	}

	if x != 42 {
		t.Error("Expected x to be 42 got", x)
	}

}

func TestThen(t *testing.T) {

	x, err := Create(func() (interface{}, interface{}) {
		return 42, nil
	}).Then(func(ret interface{}) (interface{}, interface{}) {
		x := ret.(int)
		return x * 2, nil
	}).Wait()

	if err != nil {
		t.Error("Expected err to be nil got", err)
	}

	if x != 84 {
		t.Error("Expected x to be 84 got", x)
	}

}

func TestNestedCreate(t *testing.T) {

	x, err := Create(func() (interface{}, interface{}) {
		return Create(func() (interface{}, interface{}) { return 42, nil }), nil
	}).Then(func(ret interface{}) (interface{}, interface{}) {
		x := ret.(int)
		return x * 2, nil
	}).Wait()

	if err != nil {
		t.Error("Expected err to be nil got", err)
	}

	if x != 84 {
		t.Error("Expected x to be 84 got", x)
	}

}

func TestNestedThen(t *testing.T) {

	x, err := Create(func() (interface{}, interface{}) {
		return 42, nil
	}).Then(func(ret interface{}) (interface{}, interface{}) {
		x := ret.(int)
		return Create(func() (interface{}, interface{}) { return x * 2, nil }), nil
	}).Wait()

	if err != nil {
		t.Error("Expected err to be nil got", err)
	}

	if x != 84 {
		t.Error("Expected x to be 84 got", x)
	}

}

func TestDelayedThen(t *testing.T) {

	var p *SPromise

	p = Create(func() (interface{}, interface{}) {
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

func TestError(t *testing.T) {

	x, err := Create(func() (interface{}, interface{}) {
		return true, nil

	}).Then(func(ret interface{}) (interface{}, interface{}) {
		return true, true

	}).Then(func(ret interface{}) (interface{}, interface{}) {
		t.Error("Then should not execute on error")
		return true, nil
	}).Wait()

	if err != true {
		t.Error("Expected err to be true got", err)
	}

	if x != nil {
		t.Error("Expected x to be nil got", x)
	}

}

func TestDelayedWait(t *testing.T) {

	p := Create(func() (interface{}, interface{}) {
		return 42, nil
	})

	time.Sleep(100 * time.Millisecond)

	x, err := p.Wait()

	if err != nil {
		t.Error("Expected err to be nil got", err)
	}

	if x != 42 {
		t.Error("Expected x to be 42 got", x)
	}

}

func TestAll(t *testing.T) {
	x, err := Create(func() (interface{}, interface{}) {
		return [...]*SPromise{
			Create(func() (interface{}, interface{}) { return 1, nil }),
			Create(func() (interface{}, interface{}) { return 2, nil }),
			Create(func() (interface{}, interface{}) { return 3, nil }),
		}, nil
	}).All().Wait()

	if err != nil {
		t.Error("Expected err to be nil got", err)
	}

	xi := x.([]interface{})

	if len(xi) != 3 {
		t.Error("Expected length of 3 got", len(xi))
	}

	var sum int

	for i := 0; i < len(xi); i++ {
		sum = sum + xi[i].(int)
	}

	if sum != 6 {
		t.Error("Expected sum to be 6 got", sum)
	}
}
