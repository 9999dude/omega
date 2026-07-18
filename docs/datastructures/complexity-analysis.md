# Complexity Analysis in Data Structures

Complexity analysis answers one basic question:

> **As the input becomes bigger, how much more work or memory does my program need?**

It does not usually calculate the exact number of milliseconds. It describes how the program **grows**.

---

# 1. Explain It Like a Child

Imagine you have toy boxes.

## Problem A: Find the first toy

The first toy is always at the front.

Whether there are:

* 10 toys
* 1,000 toys
* 1 million toys

You directly pick the first one.

```text
Work stays the same.
```

This is called:

[
O(1)
]

## Problem B: Find your red car

The toys are not arranged. You inspect them one by one.

```text
Toy 1 → not red
Toy 2 → not red
Toy 3 → not red
...
```

If there are twice as many toys, you may need twice as much work.

This is:

[
O(n)
]

Here, `n` means the number of toys.

## Problem C: Compare every toy with every other toy

For every toy, you inspect all the toys again.

```text
Toy 1 compares with all toys
Toy 2 compares with all toys
Toy 3 compares with all toys
...
```

If you have `n` toys, you perform roughly:

[
n \times n
]

operations.

This is:

[
O(n^2)
]

---

# 2. Why Complexity Analysis Is Needed

Suppose two programs solve the same problem.

### Program A

```text
For 10 users:        100 operations
For 1,000 users:     1,000,000 operations
For 1 million users: 1,000,000,000,000 operations
```

### Program B

```text
For 10 users:        around 4 operations
For 1,000 users:     around 10 operations
For 1 million users: around 20 operations
```

Both might appear fast with 10 users. Their difference becomes enormous at scale.

Complexity analysis helps you answer:

* Will this solution work for one million records?
* Which data structure should I use?
* Will memory usage become too large?
* Why is the system becoming slow?
* Can this algorithm pass an interview constraint?
* Can this production system scale?

## Real-life example

You need to check whether a username exists.

### Using a list

```text
["alice", "bob", "charlie", ...]
```

You may inspect every username.

[
O(n)
]

### Using a hash set

```text
{"alice", "bob", "charlie"}
```

You can usually look it up directly.

[
O(1)
]

Choosing the correct data structure changes the performance of the entire solution.

---

# 3. What Does `n` Mean?

`n` represents the size of the input.

Examples:

| Problem                | Meaning of `n`             |
| ---------------------- | -------------------------- |
| Search an array        | Number of elements         |
| Process users          | Number of users            |
| Traverse a linked list | Number of nodes            |
| Process a string       | Number of characters       |
| Traverse a tree        | Number of tree nodes       |
| Traverse a graph       | Usually vertices and edges |
| Process a matrix       | Rows and columns           |

For graphs, complexity is often written using:

* `V` = number of vertices
* `E` = number of edges

For example, graph traversal is:

[
O(V + E)
]

For a matrix with `m` rows and `n` columns:

[
O(mn)
]

---

# 4. Time Complexity and Space Complexity

## Time complexity

Time complexity measures how the number of operations grows.

It does not directly mean seconds.

```go
for i := 0; i < n; i++ {
    fmt.Println(i)
}
```

The loop runs `n` times:

[
O(n)
]

## Space complexity

Space complexity measures how much additional memory is required.

```go
result := make([]int, n)
```

The program creates an array containing `n` elements:

[
O(n)
]

This program uses constant additional memory:

```go
sum := 0

for _, value := range numbers {
    sum += value
}
```

Only `sum` and the loop variables are needed:

[
O(1)
]

The input array usually does not count as auxiliary space because it already exists.

---

# 5. What Is Big O?

Big O describes the **upper growth rate** of an algorithm.

It answers:

> How badly can the amount of work grow as the input becomes large?

Common complexities:

| Complexity   | Name         | Growth                                  |
| ------------ | ------------ | --------------------------------------- |
| `O(1)`       | Constant     | Excellent                               |
| `O(log n)`   | Logarithmic  | Excellent                               |
| `O(n)`       | Linear       | Good                                    |
| `O(n log n)` | Linearithmic | Usually good                            |
| `O(n²)`      | Quadratic    | Often acceptable only for smaller input |
| `O(2ⁿ)`      | Exponential  | Very expensive                          |
| `O(n!)`      | Factorial    | Extremely expensive                     |

---

# 6. The Most Important Mental Model

