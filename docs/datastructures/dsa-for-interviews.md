For software engineering interviews, you should know these **data structures and algorithms** well enough to explain, implement, and analyze time/space complexity.

## 1. Complexity Analysis

You must be comfortable with:

| Topic                       | What to know                                      |
| --------------------------- | ------------------------------------------------- |
| Big-O notation              | `O(1)`, `O(log n)`, `O(n)`, `O(n log n)`, `O(n²)` |
| Time complexity             | How runtime grows with input size                 |
| Space complexity            | Extra memory used by your solution                |
| Best / average / worst case | Especially for sorting, hashing, trees            |
| Amortized analysis          | Example: dynamic array append                     |

Common interview expectation:

```text
Can you explain why your solution is O(n)?
Can you reduce it from O(n²) to O(n)?
Can you trade memory for speed?
```

---

# Data Structures

## 2. Arrays / Slices

Very important.

Know:

| Topic                 | Examples                            |
| --------------------- | ----------------------------------- |
| Traversal             | Loop through array                  |
| Two pointers          | Sorted array, palindrome, pair sum  |
| Sliding window        | Longest substring, max sum subarray |
| Prefix sum            | Range sum, subarray sum             |
| Difference array      | Range update problems               |
| Kadane’s algorithm    | Maximum subarray sum                |
| In-place modification | Remove duplicates, rotate array     |

Common problems:

```text
Two Sum
Best Time to Buy and Sell Stock
Maximum Subarray
Rotate Array
Merge Sorted Array
Move Zeroes
Product of Array Except Self
Container With Most Water
```

---

## 3. Strings

Very important.

Know:

| Topic               | Examples                         |
| ------------------- | -------------------------------- |
| Character frequency | Anagram, duplicates              |
| String scanning     | Substring search                 |
| Sliding window      | Longest substring without repeat |
| Palindrome checks   | Valid palindrome                 |
| String builder      | Efficient string construction    |
| Trie basics         | Prefix search, autocomplete      |

Common problems:

```text
Valid Anagram
Longest Substring Without Repeating Characters
Valid Palindrome
Group Anagrams
Longest Palindromic Substring
Minimum Window Substring
String Compression
```

---

## 4. Hash Map / Hash Set

Extremely important.

Know:

| Topic               | What to know               |
| ------------------- | -------------------------- |
| Hash map lookup     | Average `O(1)`             |
| Hash set membership | Duplicate detection        |
| Frequency counting  | Count chars/numbers        |
| Index mapping       | Two Sum style problems     |
| Collision basics    | Why worst case can degrade |

Common problems:

```text
Two Sum
Contains Duplicate
Group Anagrams
Subarray Sum Equals K
Longest Consecutive Sequence
First Unique Character
```

---

## 5. Linked List

Important, but less used in real systems.

Know:

| Topic                 | Examples                     |
| --------------------- | ---------------------------- |
| Singly linked list    | Basic traversal              |
| Doubly linked list    | LRU cache                    |
| Fast and slow pointer | Cycle detection, middle node |
| Reverse linked list   | Iterative and recursive      |
| Dummy node            | Cleaner insert/delete logic  |
| Merge lists           | Sorted linked lists          |

Common problems:

```text
Reverse Linked List
Detect Cycle
Merge Two Sorted Lists
Remove Nth Node From End
Middle of Linked List
Add Two Numbers
LRU Cache
```

---

## 6. Stack

Very important.

Know:

| Topic                 | Examples             |
| --------------------- | -------------------- |
| LIFO behavior         | Last in, first out   |
| Parentheses matching  | Valid brackets       |
| Monotonic stack       | Next greater element |
| DFS implementation    | Graph/tree traversal |
| Expression evaluation | Calculator problems  |

Common problems:

```text
Valid Parentheses
Min Stack
Daily Temperatures
Next Greater Element
Evaluate Reverse Polish Notation
Largest Rectangle in Histogram
```

---

## 7. Queue / Deque

Important.

Know:

| Topic          | Examples                      |
| -------------- | ----------------------------- |
| FIFO behavior  | First in, first out           |
| BFS            | Level-order traversal         |
| Deque          | Sliding window maximum        |
| Circular queue | Fixed-size queue              |
| Priority queue | Usually implemented with heap |

