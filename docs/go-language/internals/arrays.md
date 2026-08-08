# Go array resource allocation

## Short summary

```go
var values [1000]int64
```

An array stores all elements inline. On this machine:

```text
unsafe.Sizeof(values) = 1,000 × 8 = 8,000 bytes
```

There is no array header, stored length, capacity field, or separate backing
array. The length is part of the compile-time type.

## 1. Inline layout

For `[N]T`, the ordinary layout is:

```text
array size = N × unsafe.Sizeof(T)
```

Elements are contiguous and indexed with a constant stride:

```text
[element 0][element 1][element 2] ... [element N-1]
```

The element's size already includes any padding needed to place consecutive
values correctly. For example, an array of padded structs repeats the complete
padded struct layout.

## 2. The array value owns its storage

Unlike a slice, an array is the data itself:

```go
a := [3]int{10, 20, 30}
b := a // copies all three integers
b[0] = 99
```

Changing `b` does not change `a`. Assignment, argument passing, and return by
value conceptually copy the entire array, although the compiler may optimize
copies when the result is unchanged.

## 3. Stack versus heap

A local array can live in a goroutine's stack frame. It may move to the heap
when it escapes or is too large for the compiler/runtime's stack-allocation
decisions:

```go
func buffer() *[4096]byte {
    var b [4096]byte
    return &b
}
```

The returned pointer is one machine word, while the pointed-to array has 4,096
bytes of payload. Escape analysis decides placement; source syntax alone is
not a reliable heap-allocation count.

## 4. Zero-length and zero-size elements

```go
var empty [0]int64
```

`unsafe.Sizeof(empty)` is zero, though the type retains the element alignment.

Arrays of zero-size types also have zero logical size:

```go
var markers [1_000_000]struct{}
```

No byte is needed for each marker. Addresses of distinct zero-size variables
are not guaranteed to be distinct.

## 5. Array length changes the type

`[4]int` and `[5]int` are different types. The compiler already knows the
length, so the value needs no runtime length word. A pointer to an array is
still only 8 bytes on this machine regardless of array length:

```go
var p *[1_000_000]int // p is 8 bytes; the array is separate
```

## 6. Alignment and containing structs

An array's alignment is the alignment of its element type. When embedded in a
struct, padding may appear before or after the array because of neighboring
fields and the containing struct's final alignment.

Use these separately:

```go
unsafe.Sizeof(value)   // complete array or struct size
unsafe.Alignof(value)  // required alignment
unsafe.Offsetof(field) // field start within a struct
```

## 7. Garbage collection

An array of pointer-free elements is non-scannable data. An array containing
pointers, strings, slices, interfaces, maps, channels, or functions must have
the relevant words scanned when it resides in the heap.

The GC cost therefore depends on element pointer content, not only total byte
size.

## 8. CPU and measurement

Indexing is constant-time. Bounds checks may be removed when the compiler can
prove an index is valid. Copying cost is proportional to the array's byte size.

```bash
go test -bench=. -benchmem
go test -gcflags=-m=2
```

Use a benchmark when deciding between an array value, an array pointer, and a
slice; compiler optimizations and escape behavior materially affect the result.