Ask:

> **How many times does my code touch the data?**

## One direct touch

```go
value := numbers[5]
```

[
O(1)
]

## Touch every element once

```go
for _, value := range numbers {
    fmt.Println(value)
}
```

[
O(n)
]

## Touch every element for every element

```go
for _, first := range numbers {
    for _, second := range numbers {
        fmt.Println(first, second)
    }
}
```

[
O(n^2)
]

## Throw away half the data repeatedly

```text
1,000 items
500 items
250 items
125 items
...
```

[
O(\log n)
]

This is the idea behind binary search.

---

# 7. Understanding Common Complexities

## O(1): Constant time

The amount of work does not depend on input size.

```go
func first(numbers []int) int {
    return numbers[0]
}
```

Whether the array has 10 or 10 million elements, accessing one index takes approximately the same amount of work.

### Mental model

> Go directly to one known location.

Examples:

* Array index access
* Stack push
* Stack pop
* Hash map lookup on average
* Reading a variable

---

## O(log n): Logarithmic time

The input becomes smaller by a fixed ratio during every step.

Binary search removes half of the remaining search space.

```text
1,024
512
256
128
64
32
16
8
4
2
1
```

Only 10 steps are needed to reduce 1,024 items to one item.

Because:

[
2^{10} = 1024
]

Therefore:

[
\log_2 1024 = 10
]

### Mental model

> Keep cutting the problem in half.

Examples:

* Binary search
* Balanced binary-search-tree lookup
* Heap insertion
* Heap deletion

---

## O(n): Linear time

Work grows directly with input size.

```go
func contains(numbers []int, target int) bool {
    for _, number := range numbers {
        if number == target {
            return true
        }
    }

    return false
}
```

In the worst case, every element must be checked.

[
O(n)
]

### Mental model

> Visit everyone once.

Examples:

* Linear search
* Finding maximum in an unsorted array
* Linked-list traversal
* Counting characters

---

## O(n log n): Linearithmic time

You perform logarithmic work for each level, while processing all `n` elements across the level.

Common sorting algorithms have this complexity:

* Merge sort
* Heap sort
* Average-case quicksort

### Mental model

> Split the input into levels, but process all elements at each level.

For merge sort:

```text
Level 1: process n elements
Level 2: process n elements
Level 3: process n elements
...
Number of levels: log n
```

Therefore:

[
n \times \log n = O(n \log n)
]

---

## O(n²): Quadratic time

Usually caused by two full nested loops.

```go
for i := 0; i < n; i++ {
    for j := 0; j < n; j++ {
        fmt.Println(i, j)
    }
}
```

The inner loop runs `n` times for each of the `n` outer iterations:

[
n \times n = n^2
]

### Mental model

> Everyone meets everyone.

Examples:

* Comparing every pair
* Bubble sort
* Selection sort
* Insertion sort worst case
* Brute-force duplicate detection

---

## O(2ⁿ): Exponential time

Every decision creates multiple new possibilities.

For a set containing `n` elements, there are:

[
2^n
]

possible subsets.

```text
n = 3  → 8 possibilities
n = 10 → 1,024 possibilities
n = 20 → 1,048,576 possibilities
n = 40 → over 1 trillion possibilities
```

### Mental model

> Every choice creates two more worlds.

Examples:

* Generating every subset
* Naive Fibonacci recursion
* Some backtracking problems

---

## O(n!): Factorial time

Used when examining every possible ordering.

For three objects:

```text
ABC
ACB
BAC
BCA
CAB
CBA
```

There are:

[
3! = 6
]

arrangements.

For 10 objects:

[
10! = 3,628,800
]

### Mental model

> Try every possible arrangement.

Examples:

* Generating every permutation
* Brute-force travelling salesperson problem

---

# 8. Growth Comparison

Suppose `n = 1,000`.

| Complexity   | Approximate operations |
| ------------ | ---------------------: |
| `O(1)`       |                      1 |
| `O(log n)`   |                     10 |
| `O(n)`       |                  1,000 |
| `O(n log n)` |                 10,000 |
| `O(n²)`      |              1,000,000 |
| `O(2ⁿ)`      |       Impossibly large |
| `O(n!)`      |       Impossibly large |

This is why complexity matters more than small code-level optimizations.

Changing an algorithm from `O(n²)` to `O(n log n)` is usually much more valuable than making each individual operation slightly faster.

---

