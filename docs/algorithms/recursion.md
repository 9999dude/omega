# Recursion — explained like a baby

Recursion means:

> **A function solves a problem by asking a smaller version of itself to solve part of the problem.**

Imagine opening a Russian nesting doll:

1. Open the current doll.
2. Find a smaller doll.
3. Open that smaller doll.
4. Continue until there is no smaller doll.
5. Close them in reverse order.

That is recursion.

```mermaid
flowchart TD
    A[Big problem] --> B[Smaller problem]
    B --> C[Even smaller problem]
    C --> D[Tiny problem]
    D --> E[Base case: answer directly]
    E --> F[Return answer upward]
    F --> G[Combine answers]
```

---

# 1. The smallest possible example

Suppose we want to count down from `3`.

```go
// Exact question: How does `countdown` solve The smallest possible example?
//
// Possible answer: Use `countdown` to format and write the demonstrated value to standard output.
//
// Output format: `countdown` has no return value; its observable result is the mutation or output performed in the function body.
//
// Inline descriptions:
// - `n` is the int input used by this example.
//
// Boundary checks:
// - `n == 0` handles the smallest valid state or recursive base case.
//
// Key variables:
// - `n` is the int input used by this example.
//
// Logic:
// 1. Recursively reduce the current problem to smaller calls until a base condition is reached.
// 2. Return the value produced after the state updates are complete.
func countdown(n int) {
	if n == 0 {
		fmt.Println("Done")
		return
	}

	fmt.Println(n)
	countdown(n - 1)
}

// time complexity: O(n) -> the algorithm visits each of the `n` input elements or states once.
// space complexity: O(n) -> the recursion stack and any cache can grow to one entry per input state.
```

Calling:

```go
// Exact question: How does this Go example demonstrate The smallest possible example?
//
// Possible answer: Use the fragment to execute the shown state update directly from top to bottom.
//
// Output format: This fragment demonstrates syntax or state updates and does not define a standalone output value.
//
// Inline descriptions:
// - The comments in this preface describe how the important expressions and state changes are used.
//
// Boundary checks:
// - No explicit boundary branch appears in this fragment; its caller or surrounding example supplies valid inputs.
//
// Key variables:
// - This fragment operates directly on the values named in each statement; it introduces no separate data structure.
//
// Logic:
// 1. Execute the statements from top to bottom to perform the demonstrated operation.
countdown(3)

// time complexity: O(1) -> the snippet performs a fixed number of operations independent of input size.
// space complexity: O(1) -> only a fixed number of scalar variables or references is kept.
```

Produces:

```text
3
2
1
Done
```

The function keeps calling itself with a smaller number:

```text
countdown(3)
countdown(2)
countdown(1)
countdown(0)
```

At `0`, it stops.

---

# 2. The two required parts of recursion

Every recursive solution needs two things.

## Part 1: Base case

The condition that stops recursion.

```go
// Exact question: How does this Go example demonstrate Part 1: Base case?
//
// Possible answer: Use the fragment to compare the current values and return or update state when the condition matches.
//
// Output format: This is a partial Go fragment; its surrounding function determines the final returned value.
//
// Inline descriptions:
// - The comments in this preface describe how the important expressions and state changes are used.
//
// Boundary checks:
// - `n == 0` handles the smallest valid state or recursive base case.
//
// Key variables:
// - This fragment operates directly on the values named in each statement; it introduces no separate data structure.
//
// Logic:
// 1. Return the value produced after the state updates are complete.
if n == 0 {
	return
}

// time complexity: O(1) -> the snippet performs a fixed number of operations independent of input size.
// space complexity: O(1) -> only a fixed number of scalar variables or references is kept.
```

Without a base case, the function keeps calling itself until the program runs out of stack memory.

## Part 2: Recursive case

The function calls itself with a smaller or simpler problem.

```go
// Exact question: How does this Go example demonstrate Part 2: Recursive case?
//
// Possible answer: Use the fragment to execute the shown state update directly from top to bottom.
//
// Output format: This fragment demonstrates syntax or state updates and does not define a standalone output value.
//
// Inline descriptions:
// - The comments in this preface describe how the important expressions and state changes are used.
//
// Boundary checks:
// - No explicit boundary branch appears in this fragment; its caller or surrounding example supplies valid inputs.
//
// Key variables:
// - This fragment operates directly on the values named in each statement; it introduces no separate data structure.
//
// Logic:
// 1. Execute the statements from top to bottom to perform the demonstrated operation.
countdown(n - 1)

// time complexity: O(1) -> the snippet performs a fixed number of operations independent of input size.
// space complexity: O(1) -> only a fixed number of scalar variables or references is kept.
```

The input must move toward the base case.

```mermaid
flowchart TD
    A[Start with input n] --> B{Base case reached?}
    B -- Yes --> C[Return direct answer]
    B -- No --> D[Reduce the problem]
    D --> E[Call function recursively]
    E --> F[Use returned result]
    F --> G[Return final result]
```

A useful formula is:

```text
Recursive solution =
Base case
+
Smaller recursive problem
+
Work done at the current level
```

---

# 3. Why is recursion needed?

Recursion is useful when the problem itself contains smaller versions of the same problem.

For example:

* A folder contains files and more folders.
* A tree node contains smaller trees.
* A manager has employees who may manage other employees.
* A decision produces more decisions.
* A maze contains smaller paths to explore.
* A string can be split into smaller strings.

These structures are naturally recursive.

## Folder example

```text
Documents
├── Resume.pdf
├── Work
│   ├── ProjectA
│   │   └── Design.md
│   └── ProjectB
└── Personal
    └── Photos
```

To search all files:

1. Examine the current folder.
2. For every subfolder, perform the same search.
3. Stop when a folder has no subfolders.

That is recursion.

---

# 4. The most important mental model

Think of a recursive function as making a promise:

> “Give me the answer for the smaller problem, and I will use it to construct the answer for the current problem.”

For factorial:

```text
factorial(5)
```

The function thinks:

> “I do not know `5!`, but if somebody gives me `4!`, I can multiply it by `5`.”

So:

```text
5! = 5 × 4!
4! = 4 × 3!
3! = 3 × 2!
2! = 2 × 1!
1! = 1
```

The recursive function does not need to understand the entire chain at once.

It only needs to know:

1. What is the smallest answer?
2. How can I reduce the current problem?
3. How do I use the smaller answer?

---

# 5. How recursion is calculated

Let us calculate factorial.

```text
5! = 5 × 4 × 3 × 2 × 1 = 120
```

Recursive definition:

```text
factorial(n) = n × factorial(n - 1)
factorial(1) = 1
```

Go implementation:

```go
// Exact question: How recursion is calculated?
//
// Possible answer: Use `factorial` to compare the current values and return or update state when the condition matches.
//
// Output format: Return an `int` value from `factorial`; the function does not print the answer.
//
// Inline descriptions:
// - `n` is the int input used by this example.
//
// Boundary checks:
// - `n <= 1` handles the smallest valid state or recursive base case.
//
// Key variables:
// - `n` is the int input used by this example.
//
// Logic:
// 1. Recursively reduce the current problem to smaller calls until a base condition is reached.
// 2. Return the value produced after the state updates are complete.
func factorial(n int) int {
	if n <= 1 {
		return 1
	}

	return n * factorial(n-1)
}

// time complexity: O(n) -> the algorithm visits each of the `n` input elements or states once.
// space complexity: O(n) -> the recursion stack and any cache can grow to one entry per input state.
```

Calling:

```go
// Exact question: How recursion is calculated?
//
// Possible answer: Use the fragment to execute the shown state update directly from top to bottom.
//
// Output format: This fragment demonstrates syntax or state updates and does not define a standalone output value.
//
// Inline descriptions:
// - The comments in this preface describe how the important expressions and state changes are used.
//
// Boundary checks:
// - No explicit boundary branch appears in this fragment; its caller or surrounding example supplies valid inputs.
//
// Key variables:
// - This fragment operates directly on the values named in each statement; it introduces no separate data structure.
//
// Logic:
// 1. Execute the statements from top to bottom to perform the demonstrated operation.
factorial(4)

// time complexity: O(1) -> the snippet performs a fixed number of operations independent of input size.
// space complexity: O(1) -> only a fixed number of scalar variables or references is kept.
```

First, calls go downward:

```text
factorial(4)
= 4 × factorial(3)

factorial(3)
= 3 × factorial(2)

factorial(2)
= 2 × factorial(1)

factorial(1)
= 1
```

Then results return upward:

```text
factorial(1) = 1
factorial(2) = 2 × 1 = 2
factorial(3) = 3 × 2 = 6
factorial(4) = 4 × 6 = 24
```

```mermaid
flowchart TD
    A["factorial(4)"] --> B["4 × factorial(3)"]
    B --> C["3 × factorial(2)"]
    C --> D["2 × factorial(1)"]
    D --> E["factorial(1) = 1"]
    E --> F["factorial(2) = 2"]
    F --> G["factorial(3) = 6"]
    G --> H["factorial(4) = 24"]
```

Recursion has two phases:

```text
Calling downward
Unwinding upward
```

---

# 6. The call stack

Each function call needs to remember:

* Its arguments
* Its local variables
* Where execution should continue
* What result it is waiting for

This information is stored in the **call stack**.

For:

```go
// Exact question: How does this Go example demonstrate The call stack?
//
// Possible answer: Use the fragment to execute the shown state update directly from top to bottom.
//
// Output format: This fragment demonstrates syntax or state updates and does not define a standalone output value.
//
// Inline descriptions:
// - The comments in this preface describe how the important expressions and state changes are used.
//
// Boundary checks:
// - No explicit boundary branch appears in this fragment; its caller or surrounding example supplies valid inputs.
//
// Key variables:
// - This fragment operates directly on the values named in each statement; it introduces no separate data structure.
//
// Logic:
// 1. Execute the statements from top to bottom to perform the demonstrated operation.
factorial(4)

// time complexity: O(1) -> the snippet performs a fixed number of operations independent of input size.
// space complexity: O(1) -> only a fixed number of scalar variables or references is kept.
```

The stack grows like this:

```text
Top
┌──────────────────────┐
│ factorial(1), n = 1  │
├──────────────────────┤
│ factorial(2), n = 2  │
├──────────────────────┤
│ factorial(3), n = 3  │
├──────────────────────┤
│ factorial(4), n = 4  │
└──────────────────────┘
Bottom
```

When `factorial(1)` returns, its stack frame is removed.

Then `factorial(2)` continues.

```mermaid
sequenceDiagram
    participant Main
    participant F4 as factorial(4)
    participant F3 as factorial(3)
    participant F2 as factorial(2)
    participant F1 as factorial(1)

    Main->>F4: call
    F4->>F3: call
    F3->>F2: call
    F2->>F1: call
    F1-->>F2: return 1
    F2-->>F3: return 2
    F3-->>F4: return 6
    F4-->>Main: return 24
```

