package main

import (
	"container/list"
	"fmt"
	br "github.com/arnauddri/algorithms/data-structures/binary-tree"
)

func main() {
	n := br.NewNode(8)
	tree := br.NewTree(n)
	tree.Insert(4)
	tree.Insert(2)
	tree.Insert(9)
	tree.Insert(5)
	tree.Insert(3)
	tree.Insert(1)
	tree.Insert(7)
	tree.Insert(10)

	l := list.New()
	Dofind(tree.Head, 0, 24, l)

	var a = []int{0, 1, 0, 2, 1, 0, 1, 3, 2, 1, 2, 1}
	fmt.Println("---", traprain(a))
}

//递归方式实现，本质是二叉树的遍历 ：只要存在左节点，会一直递归（压栈），直到该节点没有子节点，然后最后一个左节点出栈，回到父节点上
//然后走右节点的递归（压栈)
func Dofind(t *br.Node, current, sum int, l *list.List) {
	if t == nil {
		return
	}
	fmt.Println("push ", t.Value)
	l.PushBack(t.Value)
	current += t.Value

	if t.Left != nil {
		Dofind(t.Left, current, sum, l)
	}
	if t.Right != nil {
		Dofind(t.Right, current, sum, l)
	}

	if t.Left == nil && t.Right == nil && current == sum {
		show(l)
	}

	if e := l.Back(); e != nil {
		fmt.Println("pop ", e.Value)
		l.Remove(e)
	}
	current -= t.Value
}

func show(l *list.List) {
	for e := l.Front(); e != nil; e = e.Next() {
		fmt.Printf("  -> %d ", e.Value)
	}
	fmt.Println()
}

func traprain(a []int) int {
	left := 0
	right := len(a) - 1
	ans := 0
	l_max := 0
	r_max := 0

	//思想：两头比较，向中间靠近，直到重合退出
	for left < right {
		//若左边的小，比较左边的值与左边最大值，若小于最大值，差距记为存量；否则，修改左边最大值。右移一位，从新进入循环
		if a[left] < a[right] {
			if a[left] >= l_max {
				l_max = a[left]
			} else {
				ans += (l_max - a[left])
			}
			left++
		} else {
			//若右边的小，比较右边的值与右边最大值，若小于最大值，差距记为存量；否则，修改左边最大值。左移一位，从新进入循环
			if a[right] >= r_max {
				r_max = a[right]
			} else {
				ans += (r_max - a[right])
			}
			right--
		}
	}
	return ans
}