# 9. How to Calculate Complexity

Use these rules.

---

## Rule 1: Simple statements are O(1)

```go
x := 10
y := x + 20
numbers[0] = y
```

Each statement takes constant time.

```text
O(1) + O(1) + O(1)
```

This becomes:

[
O(1)
]

Constants are ignored.

---

## Rule 2: Consecutive operations are added

```go
for i := 0; i < n; i++ {
    fmt.Println(i)
}

for i := 0; i < n; i++ {
    fmt.Println(i)
}
```

The first loop is `O(n)`.

The second loop is `O(n)`.

```text
O(n) + O(n) = O(2n)
```

Ignore the constant `2`:

[
O(n)
]

### Important

Consecutive loops are usually added, not multiplied.

---

## Rule 3: Nested loops are multiplied

```go
for i := 0; i < n; i++ {
    for j := 0; j < n; j++ {
        fmt.Println(i, j)
    }
}
```

```text
Outer loop: n
Inner loop: n
```

Therefore:

[
n \times n = O(n^2)
]

Three full nested loops would usually be:

[
O(n^3)
]

---

## Rule 4: Keep the dominant term

Suppose the number of operations is:

[
n^2 + 5n + 100
]

For large `n`, `n²` dominates everything else.

Therefore:

[
O(n^2)
]

Examples:

```text
O(n + 10)       → O(n)
O(5n)           → O(n)
O(n² + n)       → O(n²)
O(n log n + n)  → O(n log n)
```

---

## Rule 5: Halving or doubling often means O(log n)

```go
for i := 1; i < n; i *= 2 {
    fmt.Println(i)
}
```

Values of `i`:

```text
1, 2, 4, 8, 16, 32...
```

The question is:

> How many times can we multiply by two before reaching `n`?

The answer is approximately `log₂ n`.

Therefore:

[
O(\log n)
]

The same applies when dividing:

```go
for n > 1 {
    n /= 2
}
```

---

## Rule 6: For conditionals, take the more expensive branch

```go
if condition {
    fmt.Println("hello")
} else {
    for i := 0; i < n; i++ {
        fmt.Println(i)
    }
}
```

The first branch is `O(1)`.

The second branch is `O(n)`.

Worst-case complexity:

[
O(n)
]

---

## Rule 7: Different inputs need different variables

```go
for _, user := range users {
    fmt.Println(user)
}

for _, product := range products {
    fmt.Println(product)
}
```

If there are `n` users and `m` products:

[
O(n + m)
]

Do not automatically call it `O(n)` because the collections may have different sizes.

Nested version:

```go
for _, user := range users {
    for _, product := range products {
        fmt.Println(user, product)
    }
}
```

Complexity:

[
O(nm)
]

---

# 10. A Common Nested-Loop Trap

Not every nested loop is `O(n²)`.

Consider:

```go
for i := 0; i < n; i++ {
    for j := i; j < n; j++ {
        fmt.Println(i, j)
    }
}
```

The inner loop runs:

```text
n times
n - 1 times
n - 2 times
...
1 time
```

Total:

[
n + (n-1) + (n-2) + ... + 1
]

This equals:

[
\frac{n(n+1)}{2}
]

After removing constants and lower-order terms:

[
O(n^2)
]

However, consider this:

```go
for i := 0; i < n; i++ {
    j := i

    for j < n {
        j *= 2
    }
}
```

Here, the inner loop may be logarithmic, making the total closer to:

[
O(n \log n)
]

You must analyze how the loop variables change rather than merely counting loop keywords.

---

# 11. Best, Average and Worst Case

Consider linear search:

```go
func find(numbers []int, target int) bool {
    for _, number := range numbers {
        if number == target {
            return true
        }
    }

    return false
}
```

## Best case

The target is the first element.

[
O(1)
]

## Worst case

The target is last or missing.

[
O(n)
]

## Average case

The target is somewhere around the middle.

Approximately `n/2` checks:

[
O(n)
]

Because constants are ignored.

In interviews, complexity usually means **worst-case complexity**, unless otherwise specified.

---

# 12. Amortized Complexity

Some operations are occasionally expensive but cheap on average.

Consider appending to a dynamic array.

Most append operations place the element in available space:

[
O(1)
]

When the array becomes full, it must:

1. Allocate a larger array.
2. Copy the old elements.
3. Insert the new element.

That particular append is:

[
O(n)
]