## Important consequence

Recursive algorithms use stack memory.

If recursion depth is `n`, the space complexity is often:

```text
O(n)
```

Even when the function uses no explicit array or map.

---

# 7. Recursion versus iteration

A recursive countdown:

```go
// Exact question: How does `countdown` solve Recursion versus iteration?
//
// Possible answer: Use `countdown` to format and write the demonstrated value to standard output.
//
// Output format: `countdown` has no return value; its observable result is the mutation or output performed in the function body.
//
// Inline descriptions:
// - `n` is the int input used by this example.
//
// Boundary checks:
// - `n == 0` handles the smallest valid state or recursive base case.
//
// Key variables:
// - `n` is the int input used by this example.
//
// Logic:
// 1. Recursively reduce the current problem to smaller calls until a base condition is reached.
// 2. Return the value produced after the state updates are complete.
func countdown(n int) {
	if n == 0 {
		return
	}

	fmt.Println(n)
	countdown(n - 1)
}

// time complexity: O(n) -> the algorithm visits each of the `n` input elements or states once.
// space complexity: O(n) -> the recursion stack and any cache can grow to one entry per input state.
```

An iterative countdown:

```go
// Exact question: How does `countdown` solve Recursion versus iteration?
//
// Possible answer: Use `countdown` to iterate through the input once, updating the running state for each element.
//
// Output format: `countdown` has no return value; its observable result is the mutation or output performed in the function body.
//
// Inline descriptions:
// - `n` is the int input used by this example.
//
// Boundary checks:
// - `n > 0` decides whether the branch or loop should continue for the current input.
//
// Key variables:
// - `n` is the int input used by this example.
//
// Logic:
// 1. Iterate through the required elements or states in the order shown.
func countdown(n int) {
	for n > 0 {
		fmt.Println(n)
		n--
	}
}

// time complexity: O(n) -> the algorithm visits each of the `n` input elements or states once.
// space complexity: O(1) -> only a fixed number of scalar variables or references is kept.
```

Both solve the same problem.

| Recursion                                | Iteration                           |
| ---------------------------------------- | ----------------------------------- |
| Uses function calls                      | Uses loops                          |
| Uses call-stack memory                   | Usually uses constant extra memory  |
| Often clearer for trees and decisions    | Often clearer for linear repetition |
| Can cause stack overflow                 | Usually avoids stack overflow       |
| Naturally expresses DFS and backtracking | May require an explicit stack       |

Use recursion when it makes the problem easier to understand.

Do not use recursion merely because it looks clever.

---

# 8. The five questions to ask before writing recursion

For every recursive problem, ask:

## 1. What does my function mean?

Example:

```go
// Exact question: What does my function mean?
//
// Possible answer: Use the fragment to execute the shown state update directly from top to bottom.
//
// Output format: This fragment demonstrates syntax or state updates and does not define a standalone output value.
//
// Inline descriptions:
// - The comments in this preface describe how the important expressions and state changes are used.
//
// Boundary checks:
// - No explicit boundary branch appears in this fragment; its caller or surrounding example supplies valid inputs.
//
// Key variables:
// - This fragment operates directly on the values named in each statement; it introduces no separate data structure.
//
// Logic:
// 1. Execute the statements from top to bottom to perform the demonstrated operation.
factorial(n)

// time complexity: O(1) -> the snippet performs a fixed number of operations independent of input size.
// space complexity: O(1) -> only a fixed number of scalar variables or references is kept.
```

means:

> Return the factorial of `n`.

For a tree:

```go
// Exact question: What does my function mean?
//
// Possible answer: Use the fragment to execute the shown state update directly from top to bottom.
//
// Output format: This fragment demonstrates syntax or state updates and does not define a standalone output value.
//
// Inline descriptions:
// - The comments in this preface describe how the important expressions and state changes are used.
//
// Boundary checks:
// - No explicit boundary branch appears in this fragment; its caller or surrounding example supplies valid inputs.
//
// Key variables:
// - This fragment operates directly on the values named in each statement; it introduces no separate data structure.
//
// Logic:
// 1. Execute the statements from top to bottom to perform the demonstrated operation.
maxDepth(node)

// time complexity: O(1) -> the snippet performs a fixed number of operations independent of input size.
// space complexity: O(1) -> only a fixed number of scalar variables or references is kept.
```

means:

> Return the maximum depth of the tree rooted at `node`.

This contract must be precise.

## 2. What is the smallest problem?

For factorial:

```text
factorial(0) = 1
factorial(1) = 1
```

For a tree:

```text
Depth of an empty tree = 0
```

## 3. How do I make the problem smaller?

Examples:

```text
n → n - 1
array → remaining array
tree → left and right subtrees
string → remaining characters
```

## 4. How do I use the smaller answer?

For factorial:

```text
n × factorial(n - 1)
```

For tree depth:

```text
1 + max(leftDepth, rightDepth)
```

## 5. Does every path reach the base case?

This prevents infinite recursion.

---

# 9. A reusable recursion template