Common problems:

```text
Binary Tree Level Order Traversal
Number of Islands
Rotting Oranges
Sliding Window Maximum
Design Circular Queue
```

---

## 8. Heap / Priority Queue

Very important for medium/hard interviews.

Know:

| Topic                | Examples                    |
| -------------------- | --------------------------- |
| Min heap             | Get smallest item           |
| Max heap             | Get largest item            |
| Top K                | K largest/frequent elements |
| Merge K sorted lists | Heap-based merge            |
| Scheduling           | Meeting rooms, CPU tasks    |

Common problems:

```text
Kth Largest Element
Top K Frequent Elements
Merge K Sorted Lists
Find Median from Data Stream
Meeting Rooms II
Task Scheduler
```

---

## 9. Trees

Very important.

Know:

| Topic                  | Examples                     |
| ---------------------- | ---------------------------- |
| Binary tree traversal  | Preorder, inorder, postorder |
| BFS traversal          | Level order                  |
| DFS traversal          | Recursive and iterative      |
| Tree height/depth      | Max depth                    |
| Balanced tree          | Height-balanced check        |
| Lowest common ancestor | Common ancestor problems     |

Common problems:

```text
Maximum Depth of Binary Tree
Invert Binary Tree
Same Tree
Subtree of Another Tree
Binary Tree Level Order Traversal
Lowest Common Ancestor
Diameter of Binary Tree
```

---

## 10. Binary Search Tree

Important.

Know:

| Topic                | Examples                 |
| -------------------- | ------------------------ |
| BST property         | Left < root < right      |
| Search               | `O(log n)` if balanced   |
| Insert/delete basics | Structural changes       |
| Inorder traversal    | Gives sorted order       |
| Validate BST         | Boundary-based recursion |

Common problems:

```text
Validate Binary Search Tree
Kth Smallest Element in BST
Lowest Common Ancestor of BST
Search in BST
Insert into BST
```

---

## 11. Trie

Useful for string-heavy interviews.

Know:

| Topic         | Examples             |
| ------------- | -------------------- |
| Prefix tree   | Fast prefix lookup   |
| Insert/search | Word dictionary      |
| StartsWith    | Autocomplete         |
| DFS with trie | Word search problems |

Common problems:

```text
Implement Trie
Word Search II
Replace Words
Design Add and Search Words Data Structure
Autocomplete System
```

---

## 12. Graphs

Very important.

Know:

| Topic                | Examples                          |
| -------------------- | --------------------------------- |
| Graph representation | Adjacency list, matrix            |
| BFS                  | Shortest path in unweighted graph |
| DFS                  | Connected components              |
| Cycle detection      | Directed and undirected           |
| Topological sort     | Dependency ordering               |
| Union-Find           | Connected components              |
| Dijkstra             | Shortest path with weights        |

Common problems:

```text
Number of Islands
Clone Graph
Course Schedule
Pacific Atlantic Water Flow
Rotting Oranges
Network Delay Time
Redundant Connection
Evaluate Division
```

---

## 13. Union-Find / Disjoint Set

Important for graph connectivity problems.

Know:

| Topic              | Examples         |
| ------------------ | ---------------- |
| Find               | Find parent/root |
| Union              | Merge components |
| Path compression   | Optimization     |
| Union by rank/size | Optimization     |
| Cycle detection    | Undirected graph |

Common problems:

```text
Number of Connected Components
Redundant Connection
Accounts Merge
Friend Circles / Number of Provinces
Graph Valid Tree
```

---

# Algorithms

## 14. Sorting

Must know.

Know:

| Algorithm      |     Average Time | Notes                       |
| -------------- | ---------------: | --------------------------- |
| Bubble sort    |          `O(n²)` | Mostly for basics           |
| Selection sort |          `O(n²)` | Simple but inefficient      |
| Insertion sort |          `O(n²)` | Good for nearly sorted data |
| Merge sort     |     `O(n log n)` | Stable, uses extra space    |
| Quick sort     | `O(n log n)` avg | Worst case `O(n²)`          |
| Heap sort      |     `O(n log n)` | In-place, not stable        |
| Counting sort  |       `O(n + k)` | When value range is small   |

Common problems:

