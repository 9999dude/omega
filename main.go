package main

import (
	"fmt"
)

// Node represents a node in the AVL Tree
type Node struct {
	Value  int
	Left   *Node
	Right  *Node
	Height int // Keep track of height for balance
}

// Helper function to get the height of a node
func getHeight(node *Node) int {
	if node == nil {
		return 0
	}
	return node.Height
}

// Helper function to update the height of a node
func updateHeight(node *Node) {
	if node != nil {
		node.Height = 1 + max(getHeight(node.Left), getHeight(node.Right))
	}
}

// Helper function to find the maximum of two integers
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// Helper function to perform a right rotation
func rightRotate(y *Node) *Node {
	x := y.Left
	T2 := x.Right

	// Perform rotation
	x.Right = y
	y.Left = T2

	// Update heights
	updateHeight(y)
	updateHeight(x)

	return x
}

// Helper function to perform a left rotation
func leftRotate(x *Node) *Node {
	y := x.Right
	T2 := y.Left

	// Perform rotation
	y.Left = x
	x.Right = T2

	// Update heights
	updateHeight(x)
	updateHeight(y)

	return y
}

// Helper function to get the balance factor of a node
func getBalanceFactor(node *Node) int {
	if node == nil {
		return 0
	}
	return getHeight(node.Left) - getHeight(node.Right)
}

// Insert a value into the AVL Tree
func (node *Node) Insert(value int) *Node {
	if node == nil {
		return &Node{Value: value, Height: 1}
	}

	if value < node.Value {
		node.Left = node.Left.Insert(value)
	} else if value > node.Value {
		node.Right = node.Right.Insert(value)
	} else {
		return node // Duplicate values are not allowed
	}

	// Update height of the current node
	updateHeight(node)

	// Balance the tree
	balanceFactor := getBalanceFactor(node)

	// Left Left Case
	if balanceFactor > 1 && value < node.Left.Value {
		return rightRotate(node)
	}

	// Right Right Case
	if balanceFactor < -1 && value > node.Right.Value {
		return leftRotate(node)
	}

	// Left Right Case
	if balanceFactor > 1 && value > node.Left.Value {
		node.Left = leftRotate(node.Left)
		return rightRotate(node)
	}

	// Right Left Case
	if balanceFactor < -1 && value < node.Right.Value {
		node.Right = rightRotate(node.Right)
		return leftRotate(node)
	}

	return node
}

// Delete a value from the AVL Tree
func (node *Node) Delete(value int) *Node {
	if node == nil {
		return node
	}

	if value < node.Value {
		node.Left = node.Left.Delete(value)
	} else if value > node.Value {
		node.Right = node.Right.Delete(value)
	} else {
		if node.Left == nil || node.Right == nil {
			var temp *Node
			if temp == node.Left {
				temp = node.Right
			} else {
				temp = node.Left
			}
			if temp == nil {
				temp = node
				node = nil
			} else {
				*node = *temp
			}
		} else {
			temp := minValueNode(node.Right)
			node.Value = temp.Value
			node.Right = node.Right.Delete(temp.Value)
		}
	}

	if node == nil {
		return node
	}

	// Update height of the current node
	updateHeight(node)

	// Balance the tree
	balanceFactor := getBalanceFactor(node)

	// Left Left Case
	if balanceFactor > 1 && getBalanceFactor(node.Left) >= 0 {
		return rightRotate(node)
	}

	// Left Right Case
	if balanceFactor > 1 && getBalanceFactor(node.Left) < 0 {
		node.Left = leftRotate(node.Left)
		return rightRotate(node)
	}

	// Right Right Case
	if balanceFactor < -1 && getBalanceFactor(node.Right) <= 0 {
		return leftRotate(node)
	}

	// Right Left Case
	if balanceFactor < -1 && getBalanceFactor(node.Right) > 0 {
		node.Right = rightRotate(node.Right)
		return leftRotate(node)
	}

	return node
}

// Helper function to find the minimum value node in the subtree
func minValueNode(node *Node) *Node {
	current := node
	for current.Left != nil {
		current = current.Left
	}
	return current
}

// Search operation to find a node in the AVL Tree
func (node *Node) Search(value int) *Node {
	if node == nil || node.Value == value {
		return node
	}

	if value < node.Value {
		return node.Left.Search(value)
	} else {
		return node.Right.Search(value)
	}
}

// In-order traversal of the AVL Tree
func (node *Node) InOrder() {
	if node != nil {
		node.Left.InOrder()
		fmt.Printf("%d ", node.Value)
		node.Right.InOrder()
	}
}

// Pre-order traversal of the AVL Tree
func (node *Node) PreOrder() {
	if node != nil {
		fmt.Printf("%d ", node.Value)
		node.Left.PreOrder()
		node.Right.PreOrder()
	}
}

// Post-order traversal of the AVL Tree
func (node *Node) PostOrder() {
	if node != nil {
		node.Left.PostOrder()
		node.Right.PostOrder()
		fmt.Printf("%d ", node.Value)
	}
}

func main() {
	root := &Node{} // Create the root node
	root.Insert(10)
	root.Insert(20)
	root.Insert(30)
	root.Insert(40)
	root.Insert(50)

	fmt.Println("Tree after inserting values:")
	fmt.Print("In-order: ")
	root.InOrder()
	fmt.Println()
	fmt.Print("Pre-order: ")
	root.PreOrder()
	fmt.Println()
	fmt.Print("Post-order: ")
	root.PostOrder()
	fmt.Println()

	// Search for a value
	fmt.Println("Searching for 30:", root.Search(30))
	fmt.Println("Searching for 60:", root.Search(60))

	// Delete a value
	root = root.Delete(30)

	fmt.Println("Tree after deleting 30:")
	fmt.Print("In-order: ")
	root.InOrder()
	fmt.Println()
	fmt.Print("Pre-order: ")
	root.PreOrder()
	fmt.Println()
	fmt.Print("Post-order: ")
	root.PostOrder()
	fmt.Println()
}
