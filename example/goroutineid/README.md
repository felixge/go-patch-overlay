# goroutineid

An example for using the go-patch-overlay tool to hack a `runtime.Getid()` function into the runtime that returns the current goroutine id.

To run the code, you need to invoke `go` like this:

```
$ go run -overlay="$(go-patch-overlay ./patches)"
```
