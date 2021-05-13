# go-patch-overlay

An experimental way to apply patches to the Go runtime at build time.

Assuming you have a directory of [patches](./example/goroutineid/patches) to apply to the Go source tree, you can use it like this:

```
$ go build -overlay="$(go-patch-overlay ./patches)"
```
