# Go map resource allocation

## Short summary
```go
memo := make(map[int]int)

// When empty, it consumes very little memory:
// memo itself: typically one machine word—8 bytes on a 64-bit machine.
// A small runtime map descriptor may also exist, depending on Go version and compiler escape analysis.
// Entry storage is generally allocated lazily when you insert values.
// CPU usage is effectively zero until the map is accessed.
// After insertion, each int → int pair contains 16 bytes of raw data, plus hash-table metadata and unused capacity. A practical estimate is roughly 20–30 bytes per entry, though it varies.

// So approximately:
// 1 key int        = 8 bytes
// 1 value int      = 8 bytes
// raw total      = 16 bytes
// map overhead   ≈ 4–14 bytes per entry
// estimated total ≈ 20–30 bytes per entry

// For 1,000 entries
// Lower estimate: 1,000 × 20 bytes = 20,000 bytes ≈ 20 KB
// Upper estimate: 1,000 × 30 bytes = 30,000 bytes ≈ 30 KB


// To measure accurately for your Go version and usage pattern, use a benchmark with go test -bench=. -benchmem, because allocation behavior changes depending on whether the map escapes to the heap.
```

Yes. First, one correction: **20–30 bytes per entry is only a rough dense-map estimate**, not a guaranteed number. Your machine currently has Go **1.26.5**, whose map uses a Swiss-table design.

For:

```go
memo := make(map[int]int)
```

## Detailed allocation
There are several kinds of memory overhead.

## 1. The map variable

`memo` is not the whole map. It is effectively a pointer to the runtime map structure.

On your 64-bit machine:

```go
unsafe.Sizeof(memo) // 8 bytes
```

This reports only the map reference, not the internal storage.

## 2. The map header

Go maintains an internal structure containing information such as:

- Number of entries
- Random hash seed
- Pointer to the map’s tables
- Directory length
- Growth information
- Concurrent-write detection state
- Deletion/tombstone state
- Iteration/clear state

In Go 1.26.5, this header is approximately **48 bytes** on a 64-bit machine. Its fields are visible in the installed runtime source: [map.go](/Users/9999dude/.asdf/installs/golang/1.26.5/go/src/internal/runtime/maps/map.go:195).

The compiler may:

- Place this structure on the stack
- Allocate it on the heap
- Eliminate it entirely if the map is unused

So `make(map[int]int)` does not necessarily mean a heap allocation.

## 3. Entry storage

For `map[int]int`, each slot contains:

```text
key:   int = 8 bytes
value: int = 8 bytes
---------------------
slot:       16 bytes
```

But Go does not allocate entries individually. It stores them in groups of eight slots.

Conceptually, one group looks like:

```text
8-byte control word
8 × 16-byte key/value slots
```

Therefore:

```text
Control metadata        8 bytes
Eight slots     8 × 16 = 128 bytes
---------------------------------
Logical group size     136 bytes
```

The group layout is defined in [group.go](/Users/9999dude/.asdf/installs/golang/1.26.5/go/src/internal/runtime/maps/group.go:241).

For a small map, the allocator may round that 136-byte allocation to the next size class—**144 bytes** on this Go installation.

## 4. Control-byte overhead

Every slot has one control byte. For eight slots, these eight bytes form the control word.

A control byte records whether the slot is:

- Empty
- Occupied
- Deleted—a tombstone

For occupied slots, it also stores seven lower bits of the key’s hash. That lets Go reject most nonmatching slots before comparing full keys.

Therefore, for `map[int]int`:

```text
Raw key/value data = 16 bytes per slot
Control metadata   =  1 byte per slot
```

That gives **17 logical bytes per available slot**, before unused capacity and other metadata.

## 5. Unused capacity

A hash map cannot remain completely full without lookups becoming slow. Go normally grows a table when it reaches approximately a **7/8 load factor**: [group.go](/Users/9999dude/.asdf/installs/golang/1.26.5/go/src/internal/runtime/maps/group.go:15).

That means approximately:

```text
8 allocated slots
7 occupied slots
1 intentionally unused slot
```

