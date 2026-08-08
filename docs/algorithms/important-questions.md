
```go
// Exact question: Given `n` services and their dependencies, return one valid deployment order or an empty slice when a cycle makes deployment impossible.
//
// Possible answer: Use `deploymentOrder` with a FIFO queue to process items in discovery order and record each result.
//
// Output format: Return the `([]int, error)` value from `deploymentOrder`; the function does not print the answer.
//
// Inline descriptions:
// - `serviceCount` is the int input used by this example.
// - `dependencies` is a two-dimensional slice: row and column indexes identify positions, and each cell stores a value.
//
// Boundary checks:
// - `serviceCount < 0` rejects or terminates a state that has moved below the valid range.
// - `serviceCount == 0` handles the smallest valid state or recursive base case.
// - `len(dependency) != 2` decides whether the branch or loop should continue for the current input.
//
// Key variables:
// - `serviceCount` is the int input used by this example.
// - `dependencies` is a two-dimensional slice: row and column indexes identify positions, and each cell stores a value.
// - `adjacencyList` is a two-dimensional table; indexes select a row and column, while cells hold computed values.
// - `inDegree` is a slice; indexes identify positions and elements hold the corresponding values.
// - `seenEdges` is a map whose keys are int lookups and whose values are stored map[int]struct{} results.
//
// Logic:
// 1. Create or use a map to associate each lookup key with its stored value.
// 2. Create or use a slice so indexes identify positions and elements store their data or state.
// 3. Iterate through the required elements or states in the order shown.
package main

import "fmt"

// Question:
// You are given n services numbered from 0 to n-1.
//
// Each dependency is represented as:
//
//	[service, prerequisite]
//
// Example:
//
//	[3, 1]
//
// means:
//
//	Service 3 depends on service 1.
//	Service 1 must be deployed before service 3.
//
// Return one valid deployment order.
//
// If no valid deployment order exists because of a cyclic dependency,
// return an empty slice.
//
// Input to function:
//
//	serviceCount int
//		The total number of services.
//		Valid service IDs are from 0 to serviceCount-1.
//
//	dependencies [][]int
//		Each element must contain exactly two integers:
//
//			[service, prerequisite]
//
// Output from the function:
//
//	[]int
//		A valid order in which the services can be deployed.
//
//		If the dependency graph contains a cycle or the input is invalid,
//		return an empty slice.
//
// Examples:
//
// Example 1:
//
//	serviceCount = 4
//
//	dependencies = [
//	    [1, 0],
//	    [2, 0],
//	    [3, 1],
//	    [3, 2],
//	]
//
// Meaning:
//
//	Service 1 depends on service 0.
//	Service 2 depends on service 0.
//	Service 3 depends on service 1.
//	Service 3 depends on service 2.
//
// One valid output:
//
//	[0, 1, 2, 3]
//
// Another valid output:
//
//	[0, 2, 1, 3]
//
// Example 2:
//
//	serviceCount = 3
//
//	dependencies = [
//	    [1, 0],
//	    [2, 1],
//	    [0, 2],
//	]
//
// Graph:
//
//	0 -> 1 -> 2 -> 0
//
// This graph contains a cycle.
//
// Output:
//
//	[]
//
// -----------------------------------------------------------------------------
//
// Answer:
//
// Approach:
//
// Model the problem as a directed graph.
//
// For a dependency:
//
//	[service, prerequisite]
//
// create the directed edge:
//
//	prerequisite -> service
//
// This direction means:
//
//	"After prerequisite is deployed, service may become eligible."
//
// Use Kahn's algorithm for topological sorting.
//
// Kahn's algorithm:
//
//  1. Build an adjacency list.
//  2. Count the incoming edges for each service using an indegree array.
//  3. Add every service with indegree 0 to a queue.
//  4. Repeatedly remove a service from the queue.
//  5. Add it to the deployment order.
//  6. Reduce the indegree of services that depend on it.
//  7. When another service's indegree becomes 0, add it to the queue.
//  8. If every service is processed, the graph is acyclic.
//  9. If fewer than serviceCount services are processed, a cycle exists.
//
// Why use a directed graph?
//
// Dependencies have direction.
//
// If service 3 depends on service 1:
//
//	1 must happen before 3.
//
// Therefore:
//
//	1 -> 3
//
// Why use topological sorting?
//
// Topological sorting returns an ordering of directed graph nodes where every
// prerequisite appears before every dependent node.
//
// Why use Kahn's algorithm?
//
// Kahn's algorithm naturally provides:
//
//	- A valid deployment order.
//	- Cycle detection.
//	- Services that are ready to deploy.
//	- A foundation for parallel deployment waves.
//	- A foundation for checking whether the order is unique.
//
// Questions to ask the interviewer before coding:
//
//  1. Are service IDs guaranteed to be between 0 and serviceCount-1?
//
//  2. Can dependencies contain malformed entries such as:
//
//	[1]
//	[1, 2, 3]
//
//  3. Can a service depend on itself?
//
//	[2, 2]
//
//	This is an immediate cycle.
//
//  4. Can the same dependency appear more than once?
//
//	[1, 0]
//	[1, 0]
//
//	If duplicates are allowed, should they be treated as one dependency
//	or two separate edges?
//
//  5. Should the output be deterministic?
//
//	For example, if services 1 and 2 are both ready, should the smaller
//	service ID always be chosen first?
//
//	This implementation uses a FIFO queue and is deterministic for integer
//	service IDs because the initial scan is from 0 to serviceCount-1 and
//	adjacency lists preserve insertion order.
//
//  6. Should an empty graph return:
//
//	[]
//
//	or should it be considered invalid?
//
//	This implementation returns an empty slice.
//
//  7. Should invalid input return an error instead of an empty slice?
//
//	This implementation returns an error so callers can distinguish:
//
//	- Invalid input.
//	- Cyclic dependency.
//	- Valid empty graph.
//
//  8. Do we need only one valid order, or all possible valid orders?
//
//	This implementation returns one valid order.
//
//  9. Do we need to return the exact cycle?
//
//	Kahn's algorithm detects that a cycle exists but does not directly
//	reconstruct the cycle.
//
// 10. Will this function be called concurrently?
//
//	This implementation uses only local variables, so separate calls are
//	safe to execute concurrently.
//
// How to plan before coding:
//
// Step 1:
//
// Define the edge direction carefully.
//
// Given:
//
//	[service, prerequisite]
//
// use:
//
//	prerequisite -> service
//
// A common mistake is reversing this edge.
//
// Step 2:
//
// Build:
//
//	adjacencyList[prerequisite] = list of dependent services
//
// Step 3:
//
// Build:
//
//	inDegree[service] = number of prerequisites not yet processed
//
// Step 4:
//
// Put every service with inDegree == 0 into the queue.
//
// Step 5:
//
// Process the queue using a head index.
//
// Avoid repeatedly removing queue[0] because that can require shifting elements.
//
// Step 6:
//
// For each processed service:
//
//	- Add it to order.
//	- Visit all dependent services.
//	- Decrease their indegree.
//	- Add a dependent service to the queue when its indegree becomes 0.
//
// Step 7:
//
// Detect a cycle by comparing:
//
//	len(order)
//
// with:
//
//	serviceCount
//
// If they are different, at least one service could never be processed.
//
// Time complexity:
//
//	O(V + E)
//
// Where:
//
//	V = number of services
//	E = number of dependency edges
//
// Why:
//
//	- Every service is added to and removed from the queue at most once.
//	- Every dependency edge is processed once.
//
// Space complexity:
//
//	O(V + E)
//
// Used by:
//
//	- Adjacency list: O(V + E)
//	- Indegree array: O(V)
//	- Queue: O(V)
//	- Output order: O(V)
//
// Race conditions:
//
// This function is safe for concurrent calls because:
//
//	- It does not modify global state.
//	- It does not share mutable data between calls.
//	- Every internal slice and map is created inside the function.
//
// However:
//
//	- The caller must not mutate the dependencies slice concurrently while
//	  this function is reading it.
//	- Concurrent reads are safe only when no goroutine modifies the same data.
//	- If dependencies can be mutated concurrently, the caller must protect it
//	  with synchronization or pass an immutable copy.
//
// -----------------------------------------------------------------------------

func deploymentOrder(
	serviceCount int,
	dependencies [][]int,
) ([]int, error) {
	// Boundary check:
	//
	// A negative service count is invalid because a collection cannot contain
	// a negative number of services.
	if serviceCount < 0 {
		return nil, fmt.Errorf(
			"service count cannot be negative: %d",
			serviceCount,
		)
	}

	// Boundary check:
	//
	// No services means there is nothing to deploy.
	//
	// This is a valid empty graph, so return an empty deployment order.
	if serviceCount == 0 {
		return []int{}, nil
	}

	// adjacencyList is indexed by service ID.
	//
	// The index represents a prerequisite service.
	//
	// The values stored at that index are dependent service IDs.
	//
	// Example:
	//
	//	dependencies = [[1,0], [2,0]]
	//
	//	adjacencyList[0] = [1,2]
	//
	// Meaning:
	//
	//	After service 0 is deployed, services 1 and 2 each have one
	// prerequisite satisfied.
	//
	// Important distinction:
	//
	//	- Index = prerequisite service ID.
	//	- Stored values = dependent service IDs.
	adjacencyList := make([][]int, serviceCount)

	// inDegree is indexed by service ID.
	//
	// The value at inDegree[i] is the number of unresolved prerequisites
	// service i currently has.
	//
	// Example:
	//
	//	dependencies = [[1,0], [2,0], [3,1], [3,2]]
	//
	//	inDegree[0] = 0
	//	inDegree[1] = 1
	//	inDegree[2] = 1
	//	inDegree[3] = 2
	//
	// Important distinction:
	//
	//	- Index = service ID.
	//	- Stored value = count of unresolved prerequisites.
	inDegree := make([]int, serviceCount)

	// seenEdges prevents duplicate dependency edges from being counted twice.
	//
	// Without duplicate detection:
	//
	//	dependencies = [[1,0], [1,0]]
	//
	// would incorrectly produce:
	//
	//	inDegree[1] = 2
	//
	// even though logically service 1 has only one unique prerequisite:
	//
	//	service 0
	//
	// The outer map key is the prerequisite service.
	// The inner map key is the dependent service.
	//
	// Example:
	//
	//	seenEdges[0][1] = true
	//
	// means the edge:
	//
	//	0 -> 1
	//
	// has already been recorded.
	seenEdges := make(map[int]map[int]struct{})

	for dependencyIndex, dependency := range dependencies {
		// Boundary check:
		//
		// Each dependency must contain exactly:
		//
		//	[service, prerequisite]
		if len(dependency) != 2 {
			return nil, fmt.Errorf(
				"dependency at index %d must contain exactly two values",
				dependencyIndex,
			)
		}

		// service is a service ID value, not an array position in dependencies.
		//
		// For:
		//
		//	dependency = [3,1]
		//
		//	service = 3
		service := dependency[0]

		// prerequisite is a service ID value.
		//
		// For:
		//
		//	dependency = [3,1]
		//
		//	prerequisite = 1
		prerequisite := dependency[1]

		// Boundary check:
		//
		// Valid service IDs are:
		//
		//	0 through serviceCount-1
		if service < 0 || service >= serviceCount {
			return nil, fmt.Errorf(
				"service ID %d at dependency index %d is outside range [0,%d]",
				service,
				dependencyIndex,
				serviceCount-1,
			)
		}

		// Boundary check:
		//
		// The prerequisite must also refer to an existing service.
		if prerequisite < 0 || prerequisite >= serviceCount {
			return nil, fmt.Errorf(
				"prerequisite ID %d at dependency index %d is outside range [0,%d]",
				prerequisite,
				dependencyIndex,
				serviceCount-1,
			)
		}

		// Cyclic dependency check:
		//
		// A self-dependency is an immediate cycle.
		//
		// Example:
		//
		//	[2,2]
		//
		// means:
		//
		//	Service 2 cannot deploy until service 2 has already deployed.
		if service == prerequisite {
			return []int{}, nil
		}

		// Initialize the inner set for this prerequisite when necessary.
		if _, exists := seenEdges[prerequisite]; !exists {
			seenEdges[prerequisite] = make(map[int]struct{})
		}

		// Duplicate dependency check:
		//
		// Ignore an edge that has already been added.
		//
		// This prevents duplicate edges from incorrectly increasing indegree.
		if _, duplicate := seenEdges[prerequisite][service]; duplicate {
			continue
		}

		seenEdges[prerequisite][service] = struct{}{}

		// Build the directed edge:
		//
		//	prerequisite -> service
		//
		// Example:
		//
		//	[3,1]
		//
		// becomes:
		//
		//	adjacencyList[1] = append(adjacencyList[1], 3)
		adjacencyList[prerequisite] = append(
			adjacencyList[prerequisite],
			service,
		)

		// service has one more unresolved prerequisite.
		//
		// Example:
		//
		//	[3,1]
		//
		// changes:
		//
		//	inDegree[3] from 0 to 1
		inDegree[service]++
	}

	// queue stores service ID values that are currently ready to deploy.
	//
	// It does not store dependency indexes or positions from the input.
	//
	// A service enters the queue only when:
	//
	//	inDegree[service] == 0
	//
	// Meaning:
	//
	//	The service has no remaining unresolved prerequisites.
	queue := make([]int, 0, serviceCount)

	for serviceID := 0; serviceID < serviceCount; serviceID++ {
		if inDegree[serviceID] == 0 {
			queue = append(queue, serviceID)
		}
	}

	// order stores service ID values in deployment order.
	//
	// Example:
	//
	//	order = [0,1,2]
	//
	// means services 0, 1, and 2 have been processed in that order.
	order := make([]int, 0, serviceCount)

	// head is an index into the queue slice.
	//
	// It is not a service ID.
	//
	// queue[head] is the next service ID to process.
	//
	// Example:
	//
	//	queue = [0,1,2,3]
	//	head  = 2
	//
	// means:
	//
	//	- queue[0] and queue[1] have already been processed.
	//	- queue[2], whose value is service ID 2, is next.
	head := 0

	// Example trace:
	//
	// serviceCount = 4
	//
	// dependencies:
	//
	//	[1,0]
	//	[2,0]
	//	[3,1]
	//	[3,2]
	//
	// Graph:
	//
	//	    0
	//	   / \
	//	  1   2
	//	   \ /
	//	    3
	//
	// Initial values:
	//
	//	adjacencyList[0] = [1,2]
	//	adjacencyList[1] = [3]
	//	adjacencyList[2] = [3]
	//	adjacencyList[3] = []
	//
	//	inDegree = [0,1,1,2]
	//
	//	queue = [0]
	//	head  = 0
	//	order = []
	//
	// Iteration 1:
	//
	//	currentService = queue[0] = 0
	//	head changes from 0 to 1
	//	order becomes [0]
	//
	//	Process dependent service 1:
	//
	//		inDegree[1] changes from 1 to 0
	//		service 1 is appended to queue
	//
	//	Process dependent service 2:
	//
	//		inDegree[2] changes from 1 to 0
	//		service 2 is appended to queue
	//
	//	Values after iteration 1:
	//
	//		inDegree = [0,0,0,2]
	//		queue     = [0,1,2]
	//		head      = 1
	//		order     = [0]
	//
	// Iteration 2:
	//
	//	currentService = queue[1] = 1
	//	head changes from 1 to 2
	//	order becomes [0,1]
	//
	//	Process dependent service 3:
	//
	//		inDegree[3] changes from 2 to 1
	//
	//	Service 3 is not appended because one unresolved prerequisite remains.
	//
	//	Values after iteration 2:
	//
	//		inDegree = [0,0,0,1]
	//		queue     = [0,1,2]
	//		head      = 2
	//		order     = [0,1]
	//
	// Iteration 3:
	//
	//	currentService = queue[2] = 2
	//	head changes from 2 to 3
	//	order becomes [0,1,2]
	//
	//	Process dependent service 3:
	//
	//		inDegree[3] changes from 1 to 0
	//		service 3 is appended to queue
	//
	//	Values after iteration 3:
	//
	//		inDegree = [0,0,0,0]
	//		queue     = [0,1,2,3]
	//		head      = 3
	//		order     = [0,1,2]
	//
	// Iteration 4:
	//
	//	currentService = queue[3] = 3
	//	head changes from 3 to 4
	//	order becomes [0,1,2,3]
	//
	//	Service 3 has no dependent services.
	//
	// Loop ends because:
	//
	//	head == len(queue)
	//
	//	4 == 4
	for head < len(queue) {
		// currentService holds a service ID value.
		//
		// It does not hold the queue index.
		//
		// head is the queue index.
		// queue[head] is the service ID value.
		currentService := queue[head]

		// Move head to the next queue position.
		head++

		// currentService is now safe to deploy because all its prerequisites
		// have already been processed.
		order = append(order, currentService)

		// adjacencyList[currentService] contains service ID values that depend
		// directly on currentService.
		for _, dependentService := range adjacencyList[currentService] {
			// One prerequisite of dependentService has now been satisfied.
			inDegree[dependentService]--

			// Defensive boundary check:
			//
			// An indegree should never become negative when duplicate edges
			// have been removed correctly.
			//
			// A negative value would indicate corrupted graph construction.
			if inDegree[dependentService] < 0 {
				return nil, fmt.Errorf(
					"internal error: negative indegree for service %d",
					dependentService,
				)
			}

			// The dependent service becomes eligible only when every
			// prerequisite has been processed.
			if inDegree[dependentService] == 0 {
				queue = append(queue, dependentService)
			}
		}
	}

	// Cyclic dependency detection:
	//
	// In an acyclic graph:
	//
	//	Every service eventually reaches indegree 0.
	//	Every service enters the queue.
	//	Every service is appended to order.
	//
	// Therefore:
	//
	//	len(order) == serviceCount
	//
	// In a cyclic graph:
	//
	//	Services in the cycle continue waiting for each other.
	//	Their indegree values never become 0.
	//	They never enter the queue.
	//
	// Example:
	//
	//	dependencies = [
	//	    [1,0],
	//	    [2,1],
	//	    [0,2],
	//	]
	//
	// Graph:
	//
	//	0 -> 1 -> 2 -> 0
	//
	// Initial indegrees:
	//
	//	inDegree = [1,1,1]
	//
	// No service has indegree 0.
	//
	// Therefore:
	//
	//	queue = []
	//	order = []
	//
	// Since:
	//
	//	len(order) != serviceCount
	//
	// a cycle exists.
	if len(order) != serviceCount {
		return []int{}, nil
	}

	return order, nil
}

func main() {
	serviceCount := 4

	dependencies := [][]int{
		{1, 0},
		{2, 0},
		{3, 1},
		{3, 2},
	}

	order, err := deploymentOrder(serviceCount, dependencies)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	if len(order) == 0 && serviceCount > 0 {
		fmt.Println("no valid deployment order: cyclic dependency")
		return
	}

	fmt.Println(order)

	// Expected output:
	//
	//	[0 1 2 3]
	//
	// Another valid topological order could be:
	//
	//	[0 2 1 3]
}

// Alternative questions that can be solved using this approach:
//
// 1. Course Schedule:
//
//	Given courses and prerequisites, determine whether all courses can
//	be completed.
//
//	Variation:
//
//	Return true or false instead of returning the complete order.
//
// 2. Course Schedule II:
//
//	Given courses and prerequisites, return one valid course-completion order.
//
//	This is structurally identical to the deployment-order problem.
//
// 3. Package Installation Order:
//
//	Given software packages and package dependencies, determine a safe
//	installation order.
//
//	Example:
//
//	logging -> framework -> application
//
// 4. Build-System Compilation Order:
//
//	Given source modules and import dependencies, determine the order in
//	which modules should be compiled.
//
// 5. CI/CD Pipeline Job Ordering:
//
//	Given pipeline jobs and prerequisites, determine a valid execution order.
//
//	Example:
//
//	build -> test -> package -> deploy
//
// 6. Kubernetes Resource Application Order:
//
//	Determine a safe order for applying resources such as:
//
//	Namespace -> ServiceAccount -> Deployment
//	CRD -> Custom Resource
//	Service -> Ingress
//
//	Production note:
//
//	Kubernetes readiness is eventually consistent, so object creation order
//	alone may not guarantee that the dependency is ready.
//
// 7. Terraform Resource Provisioning:
//
//	Determine the order for creating infrastructure resources.
//
//	Example:
//
//	VPC -> Subnet -> Load Balancer
//	IAM Role -> Compute Instance
//
// 8. Service Startup Order:
//
//	Determine the order in which system services should start.
//
//	Example:
//
//	storage -> database -> application
//
// 9. Safe Shutdown Order:
//
//	Compute the normal startup topological order and reverse it.
//
//	Example:
//
//	Startup:
//
//		database -> application
//
//	Shutdown:
//
//		application -> database
//
// 10. Detect Circular Configuration References:
//
//	Example:
//
//	A references B
//	B references C
//	C references A
//
//	The graph contains a cycle.
//
// 11. Workflow Dependency Resolution:
//
//	Determine whether a workflow contains valid task dependencies and
//	return a valid execution sequence.
//
// 12. Parallel Deployment Waves:
//
//	Use level-order Kahn's algorithm.
//
//	Instead of processing one queue entry at a time, process all services
//	currently in the queue as one deployment wave.
//
//	Example:
//
//	Wave 1: [0]
//	Wave 2: [1,2]
//	Wave 3: [3]
//
// 13. Minimum Number of Deployment Stages:
//
//	Process the graph level by level.
//
//	The number of topological levels is the minimum number of sequential
//	stages when independent services can run in parallel.
//
// 14. Determine Whether the Order Is Unique:
//
//	During Kahn's algorithm, check how many unprocessed nodes are currently
//	available in the queue.
//
//	If more than one service is ready at the same time, more than one valid
//	topological ordering exists.
//
// 15. Find Services Blocked by Cyclic Dependencies:
//
//	After Kahn's algorithm finishes, services whose indegree remains greater
//	than 0 are unresolved.
//
//	Important:
//
//	This unresolved set may include:
//
//	- Services directly inside a cycle.
//	- Services that depend on a cycle.
//
//	To identify the exact cycle members, use:
//
//	- DFS cycle reconstruction.
//	- Tarjan's strongly connected components algorithm.
//	- Kosaraju's strongly connected components algorithm.
//
// 16. Alien Dictionary:
//
//	Infer ordering constraints between characters by comparing adjacent words,
//	build a directed graph, and run topological sort.
//
// 17. Dependency-Aware Cluster Upgrade:
//
//	Model cluster components or add-ons as graph nodes.
//
//	Example:
//
//	CRD upgrade -> controller upgrade -> custom resource migration
//
//	Return a safe upgrade sequence and reject cyclic plans.
//
// 18. Data-Pipeline Stage Ordering:
//
//	Determine the execution order for extraction, validation, transformation,
//	aggregation, and publishing stages.
//
// Recognition pattern:
//
// When a question contains phrases such as:
//
//	- depends on
//	- prerequisite
//	- must happen before
//	- deployment order
//	- installation order
//	- startup order
//	- build order
//	- workflow ordering
//	- circular dependency
//	- execution stages
//	- parallel waves
//
// consider:
//
//	Directed graph
//	    +
//	Topological sort
//	    +
//	Cycle detection

// time complexity: O(V + E) -> each reachable vertex is processed once and each edge is examined once.
// space complexity: O(V + E) -> the graph representation stores vertices together with their adjacency edges.
```