But resizing does not happen on every append. Across many append operations, the average cost per append is:

[
O(1) \text{ amortized}
]

### Mental model

Imagine paying SGD 1 into a repair fund every day. Most days, no repair happens. Occasionally, the accumulated money pays for an expensive repair.

The occasional cost is distributed across many cheap operations.

---

# 13. Common Data Structure Complexities

## Arrays and dynamic arrays

| Operation                              |       Complexity |
| -------------------------------------- | ---------------: |
| Access by index                        |           `O(1)` |
| Update by index                        |           `O(1)` |
| Search unsorted array                  |           `O(n)` |
| Search sorted array with binary search |       `O(log n)` |
| Append to dynamic array                | `O(1)` amortized |
| Insert at beginning                    |           `O(n)` |
| Delete from beginning                  |           `O(n)` |
| Insert in middle                       |           `O(n)` |
| Delete from middle                     |           `O(n)` |

Why is inserting at the beginning `O(n)`?

```text
Before: [A, B, C, D]
Insert X
After:  [X, A, B, C, D]
```

All existing elements must shift.

---

## Linked lists

| Operation                     | Complexity |
| ----------------------------- | ---------: |
| Access by position            |     `O(n)` |
| Search                        |     `O(n)` |
| Insert at head                |     `O(1)` |
| Delete head                   |     `O(1)` |
| Insert after a known node     |     `O(1)` |
| Delete a known node           |     `O(1)` |
| Find the node before deleting |     `O(n)` |

Important interview distinction:

> Deleting a node is `O(1)` only when you already have the required node reference and any required predecessor information.

---

## Stack

| Operation | Complexity |
| --------- | ---------: |
| Push      |     `O(1)` |
| Pop       |     `O(1)` |
| Peek      |     `O(1)` |
| Search    |     `O(n)` |

Mental model:

> A stack is a pile of plates. You interact with the top.

---

## Queue

| Operation | Complexity |
| --------- | ---------: |
| Enqueue   |     `O(1)` |
| Dequeue   |     `O(1)` |
| Peek      |     `O(1)` |
| Search    |     `O(n)` |

Mental model:

> A queue is a line of people. Join at the back and leave from the front.

This assumes an appropriate queue implementation. Removing the first element from a basic dynamic array may require shifting and become `O(n)`.

---

## Hash map and hash set

| Operation | Average |  Worst |
| --------- | ------: | -----: |
| Insert    |  `O(1)` | `O(n)` |
| Search    |  `O(1)` | `O(n)` |
| Delete    |  `O(1)` | `O(n)` |

Worst-case `O(n)` can happen because of excessive collisions.

Mental model:

> A hash function gives you a drawer number, so you usually go directly to the correct drawer.

---

## Binary search tree

### Balanced BST

| Operation | Complexity |
| --------- | ---------: |
| Search    | `O(log n)` |
| Insert    | `O(log n)` |
| Delete    | `O(log n)` |

### Unbalanced BST

A badly shaped tree may look like a linked list:

```text
1
 \
  2
   \
    3
     \
      4
```

Search may become:

[
O(n)
]

Balanced trees such as AVL and red-black trees maintain approximately logarithmic height.

---

## Heap

| Operation                 | Complexity |
| ------------------------- | ---------: |
| Read minimum or maximum   |     `O(1)` |
| Insert                    | `O(log n)` |
| Remove minimum or maximum | `O(log n)` |
| Search arbitrary value    |     `O(n)` |
| Build heap                |     `O(n)` |

Important:

A heap lets you find the highest- or lowest-priority element efficiently. It does not provide efficient arbitrary search.

---

## Graphs

Using an adjacency list:

| Operation | Complexity |
| --------- | ---------: |
| BFS       | `O(V + E)` |
| DFS       | `O(V + E)` |
| Space     | `O(V + E)` |

Where:

* `V` = vertices
* `E` = edges

Every vertex and edge may need to be examined.

---

# 14. Space Complexity Mental Model

Ask:

> What new memory grows when the input grows?

## Constant space

```go
func maximum(numbers []int) int {
    maxValue := numbers[0]

    for _, number := range numbers {
        if number > maxValue {
            maxValue = number
        }
    }

    return maxValue
}
```

Additional memory:

```text
maxValue
number
```

The number of variables does not grow with `n`.

[
O(1)
]

## Linear space