For a densely populated map:

```text
Group size ÷ live entries
= 136 ÷ 7
≈ 19.43 bytes per entry
```

The breakdown is:

```text
Actual int key/value        16.00 bytes
Control metadata: 8 ÷ 7      1.14 bytes
Unused slot: 16 ÷ 7           2.29 bytes
---------------------------------------
Total                        19.43 bytes
```

That is where the lower end of the earlier **20–30 bytes per entry** estimate came from.

## 6. Growth overhead

Map capacity grows in powers of two. Suppose a table reaches its limit and grows from 1,024 to 2,048 slots.

Immediately after growth, it might contain only around 896 entries:

```text
Old limit = 1,024 × 7/8 = 896 entries
New capacity             = 2,048 slots
```

Its group storage is then approximately:

```text
2,048 slots × 17 bytes = 34,816 bytes
34,816 ÷ 896 entries   ≈ 38.9 bytes per entry
```

As more entries are added, the storage becomes denser and approaches approximately 19.4 bytes per entry again.

Therefore, real memory per entry fluctuates:

```text
Dense table:              approximately 19–21 bytes/entry
Recently grown table:     potentially around 39 bytes/entry
Plus top-level metadata and allocator overhead
```

Go 1.26 splits large maps into multiple smaller tables, so the entire map does not necessarily double simultaneously.

## 7. Table and directory metadata

Larger maps have additional structures:

```text
Map header
    │
    ▼
Directory: array of table pointers
    │
    ├── Table metadata
    │      └── Groups containing slots
    │
    └── Table metadata
           └── Groups containing slots
```

Each table records:

- Used slot count
- Capacity
- Remaining growth allowance
- Directory depth and index
- Pointer to its groups
- Group count/mask

The table structure is approximately **32 bytes** on your machine: [table.go](/Users/9999dude/.asdf/installs/golang/1.26.5/go/src/internal/runtime/maps/table.go:33).

Each directory entry is another 8-byte pointer. This overhead is normally small compared with millions of key/value slots.

## 8. Deletion overhead

Deleting an entry does not normally reduce the allocated capacity.

Some deleted slots become tombstones:

```text
Before deletion: occupied
After deletion:  deleted/tombstone
```

A tombstone is needed so lookup probing does not stop prematurely. Tombstones can be reused by later insertions and are cleaned up during rehashing.

Consequently:

```go
for k := range memo {
    delete(memo, k)
}
```

can leave an empty map that still owns most of its previously allocated storage. The number of live entries becomes zero, but retained capacity may remain high.

## 9. Allocation rounding

The Go allocator uses fixed size classes. An object requesting 136 bytes may receive a 144-byte allocation.

For large group arrays, the runtime allocates many groups together, so rounding overhead becomes proportionally smaller. There is also allocator bookkeeping outside the object’s directly visible bytes.

## 10. Garbage-collector overhead

For `map[int]int`, the keys and values do not contain pointers. That is helpful because the garbage collector does not need to scan every integer as a possible heap pointer.

For something like:

```go
map[int]*LargeObject
```

the map stores an 8-byte pointer, but the referenced `LargeObject` allocation is separate and must also be scanned and retained by the garbage collector.

Likewise:

```go
map[string]int
```

stores a 16-byte string header as the key, but the string’s underlying bytes may live in a separate allocation.

## Putting it together

For a large `map[int]int` with `N` live entries:

```text
Memory ≈ map header
       + directory metadata
       + table metadata
       + allocated slot capacity × 17 bytes
       + allocation rounding
       + temporary growth memory
```

At a dense 7/8 load:

```text
Capacity ≈ N × 8/7

Slot storage ≈ N × 8/7 × 17
             ≈ N × 19.43 bytes
```

For one million densely stored entries:

```text
1,000,000 × 19.43
≈ 19,430,000 bytes
≈ 18.5 MiB
```

But because of capacity rounding, partial tables, deletions, and recent growth, the actual retained memory may be considerably higher—sometimes closer to **30–40 MB**. So the earlier 20–30 MB figure is a useful ballpark, not a strict bound.
