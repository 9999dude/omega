# Go function-value resource allocation

## Short summary

```go
var f func(int) int
```

A function variable is **8 bytes** on this 64-bit machine. Its zero value is
`nil`. A non-nil value points to runtime/compiler-managed function-value data;
capturing closures can have an additional environment whose size depends on
the captured variables.

The runtime's variable-size `funcval` begins with a code pointer in
[runtime2.go](/Users/9999dude/.asdf/installs/golang/1.26.5/go/src/runtime/runtime2.go:179).

## 1. Plain function values

```go
func double(x int) int { return x * 2 }

f := double
```

The variable `f` is one word. The executable machine code is stored in the
program's text segment and is shared; assigning `double` does not copy that
code.

Copying `f` copies the one-word reference:

```go
g := f
```

Both values invoke the same function and normally need no new heap allocation.

## 2. Closures and captured variables

```go
func counter() func() int {
    n := 0
    return func() int {
        n++
        return n
    }
}
```

The returned function must retain `n`. Conceptually the closure contains:

```text
function code pointer
captured environment: reference to or storage for n
```

Because the closure and `n` outlive `counter`, their environment normally
needs heap storage. The exact layout is compiler-generated and can combine,
copy, or reference captured values.

## 3. Capturing by value versus shared variable

Closures capture variables, not a fresh snapshot on every reference:

```go
n := 1
f := func() int { return n }
n = 2
fmt.Println(f()) // 2
```

When mutation or lifetime requires it, the compiler moves shared captured
storage to a location accessible by all closures. Capturing a large variable
can therefore retain more memory than the function's one-word visible size
suggests.

The compiler may capture an immutable value directly when it proves that doing
so preserves semantics.

## 4. Method values

```go
f := service.Handle
```

A method value binds a receiver. Its hidden environment retains or copies that
receiver according to the method and value semantics. A large value receiver
can be copied; a pointer receiver retains the target object.

A method expression does not bind a receiver:

```go
f := (*Service).Handle
```

The receiver becomes an explicit argument, so this form generally needs no
bound-receiver environment.

## 5. Escape and allocation behavior

Not every closure allocates:

```go
func useNow() int {
    n := 3
    f := func() int { return n + 1 }
    return f()
}
```

Inlining and escape analysis may keep or eliminate the closure completely. A
closure more commonly allocates when it is returned, stored globally, placed
in a long-lived object, or passed to code whose lifetime cannot be proven.

## 6. Garbage collection

The function value is a GC-visible reference. A reachable closure keeps its
captured environment and everything reachable from that environment alive.

Setting an unused long-lived callback field to `nil` can release its capture
graph. This matters when closures capture large buffers, requests, or service
objects unintentionally.

## 7. CPU considerations

An indirect function call may prevent inlining and add dispatch overhead. A
direct call can often be inlined. The larger cost is frequently what the
function does, but closure allocation and retained captures matter in hot
callback creation paths.

## 8. Measurement

```go
unsafe.Sizeof((func())(nil)) // 8 on this machine
```

This sees only the function reference, not its environment. Use escape output
and benchmarks:

```bash
go test -gcflags=-m=2
go test -bench=. -benchmem
```

Heap profiles are the clearest way to find unexpectedly retained closure
environments.
