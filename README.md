# GooBird

[![Build Status](https://travis-ci.org/andrewcurioso/goobird.svg?branch=master)](https://travis-ci.org/andrewcurioso/goobird)

Write elegant asynchronous code using promises.

## Features

### Easily create promises

```go
p := promise.Create(func() (interface{}, interface{}) {
  time.Sleep(100 * time.Milliseconds)
  return 42, nil
})

fmt.Println("Waiting...")

x,_ := p.Wait()

fmt.Println("x =", x)
```

prints

```
Waiting...
x = 42
```

### Use of Then to chain promises

```go
p := promise.Create(func() (interface{}, interface{}) {
  return 42, nil
}).Then(func(res interface{}) (interface{}, interface{}) {
  x := res.(int)
  return x * 2, nil
})
```

### Nest promises (works in Then too)

```go
p := promise.Create(func() (interface{}, interface{}) {
  return promise.Create(func() (interface{}, interface{}) { return 42, nil }, nil
}).Then(func(res interface{}) (interface{}, interface{}) {
  x := res.(int)
  return x * 2
})
```

### Fast

Uses a single GoRoutine per promise to run serial links in the chain and uses metuxes instead of channels to synchronize the promises.

## Todo

### Features

- All
- Each
- Map
- Reduce
- Filter

### Improvements

- Documentation
- License
