package main

import (
	"log"
	"math"
	"math/rand"

	"github.com/davecgh/go-spew/spew"
)

// Do 做题
func Do() {
	p236()
}

// 字典序排序：13-> 1,10,11,12,13,2,3
func p386() {

}

// 字符串转数字
func p8() {
	str := "abc我是def"
	for _, s := range str {
		log.Println(string(s))
	}
}

func p15() {
	// 锚点，双指针
	nums := []int{-1, 0, 1, 2, -1, -4}
	l := len(nums)
	var i, j, k int
	for i = 0; i < l-2; i++ {
		for j = i + 1; j < l-1; j++ {
			for k = j + 1; k < l; k++ {
				if nums[i]+nums[j]+nums[k] == 0 {
					log.Println(nums[i], nums[j], nums[k])
				}
			}
		}
	}
}

type ListNode struct {
	Val  int
	Next *ListNode
}

func (l *ListNode) print() {
	for l != nil {
		print(l.Val)
		l = l.Next
	}
	println()
}

func (l *ListNode) init(nums []int) {
	cur := l
	for _, n := range nums {
		cur.Next = &ListNode{Val: n}
		cur = cur.Next
	}
}

// p23: 合并k个有序链表
func p23() {
	mergeKLists := func(lists []*ListNode) *ListNode {
		return nil
	}
	var lists []*ListNode
	// build data
	{
		nums := [][]int{
			{1, 4, 5},
			{1, 3, 4},
			{2, 6},
		}
		for _, nl := range nums {
			root := &ListNode{}
			cur := root
			for _, n := range nl {
				nd := &ListNode{Val: n}
				cur.Next = nd
				cur = nd
			}
			lists = append(lists, root.Next)
		}
	}
	spew.Dump(lists)
	mergeKLists(lists)
}

type TreeNode struct {
	// maxint32 无效
	Val         int
	Left, Right *TreeNode
}

func readTree(n *TreeNode) []int {
	var res []int
	buf := []*TreeNode{n}
	for len(buf) > 0 {
		x := buf[0]
		buf = buf[1:]
		if x.Val != math.MaxInt32 {
			res = append(res, x.Val)
			if x.Left != nil {
				buf = append(buf, x.Left)
			}
			if x.Right != nil {
				buf = append(buf, x.Right)
			}
		}
	}
	log.Println(res)
	return res
}

func reportDepth(n *TreeNode) int {
	if n == nil {
		return 0
	}
	l := reportDepth(n.Left)
	r := reportDepth(n.Right)
	if l > r {
		return l + 1
	}
	return r + 1
}

// p104: 二叉树最大深度
func p104() {
	// prepare data
	root := &TreeNode{Val: 3}
	{
		root.Left = &TreeNode{Val: 9}
		root.Right = &TreeNode{
			Val:  20,
			Left: &TreeNode{Val: 15},
			Right: &TreeNode{Val: 7,
				Left: &TreeNode{Val: 8}},
		}
	}
	d := reportDepth(root)
	println(d)
}

func maxProfit(prices []int) {
	var start, maxStart, maxEnd, maxProfit int
	// greedy
	for i := range prices {
		profit := prices[i] - prices[start]
		if profit < 0 {
			start = i
			continue
		}
		if profit > maxProfit {
			maxStart = start
			maxEnd = i
			maxProfit = profit
		}
	}
	if maxStart < maxEnd {
		println(maxStart, maxEnd, maxProfit)
	} else {
		println("no prifit")
	}
}

// p121: 购买股票的时机，买卖一次
func p121() {
	prices := []int{10, 2, 9, 1, 2, 1, 3, 1}
	maxProfit(prices)
}

func maxProfit2(prices []int) {
	// 贪心算法，下降就抛
	var profit int
	for i := range prices {
		if i == 0 {
			continue
		}
		daily := prices[i] - prices[i-1]
		if daily > 0 {
			profit += daily
		}
	}
	println(profit)
}

// p122: 购买股票的时机，买卖多次
func p122() {
	prices := []int{7, 1, 5, 3, 6, 4}
	maxProfit2(prices)
	prices = []int{1, 2, 3, 4, 5}
	maxProfit2(prices)
}

func buildNilTree(nums []int) *TreeNode {
	root := &TreeNode{Val: nums[0]}
	nIndex, nLength := 1, len(nums)
	nodeList := []*TreeNode{root}
	for len(nodeList) > 0 && nIndex < nLength {
		n := nodeList[0]
		nodeList = nodeList[1:]
		if nums[nIndex] != math.MaxInt32 {
			n.Left = &TreeNode{Val: nums[nIndex]}
			nodeList = append(nodeList, n.Left)
		} else {
			nodeList = append(nodeList, &TreeNode{})
		}
		nIndex++
		if nums[nIndex] != math.MaxInt32 {
			n.Right = &TreeNode{Val: nums[nIndex]}
			nodeList = append(nodeList, n.Right)
		} else {
			nodeList = append(nodeList, &TreeNode{})
		}
		nIndex++
	}
	return root
}

