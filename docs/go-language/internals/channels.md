# Go channel resource allocation

## Short summary

```go
ch := make(chan int64, 1000)
```

A channel variable is an **8-byte reference** on this 64-bit machine. `make`
creates a runtime channel descriptor plus, for a buffered channel, storage for
`capacity × element size`:

```text
channel variable                         8 bytes
runtime hchan descriptor       approximately 112 bytes on Go 1.26.5/arm64
logical int64 buffer       1,000 × 8 = 8,000 bytes
allocator rounding and blocked-waiter costs are additional
```

The current descriptor is defined in
[chan.go](/Users/9999dude/.asdf/installs/golang/1.26.5/go/src/runtime/chan.go:34).

## 1. Nil channel

```go
var ch chan int
```

The variable is one nil word and owns no channel descriptor or buffer. Sending
or receiving on a nil channel blocks forever; closing one panics.

A nil channel is often useful for disabling a `select` case without allocating
a new channel.

## 2. Unbuffered channel

```go
ch := make(chan Item)
```

The runtime allocates its channel descriptor, but there is no element queue.
A send and receive rendezvous: one side normally blocks until the other is
ready, and the element can be transferred directly between them.

Blocked goroutines use scheduler bookkeeping (`sudog` records and wait queues).
That cost is not represented by `unsafe.Sizeof(ch)` and grows with contention
and the number of waiters.

## 3. Buffered channel

```go
ch := make(chan Item, n)
```

The runtime maintains a circular queue with send and receive indices, current
element count, close state, wait queues, element type information, and a lock.
Logical queue storage is:

```text
n × unsafe.Sizeof(Item)
```

The allocator rounds the total. In Go 1.26.5, the runtime can combine the
descriptor and pointer-free element buffer into one allocation. When elements
contain pointers, it uses separate descriptor and buffer allocations so the GC
has precise element scanning information. The allocation paths are visible in
[chan.go](/Users/9999dude/.asdf/installs/golang/1.26.5/go/src/runtime/chan.go:75).

## 4. Capacity is allocated eagerly

Unlike map entries, buffered-channel capacity is reserved when `make` runs:

```go
ch := make(chan [1024]byte, 1000)
```

This requests roughly 1,024,000 bytes of queue payload even when the channel is
empty. Choose capacity from the needed burst and backpressure behavior; a huge
buffer can hide overload and retain substantial memory.

## 5. Copying a channel value

```go
other := ch
```

This copies only the 8-byte reference. Both variables refer to the same queue,
close state, and waiters. A channel is closed once, regardless of how many
references exist.

## 6. Element copying

Sending copies an element into the queue or receiver's destination according
to the channel operation:

```go
ch <- item
```

For a large struct, that copy can be significant. Sending `*Item` copies only
an 8-byte pointer but introduces sharing, potential heap escape, GC scanning,
and pointer chasing. Select based on ownership and concurrency semantics, then
measure.

## 7. Garbage collection and retention

A reachable channel keeps its buffered elements reachable. If elements contain
pointers, strings, slices, maps, interfaces, channels, or functions, queued
values can retain their entire object graphs.

Receiving removes an element, and the runtime clears pointer-containing queue
slots so their referenced objects can become unreachable. Closing a channel
does not discard buffered elements; receivers can drain them afterward.

## 8. CPU and scheduler cost

Channel operations synchronize goroutines and are more expensive than ordinary
loads/stores. The fast path checks queue and waiter state; contended operations
may lock, park the current goroutine, and wake another goroutine.

Unbuffered operations provide direct backpressure. Buffered operations avoid
blocking only while queue space/data is available; they are not lock-free
storage.

## 9. Measurement

```go
unsafe.Sizeof((chan int)(nil)) // 8 on this machine
```

This counts only the reference. Benchmark representative buffer capacities,
element types, and contention:

```bash
go test -bench=. -benchmem
```

Use heap profiles for retained buffered elements and execution traces or block
profiles for scheduler and blocking costs.
