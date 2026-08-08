# Go struct resource allocation

## Short summary

A struct stores its fields inline, with padding inserted for alignment:

```go
type Item struct {
    ID    int64
    Valid bool
}
```

On this machine, the field payload is 9 bytes, but `unsafe.Sizeof(Item{})` is
typically **16 bytes**: 8 bytes for `ID`, 1 byte for `Valid`, and 7 trailing
padding bytes.

A struct has no universal runtime header. Its size depends entirely on its
fields, their order, alignment, and trailing padding.

## 1. Field layout

Each field begins at an offset suitable for that field's alignment. The struct
is then rounded to a multiple of its largest field alignment so array elements
remain correctly aligned.

```go
type Poor struct {
    A byte  // offset 0
    B int64 // offset 8 after 7 bytes of padding
    C byte  // offset 16, then trailing padding
} // typically 24 bytes
```

Reordering can reduce retained memory:

```go
type Better struct {
    B int64 // offset 0
    A byte  // offset 8
    C byte  // offset 9
} // typically 16 bytes
```

Do not reorder fields when external binary layout, atomic alignment, generated
code, or readability requirements are more important than the saving.

## 2. Fields are inline, referenced data is not

```go
type Record struct {
    Name  string
    Data  []byte
    Next  *Record
    Count int64
}
```

The struct contains:

```text
Name string header   16 bytes
Data slice header    24 bytes
Next pointer          8 bytes
Count                 8 bytes
------------------------------
inline total         56 bytes on this machine
```

It does not include the string bytes, slice backing array, or the next record.
Those are separate reachable objects.

## 3. Stack versus heap

A struct value can be stored directly in a stack frame, in an enclosing object,
or in static data. It moves to the heap when escape analysis requires it:

```go
func newItem() *Item {
    v := Item{ID: 1, Valid: true}
    return &v
}
```

`new(Item)` and `&Item{}` do not guarantee separate heap allocations; the
compiler can place or eliminate storage according to its escape analysis and
optimizations.

## 4. Value copying

Assignment copies every inline field:

```go
b := a
```

For string, slice, map, channel, function, interface, and pointer fields, this
copies their small descriptors or references—not the objects they reach. A
struct copy can therefore be shallow for some fields and independent for
inline scalar or array fields.

Large structs passed frequently by value may increase copy cost. A pointer can
avoid the copy but may introduce aliasing, heap escape, GC scanning, and pointer
chasing. Measure the actual call path.

## 5. Empty structs

```go
type Marker struct{}
```

`unsafe.Sizeof(Marker{})` is zero. Empty structs are useful when only presence
matters, such as `map[Key]struct{}`.

For a trailing zero-size field, the compiler may add padding so its address
does not point beyond the containing object. Always measure the actual struct
instead of assuming a trailing `struct{}` is free in every layout.

## 6. Garbage collection

The compiler emits pointer-layout metadata for each struct type. A struct made
only of numbers and booleans is non-scannable. A struct with reference-bearing
fields makes the GC inspect the relevant pointer words, while ignoring known
scalar regions.

Field order can sometimes reduce the pointer-containing prefix the GC scans,
but ordinary layout clarity and measured memory wins should guide such changes.

## 7. CPU and cache behavior

Field access is a constant-offset memory operation. Smaller, denser structs can
improve cache locality when stored in large arrays. Splitting hot and cold
fields can help workloads that touch only a subset, but can also add pointer
indirection and allocations.

## 8. Measurement

```go
fmt.Println(unsafe.Sizeof(Item{}))
fmt.Println(unsafe.Alignof(Item{}))
fmt.Println(unsafe.Offsetof(Item{}.Valid))
```

`unsafe.Sizeof` counts inline representation only. Use heap profiles and
`go test -bench=. -benchmem` to include reachable data, escape behavior, and
allocator rounding.