```go
func copyNumbers(numbers []int) []int {
    copied := make([]int, len(numbers))
    copy(copied, numbers)
    return copied
}
```

The new array grows with the input:

[
O(n)
]

## Recursive space

```go
func countdown(n int) {
    if n == 0 {
        return
    }

    countdown(n - 1)
}
```

There are `n` recursive calls on the call stack:

[
O(n)
]

Even though no explicit array was created, the call stack consumes memory.

---

# 15. Time-Space Trade-off

You can often use additional memory to improve execution time.

## Find duplicates without extra memory

```go
for i := 0; i < len(numbers); i++ {
    for j := i + 1; j < len(numbers); j++ {
        if numbers[i] == numbers[j] {
            return true
        }
    }
}
```

Time:

[
O(n^2)
]

Space:

[
O(1)
]

## Find duplicates using a hash set

```go
func hasDuplicate(numbers []int) bool {
    seen := make(map[int]struct{})

    for _, number := range numbers {
        if _, exists := seen[number]; exists {
            return true
        }

        seen[number] = struct{}{}
    }

    return false
}
```

Time:

[
O(n)
]

Space:

[
O(n)
]

You traded more memory for less execution time.

---

# 16. The Four-Step Interview Mental Model

When asked to calculate complexity, use this sequence.

## Step 1: Define the input size

Say:

> Let `n` be the number of elements in the input array.

For two collections:

> Let `n` be the number of users and `m` be the number of products.

## Step 2: Find repeated work

Identify:

* Loops
* Nested loops
* Recursion
* Sorting
* Data-structure operations
* String concatenation
* Copying
* Hash-map operations

## Step 3: Count how repetition grows

Ask:

* Does it run once?
* Does it run `n` times?
* Does it halve the data?
* Does it run `n` times inside another `n` loop?
* Does it branch into multiple recursive calls?

## Step 4: Simplify

Remove:

* Constants
* Lower-order terms

Examples:

```text
3n + 10        → O(n)
n² + 100n      → O(n²)
2n log n + n   → O(n log n)
```

Then separately explain space complexity.

---

# 17. A Compact Pattern-Recognition Table

| Code pattern                                |         Likely complexity |
| ------------------------------------------- | ------------------------: |
| Direct index access                         |                    `O(1)` |
| One full loop                               |                    `O(n)` |
| Two consecutive full loops                  |                    `O(n)` |
| Two nested full loops                       |                   `O(n²)` |
| Loop variable doubles                       |                `O(log n)` |
| Loop variable halves                        |                `O(log n)` |
| Sort then scan                              |              `O(n log n)` |
| Hash-map scan                               | `O(n)` time, `O(n)` space |
| Recursive call reducing by one              |            Usually `O(n)` |
| Two recursive calls reducing by one         |             Often `O(2ⁿ)` |
| Divide into halves and process all elements |        Often `O(n log n)` |

---

# 18. Common Interview Mistakes

## Mistake 1: Saying two consecutive loops are O(n²)

```go
for i := 0; i < n; i++ {
}

for j := 0; j < n; j++ {
}
```

Correct:

[
O(n+n)=O(n)
]

They are not nested.

---

## Mistake 2: Calling every nested loop O(n²)

```go
for i := 1; i < n; i *= 2 {
    for j := 0; j < n; j++ {
    }
}
```

Outer loop:

[
O(\log n)
]

Inner loop:

[
O(n)
]

Total:

[
O(n \log n)
]

---

## Mistake 3: Ignoring library function complexity

```go
sort.Ints(numbers)
```

Sorting is not `O(1)` merely because it is one line.

It is generally:

[
O(n \log n)
]

Similarly:

```go
strings.Contains(...)
copy(...)
append(...)
```

may perform work proportional to the amount of data involved.

---

## Mistake 4: Ignoring space used by recursion

Recursive functions use the call stack.

A recursion depth of `n` generally uses:

[
O(n)
]

stack space.

---

## Mistake 5: Saying hash maps are always O(1)

The more precise answer is:

> Hash-map insertion, lookup and deletion are `O(1)` on average and `O(n)` in the worst case.

---

## Mistake 6: Giving complexity without explaining `n`

Always define it:

> Let `n` be the number of elements in the array.

---

# 19. Worked Interview Examples

## Question 1

What is the time complexity?

```go
func printFirst(numbers []int) {
    fmt.Println(numbers[0])
}
```

### Answer

Only one array access is performed.

