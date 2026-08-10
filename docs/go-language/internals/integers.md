# Go integer resource allocation

## Short summary

Integer values store their bits directly. They have no internal header or
backing allocation.

```text
int8, uint8, byte           1 byte
int16, uint16               2 bytes
int32, uint32, rune         4 bytes
int64, uint64               8 bytes
int, uint, uintptr          8 bytes on this 64-bit machine
```

`byte` is an alias for `uint8`; `rune` is an alias for `int32`. Aliases do not
add metadata.

## 1. Machine-sized integers

The sizes of `int`, `uint`, and `uintptr` depend on the target architecture:

```text
32-bit target: 4 bytes
64-bit target: 8 bytes
```

This installation is `darwin/arm64`, so all three are 8 bytes. Prefer a
fixed-width type for binary formats and protocols whose layout must not change
between targets.

## 2. Stack, registers, and heap

Local integers are usually kept in registers or directly in a stack frame.
Arithmetic does not allocate. Storage escapes only because of how a value is
used, not because it is an integer:

```go
func counter() *int64 {
    n := int64(0)
    return &n
}
```

Here `n` must remain alive after the function returns. Its logical payload is
8 bytes; the heap allocator can reserve a larger size-class slot.

## 3. Arrays and bulk storage

For an array, multiply the element count by the element size:

```text
[1000]int8   =  1,000 bytes
[1000]int32  =  4,000 bytes
[1000]int64  =  8,000 bytes
```

A slice adds a 24-byte header and retains storage based on capacity, not
length. For example, a `[]int64` with capacity 1,000 has 8,000 bytes of logical
backing storage plus its slice header and allocator rounding.

## 4. Alignment and struct padding

Integer fields are normally aligned to their natural alignment, capped by the
platform rules. Poor field ordering can introduce padding:

```go
type A struct {
    Small int8
    Large int64
    Tail  int8
}
```

On this machine, `A` is typically 24 bytes even though its fields contain only
10 bytes. Padding makes `Large` aligned and rounds the struct size so array
elements remain aligned.

Reordering can reduce it:

```go
type B struct {
    Large int64
    Small int8
    Tail  int8
} // typically 16 bytes
```

## 5. Overflow and CPU cost

Integer operations are constant-time machine operations. Ordinary unsigned
arithmetic wraps modulo 2 to the type's width; signed operations are also
deterministically represented by two's-complement arithmetic. Go does not
allocate a wider integer automatically on overflow.

Division is usually more CPU-expensive than addition, bit operations, or
multiplication, but none allocates by itself.

## 6. Garbage collection

Ordinary integer types contain no pointers. Large integer arrays can therefore
be allocated as non-scannable heap memory, reducing GC scan work compared with
equally sized pointer arrays.

`uintptr` is an integer, not a GC-visible reference. Storing an address only in
a `uintptr` does not keep the referenced object alive.

## 7. Measurement

```go
fmt.Println(unsafe.Sizeof(int8(0)))  // 1
fmt.Println(unsafe.Sizeof(int32(0))) // 4
fmt.Println(unsafe.Sizeof(int64(0))) // 8
fmt.Println(unsafe.Sizeof(int(0)))   // 8 on this machine
```

Use `unsafe.Sizeof` for representation, `unsafe.Offsetof` for field padding,
and `go test -bench=. -benchmem` for actual heap allocations.
