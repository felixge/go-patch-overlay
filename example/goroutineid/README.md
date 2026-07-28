# goroutineid

An example for using the go-patch-overlay tool to [hack](https://github.com/felixge/go-patch-overlay/blob/ec286d1fd9dc/example/goroutineid/patches/getgid.patch) a `runtime.Getgid()` function into the runtime that returns the current goroutine id.

To run the code, you need to invoke `go` like this:

```
$ go run -overlay="$(go-patch-overlay ./patches/getgid.patch)"
```