Time:

[
O(1)
]

Space:

[
O(1)
]

---

## Question 2

```go
func printAll(numbers []int) {
    for _, number := range numbers {
        fmt.Println(number)
    }
}
```

### Answer

Every element is visited once.

Time:

[
O(n)
]

Additional space:

[
O(1)
]

---

## Question 3

```go
func printPairs(numbers []int) {
    for _, first := range numbers {
        for _, second := range numbers {
            fmt.Println(first, second)
        }
    }
}
```

### Answer

For each of `n` elements, the code visits `n` elements.

[
n \times n = n^2
]

Time:

[
O(n^2)
]

Space:

[
O(1)
]

---

## Question 4

```go
func strangeLoop(n int) {
    for i := 1; i < n; i *= 2 {
        fmt.Println(i)
    }
}
```

### Answer

The value doubles:

```text
1, 2, 4, 8, 16...
```

The loop executes approximately `log₂ n` times.

Time:

[
O(\log n)
]

Space:

[
O(1)
]

---

## Question 5

```go
func mixed(numbers []int) {
    for _, number := range numbers {
        fmt.Println(number)
    }

    for _, first := range numbers {
        for _, second := range numbers {
            fmt.Println(first, second)
        }
    }
}
```

### Answer

First loop:

[
O(n)
]

Nested loops:

[
O(n^2)
]

Total:

[
O(n+n^2)
]

The dominant term is `n²`.

Time:

[
O(n^2)
]

---

## Question 6

```go
func binarySearch(numbers []int, target int) int {
    left := 0
    right := len(numbers) - 1

    for left <= right {
        middle := left + (right-left)/2

        if numbers[middle] == target {
            return middle
        }

        if numbers[middle] < target {
            left = middle + 1
        } else {
            right = middle - 1
        }
    }

    return -1
}
```

### Answer

Each iteration removes half the search space.

Time:

[
O(\log n)
]

Space:

[
O(1)
]

The input must be sorted.

---

## Question 7

```go
func recursiveSum(n int) int {
    if n == 0 {
        return 0
    }

    return n + recursiveSum(n-1)
}
```

### Answer

There are `n` recursive calls.

Time:

[
O(n)
]

Each call remains on the stack until recursion returns.

Space:

[
O(n)
]

---

## Question 8

```go
func fibonacci(n int) int {
    if n <= 1 {
        return n
    }

    return fibonacci(n-1) + fibonacci(n-2)
}
```

### Answer

Each call creates two more recursive calls, with many calculations repeated.

Time:

[
O(2^n)
]

More precisely, it grows close to `O(φⁿ)`, but `O(2ⁿ)` is commonly accepted in interviews.

Recursion depth:

[
O(n)
]

Therefore space:

[
O(n)
]

---

# 20. Mock Interview Questions

## Mock Question 1

You have an unsorted array. What is the complexity of finding a value?

### Expected answer

A linear scan may inspect every element.

* Best case: `O(1)`
* Worst case: `O(n)`
* Average case: `O(n)`
* Space: `O(1)`

---

## Mock Question 2

What must be true before binary search can be used?

### Expected answer

The data must be sorted, or there must be a monotonic condition allowing one half of the search space to be discarded.

Binary search runs in:

[
O(\log n)
]

---

## Mock Question 3

Why is array access O(1), but linked-list access O(n)?

### Expected answer

An array stores elements in contiguous memory, so the address of index `i` can be calculated directly.

A linked list stores nodes connected by pointers. To reach position `i`, the program must follow nodes one by one.

---

## Mock Question 4

Why is insertion at the beginning of an array O(n)?

### Expected answer

All existing elements must shift one position to create space for the new element.

Inserting at the head of a linked list is `O(1)` because only pointers need to change.

---

## Mock Question 5

What is the difference between O(n) and O(2n)?

### Expected answer

In Big O notation, both are `O(n)` because constants are ignored.

They have the same growth category, though one implementation may still be approximately twice as slow in practice.

---

## Mock Question 6

Why do we ignore constants?

### Expected answer

Complexity analysis focuses on how performance grows for large inputs.

For sufficiently large `n`, the growth category matters more than fixed multipliers.

However, constants still matter in real production performance and should not always be ignored during optimization.

---

## Mock Question 7

What is the complexity of sorting an array and then scanning it?

Sorting:

[
O(n \log n)
]

Scanning:

[
O(n)
]

