# Go interface resource allocation

## Short summary

```go
var x any
```

An interface value occupies **16 bytes** on this 64-bit machine: two machine
words. It stores type/method metadata in one word and data in the other.

Conceptually:

```text
empty interface:     [dynamic type pointer][data pointer]
non-empty interface: [method-table pointer][data pointer]
```

The current runtime definitions are in
[runtime2.go](/Users/9999dude/.asdf/installs/golang/1.26.5/go/src/runtime/runtime2.go:184).

## 1. Empty interfaces

For `any`—the alias of `interface{}`—the first word identifies the concrete
dynamic type and the second word represents its value:

```go
var x any = int64(42)
```

The interface header is 16 bytes. The dynamic `int64` payload is conceptually
separate from that header; compiler/runtime boxing rules determine whether it
needs fresh storage, can reference existing storage, or can use an optimized
representation.

Do not assume every interface conversion allocates. Whether it does depends on
the concrete value, escape context, compiler version, and subsequent use.

## 2. Interfaces with methods

```go
type Reader interface {
    Read([]byte) (int, error)
}
```

The first word points to an `itab` that pairs the interface type with the
concrete type and contains method dispatch information. The second word points
to the concrete data.

The `itab` is runtime/type metadata shared by values of the same pairing, not a
full per-interface copy of every method.

## 3. Nil interface versus typed nil

An interface is nil only when both words are nil:

```go
var a any       // type=nil, data=nil: a == nil
var p *Node     // p is nil
var b any = p   // type=*Node, data=nil: b != nil
```

The typed-nil interface still carries dynamic type information. This is a
semantic distinction, not extra header size.

## 4. Boxing and escape behavior

Passing a value through an interface may make it escape when the compiler
cannot prove its lifetime or use:

```go
func store(v any) { global = v }

store(LargeStruct{...})
```

The concrete value must remain reachable after the call and may require heap
storage. A large value also has to be copied into its interface representation
or backing storage. Passing a pointer copies only the pointer, but retains and
aliases its target.

Interface use that is fully inlined and does not escape can avoid heap
allocation. Inspect the actual build rather than treating "boxing" as always
allocating.

## 5. Copying interface values

```go
y := x
```

This copies the two-word interface representation. It does not deep-copy the
dynamic object. If the concrete value is a pointer, map, slice, channel, or
function, both interfaces can refer to shared state.

## 6. Type assertions and method calls

A type assertion checks runtime type information and normally does not allocate
by itself:

```go
n, ok := x.(int64)
```

An interface method call performs indirect dispatch through type/method
metadata. It may be devirtualized and inlined when the compiler can identify
the concrete type. Otherwise it adds indirection but no automatic allocation.

## 7. Garbage collection

The data word is a GC-visible reference when the dynamic representation needs
one. The type information lets the runtime know which words in the concrete
value are pointers.

An interface retaining a small view—such as a slice or substring—can indirectly
keep its complete backing allocation alive.

## 8. Measurement

```go
unsafe.Sizeof(any(nil)) // 16 on this machine
```

This counts only the interface header. Use:

```bash
go test -gcflags=-m=2
go test -bench=. -benchmem
```

to determine whether a concrete value is copied to heap storage in the real
call path.