```go
// Exact question: How can a sliding-window rate limiter allow at most `limit` requests during any `windowSeconds` period?
//
// Possible answer: Use `NewSlidingWindowLimiter`, `Allow` with a FIFO queue to process items in discovery order and record each result.
//
// Output format: Return the `(*SlidingWindowLimiter, error)` value from `NewSlidingWindowLimiter`; the function does not print the answer.
//
// Inline descriptions:
// - `limit` is the int input used by this example.
// - `windowSeconds` is the int input used by this example.
//
// Boundary checks:
// - `limit <= 0` handles the smallest valid state or recursive base case.
// - `windowSeconds <= 0` handles the smallest valid state or recursive base case.
// - `limiter.hasTimestamp && timestamp < limiter.lastTimestamp` decides whether the branch or loop should continue for the current input.
//
// Key variables:
// - `limit` is the int input used by this example.
// - `windowSeconds` is the int input used by this example.
// - `compacted` is a slice; indexes identify positions and elements hold the corresponding values.
// - `results` is a slice whose elements hold the ordered values produced or awaiting processing.
// - `cutoff` holds the intermediate value produced by `timestamp - limiter.windowSeconds`.
//
// Logic:
// 1. Create or use a slice so indexes identify positions and elements store their data or state.
// 2. Iterate through the required elements or states in the order shown.
// 3. Recursively reduce the current problem to smaller calls until a base condition is reached.
package main

import (
	"fmt"
	"sync"
)

// Question:
// Implement a sliding-window rate limiter.
//
// The rate limiter allows at most "limit" requests from a client during any
// "windowSeconds" period.
//
// Requests arrive with integer timestamps measured in seconds.
//
// For every request timestamp, return:
//
//	true  -> request is allowed
//	false -> request is rejected
//
// Input to function:
//
//	limit int
//		Maximum number of accepted requests within the active window.
//
//	windowSeconds int
//		Size of the sliding window in seconds.
//
//	timestamps []int
//		Request timestamps in non-decreasing order.
//
// Output from the function:
//
//	[]bool
//		One boolean result for each input request.
//
// Example:
//
//	limit = 3
//	windowSeconds = 10
//	timestamps = [1, 2, 3, 4, 11, 12]
//
// Output:
//
//	[true, true, true, false, true, true]
//
// Explanation:
//
//	At t=1:
//	    accepted timestamps = [1]
//	    result = true
//
//	At t=2:
//	    accepted timestamps = [1,2]
//	    result = true
//
//	At t=3:
//	    accepted timestamps = [1,2,3]
//	    result = true
//
//	At t=4:
//	    timestamps 1,2,3 are still in the active window
//	    result = false
//
//	At t=11:
//	    timestamp 1 has expired because:
//
//	        11 - 1 = 10
//
//	    This implementation treats requests exactly windowSeconds old as expired.
//
//	    active timestamps become [2,3]
//	    timestamp 11 is accepted
//
//	    active timestamps become [2,3,11]
//
// -----------------------------------------------------------------------------
//
// Answer:
//
// Approach:
//
// Use a queue containing timestamps of accepted requests that are still
// inside the active sliding window.
//
// For each request:
//
//  1. Calculate the cutoff timestamp:
//
//	    cutoff = currentTimestamp - windowSeconds
//
//  2. Remove accepted timestamps that are less than or equal to cutoff.
//
//  3. Count how many accepted timestamps remain.
//
//  4. If the count is already equal to limit, reject the request.
//
//  5. Otherwise, append the current timestamp and allow the request.
//
// Data structure used:
//
// Queue implemented using:
//
//	acceptedTimestamps []int
//	head               int
//
// Why use a queue?
//
// Timestamps expire in the same order in which they arrive.
//
// The oldest accepted timestamp is always at the front.
//
// Why use a head index?
//
// Repeatedly doing:
//
//	queue = queue[1:]
//
// changes the slice header, but the underlying array may continue retaining
// old values. Repeated element shifting would also be inefficient.
//
// Instead:
//
//	head
//
// identifies the first active timestamp.
//
// Questions to ask the interviewer:
//
//  1. Are timestamps guaranteed to arrive in non-decreasing order?
//
//  2. Does a request exactly windowSeconds old count inside or outside the window?
//
//	This implementation uses:
//
//	    (currentTimestamp-windowSeconds, currentTimestamp]
//
//	Therefore, a timestamp equal to the cutoff is expired.
//
//  3. Should rejected requests consume capacity?
//
//	This implementation stores only accepted requests.
//
//  4. Is the limiter per user, per IP, per API key, or global?
//
//	This example implements the state for one logical client.
//
//  5. Must the limiter be thread-safe?
//
//	This implementation uses a mutex so multiple goroutines may call Allow.
//
//  6. Must the limiter work across multiple servers?
//
//	This implementation is process-local.
//
//	A distributed implementation may use Redis, a central rate-limit service,
//	or token allocation among application instances.
//
//  7. What should happen if the distributed rate-limit backend is unavailable?
//
//	Possible policies:
//
//	    fail open:
//	        Allow requests to preserve availability.
//
//	    fail closed:
//	        Reject requests to preserve protection.
//
//  8. Can timestamps move backward because of clock adjustments?
//
//	This implementation rejects timestamps older than the previous timestamp.
//
//	In production, a monotonic clock should be preferred.
//
//  9. Is burst traffic allowed?
//
//	A sliding-window log precisely limits accepted requests but stores one
//	timestamp per accepted request.
//
//	A token bucket may be more memory-efficient and naturally supports bursts.
//
// Plan:
//
//  1. Validate the constructor inputs.
//  2. Store accepted timestamps in a queue.
//  3. Maintain a head index.
//  4. On every request, remove expired timestamps logically by advancing head.
//  5. Reject if active count reaches the limit.
//  6. Otherwise, append the timestamp.
//  7. Periodically compact the slice to release unused storage.
//
// Time complexity:
//
//	Amortized O(1) per request.
//
// Each accepted timestamp is:
//
//	- appended once
//	- removed once by advancing head
//
// Space complexity:
//
//	O(limit)
//
// Rejected requests are not stored.
//
// Race conditions:
//
// Without synchronization, two requests could both observe:
//
//	activeCount = limit - 1
//
// and both could be accepted, exceeding the configured limit.
//
// A mutex protects:
//
//	- acceptedTimestamps
//	- head
//	- lastTimestamp
//
// This mutex protects one limiter instance.
//
// A process-local mutex does not coordinate multiple servers.

type SlidingWindowLimiter struct {
	limit         int
	windowSeconds int

	// acceptedTimestamps holds timestamp VALUES.
	//
	// Example:
	//
	//	acceptedTimestamps = [1,2,3,11]
	//
	// These are actual request timestamps, not indexes.
	acceptedTimestamps []int

	// head holds an INDEX into acceptedTimestamps.
	//
	// Example:
	//
	//	acceptedTimestamps = [1,2,3,11]
	//	head = 1
	//
	// Logical active queue:
	//
	//	[2,3,11]
	//
	// acceptedTimestamps[0] is expired and ignored.
	head int

	// lastTimestamp holds the VALUE of the latest timestamp seen.
	//
	// It is used to detect out-of-order timestamps.
	lastTimestamp int

	// hasTimestamp distinguishes:
	//
	//	no timestamp received yet
//
// from:
//
//	lastTimestamp == 0
	hasTimestamp bool

	mu sync.Mutex
}

func NewSlidingWindowLimiter(
	limit int,
	windowSeconds int,
) (*SlidingWindowLimiter, error) {
	// Boundary check:
	//
	// A rate limit of zero or less cannot allow meaningful traffic.
	if limit <= 0 {
		return nil, fmt.Errorf("limit must be greater than zero")
	}

	// Boundary check:
	//
	// A non-positive time window is invalid.
	if windowSeconds <= 0 {
		return nil, fmt.Errorf("windowSeconds must be greater than zero")
	}

	return &SlidingWindowLimiter{
		limit:              limit,
		windowSeconds:      windowSeconds,
		acceptedTimestamps: make([]int, 0, limit),
	}, nil
}

func (limiter *SlidingWindowLimiter) Allow(timestamp int) (bool, error) {
	// Lock protects this limiter's mutable state.
	//
	// Without this lock, concurrent requests could exceed the limit.
	limiter.mu.Lock()
	defer limiter.mu.Unlock()

	// Boundary check:
	//
	// This implementation requires non-decreasing timestamps.
	//
	// Example invalid sequence:
	//
	//	10, 12, 9
	//
	// Processing timestamp 9 after timestamp 12 would invalidate the queue's
	// expiration ordering.
	if limiter.hasTimestamp && timestamp < limiter.lastTimestamp {
		return false, fmt.Errorf(
			"timestamps must be non-decreasing: received %d after %d",
			timestamp,
			limiter.lastTimestamp,
		)
	}

	limiter.lastTimestamp = timestamp
	limiter.hasTimestamp = true

	// cutoff is a timestamp VALUE.
	//
	// Example:
	//
	//	timestamp = 11
	//	windowSeconds = 10
	//	cutoff = 1
	//
	// Timestamps less than or equal to 1 are expired.
	cutoff := timestamp - limiter.windowSeconds

	// Remove expired timestamps by advancing the head index.
	for limiter.head < len(limiter.acceptedTimestamps) {
		// oldestTimestamp is a timestamp VALUE.
		//
		// limiter.head is an INDEX.
		oldestTimestamp := limiter.acceptedTimestamps[limiter.head]

		// Active window:
		//
		//	(cutoff, timestamp]
		//
		// Therefore, oldestTimestamp == cutoff is expired.
		if oldestTimestamp > cutoff {
			break
		}

		limiter.head++
	}

	// activeRequestCount is a COUNT.
	//
	// It is not an index and not a timestamp.
	//
	// Example:
	//
	//	acceptedTimestamps = [1,2,3,11]
	//	head = 1
	//
	//	activeRequestCount = 4 - 1 = 3
	activeRequestCount :=
		len(limiter.acceptedTimestamps) - limiter.head

	// Important iteration trace:
	//
	// limit = 3
	// windowSeconds = 10
	//
	// Request t=1:
	//
	//	cutoff = -9
	//	activeRequestCount = 0
	//	acceptedTimestamps becomes [1]
	//	result = true
	//
	// Request t=2:
	//
	//	cutoff = -8
	//	activeRequestCount = 1
	//	acceptedTimestamps becomes [1,2]
	//	result = true
	//
	// Request t=3:
	//
	//	cutoff = -7
	//	activeRequestCount = 2
	//	acceptedTimestamps becomes [1,2,3]
	//	result = true
	//
	// Request t=4:
	//
	//	cutoff = -6
	//	no timestamps expire
	//	activeRequestCount = 3
	//	result = false
	//
	//	The rejected timestamp 4 is not stored.
	//
	// Request t=11:
	//
	//	cutoff = 1
	//
	//	timestamp 1 <= cutoff
	//	head changes from 0 to 1
	//
	//	active logical timestamps = [2,3]
	//	activeRequestCount = 2
	//
	//	timestamp 11 is appended
	//	logical queue becomes [2,3,11]
	//	result = true

	if activeRequestCount >= limiter.limit {
		return false, nil
	}

	// Only accepted requests consume rate-limit capacity.
	limiter.acceptedTimestamps = append(
		limiter.acceptedTimestamps,
		timestamp,
	)

	// Memory compaction:
	//
	// Entries before head are logically expired but may still occupy memory.
	//
	// Compact only when:
	//
	//	- at least 1024 entries are expired
	//	- expired entries represent at least half the slice
	//
	// This avoids copying on every request.
	if limiter.head >= 1024 &&
		limiter.head*2 >= len(limiter.acceptedTimestamps) {

		activeTimestamps :=
			limiter.acceptedTimestamps[limiter.head:]

		compacted := make([]int, len(activeTimestamps))
		copy(compacted, activeTimestamps)

		limiter.acceptedTimestamps = compacted
		limiter.head = 0
	}

	return true, nil
}

func processRequests(
	limit int,
	windowSeconds int,
	timestamps []int,
) ([]bool, error) {
	limiter, err := NewSlidingWindowLimiter(
		limit,
		windowSeconds,
	)
	if err != nil {
		return nil, err
	}

	results := make([]bool, 0, len(timestamps))

	for _, timestamp := range timestamps {
		allowed, err := limiter.Allow(timestamp)
		if err != nil {
			return nil, err
		}

		results = append(results, allowed)
	}

	return results, nil
}

func main() {
	results, err := processRequests(
		3,
		10,
		[]int{1, 2, 3, 4, 11, 12},
	)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	fmt.Println(results)

	// Expected output:
	//
	//	[true true true false true true]
}

// Alternative questions that can be solved using this approach:
//
// 1. Limit login attempts:
//
//	Allow at most five failed login attempts per minute.
//
// 2. Kubernetes admission request throttling:
//
//	Limit mutation or validation requests per tenant.
//
// 3. API gateway rate limiting:
//
//	Limit requests by API key, user, IP address, or endpoint.
//
// 4. Alert suppression:
//
//	Allow at most one identical alert notification during a time window.
//
// 5. Restart-loop protection:
//
//	Prevent a controller from restarting the same workload too many times
//	within a configured interval.
//
// 6. Deployment frequency guard:
//
//	Allow at most N production deployments during a rolling time window.
//
// 7. Circuit-breaker event counting:
//
//	Open a circuit when failures exceed a threshold during a rolling window.
//
// 8. Abuse detection:
//
//	Detect clients producing too many requests during a short interval.
//
// Recognition pattern:
//
// When the question contains:
//
//	- at most N events
//	- within the last X seconds
//	- rolling window
//	- recent requests
//	- recent failures
//	- expire old events
//
// consider:
//
//	Queue
//	    +
//	Sliding window
//	    +
//	Timestamp expiration

// time complexity: O(1) amortized -> the snippet performs a fixed number of operations independent of input size.
// space complexity: O(limit) -> the sliding window stores at most the configured number of accepted requests.
```

