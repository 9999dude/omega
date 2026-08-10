# Go string resource allocation

## Short summary

```go
s := "hello"
```

A string value is a two-word, read-only view of bytes:

```text
pointer to byte data   8 bytes
byte length            8 bytes
------------------------------
string header         16 bytes
```

The header does not include the bytes. `len(s)` is the byte count, not the
number of Unicode characters.

The runtime layout is visible in
[string.go](/Users/9999dude/.asdf/installs/golang/1.26.5/go/src/runtime/string.go:290).

## 1. Literal storage

String literal bytes are commonly stored in the program's read-only data:

```go
s := "hello"
```

Creating the header for a literal normally does not allocate those five bytes
on the heap. Multiple headers may refer to the same immutable literal data.

## 2. Dynamically built strings

Concatenation may allocate a new byte sequence:

```go
s := first + second
```

The compiler/runtime calculates the total byte count, obtains storage, and
copies the operands unless an optimization proves it can avoid that work.
Repeated concatenation in a loop can create many temporary strings.

`strings.Builder` or `bytes.Buffer` can amortize construction when many pieces
are appended. A builder still allocates as its internal buffer grows; it mainly
reduces repeated copying and temporary allocations.

## 3. Substrings share bytes

Slicing normally creates another 16-byte header without copying:

```go
large := loadLargeString()
small := large[:10]
large = ""
```

`small` may retain the complete allocation that contains `large`. When the
small substring is long-lived, force an independent copy if releasing the
large storage matters:

```go
small = strings.Clone(small)
```

For a non-empty string, `strings.Clone` returns a fresh copy. Its purpose is to
avoid unwanted retention of a much larger string; use it only when that saving
justifies the copy.

## 4. UTF-8 and rune processing

Strings are arbitrary byte sequences; Go source literals are normally UTF-8.
For:

```go
s := "世"
```

`len(s)` is 3, while `utf8.RuneCountInString(s)` is 1. A `range` loop decodes
runes as it advances and uses constant auxiliary space. Converting to `[]rune`
allocates a separate array of 4-byte `rune` values unless optimized away.

## 5. String and byte-slice conversions

Because strings are immutable and byte slices are mutable, these conversions
normally copy data:

```go
b := []byte(s)
t := string(b)
```

Each destination also has its own header: 24 bytes for the slice, 16 bytes for
the string. Some compiler-known, read-only uses can use temporary zero-copy
conversions, but programs must not rely on that optimization.

Unsafe zero-copy conversions transfer lifetime and immutability hazards to the
program and should not be a default memory optimization.

## 6. Copying a string value

```go
b := a
```

This copies the 16-byte header, not the bytes. Both strings view the same
immutable data, which is safe because normal Go code cannot mutate string
contents.

## 7. Garbage collection

The data pointer keeps the underlying byte allocation reachable. The bytes
contain no pointers and need no GC scanning, but the allocation cannot be
reclaimed while any string header still refers to it.

## 8. CPU and measurement

`len` and slicing are constant-time. Comparing strings is O(n) in the worst
case; unequal lengths or differing early bytes can stop sooner. Concatenation
and copying are O(total bytes).

```go
unsafe.Sizeof("")      // 16 on this machine
unsafe.Sizeof("hello") // also 16; payload is not counted
```

Use `go test -bench=. -benchmem` to measure construction and conversion paths.
