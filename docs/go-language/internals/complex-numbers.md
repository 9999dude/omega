# Go complex-number resource allocation

## Short summary

```text
complex64  = two float32 components =  8 bytes
complex128 = two float64 components = 16 bytes
```

A complex number stores its real and imaginary components inline. It has no
header or backing allocation.

## 1. Representation

```go
var a complex64  = complex(1, 2)
var b complex128 = complex(1, 2)
```

Conceptually:

```text
complex64:  [real float32][imag float32]
complex128: [real float64][imag float64]
```

The `real`, `imag`, and `complex` built-ins operate on these components and do
not allocate by themselves.

## 2. Stack, registers, and heap

Local complex values are normally held in registers or in the stack frame.
They escape only when the surrounding program requires longer-lived storage:

```go
func point() *complex128 {
    z := complex(3.0, 4.0)
    return &z
}
```

The pointee contains 16 bytes of payload. If it moves to the heap, allocator
size classes and bookkeeping determine the actual memory charged.

## 3. Arrays and slices

```text
[1000]complex64  =  8,000 bytes
[1000]complex128 = 16,000 bytes
```

A slice adds its 24-byte header and retains backing storage according to
capacity. A million-element `[]complex128` therefore has about 16 MB of raw
element storage before allocator rounding.

## 4. Alignment and copying

On this machine, `complex64` has 4-byte alignment and `complex128` has 8-byte
alignment. Struct placement can add padding around them.

Assignment copies both components. Passing a large array of complex values by
value copies the complete array; passing a slice copies only its three-word
header and shares the backing array.

## 5. CPU considerations

Addition and subtraction operate component-wise. Multiplication uses several
real multiplications and additions; division is more expensive. These
operations use constant auxiliary space and do not imply heap allocation.

Formatting a complex number or placing it in an interface may allocate because
of those surrounding operations, not because the complex type owns dynamic
storage.

## 6. Garbage collection

Complex values contain no pointers, so arrays consisting only of complex
numbers do not need their elements scanned as references by the garbage
collector.

## 7. Measurement

```go
fmt.Println(unsafe.Sizeof(complex64(0)))  // 8
fmt.Println(unsafe.Sizeof(complex128(0))) // 16
```

Use `go test -bench=. -benchmem` to distinguish raw value size from allocations
introduced by formatting, containers, or escape behavior.