```text
Merge Intervals
Sort Colors
Kth Largest Element
Meeting Rooms
Insert Interval
```

---

## 15. Binary Search

Extremely important.

Know:

| Topic                | Examples                       |
| -------------------- | ------------------------------ |
| Normal binary search | Find target                    |
| Lower bound          | First value >= target          |
| Upper bound          | First value > target           |
| Search answer space  | Minimum capacity, minimum days |
| Rotated sorted array | Modified binary search         |

Common problems:

```text
Binary Search
Search in Rotated Sorted Array
Find Minimum in Rotated Sorted Array
First Bad Version
Find First and Last Position
Median of Two Sorted Arrays
Koko Eating Bananas
Capacity to Ship Packages
```

---

## 16. Recursion

Very important.

Know:

| Topic          | Examples                           |
| -------------- | ---------------------------------- |
| Base case      | Stop condition                     |
| Recursive case | Smaller subproblem                 |
| Call stack     | Memory usage                       |
| Tree recursion | DFS                                |
| Backtracking   | Generate combinations/permutations |

Common problems:

```text
Factorial
Fibonacci
Generate Parentheses
Subsets
Permutations
Combination Sum
Tree Traversals
```

---

## 17. Backtracking

Important for medium/hard problems.

Know:

| Topic                       | Examples            |
| --------------------------- | ------------------- |
| Choose / explore / unchoose | Standard pattern    |
| Pruning                     | Avoid useless paths |
| State tracking              | Current path        |
| Constraint checking         | Valid choices       |

Common problems:

```text
Subsets
Permutations
Combination Sum
Generate Parentheses
Word Search
N-Queens
Palindrome Partitioning
```

---

## 18. Dynamic Programming

Very important for senior-level interviews.

Know:

| Topic               | Examples                      |
| ------------------- | ----------------------------- |
| Memoization         | Top-down DP                   |
| Tabulation          | Bottom-up DP                  |
| 1D DP               | Climbing stairs, house robber |
| 2D DP               | Grid paths, edit distance     |
| Knapsack            | Pick/not pick decisions       |
| Longest subsequence | LIS, LCS                      |
| State transition    | Core DP skill                 |

Common problems:

```text
Climbing Stairs
House Robber
Coin Change
Longest Increasing Subsequence
Longest Common Subsequence
Edit Distance
Unique Paths
Word Break
Partition Equal Subset Sum
```

---

## 19. Greedy Algorithms

Important.

Know:

| Topic                | Examples              |
| -------------------- | --------------------- |
| Local optimal choice | Choose best now       |
| Sorting + greedy     | Intervals, scheduling |
| Heap + greedy        | Top K, scheduling     |
| Proof intuition      | Why greedy works      |

Common problems:

```text
Jump Game
Gas Station
Non-overlapping Intervals
Meeting Rooms
Task Scheduler
Partition Labels
Minimum Number of Arrows to Burst Balloons
```

---

## 20. Sliding Window

Extremely common.

Know:

| Topic                | Examples             |
| -------------------- | -------------------- |
| Fixed-size window    | Max sum of size K    |
| Variable-size window | Longest substring    |
| Frequency map window | Anagrams, min window |
| Shrinking condition  | Keep window valid    |

Common problems:

```text
Maximum Average Subarray
Longest Substring Without Repeating Characters
Minimum Window Substring
Permutation in String
Find All Anagrams in a String
Longest Repeating Character Replacement
```

---

## 21. Two Pointers

Extremely common.

Know:

| Topic          | Examples             |
| -------------- | -------------------- |
| Opposite ends  | Pair sum, palindrome |
| Same direction | Remove duplicates    |
| Fast/slow      | Linked list cycle    |
| Merge style    | Merge sorted arrays  |

Common problems:

```text
Two Sum II
3Sum
Container With Most Water
Valid Palindrome
Remove Duplicates from Sorted Array
Linked List Cycle
Merge Sorted Array
```

---

## 22. Prefix Sum

Very useful.

Know:

| Topic                 | Examples              |
| --------------------- | --------------------- |
| Running sum           | Sum from start        |
| Range sum query       | `sum[i:j]`            |
| Hash map + prefix sum | Subarray sum problems |
| 2D prefix sum         | Matrix range sum      |

