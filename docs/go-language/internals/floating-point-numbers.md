# Go floating-point resource allocation

## Short summary

```text
float32 = 4 bytes
float64 = 8 bytes
```

Both types store an IEEE 754 value directly. They have no per-value header,
backing store, or automatic heap allocation.

## 1. Representation

The bit fields are:

```text
float32: 1 sign +  8 exponent + 23 fraction bits = 32 bits
float64: 1 sign + 11 exponent + 52 fraction bits = 64 bits
```

These encodings include finite numbers, positive and negative zero,
infinities, and NaN values. The stored fraction omits the implicit leading bit
used for normal finite values.

## 2. Stack, registers, and heap

A local floating-point value normally resides in a floating-point register or
in the goroutine's stack frame. Expressions such as addition, multiplication,
and comparison do not allocate.

Taking an address is not itself proof of a heap allocation, but returning that
address commonly makes the value escape:

```go
func ratio() *float64 {
    x := 0.5
    return &x
}
```

The logical pointee is 8 bytes. Actual heap consumption may include allocator
rounding and metadata outside the object.

## 3. Arrays, slices, and structs

Dense arrays have straightforward logical sizes:

```text
[1000]float32 = 4,000 bytes
[1000]float64 = 8,000 bytes
```

A slice additionally uses a 24-byte header and retains `capacity × element
size` bytes of logical backing storage.

A `float64` field normally requires 8-byte alignment on this machine. Smaller
neighboring fields may introduce padding; inspect the complete struct rather
than summing field sizes.

## 4. Precision is not allocation

Changing `float32` to `float64` doubles element storage, but it does not create
extra runtime objects. Conversions are arithmetic operations:

```go
small := float32(1.25)
large := float64(small)
```

The conversion may lose or preserve precision depending on direction; it does
not allocate by itself.

## 5. CPU considerations

Basic floating-point operations are constant-space and generally hardware
instructions. Division and square root are typically more expensive than
addition or multiplication. Library functions such as trigonometric operations
cost more CPU but still need not allocate.

NaN is special: `x == x` is false when `x` is NaN. This affects map-key and
comparison behavior but does not change the value's size.

## 6. Garbage collection

Floating-point values contain no pointers. Arrays containing only floats are
non-scannable by the garbage collector, although the containing heap object
must still be tracked and reclaimed.

## 7. Measurement

```go
fmt.Println(unsafe.Sizeof(float32(0))) // 4
fmt.Println(unsafe.Sizeof(float64(0))) // 8
```

Use `unsafe.Sizeof` for value size and `go test -bench=. -benchmem` to detect
allocations caused by surrounding code such as formatting or interface boxing.
