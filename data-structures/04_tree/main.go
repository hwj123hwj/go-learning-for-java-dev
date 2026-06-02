package main

import "fmt"

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

func NewTreeNode(val int) *TreeNode {
	return &TreeNode{Val: val}
}

func preOrder(root *TreeNode) {
	if root == nil {
		return
	}
	fmt.Printf("%d ", root.Val)
	preOrder(root.Left)
	preOrder(root.Right)
}

func inOrder(root *TreeNode) {
	if root == nil {
		return
	}
	inOrder(root.Left)
	fmt.Printf("%d ", root.Val)
	inOrder(root.Right)
}

func postOrder(root *TreeNode) {
	if root == nil {
		return
	}
	postOrder(root.Left)
	postOrder(root.Right)
	fmt.Printf("%d ", root.Val)
}

func levelOrder(root *TreeNode) {
	if root == nil {
		return
	}
	queue := []*TreeNode{root}
	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]
		fmt.Printf("%d ", node.Val)
		if node.Left != nil {
			queue = append(queue, node.Left)
		}
		if node.Right != nil {
			queue = append(queue, node.Right)
		}
	}
}

func treeHeight(root *TreeNode) int {
	if root == nil {
		return 0
	}
	leftH := treeHeight(root.Left)
	rightH := treeHeight(root.Right)
	if leftH > rightH {
		return leftH + 1
	}
	return rightH + 1
}

func nodeCount(root *TreeNode) int {
	if root == nil {
		return 0
	}
	return 1 + nodeCount(root.Left) + nodeCount(root.Right)
}

func leafCount(root *TreeNode) int {
	if root == nil {
		return 0
	}
	if root.Left == nil && root.Right == nil {
		return 1
	}
	return leafCount(root.Left) + leafCount(root.Right)
}

type BST struct {
	Root *TreeNode
}

func NewBST() *BST {
	return &BST{}
}

func (bst *BST) Insert(val int) {
	bst.Root = bst.insert(bst.Root, val)
}

func (bst *BST) insert(node *TreeNode, val int) *TreeNode {
	if node == nil {
		return NewTreeNode(val)
	}
	if val < node.Val {
		node.Left = bst.insert(node.Left, val)
	} else if val > node.Val {
		node.Right = bst.insert(node.Right, val)
	}
	return node
}

func (bst *BST) Search(val int) *TreeNode {
	return bst.search(bst.Root, val)
}

func (bst *BST) search(node *TreeNode, val int) *TreeNode {
	if node == nil || node.Val == val {
		return node
	}
	if val < node.Val {
		return bst.search(node.Left, val)
	}
	return bst.search(node.Right, val)
}

func (bst *BST) Delete(val int) {
	bst.Root = bst.deleteNode(bst.Root, val)
}

func (bst *BST) deleteNode(node *TreeNode, val int) *TreeNode {
	if node == nil {
		return nil
	}
	if val < node.Val {
		node.Left = bst.deleteNode(node.Left, val)
	} else if val > node.Val {
		node.Right = bst.deleteNode(node.Right, val)
	} else {
		if node.Left == nil {
			return node.Right
		}
		if node.Right == nil {
			return node.Left
		}
		successor := node.Right
		for successor.Left != nil {
			successor = successor.Left
		}
		node.Val = successor.Val
		node.Right = bst.deleteNode(node.Right, successor.Val)
	}
	return node
}

func isBalanced(root *TreeNode) bool {
	_, ok := checkHeight(root)
	return ok
}

func checkHeight(root *TreeNode) (int, bool) {
	if root == nil {
		return 0, true
	}
	leftH, leftOk := checkHeight(root.Left)
	rightH, rightOk := checkHeight(root.Right)
	if !leftOk || !rightOk {
		return 0, false
	}
	diff := leftH - rightH
	if diff < 0 {
		diff = -diff
	}
	if diff > 1 {
		return 0, false
	}
	h := leftH
	if rightH > leftH {
		h = rightH
	}
	return h + 1, true
}

func buildSampleTree() *TreeNode {
	root := NewTreeNode(1)
	root.Left = NewTreeNode(2)
	root.Right = NewTreeNode(3)
	root.Left.Left = NewTreeNode(4)
	root.Left.Right = NewTreeNode(5)
	root.Right.Left = NewTreeNode(6)
	root.Right.Right = NewTreeNode(7)
	return root
}

func main() {
	fmt.Println("=== 二叉树遍历 ===")
	fmt.Println()

	root := buildSampleTree()
	fmt.Println("        1")
	fmt.Println("       / \\")
	fmt.Println("      2   3")
	fmt.Println("     / \\ / \\")
	fmt.Println("    4  5 6  7")
	fmt.Println()

	fmt.Print("前序遍历(根左右): ")
	preOrder(root)
	fmt.Println()

	fmt.Print("中序遍历(左根右): ")
	inOrder(root)
	fmt.Println()

	fmt.Print("后序遍历(左右根): ")
	postOrder(root)
	fmt.Println()

	fmt.Print("层序遍历(BFS):    ")
	levelOrder(root)
	fmt.Println()

	fmt.Println()
	fmt.Println("--- 二叉树基本属性 ---")
	fmt.Printf("树的高度: %d\n", treeHeight(root))
	fmt.Printf("节点总数: %d\n", nodeCount(root))
	fmt.Printf("叶子节点数: %d\n", leafCount(root))

	fmt.Println()
	fmt.Println("--- 二叉排序树 (BST) ---")
	bst := NewBST()
	vals := []int{5, 3, 7, 1, 4, 6, 8}
	for _, v := range vals {
		bst.Insert(v)
	}
	fmt.Printf("插入 %v 后中序遍历: ", vals)
	inOrder(bst.Root)
	fmt.Println()

	fmt.Printf("查找4: %v\n", bst.Search(4) != nil)
	fmt.Printf("查找9: %v\n", bst.Search(9) != nil)

	bst.Delete(3)
	fmt.Print("删除3后中序遍历: ")
	inOrder(bst.Root)
	fmt.Println()

	fmt.Println()
	fmt.Println("--- 平衡二叉树判断 ---")
	fmt.Printf("示例树是否平衡: %v\n", isBalanced(root))

	unbalanced := NewTreeNode(1)
	unbalanced.Left = NewTreeNode(2)
	unbalanced.Left.Left = NewTreeNode(3)
	fmt.Printf("左斜树是否平衡: %v\n", isBalanced(unbalanced))

	fmt.Println()
	fmt.Println("=== 考研要点 ===")
	fmt.Println("1. 遍历: 前序(根左右)、中序(左根右)、后序(左右根)、层序(BFS)")
	fmt.Println("2. 由前序+中序 或 后序+中序 可唯一确定一棵二叉树")
	fmt.Println("3. BST: 中序遍历得到有序序列，查找/插入/删除 O(logn)~O(n)")
	fmt.Println("4. BST删除: 叶子直接删，单子树替代，双子树用后继(或前驱)替代")
	fmt.Println("5. 平衡因子: |左子树高度-右子树高度| <= 1")
	fmt.Println("6. 完全二叉树: 层序编号，节点i的左子2i，右子2i+1")
	fmt.Println("7. 线索二叉树: 利用空指针指向前驱/后继，n个节点有n+1个空指针")
}