Common problems:

```text
Range Sum Query
Subarray Sum Equals K
Continuous Subarray Sum
Product of Array Except Self
Maximum Size Subarray Sum Equals K
```

---

## 23. Bit Manipulation

Useful, but lower priority unless company asks.

Know:

| Topic          | Examples         |
| -------------- | ---------------- |
| AND, OR, XOR   | Basic operations |
| Check odd/even | `n & 1`          |
| Single number  | XOR trick        |
| Bit mask       | Subsets          |
| Count bits     | Hamming weight   |

Common problems:

```text
Single Number
Number of 1 Bits
Counting Bits
Missing Number
Reverse Bits
Power of Two
Subsets using Bitmask
```

---

## 24. Intervals

Very common in practical interviews.

Know:

| Topic             | Examples                  |
| ----------------- | ------------------------- |
| Sort by start/end | Standard first step       |
| Merge overlaps    | Merge intervals           |
| Insert interval   | Maintain sorted intervals |
| Meeting rooms     | Overlap detection         |
| Sweep line        | Count active intervals    |

Common problems:

```text
Merge Intervals
Insert Interval
Meeting Rooms
Meeting Rooms II
Non-overlapping Intervals
Employee Free Time
Minimum Number of Arrows
```

---

## 25. Matrix / Grid Problems

Important.

Know:

| Topic            | Examples           |
| ---------------- | ------------------ |
| BFS on grid      | Shortest path      |
| DFS on grid      | Islands            |
| Boundary checks  | Valid cell         |
| Visited tracking | Avoid loops        |
| Direction arrays | Up/down/left/right |

Common problems:

```text
Number of Islands
Max Area of Island
Rotting Oranges
Word Search
Pacific Atlantic Water Flow
Shortest Path in Binary Matrix
Set Matrix Zeroes
Spiral Matrix
```

---

# Priority Order

If you have limited time, study in this order:

## Must Know

```text
1. Arrays
2. Strings
3. Hash Map / Hash Set
4. Two Pointers
5. Sliding Window
6. Binary Search
7. Stack
8. Queue
9. Trees
10. Graph BFS / DFS
```

## Should Know

```text
11. Heap / Priority Queue
12. Recursion
13. Backtracking
14. Dynamic Programming
15. Greedy
16. Prefix Sum
17. Intervals
18. Linked List
```

## Good to Know

```text
19. Trie
20. Union-Find
21. Bit Manipulation
22. Advanced DP
23. Dijkstra
24. Topological Sort
25. Segment Tree / Fenwick Tree
```

---

# For Engineering Manager / Senior Engineer Interviews

For your profile, do not only prepare coding. You should also prepare:

```text
System design
Distributed systems
Kubernetes fundamentals
API design
Database design
Concurrency
Operational excellence
Incident handling
Trade-off discussion
Code quality
Testing strategy
```

For coding rounds, focus on:

```text
Arrays
Hash maps
Sliding window
Binary search
Trees
Graphs
Heaps
Dynamic programming basics
Concurrency problems in Go
```

For platform/Kubernetes roles, also revise:

```text
Queues
Rate limiting
Leader election
Work queues
Retries
Backoff
Caching
Idempotency
Event-driven systems
Graph traversal
Dependency resolution
Scheduling algorithms
```

---

# Suggested 6-Week Preparation Plan

## Week 1

```text
Arrays
Strings
Hash Map
Two Pointers
Sliding Window
```

## Week 2

```text
Binary Search
Sorting
Intervals
Prefix Sum
Stack
Queue
```

## Week 3

```text
Linked List
Trees
BST
Heap
Recursion
```

## Week 4

```text
Graphs
BFS
DFS
Topological Sort
Union-Find
```

## Week 5

```text
Dynamic Programming
Greedy
Backtracking
Trie
Bit Manipulation
```

## Week 6

```text
Mock interviews
Timed LeetCode medium problems
System design
Behavioral stories
Go concurrency
```

---

# Minimum Problem Count

A realistic target:

```text
Easy: 40–50
Medium: 80–120
Hard: 10–20
Total: 130–190 problems
```

For your level, prioritize **medium problems**. Hard problems are useful, but system design and senior-level tradeoff discussion will matter more.