```go
// Exact question: How does `solve` solve A reusable recursion template?
//
// Possible answer: Use `solve` to compare the current values and return or update state when the condition matches.
//
// Output format: Return a `Result` value from `solve`; the function does not print the answer.
//
// Inline descriptions:
// - `state` is the State input used by this example.
//
// Boundary checks:
// - `isBaseCase(state)` decides whether the branch or loop should continue for the current input.
//
// Key variables:
// - `state` is the State input used by this example.
// - `smallerState` holds the intermediate value produced by `makeSmaller(state`.
// - `smallerResult` holds the answer computed for the current operation.
//
// Logic:
// 1. Recursively reduce the current problem to smaller calls until a base condition is reached.
// 2. Return the value produced after the state updates are complete.
func solve(state State) Result {
	// 1. Base case
	if isBaseCase(state) {
		return baseResult
	}

	// 2. Reduce the problem
	smallerState := makeSmaller(state)

	// 3. Recursive call
	smallerResult := solve(smallerState)

	// 4. Use the smaller result
	return combine(state, smallerResult)
}

// time complexity: O(n) -> the algorithm visits each of the `n` input elements or states once.
// space complexity: O(n) -> the recursion stack and any cache can grow to one entry per input state.
```

For recursion with multiple branches:

```go
// Exact question: How does `solve` solve A reusable recursion template?
//
// Possible answer: Use `solve` to compare the current values and return or update state when the condition matches.
//
// Output format: Return a `Result` value from `solve`; the function does not print the answer.
//
// Inline descriptions:
// - `state` is the State input used by this example.
//
// Boundary checks:
// - `isBaseCase(state)` decides whether the branch or loop should continue for the current input.
//
// Key variables:
// - `state` is the State input used by this example.
// - `leftResult` holds the answer computed for the current operation.
// - `rightResult` holds the answer computed for the current operation.
//
// Logic:
// 1. Recursively reduce the current problem to smaller calls until a base condition is reached.
// 2. Return the value produced after the state updates are complete.
func solve(state State) Result {
	if isBaseCase(state) {
		return baseResult
	}

	leftResult := solve(firstChoice(state))
	rightResult := solve(secondChoice(state))

	return combine(leftResult, rightResult)
}

// time complexity: O(n) -> the algorithm visits each of the `n` input elements or states once.
// space complexity: O(n) -> the recursion stack and any cache can grow to one entry per input state.
```

---

# 10. Types of recursion

## 10.1 Linear recursion

Each call makes one recursive call.

Examples:

* Factorial
* Sum of an array
* Reverse a linked list
* Countdown

```mermaid
flowchart LR
    A["f(4)"] --> B["f(3)"]
    B --> C["f(2)"]
    C --> D["f(1)"]
```

Typical complexity:

```text
Time: O(n)
Space: O(n)
```

---

## 10.2 Tree recursion

Each call makes multiple recursive calls.

Examples:

* Fibonacci
* Tree traversal
* Generate subsets
* Backtracking

```mermaid
flowchart TD
    A["f(4)"] --> B["f(3)"]
    A --> C["f(2)"]
    B --> D["f(2)"]
    B --> E["f(1)"]
    C --> F["f(1)"]
    C --> G["f(0)"]
```

The number of calls can grow exponentially.

---

## 10.3 Divide and conquer

Split one large problem into independent smaller problems.

Examples:

* Merge sort
* Quick sort
* Binary search
* Tree algorithms

```mermaid
flowchart TD
    A[Large problem] --> B[Left half]
    A --> C[Right half]
    B --> D[Solve left]
    C --> E[Solve right]
    D --> F[Combine]
    E --> F
```

---

## 10.4 DFS recursion

Explore one path deeply before exploring another path.

Used for:

* Trees
* Graphs
* Mazes
* Islands
* File systems

```mermaid
flowchart TD
    A[A] --> B[B]
    A --> C[C]
    B --> D[D]
    B --> E[E]
    C --> F[F]

    D:::visit
    E:::visit
    B:::visit
    F:::visit
    C:::visit
    A:::visit

    classDef visit stroke-width:2px
```

A depth-first traversal might visit:

```text
A → B → D → E → C → F
```

---

## 10.5 Backtracking

Backtracking means:

> Choose something, explore it, undo the choice, and try another choice.

Examples:

* Subsets
* Permutations
* Generate Parentheses
* Combination Sum
* Sudoku
* N-Queens

```mermaid
flowchart TD
    A[Current state] --> B[Make choice A]
    A --> C[Make choice B]
    B --> D[Explore]
    D --> E[Undo choice A]
    E --> C
    C --> F[Explore]
```

Backtracking is recursion plus undoing decisions.

---

# 11. Problem 1: Factorial

## Problem

Calculate:

```text
n! = n × (n-1) × ... × 1
```

## Recursive thinking

```text
factorial(n) = n × factorial(n - 1)
```

## Code

```go
// Exact question: How does `factorial` solve Factorial?
//
// Possible answer: Use `factorial` to compare the current values and return or update state when the condition matches.
//
// Output format: Return an `int` value from `factorial`; the function does not print the answer.
//
// Inline descriptions:
// - `n` is the int input used by this example.
//
// Boundary checks:
// - `n <= 1` handles the smallest valid state or recursive base case.
//
// Key variables:
// - `n` is the int input used by this example.
//
// Logic:
// 1. Recursively reduce the current problem to smaller calls until a base condition is reached.
// 2. Return the value produced after the state updates are complete.
func factorial(n int) int {
	if n <= 1 {
		return 1
	}

	return n * factorial(n-1)
}

// time complexity: O(n) -> the algorithm visits each of the `n` input elements or states once.
// space complexity: O(n) -> the recursion stack and any cache can grow to one entry per input state.
```

## Complexity

```text
Time:  O(n)
Space: O(n)
```

The time is `O(n)` because there are `n` calls.

The space is `O(n)` because `n` stack frames exist at the deepest point.

---

# 12. Problem 2: Fibonacci

Fibonacci numbers:

```text
0, 1, 1, 2, 3, 5, 8, 13...
```

Definition:

```text
fib(n) = fib(n - 1) + fib(n - 2)
```

## Naive recursive code

```go
// Exact question: How does `fib` solve Naive recursive code?
//
// Possible answer: Use `fib` to compare the current values and return or update state when the condition matches.
//
// Output format: Return an `int` value from `fib`; the function does not print the answer.
//
// Inline descriptions:
// - `n` is the int input used by this example.
//
// Boundary checks:
// - `n <= 1` handles the smallest valid state or recursive base case.
//
// Key variables:
// - `n` is the int input used by this example.
//
// Logic:
// 1. Recursively reduce the current problem to smaller calls until a base condition is reached.
// 2. Return the value produced after the state updates are complete.
func fib(n int) int {
	if n <= 1 {
		return n
	}

	return fib(n-1) + fib(n-2)
}

// time complexity: O(2^n) -> the recursion can branch into two choices at each of `n` levels.
// space complexity: O(n) -> the recursion stack and any cache can grow to one entry per input state.
```

For `fib(5)`:

```mermaid
flowchart TD
    A["fib(5)"] --> B["fib(4)"]
    A --> C["fib(3)"]
    B --> D["fib(3)"]
    B --> E["fib(2)"]
    C --> F["fib(2)"]
    C --> G["fib(1)"]
    D --> H["fib(2)"]
    D --> I["fib(1)"]
```

Notice that `fib(3)` and `fib(2)` are calculated repeatedly.

## Complexity

```text
Time:  O(2^n)
Space: O(n)
```

The stack depth is only `n`, even though the number of calls is exponential.

## Improved with memoization

```go
// Exact question: How can Improved with memoization be solved using memoization?
//
// Possible answer: Use top-down recursion and cache each state so repeated subproblems are returned immediately.
//
// Output format: Return an `int` value from `fib`; the function does not print the answer.
//
// Inline descriptions:
// - `n` is the int input used by this example.
// - `memo` is a map: each key is int, and each value is int.
//
// Boundary checks:
// - `n <= 1` handles the smallest valid state or recursive base case.
//
// Key variables:
// - `n` is the int input used by this example.
// - `memo` is a map: each key is int, and each value is int.
//
// Logic:
// 1. Create or use a map to associate each lookup key with its stored value.
// 2. Recursively reduce the current problem to smaller calls until a base condition is reached.
// 3. Reuse stored state instead of recomputing the same subproblem.
func fib(n int, memo map[int]int) int {
	if n <= 1 {
		return n
	}

	if result, exists := memo[n]; exists {
		return result
	}

	memo[n] = fib(n-1, memo) + fib(n-2, memo)
	return memo[n]
}

// time complexity: O(n) -> each of the `n` states is calculated once and then reused.
// space complexity: O(n) -> the recursion stack and any cache can grow to one entry per input state.
```

Usage:

```go
// Exact question: How can Improved with memoization be solved using memoization?
//
// Possible answer: Use the fragment with a map whose keys identify lookups and whose values hold the associated data.
//
// Output format: This is a partial Go fragment; its surrounding function determines the final returned value.
//
// Inline descriptions:
// - The comments in this preface describe how the important expressions and state changes are used.
//
// Boundary checks:
// - No explicit boundary branch appears in this fragment; its caller or surrounding example supplies valid inputs.
//
// Key variables:
// - `result` holds the answer computed for the current operation.
//
// Logic:
// 1. Create or use a map to associate each lookup key with its stored value.
result := fib(10, make(map[int]int))

// time complexity: O(n) -> the algorithm visits each of the `n` input elements or states once.
// space complexity: O(n) -> the auxiliary slice, map, table, queue, or returned collection can grow with `n`.
```

Complexity:

```text
Time:  O(n)
Space: O(n)
```

This is the connection between recursion and dynamic programming.

---

# 13. Problem 3: Sum of an array

Input:

```text
[4, 7, 2]
```

Recursive idea:

```text
sum([4, 7, 2])
= 4 + sum([7, 2])
= 4 + 7 + sum([2])
= 4 + 7 + 2 + sum([])
```

## Go code

```go
// Exact question: How does `sum` solve Go code?
//
// Possible answer: Use `sum` with a slice whose indexes identify positions and whose elements hold their values.
//
// Output format: Return an `int` value from `sum`; the function does not print the answer.
//
// Inline descriptions:
// - `nums` is a slice: the index identifies an element or state, and the stored item has type int.
// - `index` is the int input used by this example.
//
// Boundary checks:
// - `index == len(nums)` keeps indexes or pointers within the portion of the input still being processed.
//
// Key variables:
// - `nums` is a slice: the index identifies an element or state, and the stored item has type int.
// - `index` is the int input used by this example.
//
// Logic:
// 1. Create or use a slice so indexes identify positions and elements store their data or state.
// 2. Recursively reduce the current problem to smaller calls until a base condition is reached.
// 3. Return the value produced after the state updates are complete.
func sum(nums []int, index int) int {
	if index == len(nums) {
		return 0
	}

	return nums[index] + sum(nums, index+1)
}

// time complexity: O(n) -> the algorithm visits each of the `n` input elements or states once.
// space complexity: O(n) -> the recursion stack and any cache can grow to one entry per input state.
```

## Complexity

```text
Time:  O(n)
Space: O(n)
```

A common mistake is to create a new slice during every call:

```go
// Exact question: How does this Go example demonstrate Complexity?
//
// Possible answer: Use the fragment to execute the shown state update directly from top to bottom.
//
// Output format: This is a partial Go fragment; its surrounding function determines the final returned value.
//
// Inline descriptions:
// - The comments in this preface describe how the important expressions and state changes are used.
//
// Boundary checks:
// - No explicit boundary branch appears in this fragment; its caller or surrounding example supplies valid inputs.
//
// Key variables:
// - This fragment operates directly on the values named in each statement; it introduces no separate data structure.
//
// Logic:
// 1. Return the value produced after the state updates are complete.
return nums[0] + sum(nums[1:])

// time complexity: O(1) -> the snippet performs a fixed number of operations independent of input size.
// space complexity: O(1) -> only a fixed number of scalar variables or references is kept.
```

This may create additional slice-related overhead depending on what else is done. Passing an index is usually cleaner.

---

# 14. Problem 4: Binary tree traversal

A tree is recursive by definition.

A tree contains:

* A root
* A left subtree
* A right subtree

```mermaid
flowchart TD
    A[1] --> B[2]
    A --> C[3]
    B --> D[4]
    B --> E[5]
```

## Tree node

```go
// Exact question: How does this Go example demonstrate Tree node?
//
// Possible answer: Use the fragment to declare the types or values required by the surrounding implementation.
//
// Output format: This fragment demonstrates syntax or state updates and does not define a standalone output value.
//
// Inline descriptions:
// - The comments in this preface describe how the important expressions and state changes are used.
//
// Boundary checks:
// - No explicit boundary branch appears in this fragment; its caller or surrounding example supplies valid inputs.
//
// Key variables:
// - This fragment operates directly on the values named in each statement; it introduces no separate data structure.
//
// Logic:
// 1. Execute the statements from top to bottom to perform the demonstrated operation.
type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// time complexity: O(1) -> the snippet performs a fixed number of operations independent of input size.
// space complexity: O(1) -> only a fixed number of scalar variables or references is kept.
```

## Preorder traversal

Order:

```text
Root → Left → Right
```

```go
// Exact question: How does `preorder` solve Preorder traversal?
//
// Possible answer: Use `preorder` to format and write the demonstrated value to standard output.
//
// Output format: `preorder` has no return value; its observable result is the mutation or output performed in the function body.
//
// Inline descriptions:
// - `node` points to a TreeNode value that the function reads or updates.
//
// Boundary checks:
// - `node == nil` checks whether the referenced value exists before it is used.
//
// Key variables:
// - `node` points to a TreeNode value that the function reads or updates.
//
// Logic:
// 1. Recursively reduce the current problem to smaller calls until a base condition is reached.
// 2. Return the value produced after the state updates are complete.
func preorder(node *TreeNode) {
	if node == nil {
		return
	}

	fmt.Println(node.Val)
	preorder(node.Left)
	preorder(node.Right)
}

// time complexity: O(n) -> the algorithm visits each of the `n` input elements or states once.
// space complexity: O(h) -> the recursion stack grows with the height of the traversed tree or Trie.
```

Output:

```text
1 2 4 5 3
```

## Inorder traversal

Order:

```text
Left → Root → Right
```

```go
// Exact question: How does `inorder` solve Inorder traversal?
//
// Possible answer: Use `inorder` to format and write the demonstrated value to standard output.
//
// Output format: `inorder` has no return value; its observable result is the mutation or output performed in the function body.
//
// Inline descriptions:
// - `node` points to a TreeNode value that the function reads or updates.
//
// Boundary checks:
// - `node == nil` checks whether the referenced value exists before it is used.
//
// Key variables:
// - `node` points to a TreeNode value that the function reads or updates.
//
// Logic:
// 1. Recursively reduce the current problem to smaller calls until a base condition is reached.
// 2. Return the value produced after the state updates are complete.
func inorder(node *TreeNode) {
	if node == nil {
		return
	}

	inorder(node.Left)
	fmt.Println(node.Val)
	inorder(node.Right)
}

// time complexity: O(n) -> the algorithm visits each of the `n` input elements or states once.
// space complexity: O(h) -> the recursion stack grows with the height of the traversed tree or Trie.
```

Output:

```text
4 2 5 1 3
```

## Postorder traversal

Order:

```text
Left → Right → Root
```

```go
// Exact question: How does `postorder` solve Postorder traversal?
//
// Possible answer: Use `postorder` to format and write the demonstrated value to standard output.
//
// Output format: `postorder` has no return value; its observable result is the mutation or output performed in the function body.
//
// Inline descriptions:
// - `node` points to a TreeNode value that the function reads or updates.
//
// Boundary checks:
// - `node == nil` checks whether the referenced value exists before it is used.
//
// Key variables:
// - `node` points to a TreeNode value that the function reads or updates.
//
// Logic:
// 1. Recursively reduce the current problem to smaller calls until a base condition is reached.
// 2. Return the value produced after the state updates are complete.
func postorder(node *TreeNode) {
	if node == nil {
		return
	}

	postorder(node.Left)
	postorder(node.Right)
	fmt.Println(node.Val)
}

// time complexity: O(n) -> the algorithm visits each of the `n` input elements or states once.
// space complexity: O(h) -> the recursion stack grows with the height of the traversed tree or Trie.
```

Output:

```text
4 5 2 3 1
```

## Complexity

For all three traversals:

```text
Time:  O(n)
Space: O(h)
```

Where `h` is tree height.

* Balanced tree: `h = O(log n)`
* Skewed tree: `h = O(n)`

---

# 15. Problem 5: Maximum depth of a tree

## Recursive contract

```text
maxDepth(node)
```

means:

> Return the maximum depth of the tree rooted at `node`.

For an empty tree:

```text
depth = 0
```

For a non-empty tree:

```text
depth = 1 + maximum(left depth, right depth)
```

```go
// Exact question: How does `maxDepth`, `max` solve Recursive contract?
//
// Possible answer: Use `maxDepth`, `max` to compare the current values and return or update state when the condition matches.
//
// Output format: Return an `int` value from `maxDepth`; the function does not print the answer.
//
// Inline descriptions:
// - `node` points to a TreeNode value that the function reads or updates.
//
// Boundary checks:
// - `node == nil` checks whether the referenced value exists before it is used.
// - `a > b` decides whether the branch or loop should continue for the current input.
//
// Key variables:
// - `node` points to a TreeNode value that the function reads or updates.
// - `leftDepth` marks the current left boundary or left-side value.
// - `rightDepth` marks the current right boundary or right-side value.
//
// Logic:
// 1. Recursively reduce the current problem to smaller calls until a base condition is reached.
// 2. Return the value produced after the state updates are complete.
func maxDepth(node *TreeNode) int {
	if node == nil {
		return 0
	}

	leftDepth := maxDepth(node.Left)
	rightDepth := maxDepth(node.Right)

	return 1 + max(leftDepth, rightDepth)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// time complexity: O(n) -> the algorithm visits each of the `n` input elements or states once.
// space complexity: O(h) -> the recursion stack grows with the height of the traversed tree or Trie.
```

```mermaid
flowchart TD
    A["depth(1)"] --> B["depth(2)"]
    A --> C["depth(3)"]
    B --> D["depth(4)"]
    B --> E["depth(5)"]
    C --> F["depth(nil) = 0"]
    C --> G["depth(nil) = 0"]
```

Complexity:

```text
Time:  O(n)
Space: O(h)
```

---

# 16. Backtracking mental model

Suppose you are choosing clothes:

```text
Shirt:
- Blue
- Red

Pants:
- Black
- White
```

You try:

```text
Blue + Black
Blue + White
Red + Black
Red + White
```

The process is:

1. Choose blue.
2. Choose black.
3. Record the outfit.
4. Remove black.
5. Choose white.
6. Record the outfit.
7. Remove white.
8. Remove blue.
9. Choose red.
10. Repeat.

This is backtracking.

The standard pattern is:

```go
// Exact question: How can Backtracking mental model be solved using backtracking?
//
// Possible answer: Choose an option, recurse, then undo the choice before trying the next option.
//
// Output format: `backtrack` has no return value; its observable result is the mutation or output performed in the function body.
//
// Inline descriptions:
// - `path` is a slice: the index identifies an element or state, and the stored item has type Choice.
// - `choices` is a slice: the index identifies an element or state, and the stored item has type Choice.
//
// Boundary checks:
// - `solutionComplete(path)` decides whether the branch or loop should continue for the current input.
// - `!isValid(choice, path)` decides whether the branch or loop should continue for the current input.
//
// Key variables:
// - `path` is a slice: the index identifies an element or state, and the stored item has type Choice.
// - `choices` is a slice: the index identifies an element or state, and the stored item has type Choice.
//
// Logic:
// 1. Create or use a slice so indexes identify positions and elements store their data or state.
// 2. Iterate through the required elements or states in the order shown.
// 3. Recursively reduce the current problem to smaller calls until a base condition is reached.
func backtrack(path []Choice, choices []Choice) {
	if solutionComplete(path) {
		saveCopy(path)
		return
	}

	for _, choice := range choices {
		if !isValid(choice, path) {
			continue
		}

		// Choose
		path = append(path, choice)

		// Explore
		backtrack(path, choices)

		// Undo
		path = path[:len(path)-1]
	}
}

// time complexity: O(n) -> the algorithm visits each of the `n` input elements or states once.
// space complexity: O(n) -> the recursion stack and any cache can grow to one entry per input state.
```

Remember:

```text
Choose
Explore
Undo
```

---

# 17. Problem 6: Generate all subsets

Input:

```text
[1, 2]
```

Output:

```text
[]
[1]
[2]
[1, 2]
```

For every number, there are two choices:

```text
Include it
Do not include it
```

```mermaid
flowchart TD
    A["Start: []"] --> B["Exclude 1: []"]
    A --> C["Include 1: [1]"]
    B --> D["Exclude 2: []"]
    B --> E["Include 2: [2]"]
    C --> F["Exclude 2: [1]"]
    C --> G["Include 2: [1,2]"]
```

## Go implementation

```go
// Exact question: How does `subsets` solve Generate all subsets?
//
// Possible answer: Use `subsets` with a slice whose indexes identify positions and whose elements hold their values.
//
// Output format: Return the `[][]int` value from `subsets`; the function does not print the answer.
//
// Inline descriptions:
// - `nums` is a slice: the index identifies an element or state, and the stored item has type int.
//
// Boundary checks:
// - `index == len(nums)` keeps indexes or pointers within the portion of the input still being processed.
//
// Key variables:
// - `nums` is a slice: the index identifies an element or state, and the stored item has type int.
// - `result` is a two-dimensional table; indexes select a row and column, while cells hold computed values.
// - `path` is a slice; indexes identify positions and elements hold the corresponding values.
// - `backtrack` holds the intermediate value produced by `func(index int`.
// - `subset` holds the intermediate value produced by `append([]int(nil), path...`.
//
// Logic:
// 1. Create or use a slice so indexes identify positions and elements store their data or state.
// 2. Return the value produced after the state updates are complete.
func subsets(nums []int) [][]int {
	result := make([][]int, 0)
	path := make([]int, 0)

	var backtrack func(index int)
	backtrack = func(index int) {
		if index == len(nums) {
			subset := append([]int(nil), path...)
			result = append(result, subset)
			return
		}

		// Do not include nums[index]
		backtrack(index + 1)

		// Include nums[index]
		path = append(path, nums[index])
		backtrack(index + 1)

		// Undo
		path = path[:len(path)-1]
	}

	backtrack(0)
	return result
}

// time complexity: O(2^n) -> the recursion can branch into two choices at each of `n` levels.
// space complexity: O(n) -> the recursion stack and any cache can grow to one entry per input state.
```

## Why copy the path?

This is wrong:

```go
// Exact question: Why copy the path?
//
// Possible answer: Use the fragment to execute the shown state update directly from top to bottom.
//
// Output format: This fragment demonstrates syntax or state updates and does not define a standalone output value.
//
// Inline descriptions:
// - The comments in this preface describe how the important expressions and state changes are used.
//
// Boundary checks:
// - No explicit boundary branch appears in this fragment; its caller or surrounding example supplies valid inputs.
//
// Key variables:
// - This fragment operates directly on the values named in each statement; it introduces no separate data structure.
//
// Logic:
// 1. Execute the statements from top to bottom to perform the demonstrated operation.
result = append(result, path)

// time complexity: O(1) amortized -> the snippet performs a fixed number of operations independent of input size.
// space complexity: O(1) amortized -> only a fixed number of scalar variables or references is kept.
```

The underlying array may be reused and modified later.

Instead:

```go
// Exact question: Why copy the path?
//
// Possible answer: Use the fragment with a slice whose indexes identify positions and whose elements hold their values.
//
// Output format: This is a partial Go fragment; its surrounding function determines the final returned value.
//
// Inline descriptions:
// - The comments in this preface describe how the important expressions and state changes are used.
//
// Boundary checks:
// - No explicit boundary branch appears in this fragment; its caller or surrounding example supplies valid inputs.
//
// Key variables:
// - `subset` holds the intermediate value produced by `append([]int(nil), path...`.
//
// Logic:
// 1. Create or use a slice so indexes identify positions and elements store their data or state.
subset := append([]int(nil), path...)
result = append(result, subset)

// time complexity: O(n * 2^n) -> the recursion can branch into two choices at each of `n` levels.
// space complexity: O(n) -> the auxiliary slice, map, table, queue, or returned collection can grow with `n`.
```

## Complexity

For `n` numbers, there are `2^n` subsets.

```text
Time:  O(n × 2^n)
Space: O(n) recursion depth
Output: O(n × 2^n)
```

The additional `n` factor comes from copying each subset.

---

# 18. Problem 7: Permutations

Input:

```text
[1, 2, 3]
```

Possible outputs:

```text
[1, 2, 3]
[1, 3, 2]
[2, 1, 3]
[2, 3, 1]
[3, 1, 2]
[3, 2, 1]
```

At every position, choose any number not already used.

```mermaid
flowchart TD
    A["[]"] --> B["[1]"]
    A --> C["[2]"]
    A --> D["[3]"]
    B --> E["[1,2]"]
    B --> F["[1,3]"]
    E --> G["[1,2,3]"]
    F --> H["[1,3,2]"]
```

## Go implementation

```go
// Exact question: How does `permute` solve Permutations?
//
// Possible answer: Use `permute` to iterate through the input once, updating the running state for each element.
//
// Output format: Return the `[][]int` value from `permute`; the function does not print the answer.
//
// Inline descriptions:
// - `nums` is a slice: the index identifies an element or state, and the stored item has type int.
//
// Boundary checks:
// - `len(path) == len(nums)` decides whether the branch or loop should continue for the current input.
// - `used[i]` keeps indexes or pointers within the portion of the input still being processed.
//
// Key variables:
// - `nums` is a slice: the index identifies an element or state, and the stored item has type int.
// - `result` is a two-dimensional table; indexes select a row and column, while cells hold computed values.
// - `path` is a slice; indexes identify positions and elements hold the corresponding values.
// - `used` is a slice; indexes identify positions and elements hold the corresponding values.
// - `backtrack` holds the intermediate value produced by `func(`.
//
// Logic:
// 1. Create or use a slice so indexes identify positions and elements store their data or state.
// 2. Iterate through the required elements or states in the order shown.
// 3. Return the value produced after the state updates are complete.
func permute(nums []int) [][]int {
	result := make([][]int, 0)
	path := make([]int, 0, len(nums))
	used := make([]bool, len(nums))

	var backtrack func()
	backtrack = func() {
		if len(path) == len(nums) {
			permutation := append([]int(nil), path...)
			result = append(result, permutation)
			return
		}

		for i := 0; i < len(nums); i++ {
			if used[i] {
				continue
			}

			// Choose
			used[i] = true
			path = append(path, nums[i])

			// Explore
			backtrack()

			// Undo
			path = path[:len(path)-1]
			used[i] = false
		}
	}

	backtrack()
	return result
}

// time complexity: O(n * n!) -> the search may enumerate factorially many orderings, with work at each recursion level.
// space complexity: O(n) -> the recursion stack and any cache can grow to one entry per input state.
```

## Complexity

There are `n!` permutations.

```text
Time:  O(n × n!)
Space: O(n)
Output: O(n × n!)
```

---

# 19. Problem 8: Generate Parentheses

For `n = 3`, generate:

```text
((()))
(()())
(())()
()(())
()()()
```

At every step, choose:

```text
Add "("
Add ")"
```

But the partial string must remain valid.

Rules:

```text
Opening brackets used < n
Closing brackets used < opening brackets used
```

```mermaid
flowchart TD
    A["'' open=0 close=0"] --> B["'(' 1,0"]
    B --> C["'((' 2,0"]
    B --> D["'()' 1,1"]
    C --> E["'(((' 3,0"]
    C --> F["'(()' 2,1"]
    D --> G["'()(' 2,1"]
```

## Go implementation

```go
// Exact question: How does `generateParenthesis` solve Generate Parentheses?
//
// Possible answer: Use `generateParenthesis` with a slice whose indexes identify positions and whose elements hold their values.
//
// Output format: Return the `[]string` value from `generateParenthesis`; the function does not print the answer.
//
// Inline descriptions:
// - `n` is the int input used by this example.
//
// Boundary checks:
// - `len(current) == 2*n` decides whether the branch or loop should continue for the current input.
// - `open < n` decides whether the branch or loop should continue for the current input.
// - `close < open` decides whether the branch or loop should continue for the current input.
//
// Key variables:
// - `n` is the int input used by this example.
// - `result` is a slice whose elements hold the ordered values produced or awaiting processing.
// - `backtrack` holds the intermediate value produced by `func(current []byte, open, close int`.
//
// Logic:
// 1. Create or use a slice so indexes identify positions and elements store their data or state.
// 2. Return the value produced after the state updates are complete.
func generateParenthesis(n int) []string {
	result := make([]string, 0)

	var backtrack func(current []byte, open, close int)
	backtrack = func(current []byte, open, close int) {
		if len(current) == 2*n {
			result = append(result, string(current))
			return
		}

		if open < n {
			current = append(current, '(')
			backtrack(current, open+1, close)
			current = current[:len(current)-1]
		}

		if close < open {
			current = append(current, ')')
			backtrack(current, open, close+1)
			current = current[:len(current)-1]
		}
	}

	backtrack(make([]byte, 0, 2*n), 0, 0)
	return result
}

// time complexity: O(2^n) -> the recursion can branch into two choices at each of `n` levels.
// space complexity: O(n) -> the recursion stack and any cache can grow to one entry per input state.
```

## Interview insight

Do not generate every possible sequence and validate afterward.

Prune invalid branches before exploring them.

That is one of the main goals of backtracking.

---

# 20. Problem 9: Combination Sum

Input:

```text
candidates = [2, 3, 6, 7]
target = 7
```

Output:

```text
[2, 2, 3]
[7]
```

At each candidate, you can:

```text
Use it again
Move to the next candidate
```

## Go implementation

```go
// Exact question: How does `combinationSum` solve Combination Sum?
//
// Possible answer: Use `combinationSum` to iterate through the input once, updating the running state for each element.
//
// Output format: Return the `[][]int` value from `combinationSum`; the function does not print the answer.
//
// Inline descriptions:
// - `candidates` is a slice: the index identifies an element or state, and the stored item has type int.
// - `target` is the int input used by this example.
//
// Boundary checks:
// - `remaining == 0` handles the smallest valid state or recursive base case.
// - `remaining < 0` rejects or terminates a state that has moved below the valid range.
//
// Key variables:
// - `candidates` is a slice: the index identifies an element or state, and the stored item has type int.
// - `target` is the int input used by this example.
// - `result` is a two-dimensional table; indexes select a row and column, while cells hold computed values.
// - `path` is a slice; indexes identify positions and elements hold the corresponding values.
// - `backtrack` holds the intermediate value produced by `func(start, remaining int`.
//
// Logic:
// 1. Create or use a slice so indexes identify positions and elements store their data or state.
// 2. Iterate through the required elements or states in the order shown.
// 3. Return the value produced after the state updates are complete.
func combinationSum(candidates []int, target int) [][]int {
	result := make([][]int, 0)
	path := make([]int, 0)

	var backtrack func(start, remaining int)
	backtrack = func(start, remaining int) {
		if remaining == 0 {
			combination := append([]int(nil), path...)
			result = append(result, combination)
			return
		}

		if remaining < 0 {
			return
		}

		for i := start; i < len(candidates); i++ {
			value := candidates[i]

			path = append(path, value)
			backtrack(i, remaining-value)
			path = path[:len(path)-1]
		}
	}

	backtrack(0, target)
	return result
}

// time complexity: O(n) -> the algorithm visits each of the `n` input elements or states once.
// space complexity: O(n) -> the recursion stack and any cache can grow to one entry per input state.
```

Why call:

```go
// Exact question: How does this Go example demonstrate Combination Sum?
//
// Possible answer: Use the fragment to execute the shown state update directly from top to bottom.
//
// Output format: This fragment demonstrates syntax or state updates and does not define a standalone output value.
//
// Inline descriptions:
// - The comments in this preface describe how the important expressions and state changes are used.
//
// Boundary checks:
// - No explicit boundary branch appears in this fragment; its caller or surrounding example supplies valid inputs.
//
// Key variables:
// - This fragment operates directly on the values named in each statement; it introduces no separate data structure.
//
// Logic:
// 1. Execute the statements from top to bottom to perform the demonstrated operation.
backtrack(i, remaining-value)

// time complexity: O(1) -> the snippet performs a fixed number of operations independent of input size.
// space complexity: O(1) -> only a fixed number of scalar variables or references is kept.
```

instead of:

```go
// Exact question: How does this Go example demonstrate Combination Sum?
//
// Possible answer: Use the fragment to execute the shown state update directly from top to bottom.
//
// Output format: This fragment demonstrates syntax or state updates and does not define a standalone output value.
//
// Inline descriptions:
// - The comments in this preface describe how the important expressions and state changes are used.
//
// Boundary checks:
// - No explicit boundary branch appears in this fragment; its caller or surrounding example supplies valid inputs.
//
// Key variables:
// - This fragment operates directly on the values named in each statement; it introduces no separate data structure.
//
// Logic:
// 1. Execute the statements from top to bottom to perform the demonstrated operation.
backtrack(i+1, remaining-value)

// time complexity: O(1) -> the snippet performs a fixed number of operations independent of input size.
// space complexity: O(1) -> only a fixed number of scalar variables or references is kept.
```

Because the same number may be reused.

Why use `start`?

To avoid generating duplicate orderings such as:

```text
[2, 2, 3]
[2, 3, 2]
[3, 2, 2]
```

---

# 21. Recursion complexity

Recursive complexity has two parts:

```text
Time complexity:
How many calls are made?
×
How much work happens per call?

Space complexity:
Maximum recursion depth
+
Any additional data structures
```

## Linear recursion

```go
// Exact question: How does `solve` solve Linear recursion?
//
// Possible answer: Use `solve` to compare the current values and return or update state when the condition matches.
//
// Output format: `solve` has no return value; its observable result is the mutation or output performed in the function body.
//
// Inline descriptions:
// - `n` is the int input used by this example.
//
// Boundary checks:
// - `n == 0` handles the smallest valid state or recursive base case.
//
// Key variables:
// - `n` is the int input used by this example.
//
// Logic:
// 1. Recursively reduce the current problem to smaller calls until a base condition is reached.
// 2. Return the value produced after the state updates are complete.
func solve(n int) {
	if n == 0 {
		return
	}
	solve(n - 1)
}

// time complexity: O(n) -> the algorithm visits each of the `n` input elements or states once.
// space complexity: O(n) -> the recursion stack and any cache can grow to one entry per input state.
```

Recurrence:

```text
T(n) = T(n - 1) + O(1)
```

Result:

```text
Time: O(n)
Space: O(n)
```

## Binary tree recursion

```go
// Exact question: How does `solve` solve Binary tree recursion?
//
// Possible answer: Use `solve` to compare the current values and return or update state when the condition matches.
//
// Output format: `solve` has no return value; its observable result is the mutation or output performed in the function body.
//
// Inline descriptions:
// - `n` is the int input used by this example.
//
// Boundary checks:
// - `n == 0` handles the smallest valid state or recursive base case.
//
// Key variables:
// - `n` is the int input used by this example.
//
// Logic:
// 1. Recursively reduce the current problem to smaller calls until a base condition is reached.
// 2. Return the value produced after the state updates are complete.
func solve(n int) {
	if n == 0 {
		return
	}

	solve(n - 1)
	solve(n - 1)
}

// time complexity: O(2^n) -> the recursion can branch into two choices at each of `n` levels.
// space complexity: O(n) -> the recursion stack and any cache can grow to one entry per input state.
```

Recurrence:

```text
T(n) = 2T(n - 1) + O(1)
```

Result:

```text
Time: O(2^n)
Space: O(n)
```

A critical distinction:

> Space depends on the longest active path, not the total number of calls.

---

# 22. Understanding recursion trees

Consider:

```go
// Exact question: How does `f` solve Understanding recursion trees?
//
// Possible answer: Use `f` to compare the current values and return or update state when the condition matches.
//
// Output format: `f` has no return value; its observable result is the mutation or output performed in the function body.
//
// Inline descriptions:
// - `n` is the int input used by this example.
//
// Boundary checks:
// - `n == 0` handles the smallest valid state or recursive base case.
//
// Key variables:
// - `n` is the int input used by this example.
//
// Logic:
// 1. Recursively reduce the current problem to smaller calls until a base condition is reached.
// 2. Return the value produced after the state updates are complete.
func f(n int) {
	if n == 0 {
		return
	}

	f(n - 1)
	f(n - 1)
}

// time complexity: O(n) -> the algorithm visits each of the `n` input elements or states once.
// space complexity: O(h) -> the recursion stack grows with the height of the traversed tree or Trie.
```

At each level, the number of calls doubles.

```text
Level 0: 1 call
Level 1: 2 calls
Level 2: 4 calls
Level 3: 8 calls
```

Total:

```text
1 + 2 + 4 + ... + 2^n
≈ 2^(n+1)
= O(2^n)
```

```mermaid
flowchart TD
    A["f(3)"] --> B["f(2)"]
    A --> C["f(2)"]
    B --> D["f(1)"]
    B --> E["f(1)"]
    C --> F["f(1)"]
    C --> G["f(1)"]
```

---

# 23. Tail recursion

A tail-recursive function performs the recursive call as its final operation.

```go
// Exact question: How does `factorialTail` solve Tail recursion?
//
// Possible answer: Use `factorialTail` to compare the current values and return or update state when the condition matches.
//
// Output format: Return an `int` value from `factorialTail`; the function does not print the answer.
//
// Inline descriptions:
// - `n` is the int input used by this example.
// - `accumulator` is the int input used by this example.
//
// Boundary checks:
// - `n <= 1` handles the smallest valid state or recursive base case.
//
// Key variables:
// - `n` is the int input used by this example.
// - `accumulator` is the int input used by this example.
//
// Logic:
// 1. Recursively reduce the current problem to smaller calls until a base condition is reached.
// 2. Return the value produced after the state updates are complete.
func factorialTail(n, accumulator int) int {
	if n <= 1 {
		return accumulator
	}

	return factorialTail(n-1, accumulator*n)
}

// time complexity: O(n) -> the algorithm visits each of the `n` input elements or states once.
// space complexity: O(n) -> the recursion stack and any cache can grow to one entry per input state.
```

Usage:

```go
// Exact question: How does this Go example demonstrate Tail recursion?
//
// Possible answer: Use the fragment to execute the shown state update directly from top to bottom.
//
// Output format: This is a partial Go fragment; its surrounding function determines the final returned value.
//
// Inline descriptions:
// - The comments in this preface describe how the important expressions and state changes are used.
//
// Boundary checks:
// - No explicit boundary branch appears in this fragment; its caller or surrounding example supplies valid inputs.
//
// Key variables:
// - `result` holds the answer computed for the current operation.
//
// Logic:
// 1. Execute the statements from top to bottom to perform the demonstrated operation.
result := factorialTail(5, 1)

// time complexity: O(1) -> the snippet performs a fixed number of operations independent of input size.
// space complexity: O(1) -> only a fixed number of scalar variables or references is kept.
```

Some languages optimize tail recursion and reuse the same stack frame.

Go does not guarantee tail-call optimization. Therefore, deep tail recursion can still overflow the stack or require stack growth.

For simple linear operations in Go, loops are generally preferable.

---

# 24. Common recursion mistakes

## Mistake 1: Missing the base case

```go
// Exact question: How does `countdown` solve Mistake 1: Missing the base case?
//
// Possible answer: Use `countdown` to format and write the demonstrated value to standard output.
//
// Output format: `countdown` has no return value; its observable result is the mutation or output performed in the function body.
//
// Inline descriptions:
// - `n` is the int input used by this example.
//
// Boundary checks:
// - No explicit boundary branch appears in this fragment; its caller or surrounding example supplies valid inputs.
//
// Key variables:
// - `n` is the int input used by this example.
//
// Logic:
// 1. Recursively reduce the current problem to smaller calls until a base condition is reached.
func countdown(n int) {
	fmt.Println(n)
	countdown(n - 1)
}

// time complexity: O(n) -> the algorithm visits each of the `n` input elements or states once.
// space complexity: O(n) -> the recursion stack and any cache can grow to one entry per input state.
```

This never stops.

---

## Mistake 2: Input does not move toward the base case

```go
// Exact question: How does `countdown` solve Mistake 2: Input does not move toward the base case?
//
// Possible answer: Use `countdown` to compare the current values and return or update state when the condition matches.
//
// Output format: `countdown` has no return value; its observable result is the mutation or output performed in the function body.
//
// Inline descriptions:
// - `n` is the int input used by this example.
//
// Boundary checks:
// - `n == 0` handles the smallest valid state or recursive base case.
//
// Key variables:
// - `n` is the int input used by this example.
//
// Logic:
// 1. Recursively reduce the current problem to smaller calls until a base condition is reached.
// 2. Return the value produced after the state updates are complete.
func countdown(n int) {
	if n == 0 {
		return
	}

	countdown(n)
}

// time complexity: O(n) -> the algorithm visits each of the `n` input elements or states once.
// space complexity: O(n) -> the recursion stack and any cache can grow to one entry per input state.
```

`n` never changes.

---

## Mistake 3: Wrong base-case result

For factorial:

```go
// Exact question: How does this Go example demonstrate Mistake 3: Wrong base-case result?
//
// Possible answer: Use the fragment to compare the current values and return or update state when the condition matches.
//
// Output format: This is a partial Go fragment; its surrounding function determines the final returned value.
//
// Inline descriptions:
// - The comments in this preface describe how the important expressions and state changes are used.
//
// Boundary checks:
// - `n == 0` handles the smallest valid state or recursive base case.
//
// Key variables:
// - This fragment operates directly on the values named in each statement; it introduces no separate data structure.
//
// Logic:
// 1. Return the value produced after the state updates are complete.
if n == 0 {
	return 0
}

// time complexity: O(1) -> the snippet performs a fixed number of operations independent of input size.
// space complexity: O(1) -> only a fixed number of scalar variables or references is kept.
```

This makes every factorial result zero.

Correct:

```go
// Exact question: How does this Go example demonstrate Mistake 3: Wrong base-case result?
//
// Possible answer: Use the fragment to compare the current values and return or update state when the condition matches.
//
// Output format: This is a partial Go fragment; its surrounding function determines the final returned value.
//
// Inline descriptions:
// - The comments in this preface describe how the important expressions and state changes are used.
//
// Boundary checks:
// - `n == 0` handles the smallest valid state or recursive base case.
//
// Key variables:
// - This fragment operates directly on the values named in each statement; it introduces no separate data structure.
//
// Logic:
// 1. Return the value produced after the state updates are complete.
if n == 0 {
	return 1
}

// time complexity: O(1) -> the snippet performs a fixed number of operations independent of input size.
// space complexity: O(1) -> only a fixed number of scalar variables or references is kept.
```

Because:

```text
n × 1 = n
```

---

## Mistake 4: Forgetting to use the recursive result

Wrong:

```go
// Exact question: How does `factorial` solve Mistake 4: Forgetting to use the recursive result?
//
// Possible answer: Use `factorial` to compare the current values and return or update state when the condition matches.
//
// Output format: Return an `int` value from `factorial`; the function does not print the answer.
//
// Inline descriptions:
// - `n` is the int input used by this example.
//
// Boundary checks:
// - `n <= 1` handles the smallest valid state or recursive base case.
//
// Key variables:
// - `n` is the int input used by this example.
//
// Logic:
// 1. Recursively reduce the current problem to smaller calls until a base condition is reached.
// 2. Return the value produced after the state updates are complete.
func factorial(n int) int {
	if n <= 1 {
		return 1
	}

	factorial(n - 1)
	return n
}

// time complexity: O(n) -> the algorithm visits each of the `n` input elements or states once.
// space complexity: O(n) -> the recursion stack and any cache can grow to one entry per input state.
```

Correct:

```go
// Exact question: How does this Go example demonstrate Mistake 4: Forgetting to use the recursive result?
//
// Possible answer: Use the fragment to execute the shown state update directly from top to bottom.
//
// Output format: This is a partial Go fragment; its surrounding function determines the final returned value.
//
// Inline descriptions:
// - The comments in this preface describe how the important expressions and state changes are used.
//
// Boundary checks:
// - No explicit boundary branch appears in this fragment; its caller or surrounding example supplies valid inputs.
//
// Key variables:
// - This fragment operates directly on the values named in each statement; it introduces no separate data structure.
//
// Logic:
// 1. Return the value produced after the state updates are complete.
return n * factorial(n-1)

// time complexity: O(1) -> the snippet performs a fixed number of operations independent of input size.
// space complexity: O(1) -> only a fixed number of scalar variables or references is kept.
```

---

## Mistake 5: Modifying shared state without undoing

Wrong backtracking:

```go
// Exact question: How does this Go example demonstrate Mistake 5: Modifying shared state without undoing?
//
// Possible answer: Use the fragment to execute the shown state update directly from top to bottom.
//
// Output format: This fragment demonstrates syntax or state updates and does not define a standalone output value.
//
// Inline descriptions:
// - The comments in this preface describe how the important expressions and state changes are used.
//
// Boundary checks:
// - No explicit boundary branch appears in this fragment; its caller or surrounding example supplies valid inputs.
//
// Key variables:
// - This fragment operates directly on the values named in each statement; it introduces no separate data structure.
//
// Logic:
// 1. Execute the statements from top to bottom to perform the demonstrated operation.
path = append(path, choice)
backtrack()

// time complexity: O(1) amortized -> the snippet performs a fixed number of operations independent of input size.
// space complexity: O(1) amortized -> only a fixed number of scalar variables or references is kept.
```

The choice remains in `path`.

Correct:

```go
// Exact question: How does this Go example demonstrate Mistake 5: Modifying shared state without undoing?
//
// Possible answer: Use the fragment to execute the shown state update directly from top to bottom.
//
// Output format: This fragment demonstrates syntax or state updates and does not define a standalone output value.
//
// Inline descriptions:
// - The comments in this preface describe how the important expressions and state changes are used.
//
// Boundary checks:
// - No explicit boundary branch appears in this fragment; its caller or surrounding example supplies valid inputs.
//
// Key variables:
// - This fragment operates directly on the values named in each statement; it introduces no separate data structure.
//
// Logic:
// 1. Execute the statements from top to bottom to perform the demonstrated operation.
path = append(path, choice)
backtrack()
path = path[:len(path)-1]

// time complexity: O(1) amortized -> the snippet performs a fixed number of operations independent of input size.
// space complexity: O(1) amortized -> only a fixed number of scalar variables or references is kept.
```

---

## Mistake 6: Saving the same mutable path

Wrong:

```go
// Exact question: How does this Go example demonstrate Mistake 6: Saving the same mutable path?
//
// Possible answer: Use the fragment to execute the shown state update directly from top to bottom.
//
// Output format: This fragment demonstrates syntax or state updates and does not define a standalone output value.
//
// Inline descriptions:
// - The comments in this preface describe how the important expressions and state changes are used.
//
// Boundary checks:
// - No explicit boundary branch appears in this fragment; its caller or surrounding example supplies valid inputs.
//
// Key variables:
// - This fragment operates directly on the values named in each statement; it introduces no separate data structure.
//
// Logic:
// 1. Execute the statements from top to bottom to perform the demonstrated operation.
result = append(result, path)

// time complexity: O(1) amortized -> the snippet performs a fixed number of operations independent of input size.
// space complexity: O(1) amortized -> only a fixed number of scalar variables or references is kept.
```

Correct:

```go
// Exact question: How does this Go example demonstrate Mistake 6: Saving the same mutable path?
//
// Possible answer: Use the fragment with a slice whose indexes identify positions and whose elements hold their values.
//
// Output format: This is a partial Go fragment; its surrounding function determines the final returned value.
//
// Inline descriptions:
// - The comments in this preface describe how the important expressions and state changes are used.
//
// Boundary checks:
// - No explicit boundary branch appears in this fragment; its caller or surrounding example supplies valid inputs.
//
// Key variables:
// - `copyOfPath` holds the intermediate value produced by `append([]int(nil), path...`.
//
// Logic:
// 1. Create or use a slice so indexes identify positions and elements store their data or state.
copyOfPath := append([]int(nil), path...)
result = append(result, copyOfPath)

// time complexity: O(1) amortized -> the snippet performs a fixed number of operations independent of input size.
// space complexity: O(1) amortized -> only a fixed number of scalar variables or references is kept.
```

---

## Mistake 7: Recomputing the same state

Naive Fibonacci recalculates the same values many times.

Use:

* Memoization
* Tabulation
* Another mathematical transformation

---

## Mistake 8: Ignoring stack depth

A recursive function with one million levels may be dangerous.

Consider:

* A loop
* An explicit stack
* BFS
* Iterative DFS

---

# 25. How to debug recursion

Add logging at function entry and exit.

```go
// Exact question: How to debug recursion?
//
// Possible answer: Use `factorial` to format and write the demonstrated value to standard output.
//
// Output format: Return an `int` value from `factorial`; the function does not print the answer.
//
// Inline descriptions:
// - `n` is the int input used by this example.
// - `depth` is the int input used by this example.
//
// Boundary checks:
// - `n <= 1` handles the smallest valid state or recursive base case.
//
// Key variables:
// - `n` is the int input used by this example.
// - `depth` is the int input used by this example.
// - `indent` holds the intermediate value produced by `strings.Repeat(" ", depth`.
// - `result` holds the answer computed for the current operation.
//
// Logic:
// 1. Recursively reduce the current problem to smaller calls until a base condition is reached.
// 2. Return the value produced after the state updates are complete.
func factorial(n int, depth int) int {
	indent := strings.Repeat("  ", depth)
	fmt.Printf("%senter factorial(%d)\n", indent, n)

	if n <= 1 {
		fmt.Printf("%sreturn 1\n", indent)
		return 1
	}

	result := n * factorial(n-1, depth+1)
	fmt.Printf("%sreturn %d\n", indent, result)
	return result
}

// time complexity: O(n) -> the algorithm visits each of the `n` input elements or states once.
// space complexity: O(h) -> the recursion stack grows with the height of the traversed tree or Trie.
```

Output:

```text
enter factorial(4)
  enter factorial(3)
    enter factorial(2)
      enter factorial(1)
      return 1
    return 2
  return 6
return 24
```

This makes the call stack visible.

---

# 26. How to recognize a recursion problem

Look for these signals:

* The problem contains nested structures.
* The input can be divided into smaller versions of itself.
* You need all combinations or permutations.
* You need to explore choices.
* The problem mentions trees or graphs.
* The problem asks for paths.
* The problem asks to generate all valid arrangements.
* The problem can be expressed using previous states.

Typical phrases:

```text
Generate all...
Find every...
Explore every path...
All possible combinations...
For each child...
For each neighboring node...
Left subtree and right subtree...
```

---

# 27. Recursion patterns to memorize

## Pattern 1: Process a sequence

```go
// Exact question: How does `process` solve Pattern 1: Process a sequence?
//
// Possible answer: Use `process` with a slice whose indexes identify positions and whose elements hold their values.
//
// Output format: `process` has no return value; its observable result is the mutation or output performed in the function body.
//
// Inline descriptions:
// - `nums` is a slice: the index identifies an element or state, and the stored item has type int.
// - `index` is the int input used by this example.
//
// Boundary checks:
// - `index == len(nums)` keeps indexes or pointers within the portion of the input still being processed.
//
// Key variables:
// - `nums` is a slice: the index identifies an element or state, and the stored item has type int.
// - `index` is the int input used by this example.
//
// Logic:
// 1. Create or use a slice so indexes identify positions and elements store their data or state.
// 2. Recursively reduce the current problem to smaller calls until a base condition is reached.
// 3. Return the value produced after the state updates are complete.
func process(nums []int, index int) {
	if index == len(nums) {
		return
	}

	// Process nums[index]
	process(nums, index+1)
}

// time complexity: O(n) -> the algorithm visits each of the `n` input elements or states once.
// space complexity: O(n) -> the recursion stack and any cache can grow to one entry per input state.
```

## Pattern 2: Return an answer

```go
// Exact question: How does `solve` solve Pattern 2: Return an answer?
//
// Possible answer: Use `solve` to compare the current values and return or update state when the condition matches.
//
// Output format: Return an `int` value from `solve`; the function does not print the answer.
//
// Inline descriptions:
// - `n` is the int input used by this example.
//
// Boundary checks:
// - `n == 0` handles the smallest valid state or recursive base case.
//
// Key variables:
// - `n` is the int input used by this example.
// - `smallerAnswer` tracks the best or final answer found so far.
//
// Logic:
// 1. Recursively reduce the current problem to smaller calls until a base condition is reached.
// 2. Return the value produced after the state updates are complete.
func solve(n int) int {
	if n == 0 {
		return baseValue
	}

	smallerAnswer := solve(n - 1)
	return combine(n, smallerAnswer)
}

// time complexity: O(n) -> the algorithm visits each of the `n` input elements or states once.
// space complexity: O(n) -> the recursion stack and any cache can grow to one entry per input state.
```

## Pattern 3: Binary choice

```go
// Exact question: How does `choose` solve Pattern 3: Binary choice?
//
// Possible answer: Use `choose` to compare the current values and return or update state when the condition matches.
//
// Output format: `choose` has no return value; its observable result is the mutation or output performed in the function body.
//
// Inline descriptions:
// - `index` is the int input used by this example.
//
// Boundary checks:
// - `index == n` keeps indexes or pointers within the portion of the input still being processed.
//
// Key variables:
// - `index` is the int input used by this example.
//
// Logic:
// 1. Recursively reduce the current problem to smaller calls until a base condition is reached.
// 2. Return the value produced after the state updates are complete.
func choose(index int) {
	if index == n {
		saveAnswer()
		return
	}

	// Exclude
	choose(index + 1)

	// Include
	path = append(path, values[index])
	choose(index + 1)
	path = path[:len(path)-1]
}

// time complexity: O(2^n) -> the recursion can branch into two choices at each of `n` levels.
// space complexity: O(n) -> the recursion stack and any cache can grow to one entry per input state.
```

## Pattern 4: Loop over choices

```go
// Exact question: How does `backtrack` solve Pattern 4: Loop over choices?
//
// Possible answer: Use `backtrack` to iterate through the input once, updating the running state for each element.
//
// Output format: `backtrack` has no return value; its observable result is the mutation or output performed in the function body.
//
// Inline descriptions:
// - `start` is the int input used by this example.
//
// Boundary checks:
// - `complete()` decides whether the branch or loop should continue for the current input.
//
// Key variables:
// - `start` is the int input used by this example.
//
// Logic:
// 1. Iterate through the required elements or states in the order shown.
// 2. Recursively reduce the current problem to smaller calls until a base condition is reached.
// 3. Return the value produced after the state updates are complete.
func backtrack(start int) {
	if complete() {
		saveAnswer()
		return
	}

	for i := start; i < len(choices); i++ {
		makeChoice(i)
		backtrack(nextStart(i))
		undoChoice(i)
	}
}

// time complexity: O(n) -> the algorithm visits each of the `n` input elements or states once.
// space complexity: O(n) -> the recursion stack and any cache can grow to one entry per input state.
```

## Pattern 5: Tree recursion

```go
// Exact question: How does `solve` solve Pattern 5: Tree recursion?
//
// Possible answer: Use `solve` to compare the current values and return or update state when the condition matches.
//
// Output format: Return a `Result` value from `solve`; the function does not print the answer.
//
// Inline descriptions:
// - `node` points to a TreeNode value that the function reads or updates.
//
// Boundary checks:
// - `node == nil` checks whether the referenced value exists before it is used.
//
// Key variables:
// - `node` points to a TreeNode value that the function reads or updates.
// - `left` marks the current left boundary or left-side value.
// - `right` marks the current right boundary or right-side value.
//
// Logic:
// 1. Recursively reduce the current problem to smaller calls until a base condition is reached.
// 2. Return the value produced after the state updates are complete.
func solve(node *TreeNode) Result {
	if node == nil {
		return baseResult
	}

	left := solve(node.Left)
	right := solve(node.Right)

	return combine(node, left, right)
}

// time complexity: O(n) -> the algorithm visits each of the `n` input elements or states once.
// space complexity: O(h) -> the recursion stack grows with the height of the traversed tree or Trie.
```

## Pattern 6: Graph DFS

```go
// Exact question: How can Pattern 6: Graph DFS be solved using depth-first search?
//
// Possible answer: Follow one path at a time with recursion or a stack, marking nodes to avoid revisits.
//
// Output format: `dfs` has no return value; its observable result is the mutation or output performed in the function body.
//
// Inline descriptions:
// - `node` is the int input used by this example.
// - `graph` is a two-dimensional slice: row and column indexes identify positions, and each cell stores a value.
// - `visited` is a slice: the index identifies an element or state, and the stored item has type bool.
//
// Boundary checks:
// - `visited[node]` decides whether the branch or loop should continue for the current input.
//
// Key variables:
// - `node` is the int input used by this example.
// - `graph` is a two-dimensional slice: row and column indexes identify positions, and each cell stores a value.
//
// Logic:
// 1. Create or use a slice so indexes identify positions and elements store their data or state.
// 2. Iterate through the required elements or states in the order shown.
// 3. Recursively reduce the current problem to smaller calls until a base condition is reached.
func dfs(node int, graph [][]int, visited []bool) {
	if visited[node] {
		return
	}

	visited[node] = true

	for _, neighbor := range graph[node] {
		dfs(neighbor, graph, visited)
	}
}

// time complexity: O(V + E) -> each reachable vertex is processed once and each edge is examined once.
// space complexity: O(V) -> the visited state, queue, stack, or result can hold one entry per vertex.
```

---

# 28. Most important interview problems

## Foundation

1. Factorial
2. Fibonacci
3. Sum of an array
4. Reverse a string
5. Binary search recursively
6. Power of a number

## Trees and DFS

1. Maximum Depth of Binary Tree
2. Same Tree
3. Invert Binary Tree
4. Path Sum
5. Diameter of Binary Tree
6. Lowest Common Ancestor
7. Validate Binary Search Tree
8. Tree traversals

## Backtracking

1. Subsets
2. Permutations
3. Combination Sum
4. Generate Parentheses
5. Letter Combinations of a Phone Number
6. Palindrome Partitioning
7. Word Search
8. N-Queens

## Recursion with dynamic programming

1. Climbing Stairs
2. House Robber
3. Coin Change
4. Longest Common Subsequence
5. Edit Distance
6. Partition Equal Subset Sum

---

# 29. Mock interview questions

## Question 1: What is recursion?

A strong answer:

> Recursion is a technique where a function solves a problem by calling itself on a smaller subproblem. It requires a base case to stop and a recursive case that makes progress toward the base case. Each active call occupies a stack frame, so recursion depth affects space complexity.

---

## Question 2: What happens internally during recursion?

A strong answer:

> Every recursive call creates a stack frame containing its parameters, local variables, and return address. Frames are added while recursion goes deeper and removed when calls return. This return process is called stack unwinding.

---

## Question 3: Why is a base case necessary?

A strong answer:

> The base case defines the smallest directly solvable problem and prevents infinite recursive calls. Every recursive path must eventually reach it.

---

## Question 4: What causes stack overflow?

A strong answer:

> Stack overflow happens when recursion becomes too deep or never reaches its base case, causing too many active stack frames to accumulate.

---

## Question 5: What is the space complexity of factorial recursion?

```text
O(n)
```

Although each call performs constant work, `n` calls are simultaneously active before unwinding.

---

## Question 6: Why is naive Fibonacci slow?

A strong answer:

> It repeatedly calculates the same subproblems. Its recursion tree has exponentially many calls, producing approximately `O(2^n)` time. Memoization reduces this to `O(n)` by caching each state.

---

## Question 7: What is the difference between recursion and backtracking?

A strong answer:

> Recursion is a mechanism where a function calls itself. Backtracking is an algorithmic technique that uses recursion to explore choices, undo a choice, and explore alternatives.

---

## Question 8: What is the difference between DFS and backtracking?

A strong answer:

> DFS describes an exploration order: go as deep as possible before returning. Backtracking usually uses DFS but additionally maintains a candidate solution, makes choices, validates or prunes them, and undoes choices.

---

## Question 9: Why must paths be copied in backtracking?

A strong answer:

> The same mutable path is reused while exploring multiple branches. Saving its reference would allow later modifications to change previously stored results. A copy captures the path’s current state.

---

## Question 10: When should recursion be avoided?

A strong answer:

> Avoid it when recursion depth may be extremely large, when an iterative solution is substantially simpler, or when stack overhead is a concern. An explicit stack can preserve DFS behavior without relying on the language call stack.

---

# 30. Senior-level interview questions

## 1. What determines recursive space complexity?

Not the total number of calls.

It is primarily the maximum number of simultaneously active calls.

For a balanced tree:

```text
Depth = O(log n)
```

For a skewed tree:

```text
Depth = O(n)
```

---

## 2. How do you identify overlapping subproblems?

Check whether the recursion tree reaches the same state multiple times.

For example:

```text
fib(5)
```

calculates `fib(3)` more than once.

When repeated calls depend on the same parameters, memoization may help.

---

## 3. When does memoization not help?

Memoization does not significantly help when:

* Each recursive state is unique.
* The state space is too large.
* Cached results are unlikely to be reused.
* State hashing or storage costs exceed the benefit.

For example, generating every permutation requires producing `n!` unique outputs. Memoization cannot remove the output requirement.

---

## 4. How do you prevent cycles in recursive graph traversal?

Maintain a visited set.

```go
// Exact question: How do you prevent cycles in recursive graph traversal?
//
// Possible answer: Use `dfs` to iterate through the input once, updating the running state for each element.
//
// Output format: `dfs` has no return value; its observable result is the mutation or output performed in the function body.
//
// Inline descriptions:
// - `node` is the int input used by this example.
// - `graph` is a two-dimensional slice: row and column indexes identify positions, and each cell stores a value.
// - `visited` is a slice: the index identifies an element or state, and the stored item has type bool.
//
// Boundary checks:
// - `visited[node]` decides whether the branch or loop should continue for the current input.
//
// Key variables:
// - `node` is the int input used by this example.
// - `graph` is a two-dimensional slice: row and column indexes identify positions, and each cell stores a value.
//
// Logic:
// 1. Create or use a slice so indexes identify positions and elements store their data or state.
// 2. Iterate through the required elements or states in the order shown.
// 3. Recursively reduce the current problem to smaller calls until a base condition is reached.
func dfs(node int, graph [][]int, visited []bool) {
	if visited[node] {
		return
	}

	visited[node] = true

	for _, neighbor := range graph[node] {
		dfs(neighbor, graph, visited)
	}
}

// time complexity: O(2^n) -> the recursion can branch into two choices at each of `n` levels.
// space complexity: O(n) -> the auxiliary slice, map, table, queue, or returned collection can grow with `n`.
```

Without `visited`, a cycle such as:

```text
A → B → C → A
```

causes infinite recursion.

---

## 5. What is the difference between recursion-tree size and recursion depth?

For naive Fibonacci:

```text
Number of calls: O(2^n)
Maximum depth:   O(n)
```

Therefore:

```text
Time:  O(2^n)
Space: O(n)
```

---

## 6. How do you optimize a backtracking problem?

Common techniques:

* Prune invalid branches early.
* Sort input to enable early termination.
* Skip duplicate choices.
* Use efficient mutable state.
* Avoid unnecessary copying except when saving results.
* Use bitmasks for small state spaces.
* Cache states when branches overlap.

---

# 31. Interview walkthrough: Subsets

Suppose the interviewer asks:

> Generate all subsets of an integer array.

A structured answer:

## Step 1: Define the state

```text
index = which number we are deciding about
path = currently selected numbers
```

## Step 2: Define choices

For each number:

```text
Exclude it
Include it
```

## Step 3: Define the base case

```text
index == len(nums)
```

All decisions have been made.

## Step 4: Save a copy

```go
// Exact question: How does this Go example demonstrate Step 4: Save a copy?
//
// Possible answer: Use the fragment with a slice whose indexes identify positions and whose elements hold their values.
//
// Output format: This is a partial Go fragment; its surrounding function determines the final returned value.
//
// Inline descriptions:
// - The comments in this preface describe how the important expressions and state changes are used.
//
// Boundary checks:
// - No explicit boundary branch appears in this fragment; its caller or surrounding example supplies valid inputs.
//
// Key variables:
// - `subset` holds the intermediate value produced by `append([]int(nil), path...`.
//
// Logic:
// 1. Create or use a slice so indexes identify positions and elements store their data or state.
subset := append([]int(nil), path...)

// time complexity: O(n * 2^n) -> the recursion can branch into two choices at each of `n` levels.
// space complexity: O(n) -> the auxiliary slice, map, table, queue, or returned collection can grow with `n`.
```

## Step 5: Analyze complexity

There are `2^n` subsets, and copying each can take `O(n)`.

```text
Time: O(n × 2^n)
Auxiliary recursion space: O(n)
Output space: O(n × 2^n)
```

This explanation is as important as the code.

---

# 32. Interview walkthrough: Permutations

State:

```text
path = current permutation
used[i] = whether nums[i] is already selected
```

Base case:

```text
len(path) == len(nums)
```

Choices:

```text
Every unused number
```

Backtracking:

```text
Mark used
Append
Recurse
Remove
Mark unused
```

Complexity:

```text
n! outputs
Each output has n elements

Time: O(n × n!)
Auxiliary space: O(n)
```

---

# 33. How recursion connects to other topics

```mermaid
flowchart TD
    A[Recursion] --> B[Tree DFS]
    A --> C[Graph DFS]
    A --> D[Backtracking]
    A --> E[Divide and Conquer]
    A --> F[Dynamic Programming]

    D --> G[Subsets]
    D --> H[Permutations]
    D --> I[Combination Sum]

    E --> J[Merge Sort]
    E --> K[Quick Sort]
    E --> L[Binary Search]

    F --> M[Memoization]
    F --> N[Optimal substructure]
```

Recursion is not an isolated topic.

It is the foundation of:

* Tree algorithms
* Graph DFS
* Backtracking
* Divide and conquer
* Memoized dynamic programming
* Parsing
* Compiler syntax trees
* File-system traversal

---

# 34. The simplest mental checklist

Before coding, say this aloud:

```text
1. What does my function return?
2. What is the smallest input?
3. What smaller problem will I call?
4. How will I use the smaller answer?
5. Does every path stop?
6. How deep can the stack become?
7. Am I recalculating the same state?
8. Do I need to undo mutable state?
```

For backtracking:

```text
State
Choices
Constraints
Base case
Choose
Explore
Undo
```

---

# 35. Final cheat sheet

## Basic recursion

```go
// Exact question: How does `solve` solve Basic recursion?
//
// Possible answer: Use `solve` to compare the current values and return or update state when the condition matches.
//
// Output format: Return an `int` value from `solve`; the function does not print the answer.
//
// Inline descriptions:
// - `n` is the int input used by this example.
//
// Boundary checks:
// - `n == 0` handles the smallest valid state or recursive base case.
//
// Key variables:
// - `n` is the int input used by this example.
//
// Logic:
// 1. Recursively reduce the current problem to smaller calls until a base condition is reached.
// 2. Return the value produced after the state updates are complete.
func solve(n int) int {
	if n == 0 {
		return base
	}

	return combine(n, solve(n-1))
}

// time complexity: O(n) -> the algorithm visits each of the `n` input elements or states once.
// space complexity: O(n) -> the recursion stack and any cache can grow to one entry per input state.
```

## Tree recursion

```go
// Exact question: How does `solve` solve Tree recursion?
//
// Possible answer: Use `solve` to compare the current values and return or update state when the condition matches.
//
// Output format: Return an `int` value from `solve`; the function does not print the answer.
//
// Inline descriptions:
// - `node` points to a TreeNode value that the function reads or updates.
//
// Boundary checks:
// - `node == nil` checks whether the referenced value exists before it is used.
//
// Key variables:
// - `node` points to a TreeNode value that the function reads or updates.
// - `left` marks the current left boundary or left-side value.
// - `right` marks the current right boundary or right-side value.
//
// Logic:
// 1. Recursively reduce the current problem to smaller calls until a base condition is reached.
// 2. Return the value produced after the state updates are complete.
func solve(node *TreeNode) int {
	if node == nil {
		return 0
	}

	left := solve(node.Left)
	right := solve(node.Right)

	return combine(node, left, right)
}

// time complexity: O(n) -> the algorithm visits each of the `n` input elements or states once.
// space complexity: O(h) -> the recursion stack grows with the height of the traversed tree or Trie.
```

## Backtracking

```go
// Exact question: How does `backtrack` solve Backtracking?
//
// Possible answer: Choose an option, recurse, then undo the choice before trying the next option.
//
// Output format: `backtrack` has no return value; its observable result is the mutation or output performed in the function body.
//
// Inline descriptions:
// - The comments in this preface describe how the important expressions and state changes are used.
//
// Boundary checks:
// - `complete()` decides whether the branch or loop should continue for the current input.
// - `!valid(choice)` decides whether the branch or loop should continue for the current input.
//
// Key variables:
// - This fragment operates directly on the values named in each statement; it introduces no separate data structure.
//
// Logic:
// 1. Iterate through the required elements or states in the order shown.
// 2. Recursively reduce the current problem to smaller calls until a base condition is reached.
// 3. Return the value produced after the state updates are complete.
func backtrack() {
	if complete() {
		saveCopy()
		return
	}

	for _, choice := range choices {
		if !valid(choice) {
			continue
		}

		makeChoice(choice)
		backtrack()
		undoChoice(choice)
	}
}

// time complexity: O(n) -> the algorithm visits each of the `n` input elements or states once.
// space complexity: O(n) -> the recursion stack and any cache can grow to one entry per input state.
```

## Complexity

```text
Time =
number of calls
×
work per call

Space =
maximum recursion depth
+
additional state
```

## Core mental model

> Trust the recursive call to correctly solve the smaller problem.

Your responsibility is only to define:

```text
Base case
Smaller problem
How to combine the answer
```