Total:

[
O(n \log n+n)
]

Dominant term:

[
O(n \log n)
]

---

## Mock Question 8

What is the complexity of checking duplicates using two nested loops?

### Expected answer

Each element may be compared against every other element.

Time:

[
O(n^2)
]

Space:

[
O(1)
]

Using a hash set can improve time to `O(n)` at the cost of `O(n)` space.

---

## Mock Question 9

What is amortized O(1)?

### Expected answer

An operation may occasionally be expensive, but when the total cost is distributed across a sequence of operations, the average cost per operation is constant.

Dynamic-array append is a standard example.

---

## Mock Question 10

Is a hash-map lookup always O(1)?

### Expected answer

No.

It is:

* `O(1)` on average
* `O(n)` in the worst case

The worst case can result from severe collisions or pathological behavior.

---

## Mock Question 11

Why is BFS O(V + E)?

### Expected answer

BFS may visit every vertex once and inspect every edge once.

Therefore:

[
O(V+E)
]

With an adjacency list, space is also generally:

[
O(V+E)
]

including graph storage, while the BFS queue and visited set use up to `O(V)` auxiliary space.

---

## Mock Question 12

How do you decide between an array and a linked list?

### Expected answer

Use an array when:

* Random access is important.
* Cache locality matters.
* Appending is common.
* Middle insertion and deletion are uncommon.

Use a linked list when:

* Frequent insertion or deletion occurs at known nodes.
* Random access is not needed.
* Pointer and memory overhead are acceptable.

In practice, arrays are often preferred because of cache locality and simpler memory management.

---

## Mock Question 13

What is the complexity of finding the minimum element in a min-heap?

### Expected answer

The minimum element is at the root.

[
O(1)
]

Removing it requires restoring the heap property:

[
O(\log n)
]

Searching for an arbitrary value is:

[
O(n)
]

---

## Mock Question 14

Can an O(n²) solution ever be acceptable?

### Expected answer

Yes. It depends on the input constraint.

For example:

* `n = 20`: usually harmless
* `n = 1,000`: one million operations may be acceptable
* `n = 100,000`: ten billion operations are generally unacceptable

Complexity must be evaluated together with input size, constants and system requirements.

---

# 21. How to Estimate Acceptable Complexity from Constraints

A rough interview guideline:

|            Maximum `n` | Often acceptable                         |
| ---------------------: | ---------------------------------------- |
|               `n ≤ 10` | `O(n!)` may be possible                  |
|               `n ≤ 20` | `O(2ⁿ)` may be possible                  |
|            `n ≤ 1,000` | `O(n²)` may be possible                  |
|          `n ≤ 100,000` | Usually `O(n log n)` or `O(n)`           |
|        `n ≤ 1,000,000` | Usually `O(n)` or `O(log n)`             |
| Very large / streaming | Often `O(log n)` or `O(1)` per operation |

These are approximations, not strict rules.

---

# 22. A Strong Interview Answer Template

When asked to explain your solution, say:

> Let `n` represent the number of elements in the input. I traverse the input once, so the loop performs `n` iterations. Each hash-map lookup is `O(1)` on average. Therefore, the overall average time complexity is `O(n)`. The hash map may store up to `n` elements, so the auxiliary space complexity is `O(n)`.

This answer is strong because it:

1. Defines `n`.
2. Identifies the repeated operation.
3. Mentions the data-structure operation.
4. Gives time complexity.
5. Gives space complexity.
6. States average or worst case precisely.

---

# 23. Final Mental Model

Remember these pictures:

```text
O(1)       → Pick one known toy
O(log n)   → Keep cutting the toy pile in half
O(n)       → Check every toy once
O(n log n) → Divide into levels and process all toys per level
O(n²)      → Every toy meets every toy
O(2ⁿ)      → Every decision creates two worlds
O(n!)      → Try every possible arrangement
```

And remember this analysis formula:

```text
Consecutive work → Add
Nested work      → Multiply
Halving input    → Logarithm
Recursive calls  → Draw the call tree
Final answer     → Keep the fastest-growing term
```

The most useful question is:

> **When the input doubles, what happens to the amount of work?**

```text
O(1)       → stays approximately the same
O(log n)   → increases slightly
O(n)       → approximately doubles
O(n log n) → slightly more than doubles
O(n²)      → approximately quadruples
O(2ⁿ)      → grows catastrophically
```