```go
// Exact question: How can a fixed-capacity LRU cache support average O(1) `Get` and `Put` operations?
//
// Possible answer: Use `NewLRUCache`, `Get` with a map whose keys identify lookups and whose values hold the associated data.
//
// Output format: Return the `(*LRUCache, error)` value from `NewLRUCache`; the function does not print the answer.
//
// Inline descriptions:
// - `capacity` is the int input used by this example.
//
// Boundary checks:
// - `capacity <= 0` handles the smallest valid state or recursive base case.
// - `!exists` decides whether the branch or loop should continue for the current input.
// - `len(cache.items) > cache.capacity` decides whether the branch or loop should continue for the current input.
//
// Key variables:
// - `capacity` is the int input used by this example.
// - `head` holds the intermediate value produced by `&Node{}`.
// - `tail` holds the intermediate value produced by `&Node{}`.
// - `newNode` holds the intermediate value produced by `&Node`.
// - `leastRecentlyUsed` holds the intermediate value produced by `cache.tail.prev`.
//
// Logic:
// 1. Create or use a map to associate each lookup key with its stored value.
// 2. Recursively reduce the current problem to smaller calls until a base condition is reached.
// 3. Return the value produced after the state updates are complete.
package main

import (
	"fmt"
	"sync"
)

// Question:
// Implement a Least Recently Used cache.
//
// The cache has a fixed capacity.
//
// Operations:
//
//	Get(key)
//	    Return the value when the key exists.
//	    Mark the key as most recently used.
//
//	Put(key, value)
//	    Insert a new key or update an existing key.
//	    Mark the key as most recently used.
//
// When inserting a new key causes the cache to exceed capacity:
//
//	Remove the least recently used key.
//
// Both Get and Put must run in O(1) average time.
//
// Input:
//
//	capacity int
//		Maximum number of entries.
//
//	Get(key int)
//
//	Put(key int, value int)
//
// Output:
//
//	Get returns:
//
//	    value int
//	    found bool
//
// Example:
//
//	capacity = 2
//
//	Put(1,100)
//	Put(2,200)
//	Get(1)       -> 100, true
//	Put(3,300)   -> evicts key 2
//	Get(2)       -> 0, false
//	Get(3)       -> 300, true
//
// Usage order:
//
//	After Put(1,100):
//
//	    most recent -> [1] <- least recent
//
//	After Put(2,200):
//
//	    most recent -> [2,1] <- least recent
//
//	After Get(1):
//
//	    most recent -> [1,2] <- least recent
//
//	After Put(3,300):
//
//	    temporary order = [3,1,2]
//
//	Key 2 is least recently used and is evicted.
//
//	Final:
//
//	    [3,1]
//
// -----------------------------------------------------------------------------
//
// Answer:
//
// Approach:
//
// Combine:
//
//  1. Hash map
//  2. Doubly linked list
//
// Hash map:
//
//	map[key]*Node
//
// provides O(1) average lookup.
//
// Doubly linked list:
//
// tracks recency order.
//
// The node nearest the head is most recently used.
//
// The node nearest the tail is least recently used.
//
// Why not use only a hash map?
//
// A hash map does not maintain recency order.
//
// Why not use only a linked list?
//
// Finding a node by key would require O(n) traversal.
//
// Why a doubly linked list instead of a singly linked list?
//
// A doubly linked list allows removing a known node in O(1).
//
// With a singly linked list, removing a node generally requires finding its
// previous node first.
//
// Questions to ask the interviewer:
//
//  1. What should happen when capacity is zero or negative?
//
//	This implementation returns an error.
//
//  2. Can key or value be zero or negative?
//
//	Yes. The found boolean distinguishes a missing key from a stored zero.
//
//  3. Does updating an existing key count as recent usage?
//
//	This implementation treats an update as recent usage.
//
//  4. Does Get modify the cache order?
//
//	Yes. A successful Get moves the entry to the front.
//
//  5. Must the cache be thread-safe?
//
//	This implementation uses a mutex.
//
//  6. Is an eviction callback required?
//
//	A production cache may need to close files, release connections, or
//	update metrics when an item is evicted.
//
//  7. Is time-based expiration also required?
//
//	LRU and TTL are separate policies.
//
//	A cache may combine recency ordering with expiration timestamps.
//
//  8. Should the cache support generic keys and values?
//
//	This interview implementation uses integers for clarity.
//
//  9. What should happen if a cached value is updated while another goroutine
//	reads it?
//
//	The mutex serializes access.
//
// Plan:
//
//  1. Create Node with key, value, prev, and next.
//  2. Create a map from key to node.
//  3. Create dummy head and tail sentinel nodes.
//  4. Keep most recently used nodes after head.
//  5. Keep least recently used nodes before tail.
//  6. Get:
//
//	    - look up node in map
//	    - remove it from its current position
//	    - add it after head
//
//  7. Put existing:
//
//	    - update value
//	    - move node after head
//
//  8. Put new:
//
//	    - create node
//	    - add to map
//	    - add after head
//	    - evict tail.prev when capacity is exceeded
//
// Time complexity:
//
//	Get: O(1) average
//	Put: O(1) average
//
// Space complexity:
//
//	O(capacity)
//
// Race conditions:
//
// Both Get and Put modify recency order.
//
// Even Get is a write operation because it moves a node.
//
// Therefore, an RWMutex does not provide much benefit for normal successful
// reads. This implementation uses a regular mutex.
//
// Sentinel nodes:
//
// head and tail are dummy nodes.
//
// They simplify boundary cases:
//
//	- empty cache
//	- one-element cache
//	- removing first real node
//	- removing last real node
//
// Linked-list invariants:
//
//	head.prev == nil
//	tail.next == nil
//	head.next is either tail or the most recent real node
//	tail.prev is either head or the least recent real node
//	node.prev.next == node
//	node.next.prev == node

type Node struct {
	// key and value hold actual cache VALUES.
	key   int
	value int

	// prev and next hold pointers to adjacent nodes.
	prev *Node
	next *Node
}

type LRUCache struct {
	capacity int

	// items maps a key VALUE to the node containing that key.
	//
	// Example:
	//
	//	items[3] -> Node{key:3, value:300}
	items map[int]*Node

	// head and tail are sentinel nodes.
	//
	// Recency order:
	//
	//	head <-> most recent ... least recent <-> tail
	head *Node
	tail *Node

	mu sync.Mutex
}

func NewLRUCache(capacity int) (*LRUCache, error) {
	// Boundary check:
	//
	// A zero-capacity cache cannot retain values.
	if capacity <= 0 {
		return nil, fmt.Errorf(
			"capacity must be greater than zero",
		)
	}

	head := &Node{}
	tail := &Node{}

	// Empty-list state:
	//
	//	head <-> tail
	head.next = tail
	tail.prev = head

	return &LRUCache{
		capacity: capacity,
		items:    make(map[int]*Node, capacity),
		head:     head,
		tail:     tail,
	}, nil
}

func (cache *LRUCache) Get(key int) (int, bool) {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	// node is a pointer to the node VALUE stored in the map.
	node, exists := cache.items[key]
	if !exists {
		// Returning found=false distinguishes a missing key from a stored value 0.
		return 0, false
	}

	// A successful access makes this node most recently used.
	cache.removeNode(node)
	cache.addAfterHead(node)

	return node.value, true
}

func (cache *LRUCache) Put(key int, value int) {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	// Case 1:
	//
	// The key already exists.
	//
	// Update the value and move the node to the most-recent position.
	if existingNode, exists := cache.items[key]; exists {
		existingNode.value = value

		cache.removeNode(existingNode)
		cache.addAfterHead(existingNode)
		return
	}

	// Case 2:
	//
	// The key is new.
	newNode := &Node{
		key:   key,
		value: value,
	}

	// Store the node in the map for O(1) lookup.
	cache.items[key] = newNode

	// New entries are most recently used.
	cache.addAfterHead(newNode)

	// Boundary check:
	//
	// Evict exactly one node when capacity is exceeded.
	if len(cache.items) > cache.capacity {
		// tail.prev is the least recently used real node.
		//
		// It cannot be head here because:
		//
		//	- capacity is positive
		//	- len(items) exceeded capacity
		leastRecentlyUsed := cache.tail.prev

		cache.removeNode(leastRecentlyUsed)
		delete(cache.items, leastRecentlyUsed.key)
	}
}

func (cache *LRUCache) removeNode(node *Node) {
	// Defensive boundary check:
	//
	// Sentinel nodes must never be removed.
	if node == nil || node == cache.head || node == cache.tail {
		return
	}

	// previousNode and nextNode are node POINTERS.
	previousNode := node.prev
	nextNode := node.next

	// Before:
	//
	//	previousNode <-> node <-> nextNode
	//
	// After:
	//
	//	previousNode <-> nextNode
	previousNode.next = nextNode
	nextNode.prev = previousNode

	// Clear stale pointers.
	//
	// This is not required for garbage collection, but it makes accidental
// reuse easier to detect.
	node.prev = nil
	node.next = nil
}

func (cache *LRUCache) addAfterHead(node *Node) {
	// currentMostRecent is a node POINTER.
	//
	// It may point to:
	//
	//	- the current most-recent real node
	//	- tail when the cache is empty
	currentMostRecent := cache.head.next

	// Before:
	//
	//	head <-> currentMostRecent
	//
	// After:
	//
	//	head <-> node <-> currentMostRecent
	cache.head.next = node
	node.prev = cache.head

	node.next = currentMostRecent
	currentMostRecent.prev = node
}

func main() {
	cache, err := NewLRUCache(2)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	// Iteration 1:
	//
	// Put(1,100)
	//
	// items:
	//
	//	1 -> node1
	//
	// list:
	//
	//	head <-> 1 <-> tail
	cache.Put(1, 100)

	// Iteration 2:
	//
	// Put(2,200)
	//
	// items:
	//
	//	1 -> node1
	//	2 -> node2
	//
	// list:
	//
	//	head <-> 2 <-> 1 <-> tail
	//
	// 2 is most recent.
	// 1 is least recent.
	cache.Put(2, 200)

	// Iteration 3:
	//
	// Get(1)
	//
	// node1 is removed from its old position and added after head.
	//
	// list changes:
	//
	//	head <-> 2 <-> 1 <-> tail
	//
	// to:
	//
	//	head <-> 1 <-> 2 <-> tail
	//
	// 1 becomes most recent.
	// 2 becomes least recent.
	value, found := cache.Get(1)
	fmt.Println(value, found)

	// Expected:
	//
	//	100 true

	// Iteration 4:
	//
	// Put(3,300)
	//
	// Temporary list:
	//
	//	head <-> 3 <-> 1 <-> 2 <-> tail
	//
	// len(items) becomes 3.
	// capacity is 2.
	//
	// tail.prev is node2.
	// Key 2 is evicted.
	//
	// Final list:
	//
	//	head <-> 3 <-> 1 <-> tail
	cache.Put(3, 300)

	value, found = cache.Get(2)
	fmt.Println(value, found)

	// Expected:
	//
	//	0 false

	value, found = cache.Get(3)
	fmt.Println(value, found)

	// Expected:
	//
	//	300 true
}

// Alternative questions that can be solved using this approach:
//
// 1. Metadata cache:
//
//	Cache recently requested Kubernetes object metadata.
//
// 2. DNS cache:
//
//	Retain recently resolved hostnames with bounded memory.
//
// 3. Connection cache:
//
//	Keep recently used reusable connections.
//
// 4. Page replacement:
//
//	Simulate operating-system LRU page eviction.
//
// 5. Recently used configuration cache:
//
//	Cache parsed configurations and evict inactive entries.
//
// 6. Bounded query-result cache:
//
//	Cache expensive query results while controlling memory usage.
//
// 7. API response cache:
//
//	Store recently accessed responses.
//
// 8. LRU cache with TTL:
//
//	Add expiration timestamps and remove expired entries before returning them.
//
// 9. LRU cache with eviction callback:
//
//	Run cleanup logic whenever an entry is removed.
//
// 10. Least Frequently Used cache:
//
//	This requires a different ordering structure:
//
//	    frequency buckets
//	    +
//	    linked lists
//	    +
//	    key-to-node map
//
// Recognition pattern:
//
// When the question requires:
//
//	- O(1) lookup
//	- fixed capacity
//	- remove least recently used
//	- maintain access order
//
// consider:
//
//	Hash map
//	    +
//	Doubly linked list

// time complexity: O(1) -> the snippet performs a fixed number of operations independent of input size.
// space complexity: O(capacity) -> the cache stores at most one map entry and one list node per capacity slot.
```

