package main

import (
	"log"

	"github.com/davecgh/go-spew/spew"
)

func Do() {
	p23()
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

// p23: 合并k个有序链表
func p23() {
	mergeKLists := func(lists []*ListNode) *ListNode {
		return nil
	}
	var lists []*ListNode
	// build data
	{
		nums := [][]int{
			[]int{1, 4, 5},
			[]int{1, 3, 4},
			[]int{2, 6},
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
