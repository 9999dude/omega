# Go slice resource allocation

## Short summary

```go
values := make([]int64, 1000, 1200)
```

A slice value is a three-word descriptor:

```text
pointer to backing array   8 bytes
length                     8 bytes
capacity                   8 bytes
---------------------------------
slice header              24 bytes
```

The logical backing storage is `capacity × element size`, so this example
retains at least `1,200 × 8 = 9,600` bytes plus allocator rounding. The
24-byte header does not include that backing array.

The runtime representation is visible in
[slice.go](/Users/9999dude/.asdf/installs/golang/1.26.5/go/src/runtime/slice.go:16).

## 1. Nil and empty slices

```go
var a []int          // nil:   len=0, cap=0, pointer=nil
b := []int{}         // empty: len=0, cap=0, pointer may be non-nil
c := make([]int, 0)  // empty: len=0, cap=0
```

Each slice variable is still a 24-byte value when addressable. These forms
usually need no element backing allocation, though the exact pointer and
compiler treatment can differ.

## 2. Length versus capacity

The length is the accessible element count. Capacity is how far the slice can
grow from its starting position before `append` needs different storage:

```text
header: len=3, cap=5
             │
             ▼
backing: [10][20][30][  ][  ]
          accessible     spare
```

Memory retention follows capacity and the reachable backing object, not just
length.

## 3. `make` allocation

```go
s := make([]T, length, capacity)
```

Conceptually reserves `capacity × unsafe.Sizeof(T)` bytes and zeroes the
elements exposed by the language. The allocator rounds the requested byte
count to a size class. Escape analysis can sometimes keep small, fixed-size
backing arrays on the stack.

The header and backing array need not have the same lifetime or location.

## 4. `append` growth

If spare capacity exists, `append` writes into the existing backing array and
usually allocates nothing. If capacity is insufficient, the runtime:

1. Chooses a larger capacity.
2. Allocates a new backing array.
3. Copies the old elements.
4. Appends the new elements.
5. Returns a header pointing to the new array.

The exact growth formula is an implementation detail and changes with element
size, old capacity, requested length, and allocator size classes. It is
implemented in
[slice.go](/Users/9999dude/.asdf/installs/golang/1.26.5/go/src/runtime/slice.go:178).

Preallocating a realistic capacity can avoid repeated allocation and copying:

```go
out := make([]Item, 0, expected)
```

Excessive preallocation can waste and retain memory, so larger is not always
better.

## 5. Slicing shares storage

```go
large := make([]byte, 1<<20)
small := large[:10]
large = nil
```

`small` can keep the entire 1 MiB backing array reachable even though its
length is ten. Copy the needed data when releasing the large object matters:

```go
small = append([]byte(nil), small...)
```

A full slice expression can restrict future append capacity:

```go
view := large[:10:10]
```

It does not shrink the existing allocation, but a subsequent append must use
new backing storage instead of overwriting the following elements.

## 6. Copying a slice

```go
b := a
```

This copies only the 24-byte header. Both slices initially refer to the same
backing array. Use `copy` or an append-to-nil pattern for independent elements.

## 7. Zero-size elements

For `[]struct{}`, the logical element size is zero. Very large lengths need not
reserve one byte per element. The slice still records length and capacity, and
distinct zero-size elements may have equal addresses.

## 8. Garbage collection

The header contains a pointer that keeps the backing array alive. If `T`
contains pointers, the relevant backing-array slots are GC-scanned. Pointer-free
element arrays can be allocated as non-scannable memory.

Clearing unused pointer elements can release referenced objects while reusing
the slice capacity:

```go
clear(s[:cap(s)])
s = s[:0]
```

This releases the referenced objects, not the backing array itself.

## 9. CPU and measurement

Indexing is constant-time. Reslicing is normally header arithmetic. Growth is
O(number of copied elements), which is amortized across appends.

```go
unsafe.Sizeof([]int(nil)) // 24 on this machine
```

Use `go test -bench=. -benchmem` for allocation counts and memory per operation;
`unsafe.Sizeof` sees only the header.