func buildTree(nums []int) *TreeNode {
	if len(nums) < 1 {
		return nil
	}
	root := &TreeNode{}
	nodeList := []*TreeNode{root}
	var nodeIndex int
	for _, i := range nums {
		nodeList[nodeIndex].Val = i
		if len(nodeList) < len(nums) {
			nodeList[nodeIndex].Left = &TreeNode{}
			nodeList[nodeIndex].Right = &TreeNode{}
			nodeList = append(nodeList, nodeList[nodeIndex].Left, nodeList[nodeIndex].Right)
		}
		nodeIndex++
	}
	return root
}

func maxTreePath(root *TreeNode) {
	// 遍历树
}

// p124: 二叉树最大路径和
func p124() {
	nums := []int{-10, 9, 20, 0, 0, 15, 7}
	root := buildTree(nums)
	maxTreePath(root)
}

func isLoopList(root *ListNode) {
	slow := root
	fast := root.Next
	for fast != nil {
		log.Println(slow.Val, fast.Val)
		if slow == fast {
			log.Println("loop")
			return
		}
		slow = slow.Next
		// 每次走2步
		fast = fast.Next
		if fast == nil {
			break
		}
		fast = fast.Next
	}
	log.Println("no loop")
}

func whereLoopList(root *ListNode) {
	// hash
	m := make(map[*ListNode]int)
	var index int
	for root != nil {
		// 比较指针地址
		if v, ok := m[root]; ok {
			log.Println("loop at", v)
			return
		}
		m[root] = index
		index++
		root = root.Next
	}
	log.Println("no loop")
}

// 链表是否有环
func p142() {
	// prepare data
	loop := &ListNode{Val: 4}
	root := &ListNode{
		Val: 3,
		Next: &ListNode{
			Val: 2,
			Next: &ListNode{
				Val:  0,
				Next: loop,
			},
		},
	}
	loop.Next = root.Next
	whereLoopList(root)
}

type TwoWayListNode struct {
	k, v       int
	Prev, Next *TwoWayListNode
}

func (n *TwoWayListNode) Key() int {
	return n.k
}

func (n *TwoWayListNode) Value() int {
	return n.v
}

// LRUCache 最近最少使用
// 数据是个队列，容量满了删除
type LRUCache struct {
	Capacity, Size int
	m              map[int]*TwoWayListNode
	head, tail     *TwoWayListNode
}

func constructor(capacity int) *LRUCache {
	cache := &LRUCache{
		Capacity: capacity,
		m:        make(map[int]*TwoWayListNode),
	}
	return cache
}

// get 增加数据的权重
func (this *LRUCache) get(key int) int {
	if v, ok := this.m[key]; ok {
		if v == this.head {
			this.head = v.Next
			this.head.Prev = nil
		}
		this.tail.Next = v
		v.Prev = this.tail
		v.Next = nil
		this.tail = v
		log.Println("get", key, v.Value())
		return v.Value()
	}
	log.Println("get", key, -1)
	return -1
}

// put
func (this *LRUCache) put(key int, value int) {
	var curPos *TwoWayListNode
	if v, ok := this.m[key]; ok {
		curPos = v
	}
	if curPos == nil {
		curPos = &TwoWayListNode{
			k: key,
			v: value,
		}
		this.m[key] = curPos
	}
	if len(this.m) == 1 {
		this.head, this.tail = curPos, curPos
		return
	}
	if curPos == this.head {
		this.head = curPos.Next
		this.head.Prev = nil
	}
	this.tail.Next = curPos
	curPos.Prev = this.tail
	curPos.Next = nil
	this.tail = curPos
	if len(this.m) > this.Capacity {
		this.head.Next.Prev = nil
		delete(this.m, this.head.Key())
		this.head = this.head.Next
	}
}

// LRU
func p146() {
	cache := constructor(2)
	cache.put(1, 1)
	cache.put(2, 2)
	cache.get(1)
	cache.put(3, 3)
	cache.get(2)
	cache.put(4, 4)
	cache.get(1)
	cache.get(3)
	cache.get(4)
}

func mergeSortList(n *ListNode, l *ListNode) *ListNode {
	if l == nil {
		return n
	}
	var prev *ListNode
	cur := l
	for cur != nil {
		if n.Val < cur.Val {
			if prev == nil {
				n.Next = l
				return n
			}
			prev.Next = n
			n.Next = cur
			return l
		}
		prev = cur
		cur = cur.Next
	}
	// n比l的都大
	n.Next = nil
	prev.Next = n
	return l
}