```go
// Exact question: Given maintenance windows as `[start, end]`, return a new list in which all overlapping or touching windows are merged.
//
// Possible answer: Use `mergeWindows` to sort the values first, then scan or compare them in the required order.
//
// Output format: Return the `([][]int, error)` value from `mergeWindows`; the function does not print the answer.
//
// Inline descriptions:
// - `windows` is a two-dimensional slice: row and column indexes identify positions, and each cell stores a value.
//
// Boundary checks:
// - `len(windows) == 0` handles empty input before any element is accessed.
// - `len(window) != 2` decides whether the branch or loop should continue for the current input.
// - `start > end` keeps indexes or pointers within the portion of the input still being processed.
//
// Key variables:
// - `windows` is a two-dimensional slice: row and column indexes identify positions, and each cell stores a value.
// - `sortedWindows` is a two-dimensional table; indexes select a row and column, while cells hold computed values.
// - `merged` is a two-dimensional table; indexes select a row and column, while cells hold computed values.
// - `start` holds the intermediate value produced by `window[0]`.
// - `end` holds the intermediate value produced by `window[1]`.
//
// Logic:
// 1. Create or use a slice so indexes identify positions and elements store their data or state.
// 2. Iterate through the required elements or states in the order shown.
// 3. Recursively reduce the current problem to smaller calls until a base condition is reached.
package main

import (
	"fmt"
	"sort"
)

// Question:
// Merge overlapping maintenance windows.
//
// Each maintenance window is represented as:
//
//	[start, end]
//
// Return a new list where all overlapping or touching windows are merged.
//
// Input to function:
//
//	windows [][]int
//
// Each inner slice must contain exactly:
//
//	[start, end]
//
// Output from function:
//
//	[][]int
//
// A list of non-overlapping intervals sorted by start time.
//
// Example:
//
//	Input:
//
//	[
//	    [1,4],
//	    [2,6],
//	    [8,10],
//	    [9,12],
//	]
//
//	Output:
//
//	[
//	    [1,6],
//	    [8,12],
//	]
//
// Explanation:
//
//	[1,4] overlaps [2,6], so they become [1,6].
//
//	[8,10] overlaps [9,12], so they become [8,12].
//
// Touching interval example:
//
//	[1,4]
//	[4,7]
//
// This implementation treats them as overlapping and returns:
//
//	[1,7]
//
// -----------------------------------------------------------------------------
//
// Answer:
//
// Approach:
//
//  1. Validate all intervals.
//  2. Sort intervals by start time.
//  3. Add the first interval to the result.
//  4. For each remaining interval:
//
//	    - Compare it with the last merged interval.
//	    - If they overlap, extend the last merged interval.
//	    - Otherwise, append a new interval.
//
// Why sort first?
//
// Before sorting, an overlapping interval may appear anywhere.
//
// After sorting by start time, the current interval can overlap only with the
// last interval in the merged result.
//
// Data structure:
//
// A slice stores the merged intervals.
//
// No special graph or tree is needed.
//
// Questions to ask the interviewer:
//
//  1. Are intervals closed, open, or half-open?
//
//	This implementation assumes closed intervals:
//
//	    [start, end]
//
//  2. Should touching intervals merge?
//
//	This implementation merges:
//
//	    [1,4] and [4,7]
//
//	into:
//
//	    [1,7]
//
//  3. Can start be greater than end?
//
//	This implementation treats that as invalid input.
//
//  4. Can intervals contain negative values?
//
//	Yes. Negative timestamps or offsets are valid as long as start <= end.
//
//  5. Can duplicate intervals appear?
//
//	Yes. They are naturally merged.
//
//  6. Should the function modify the input?
//
//	This implementation creates a copy before sorting, so the caller's input
//	order is preserved.
//
//  7. What should empty input return?
//
//	This implementation returns an empty slice.
//
//  8. Are timestamps integers, Unix timestamps, or time.Time values?
//
//	The algorithm is identical for any sortable time representation.
//
//  9. Can end values overflow when computing durations?
//
//	This algorithm does not subtract interval endpoints, so it avoids duration
//	overflow during merging.
//
// Plan:
//
//  1. Handle empty input.
//  2. Validate interval length.
//  3. Validate start <= end.
//  4. Copy intervals.
//  5. Sort by start, then end.
//  6. Add the first sorted interval.
//  7. Iterate through the remainder.
//  8. Merge when current.start <= last.end.
//  9. Otherwise append a separate interval.
//
// Time complexity:
//
//	O(n log n)
//
// Sorting dominates the complexity.
//
// The merge scan is O(n).
//
// Space complexity:
//
//	O(n)
//
// A copied input and merged result are stored.
//
// Race conditions:
//
// The function uses local slices.
//
// Separate calls are safe concurrently.
//
// The function copies input values before sorting, so it does not mutate the
// caller's interval order.
//
// However, the caller must not concurrently modify the nested input slices
// while this function is copying them.

func mergeWindows(windows [][]int) ([][]int, error) {
	// Boundary check:
	//
	// No maintenance windows means no merged windows.
	if len(windows) == 0 {
		return [][]int{}, nil
	}

	// sortedWindows stores copies of interval VALUES.
	//
	// Example:
	//
	//	sortedWindows[0] = []int{1,4}
	//
	// Index 0 identifies the first interval after sorting.
	//
	// Values 1 and 4 are start and end values.
	sortedWindows := make([][]int, 0, len(windows))

	for intervalIndex, window := range windows {
		// Boundary check:
		//
		// Every interval must have exactly two values.
		if len(window) != 2 {
			return nil, fmt.Errorf(
				"window at index %d must contain start and end",
				intervalIndex,
			)
		}

		start := window[0]
		end := window[1]

		// Boundary check:
		//
		// A maintenance window cannot end before it starts.
		if start > end {
			return nil, fmt.Errorf(
				"invalid window at index %d: start %d is after end %d",
				intervalIndex,
				start,
				end,
			)
		}

		// Copy values to prevent modifications to the caller's inner slices.
		sortedWindows = append(
			sortedWindows,
			[]int{start, end},
		)
	}

	sort.Slice(sortedWindows, func(i, j int) bool {
		// i and j are interval INDEXES.
		//
		// sortedWindows[i][0] is a start VALUE.
		// sortedWindows[i][1] is an end VALUE.

		// Primary sort:
		//
		// Smaller start value first.
		if sortedWindows[i][0] != sortedWindows[j][0] {
			return sortedWindows[i][0] < sortedWindows[j][0]
		}

		// Tie-break:
		//
		// Smaller end value first for deterministic output.
		return sortedWindows[i][1] < sortedWindows[j][1]
	})

	// merged stores interval VALUES.
	merged := make([][]int, 0, len(sortedWindows))

	// Add a copy of the first sorted interval.
	merged = append(
		merged,
		[]int{
			sortedWindows[0][0],
			sortedWindows[0][1],
		},
	)

	for currentIndex := 1;
		currentIndex < len(sortedWindows);
		currentIndex++ {

		// current is the current interval VALUE.
		current := sortedWindows[currentIndex]

		// lastMergedIndex is an INDEX into merged.
		lastMergedIndex := len(merged) - 1

		// lastMerged is the most recent merged interval VALUE.
		lastMerged := merged[lastMergedIndex]

		currentStart := current[0]
		currentEnd := current[1]

		lastStart := lastMerged[0]
		lastEnd := lastMerged[1]

		// lastStart is retained for explanation and debugging.
		// The overlap decision mainly uses lastEnd.
		_ = lastStart

		// Important iteration trace:
		//
		// Sorted input:
		//
		//	[1,4], [2,6], [8,10], [9,12]
		//
		// Initial:
		//
		//	merged = [[1,4]]
		//
		// Iteration 1:
		//
		//	currentIndex = 1
		//	current = [2,6]
		//	lastMerged = [1,4]
		//
		//	currentStart = 2
		//	lastEnd = 4
		//
		//	2 <= 4, so overlap exists.
		//
		//	lastMerged end becomes max(4,6) = 6
		//
		//	merged = [[1,6]]
		//
		// Iteration 2:
		//
		//	current = [8,10]
		//	lastMerged = [1,6]
		//
		//	8 > 6, so no overlap.
		//
		//	merged = [[1,6],[8,10]]
		//
		// Iteration 3:
		//
		//	current = [9,12]
		//	lastMerged = [8,10]
		//
		//	9 <= 10, so overlap exists.
		//
		//	lastMerged end becomes max(10,12) = 12
		//
		//	merged = [[1,6],[8,12]]

		// Overlap or touching check:
		//
		// current starts before or exactly when last merged interval ends.
		if currentStart <= lastEnd {
			// Extend only when the current interval ends later.
			if currentEnd > lastEnd {
				merged[lastMergedIndex][1] = currentEnd
			}

			continue
		}

		// No overlap.
		//
		// Append a copied interval.
		merged = append(
			merged,
			[]int{currentStart, currentEnd},
		)
	}

	return merged, nil
}

func main() {
	input := [][]int{
		{1, 4},
		{2, 6},
		{8, 10},
		{9, 12},
	}

	merged, err := mergeWindows(input)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	fmt.Println(merged)

	// Expected output:
	//
	//	[[1 6] [8 12]]
}

// Alternative questions that can be solved using this approach:
//
// 1. Merge incident periods:
//
//	Combine overlapping periods during which a service was unavailable.
//
// 2. Merge on-call escalation windows:
//
//	Combine overlapping coverage schedules.
//
// 3. Merge cluster maintenance windows:
//
//	Find the actual total periods affected by multiple maintenance plans.
//
// 4. Employee free time:
//
//	Merge all busy intervals, then find gaps.
//
// 5. Insert interval:
//
//	Insert one interval into an already sorted list and merge overlaps.
//
// 6. Determine total downtime:
//
//	Merge outage intervals, then sum:
//
//	    end - start
//
// 7. Meeting-room conflict detection:
//
//	Sort intervals and check whether:
//
//	    current.start < previous.end
//
// 8. Minimum number of meeting rooms:
//
//	Use sorting plus a min heap of end times.
//
// 9. Maximum concurrent deployments:
//
//	Use sweep-line events:
//
//	    +1 at start
//	    -1 at end
//
// 10. Kubernetes disruption windows:
//
//	Merge periods when multiple dependent systems are unavailable.
//
// 11. IP address range merging:
//
//	Convert addresses to sortable integers and merge overlapping ranges.
//
// 12. Log-retention range compaction:
//
//	Merge adjacent or overlapping timestamp ranges.
//
// Recognition pattern:
//
// When the question contains:
//
//	- intervals
//	- start and end
//	- overlapping ranges
//	- maintenance periods
//	- meeting times
//	- outage windows
//	- merge schedules
//
// consider:
//
//	Sort by start time
//	    +
//	Compare with last merged interval

// time complexity: O(n log n) -> sorting dominates the remaining linear scan.
// space complexity: O(n) -> the auxiliary slice, map, table, queue, or returned collection can grow with `n`.
```

