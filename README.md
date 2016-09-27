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

- **All**: Run an array of promises
- **Each**: Run a function on all items of an array and return an array
- **Map**: Run a function on each item of a map and return a map, order not guarenteed
- **Reduce**: Run a function to reduce an array to a single value, order not guarenteed
- **Filter**: Run a function to filter an array to contain only specific values

### Improvements

- Documentation
- License