func sortList(n *ListNode) *ListNode {
	if n == nil {
		return n
	}
	return mergeSortList(n, sortList(n.Next))
}

// 排序链表
func p148() {
	data := &ListNode{
		Val: 4,
		Next: &ListNode{
			Val: 2,
			Next: &ListNode{
				Val: 1,
				Next: &ListNode{
					Val: 3,
				},
			},
		},
	}
	l := sortList(data)
	l.print()
	data2 := &ListNode{}
	data2.init([]int{-1, 5, 3, 4, 0})
	l = sortList(data2.Next)
	l.print()
}

func reverseList(node *ListNode) *ListNode {
	// node 没有排序
	// node.next后面的都排序了
	if node == nil || node.Next == nil {
		return node
	}
	r := reverseList(node.Next)
	node.Next.Next = node
	node.Next = nil
	return r
}

// 反转链表
func p206() {
	data := &ListNode{}
	data.init([]int{1, 2, 3, 4, 5})
	r := reverseList(data.Next)
	r.print()
}

func findKthLargest(nums []int, k int) {
	swap := func(a, b int) {
		nums[a], nums[b] = nums[b], nums[a]
	}
	partition := func(l, r, pindex int) int {
		pivot := nums[pindex]
		// 移到最后
		swap(pindex, r)
		// 比pivot小的都移到前面
		sindex := l
		var i int
		for i = l; i <= r; i++ {
			if nums[i] < pivot {
				swap(i, sindex)
				sindex++
			}
		}
		swap(sindex, r)
		return sindex
	}
	//
	min, max := 0, len(nums)-1
	fixK := len(nums) - k
	for {
		rd := rand.Intn(max-min+1) + min
		p := partition(min, max, rd)
		if p == fixK {
			log.Println(nums)
			log.Println(nums[fixK])
			return
		}
		if p > fixK {
			max = p - 1
		} else {
			min = p + 1
		}
	}
}

// 找数组中第k大的
func p215() {
	nums := []int{3, 2, 2, 2, 8, 1, 0, 2, 7, 5, 5, 9, 6, 3, 4, 3, 6}
	findKthLargest(nums, 4)
}

func dfsFindKthMin(n *TreeNode, k *int) {
	if n.Left != nil && n.Left.Val != math.MaxInt32 {
		dfsFindKthMin(n.Left, k)
	}
	*k--
	if *k == 0 {
		log.Println(n.Val)
		return
	}
	if *k < 0 {
		return
	}
	if n.Right != nil && n.Right.Val != math.MaxInt32 {
		dfsFindKthMin(n.Right, k)
	}
}

// 二叉树中第k小的元素
func p230() {
	tr := buildTree([]int{5, 3, 6, 2, 4, math.MaxInt32, math.MaxInt32, 1, math.MaxInt32})
	k := 3
	dfsFindKthMin(tr, &k)
}

// 2的幂
func p231() {
	v := 18
	i := 1
	for {
		if v == i {
			log.Println("yes")
			return
		}
		if i > v {
			log.Println("no")
			return
		}
		i = i << 1
	}
}

func lowestCommonSearchAncestor(root *TreeNode, p, q int) {
	if root.Val > p && root.Val > q {
		lowestCommonSearchAncestor(root.Left, p, q)
	} else if root.Val < p && root.Val < q {
		lowestCommonSearchAncestor(root.Right, p, q)
	} else {
		log.Println("parent is", root.Val)
	}
}

// 二叉搜索树的最近公共祖先
func p235() {
	// 排序的
	tr := buildNilTree([]int{6, 2, 8, 0, 4, 7, 9, math.MaxInt32, math.MaxInt32, 3, 5})
	lowestCommonSearchAncestor(tr, 2, 7)
}

func searchRoute(root *TreeNode, v int) []int {
	if root == nil {
		return nil
	}
	if root.Val == v {
		return []int{root.Val}
	}
	r := searchRoute(root.Left, v)
	if len(r) > 0 {
		return append([]int{root.Val}, r...)
	}
	r = searchRoute(root.Right, v)
	if len(r) > 0 {
		return append([]int{root.Val}, r...)
	}
	return nil
}

func lowestCommonAncestor(root *TreeNode, p, q int) {
	l := searchRoute(root, p)
	r := searchRoute(root, q)
	log.Println(l)
	log.Println(r)
}

// 二叉树的最近公共祖先
func p236() {
	tr := buildNilTree([]int{3, 5, 1, 6, 2, 0, 8, math.MaxInt32, math.MaxInt32, 7, 4})
	lowestCommonAncestor(tr, 2, 7)
}