```go
// Exact question: Given failure events, return the top `k` services producing the most failures.
//
// Possible answer: Use `Len`, `Less` with a stack to process the most recently added item first.
//
// Output format: Return an `int` value from `Len`; the function does not print the answer.
//
// Inline descriptions:
// - The comments in this preface describe how the important expressions and state changes are used.
//
// Boundary checks:
// - `h[i].count != h[j].count` keeps indexes or pointers within the portion of the input still being processed.
// - `k < 0` rejects or terminates a state that has moved below the valid range.
// - `k == 0 || len(failures) == 0` handles empty input before any element is accessed.
//
// Key variables:
// - `counts` is a map whose keys are string lookups and whose values are stored int results.
// - `item` holds the intermediate value produced by `value.(ServiceFailure`.
// - `oldHeap` holds the intermediate value produced by `*h`.
// - `lastIndex` holds the intermediate value produced by `len(oldHeap) - 1`.
// - `failureHeap` holds the intermediate value produced by `&MinFailureHeap{}`.
//
// Logic:
// 1. Create or use a map to associate each lookup key with its stored value.
// 2. Create or use a slice so indexes identify positions and elements store their data or state.
// 3. Iterate through the required elements or states in the order shown.
package main

import (
	"container/heap"
	"fmt"
	"sort"
)

// Question:
// Return the top K services producing the most failures.
//
// Input:
//
//	failures []string
//		Each value is the service name associated with one failure event.
//
//	k int
//		Maximum number of services to return.
//
// Output:
//
//	[]ServiceFailure
//
// Each result contains:
//
//	service name
//	failure count
//
// Results must be sorted by:
//
//  1. Higher failure count first.
//  2. Alphabetically smaller service name first when counts are equal.
//
// Example:
//
//	Input:
//
//	failures = [
//	    "api",
//	    "worker",
//	    "api",
//	    "database",
//	    "worker",
//	    "api",
//	    "cache",
//	]
//
//	k = 2
//
//	Output:
//
//	[
//	    {service:"api", count:3},
//	    {service:"worker", count:2},
//	]
//
// -----------------------------------------------------------------------------
//
// Answer:
//
// Approach:
//
// Use two main data structures:
//
//  1. Hash map for frequency counting.
//  2. Min heap containing at most K candidates.
//
// Step 1:
//
// Count failures:
//
//	api      -> 3
//	worker   -> 2
//	database -> 1
//	cache    -> 1
//
// Step 2:
//
// Iterate over unique services.
//
// Push each service and count into a min heap.
//
// When heap size exceeds K:
//
//	Remove the least important entry.
//
// The heap therefore retains only the K best candidates.
//
// Why use a hash map?
//
// It provides O(1) average-time frequency updates.
//
// Why use a min heap?
//
// The least important retained candidate remains at the root.
//
// When a better candidate is added and heap size becomes K+1, remove the root.
//
// Why not sort every service?
//
// Sorting all unique services costs:
//
//	O(m log m)
//
// A size-K heap costs:
//
//	O(m log k)
//
// where m is the number of unique services.
//
// The heap is useful when:
//
//	k << m
//
// Questions to ask the interviewer:
//
//  1. What is the tie-breaking rule?
//
//	This implementation uses alphabetical order.
//
//  2. Can k be larger than the number of unique services?
//
//	Yes. Return all unique services.
//
//  3. What should k=0 return?
//
//	An empty result.
//
//  4. Are empty service names valid?
//
//	This implementation rejects them.
//
//  5. Are service names case-sensitive?
//
//	This implementation treats "API" and "api" as different services.
//
//  6. Must output be deterministic?
//
//	Yes. Map iteration order in Go is not deterministic, so the final result
//	is explicitly sorted.
//
//  7. Can events arrive continuously?
//
//	For streaming data, possible approaches include:
//
//	    - periodically recompute
//	    - maintain an indexed heap
//	    - use Count-Min Sketch
//	    - use heavy-hitter algorithms
//
//  8. Must the solution work across millions of hosts?
//
//	Use hierarchical aggregation:
//
//	    node -> cluster -> region -> global
//
//  9. Is the time window bounded?
//
//	This function counts the entire input.
//
//	For recent failures only, combine frequency counting with a sliding window.
//
// Plan:
//
//  1. Validate k.
//  2. Count every service using a map.
//  3. Initialize a min heap.
//  4. Push every unique service.
//  5. If heap size exceeds k, remove one entry.
//  6. Extract remaining heap entries.
//  7. Sort final output deterministically.
//
// Time complexity:
//
// Let:
//
//	n = number of failure events
//	m = number of unique services
//
// Frequency counting:
//
//	O(n)
//
// Heap processing:
//
//	O(m log k)
//
// Final sorting:
//
//	O(k log k)
//
// Total:
//
//	O(n + m log k + k log k)
//
// Space complexity:
//
//	O(m + k)
//
// Race conditions:
//
// This function uses local state and is safe for separate concurrent calls.
//
// The caller must not mutate the failures slice concurrently.
//
// For a continuously updated shared counter map:
//
//	- use a mutex
//	- use sharded counters
//	- use atomic counters for fixed service identities
//	- aggregate events through a single consumer goroutine

type ServiceFailure struct {
	// service is the service-name VALUE.
	service string

	// count is the number of observed failure events.
	count int
}

// MinFailureHeap implements heap.Interface.
//
// The least desirable retained candidate appears at index 0.
//
// "Less desirable" means:
//
//  1. Lower count.
//  2. For equal count, alphabetically larger service name.
//
// Why alphabetically larger at the root during ties?
//
// Final desired ordering says alphabetically smaller names are preferred.
//
// Therefore, alphabetically larger names should be evicted first.
type MinFailureHeap []ServiceFailure

func (h MinFailureHeap) Len() int {
	return len(h)
}

func (h MinFailureHeap) Less(i, j int) bool {
	// i and j are heap INDEXES.
	//
	// h[i].count and h[j].count are failure-count VALUES.

	if h[i].count != h[j].count {
		return h[i].count < h[j].count
	}

	// Equal counts:
	//
	// Alphabetically larger service is considered less desirable.
	return h[i].service > h[j].service
}

func (h MinFailureHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *MinFailureHeap) Push(value any) {
	// value contains a ServiceFailure VALUE.
	item := value.(ServiceFailure)
	*h = append(*h, item)
}

func (h *MinFailureHeap) Pop() any {
	oldHeap := *h

	// lastIndex is an INDEX.
	lastIndex := len(oldHeap) - 1

	// item is a ServiceFailure VALUE.
	item := oldHeap[lastIndex]

	// Remove final slice element.
	*h = oldHeap[:lastIndex]

	return item
}

func topKFailingServices(
	failures []string,
	k int,
) ([]ServiceFailure, error) {
	// Boundary check:
	//
	// Negative K is invalid.
	if k < 0 {
		return nil, fmt.Errorf("k cannot be negative")
	}

	// Boundary check:
	//
	// No requested results or no input events.
	if k == 0 || len(failures) == 0 {
		return []ServiceFailure{}, nil
	}

	// counts maps service-name VALUES to failure-count VALUES.
	//
	// Example:
	//
	//	counts["api"] = 3
	counts := make(map[string]int)

	for eventIndex, service := range failures {
		// Boundary check:
		//
		// Empty service names are treated as malformed events.
		if service == "" {
			return nil, fmt.Errorf(
				"empty service name at failure index %d",
				eventIndex,
			)
		}

		counts[service]++
	}

	failureHeap := &MinFailureHeap{}
	heap.Init(failureHeap)

	// Important note:
	//
	// Go map iteration order is intentionally non-deterministic.
	//
	// The heap retains the correct top K because its comparison rules are
// deterministic, but extraction order is not presentation order.
	//
	// Therefore, the final result is explicitly sorted.
	for service, count := range counts {
		item := ServiceFailure{
			service: service,
			count:   count,
		}

		heap.Push(failureHeap, item)

		// Keep at most K items.
		if failureHeap.Len() > k {
			heap.Pop(failureHeap)
		}
	}

	// result length is min(k, number of unique services).
	result := make(
		[]ServiceFailure,
		0,
		failureHeap.Len(),
	)

	for failureHeap.Len() > 0 {
		item := heap.Pop(failureHeap).(ServiceFailure)
		result = append(result, item)
	}

	// Final desired ordering:
	//
	//	- larger count first
	//	- alphabetically smaller service first on ties
	sort.Slice(result, func(i, j int) bool {
		if result[i].count != result[j].count {
			return result[i].count > result[j].count
		}

		return result[i].service < result[j].service
	})

	return result, nil
}

func main() {
	failures := []string{
		"api",
		"worker",
		"api",
		"database",
		"worker",
		"api",
		"cache",
	}

	// Frequency-counting iterations:
	//
	// Event 1:
	//
	//	service = "api"
	//	counts["api"] = 1
	//
	// Event 2:
	//
	//	service = "worker"
	//	counts["worker"] = 1
	//
	// Event 3:
	//
	//	service = "api"
	//	counts["api"] changes from 1 to 2
	//
	// Event 4:
	//
	//	service = "database"
	//	counts["database"] = 1
	//
	// Event 5:
	//
	//	service = "worker"
	//	counts["worker"] changes from 1 to 2
	//
	// Event 6:
	//
	//	service = "api"
	//	counts["api"] changes from 2 to 3
	//
	// Final counts:
	//
	//	api      = 3
	//	worker   = 2
	//	database = 1
	//	cache    = 1
	//
	// For k=2, only api and worker remain in the heap after all candidates
// have been considered.

	result, err := topKFailingServices(failures, 2)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	for _, item := range result {
		fmt.Printf(
			"%s: %d\n",
			item.service,
			item.count,
		)
	}

	// Expected output:
	//
	//	api: 3
	//	worker: 2
}

// Alternative questions that can be solved using this approach:
//
// 1. Top K most frequently accessed APIs.
//
// 2. Top K Kubernetes namespaces generating errors.
//
// 3. Top K nodes with the most pod failures.
//
// 4. Top K containers producing the most logs.
//
// 5. Top K alerts by occurrence count.
//
// 6. Top K tenants consuming the most resources.
//
// 7. Top K HTTP status codes.
//
// 8. Top K frequent words.
//
// 9. Top K requested URLs.
//
// 10. Top K clusters with failed deployments.
//
// 11. Top K slowest endpoints:
//
//	Store latency as the heap priority rather than frequency.
//
// 12. K largest values:
//
//	Use a min heap of size K without a frequency map.
//
// 13. K smallest values:
//
//	Use a max heap of size K.
//
// 14. Streaming heavy hitters:
//
//	Use approximate counting when exact global counts are too expensive.
//
// Recognition pattern:
//
// When the question contains:
//
//	- top K
//	- most frequent
//	- highest count
//	- largest K
//	- most common
//	- retain only K candidates
//
// consider:
//
//	Hash map for counting
//	    +
//	Min heap of size K

// time complexity: O(n + m log k + k log k) -> counting costs `O(n)`, maintaining a size-`k` heap costs `O(m log k)`, and sorting the winners costs `O(k log k)`.
// space complexity: O(m + k) -> the frequency map stores up to `m` services and the heap/result stores up to `k` entries.
```

```go
// Exact question: Find the shortest number of network hops from a starting service to every reachable service in an unweighted graph.
//
// Possible answer: Use `shortestHopDistances` with a FIFO queue to process items in discovery order and record each result.
//
// Output format: Return the `([]int, error)` value from `shortestHopDistances`; the function does not print the answer.
//
// Inline descriptions:
// - `graph` is a two-dimensional slice: row and column indexes identify positions, and each cell stores a value.
// - `start` is the int input used by this example.
//
// Boundary checks:
// - `nodeCount == 0` handles the smallest valid state or recursive base case.
// - `start < 0 || start >= nodeCount` rejects or terminates a state that has moved below the valid range.
// - `head < len(queue)` decides whether the branch or loop should continue for the current input.
//
// Key variables:
// - `graph` is a two-dimensional slice: row and column indexes identify positions, and each cell stores a value.
// - `start` is the int input used by this example.
// - `distance` is indexed by a state and stores the computed answer for that state.
// - `nodeCount` holds the intermediate value produced by `len(graph`.
// - `queue` is a slice whose elements hold the ordered values produced or awaiting processing.
//
// Logic:
// 1. Create or use a slice so indexes identify positions and elements store their data or state.
// 2. Iterate through the required elements or states in the order shown.
// 3. Recursively reduce the current problem to smaller calls until a base condition is reached.
package main

import (
	"fmt"
)

// Question:
// Find the shortest number of network hops from a starting service to every
// reachable service in an unweighted dependency or communication graph.
//
// The graph contains services numbered from 0 to n-1.
//
// graph[i] contains the direct neighbours of service i.
//
// Input:
//
//	graph [][]int
//		Adjacency list.
//
//	start int
//		Starting service ID.
//
// Output:
//
//	[]int
//
// distance[i] contains the shortest number of edges from start to service i.
//
// Use:
//
//	-1
//
// when service i is unreachable.
//
// Example:
//
//	graph = [
//	    [1,2],
//	    [0,3],
//	    [0,3,4],
//	    [1,2,4],
//	    [2,3],
//	]
//
//	start = 0
//
// Output:
//
//	[0,1,1,2,2]
//
// Explanation:
//
//	distance[0] = 0
//	    Start node requires zero hops.
//
//	distance[1] = 1
//	    0 -> 1
//
//	distance[2] = 1
//	    0 -> 2
//
//	distance[3] = 2
//	    0 -> 1 -> 3
//
//	distance[4] = 2
//	    0 -> 2 -> 4
//
// -----------------------------------------------------------------------------
//
// Answer:
//
// Approach:
//
// Use Breadth-First Search.
//
// BFS explores nodes level by level:
//
//	Level 0:
//	    start
//
//	Level 1:
//	    direct neighbours
//
//	Level 2:
//	    neighbours of level-1 nodes
//
// In an unweighted graph, the first time BFS reaches a node is through the
// minimum number of edges.
//
// Data structures:
//
//	queue []int
//	    Holds service ID values waiting to be processed.
//
//	distance []int
//	    distance[i] holds the shortest hop count to service i.
//
// Why use BFS?
//
// BFS guarantees shortest paths in an unweighted graph.
//
// Why not DFS?
//
// DFS explores one branch deeply before exploring alternatives.
//
// The first path found by DFS is not necessarily the shortest.
//
// Questions to ask the interviewer:
//
//  1. Is the graph directed or undirected?
//
//	This function works with either representation.
//
//	For an undirected graph, both directions must appear in the adjacency list.
//
//  2. Are all service IDs guaranteed to be valid?
//
//	This implementation validates every neighbour.
//
//  3. Can the graph contain self-loops?
//
//	Yes. The visited check prevents repeated processing.
//
//  4. Can the graph contain duplicate edges?
//
//	Yes. The distance check prevents duplicate queue insertion.
//
//  5. What should unreachable nodes contain?
//
//	This implementation uses -1.
//
//  6. Do we need only distances, or also actual paths?
//
//	For actual paths, maintain:
//
//	    parent[child] = current
//
//  7. Are edges weighted?
//
//	If edges have non-negative weights, use Dijkstra's algorithm.
//
//  8. Can the graph change while traversing?
//
//	This implementation assumes a stable graph snapshot.
//
//  9. Is the graph too large for one machine?
//
//	A distributed BFS may partition frontier nodes across workers.
//
// 10. Is the target one service or every service?
//
//	This implementation calculates distances to all reachable services.
//
//	If only one target is needed, stop when that target is discovered.
//
// Plan:
//
//  1. Validate start.
//  2. Create distance array initialized to -1.
//  3. Set distance[start] = 0.
//  4. Add start to queue.
//  5. Process queue using a head index.
//  6. For every neighbour:
//
//	    - validate ID
//	    - if distance is still -1, mark it
//	    - set distance[neighbour] = distance[current] + 1
//	    - append neighbour to queue
//
//  7. Return distance.
//
// Time complexity:
//
//	O(V + E)
//
// Every node is queued at most once.
//
// Every adjacency edge is examined once.
//
// Space complexity:
//
//	O(V)
//
// Used by:
//
//	- distance array
//	- queue
//
// Cyclic dependency checks:
//
// The graph may contain cycles.
//
// Example:
//
//	0 -> 1 -> 2 -> 0
//
// BFS does not require separate cycle detection for shortest paths.
//
// The distance array also acts as the visited structure.
//
// Once distance[node] is not -1, that node is not enqueued again.
//
// This prevents infinite traversal.
//
// Race conditions:
//
// The function uses local state.
//
// Separate calls are safe concurrently.
//
// The caller must not modify graph while this function is traversing it.
//
// For a mutable shared graph, use:
//
//	- read lock
//	- immutable snapshot
//	- copy-on-write structure

func shortestHopDistances(
	graph [][]int,
	start int,
) ([]int, error) {
	nodeCount := len(graph)

	// Boundary check:
	//
	// Empty graph has no valid starting node.
	if nodeCount == 0 {
		return nil, fmt.Errorf("graph cannot be empty")
	}

	// Boundary check:
	//
	// Valid node indexes are 0 through nodeCount-1.
	if start < 0 || start >= nodeCount {
		return nil, fmt.Errorf(
			"start node %d is outside range [0,%d]",
			start,
			nodeCount-1,
		)
	}

	// distance is indexed by service ID.
	//
	// Index:
	//
	//	service ID
//
// Value:
	//
	//	shortest number of hops from start
//
// Example:
	//
	//	distance[3] = 2
//
// means service 3 is two edges away from start.
	distance := make([]int, nodeCount)

	for nodeID := 0; nodeID < nodeCount; nodeID++ {
		// -1 means the node has not been reached.
		distance[nodeID] = -1
	}

	// The starting node is zero hops from itself.
	distance[start] = 0

	// queue stores service ID VALUES.
	//
	// It does not store distances or indexes into graph adjacency lists.
	queue := []int{start}

	// head is an INDEX into queue.
	//
	// queue[head] is the next service ID VALUE to process.
	head := 0

	// Example trace:
	//
	// graph:
	//
	//	0: [1,2]
	//	1: [0,3]
	//	2: [0,3,4]
	//	3: [1,2,4]
	//	4: [2,3]
	//
	// start = 0
	//
	// Initial:
	//
	//	distance = [0,-1,-1,-1,-1]
	//	queue = [0]
	//	head = 0
	//
	// Iteration 1:
	//
	//	currentNode = queue[0] = 0
	//	head changes from 0 to 1
	//
	//	Neighbour 1:
	//
	//	    distance[1] == -1
	//	    distance[1] = distance[0] + 1 = 1
	//	    queue becomes [0,1]
	//
	//	Neighbour 2:
	//
	//	    distance[2] == -1
	//	    distance[2] = 1
	//	    queue becomes [0,1,2]
	//
	//	After iteration:
	//
	//	    distance = [0,1,1,-1,-1]
	//	    queue = [0,1,2]
	//	    head = 1
	//
	// Iteration 2:
	//
	//	currentNode = queue[1] = 1
	//
	//	Neighbour 0:
	//
	//	    distance[0] is already 0
	//	    skip
	//
	//	Neighbour 3:
	//
	//	    distance[3] == -1
	//	    distance[3] = distance[1] + 1 = 2
	//	    queue becomes [0,1,2,3]
	//
	// Iteration 3:
	//
	//	currentNode = 2
	//
	//	Neighbour 0:
	//	    already visited
	//
	//	Neighbour 3:
	//	    already visited at distance 2
	//
	//	Neighbour 4:
	//
	//	    distance[4] = distance[2] + 1 = 2
	//	    queue becomes [0,1,2,3,4]
	//
	// Final:
	//
	//	distance = [0,1,1,2,2]

	for head < len(queue) {
		currentNode := queue[head]
		head++

		for neighbourIndex, neighbourNode :=
			range graph[currentNode] {

			// neighbourIndex is an INDEX within graph[currentNode].
			//
			// neighbourNode is a service ID VALUE.
			_ = neighbourIndex

			// Boundary check:
			//
			// Every neighbour must refer to a valid graph node.
			if neighbourNode < 0 ||
				neighbourNode >= nodeCount {

				return nil, fmt.Errorf(
					"node %d contains invalid neighbour %d",
					currentNode,
					neighbourNode,
				)
			}

			// Cycle and duplicate-edge protection:
			//
			// A node is visited when distance[node] != -1.
			//
			// Skip nodes already discovered.
			if distance[neighbourNode] != -1 {
				continue
			}

			// First discovery through BFS is the shortest path.
			distance[neighbourNode] =
				distance[currentNode] + 1

			queue = append(queue, neighbourNode)
		}
	}

	return distance, nil
}

func main() {
	graph := [][]int{
		{1, 2},
		{0, 3},
		{0, 3, 4},
		{1, 2, 4},
		{2, 3},
	}

	distances, err := shortestHopDistances(graph, 0)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	fmt.Println(distances)

	// Expected output:
	//
	//	[0 1 1 2 2]
}

// Alternative questions that can be solved using this approach:
//
// 1. Shortest path in an unweighted network.
//
// 2. Minimum number of service-to-service hops.
//
// 3. Kubernetes blast-radius analysis:
//
//	Starting from a failed dependency, find services affected at each level.
//
// 4. Rotting Oranges:
//
//	Use multi-source BFS.
//
// 5. Number of Islands:
//
//	Run BFS from each unvisited land cell.
//
// 6. Shortest path through a binary matrix.
//
// 7. Word Ladder:
//
//	Each valid word transformation is one graph edge.
//
// 8. Minimum number of configuration transitions.
//
// 9. Find all nodes within K hops.
//
//	Stop expanding nodes after distance K.
//
// 10. Nearest healthy replica:
//
//	Use multi-source BFS beginning from all healthy replicas.
//
// 11. Cluster network reachability:
//
//	Determine which endpoints are reachable from a source.
//
// 12. Level-order traversal of a tree:
//
//	A tree is a special graph.
//
// Recognition pattern:
//
// When the question contains:
//
//	- shortest number of edges
//	- minimum hops
//	- nearest node
//	- level by level
//	- unweighted graph
//	- all nodes within K steps
//
// consider:
//
//	Breadth-First Search
//	    +
//	Queue
//	    +
//	Visited or distance array

// time complexity: O(V + E) -> each reachable vertex is processed once and each edge is examined once.
// space complexity: O(V) -> the visited state, queue, stack, or result can hold one entry per vertex.
```
