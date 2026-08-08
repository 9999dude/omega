# Go pointer resource allocation

## Short summary

```go
var p *int
```

A typed pointer is one machine word: **8 bytes** on this 64-bit machine. The
pointer value does not include the pointed-to object's bytes.

```text
pointer variable        8 bytes
pointed-to int           8 bytes, stored elsewhere
allocator/GC overhead    depends on placement and escape behavior
```

The zero value is `nil` and points to no object.

## 1. A pointer does not allocate by itself

```go
x := 42
p := &x
```

Creating `p` is just obtaining an address. `x` may remain in the goroutine's
stack frame when the compiler proves that the pointer cannot outlive it.

This commonly escapes:

```go
func value() *int {
    x := 42
    return &x
}
```

The compiler arranges storage that remains valid after return, usually the
heap. Go never leaves a pointer referring to a discarded stack frame; stacks
can also grow and the runtime updates stack pointers as needed.

## 2. `new(T)`

```go
p := new(int)
```

`new(T)` returns `*T` pointing to zeroed storage for one `T`. It does not
guarantee a heap allocation. The compiler can place the storage on the stack
or eliminate it when it does not escape or produce an observable effect.

The pointer itself is still 8 bytes; storage for `T` is separate.

## 3. Aliasing and lifetime

Several pointers may refer to one object:

```go
p := &item
q := p
```

Copying the pointer copies 8 bytes, not `item`. Mutations through either
pointer affect the same object. Any live GC-visible pointer can keep that
object and the objects reachable from it alive.

A pointer to a small subobject can retain the complete enclosing allocation:

```go
p := &large.Field
```

As long as `p` remains reachable, the allocator generally cannot reclaim the
rest of `large` independently.

## 4. Pointers in bulk storage

```text
[1000]*Node = 8,000 bytes of pointer slots
```

That does not include any `Node` objects. If every slot points to a separately
allocated node, actual memory includes 1,000 node allocations, allocator
rounding, and reduced locality compared with inline `[1000]Node` storage.

Pointers can reduce copies and represent optional/shared data, but they add
indirection and GC scanning. Use them for semantics first and measure layout
changes in hot collections.

## 5. `unsafe.Pointer` and `uintptr`

`unsafe.Pointer` is a GC-visible pointer-sized value. `uintptr` is a
pointer-sized integer, but it is **not** a reference as far as the garbage
collector is concerned.

Converting a pointer to `uintptr` and retaining only the integer does not keep
the object alive. Pointer arithmetic must follow the `unsafe` package's strict
same-expression and object-boundary rules; otherwise the program may become
invalid even if the numeric address looks correct.

## 6. Garbage collection cost

A heap object containing pointer fields has pointer-layout metadata. The GC
scans its pointer slots and follows non-nil references. Pointer-free objects
avoid that scan work.

Assigning a pointer into heap memory may execute a GC write barrier. The barrier
is normally small but means pointer-heavy mutation is not identical in cost to
writing integers.

## 7. CPU and measurement

Dereferencing is constant-time but can cause cache misses. Pointer chains often
have worse locality than contiguous inline values.

```go
unsafe.Sizeof((*int)(nil)) // 8 on this machine
```

Use escape diagnostics and allocation benchmarks:

```bash
go test -gcflags=-m=2
go test -bench=. -benchmem
```

Heap and allocation profiles show the pointed-to object costs that
`unsafe.Sizeof(p)` cannot see.
