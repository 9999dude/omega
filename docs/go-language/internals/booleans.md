# Go boolean resource allocation

## Short summary

```go
var enabled bool
```

On this 64-bit machine, `unsafe.Sizeof(enabled)` is **1 byte**. A standalone
`bool` has no backing store or runtime descriptor and normally needs no heap
allocation. Its zero value is `false`.

## 1. Representation

Go reserves one byte for an addressable `bool`:

```text
false = 0
true  = a nonzero byte internally
size  = 1 byte
alignment = 1 byte
```

Code should not depend on the exact byte used for `true`; the language exposes
only the two logical values.

## 2. Stack, registers, and heap

A local boolean is commonly kept in a CPU register or in the goroutine's
stack frame. It moves to the heap only when escape analysis requires the
storage to outlive that frame, for example:

```go
func flag() *bool {
    v := true
    return &v // v must remain reachable after flag returns
}
```

The returned pointer is 8 bytes on this machine. The pointee logically needs
1 byte, although a heap allocation can consume more because the allocator uses
size classes and maintains metadata outside the object.

## 3. Struct padding

A `bool` can indirectly increase a struct by more than one byte when alignment
requires padding:

```go
type Record struct {
    Ready bool  // 1 byte
    Count int64 // must start at an 8-byte-aligned address
}
```

Typical layout on this machine:

```text
Ready       1 byte
padding     7 bytes
Count       8 bytes
-------------------
total      16 bytes
```

Field order can therefore matter when millions of values are retained.

## 4. Arrays and slices

An array such as `[1000]bool` occupies exactly 1,000 bytes as a value. A
`[]bool` additionally has a 24-byte slice header and a backing array whose
capacity determines retained storage.

Go does **not** bit-pack ordinary booleans. If one bit per flag matters, use an
explicit bitmap such as `[]uint64`; that trades simpler access for lower memory
usage and extra masking operations.

## 5. Garbage collection and CPU

Boolean storage contains no pointers, so the garbage collector does not scan
its bytes as references. Reads, writes, comparisons, and logical operations
are constant-time and do not allocate by themselves.

## 6. Measure the real context

Use `unsafe.Sizeof` for the value representation and a benchmark for allocation
behavior:

```bash
go test -bench=. -benchmem
go test -gcflags=-m=2
```

`unsafe.Sizeof` does not include allocator rounding, containing-object padding,
or objects reachable through pointers.
