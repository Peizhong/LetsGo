package main

import (
	"fmt"
)

// ArrangeClass: 最大活动安排问题
func ArrangeClass() {
	// 尽可能多的课
	m := make(map[string][]float32)
	m["art"] = []float32{9, 10}
	m["eng"] = []float32{9.3, 10.3}
	m["math"] = []float32{10, 11}
	m["it"] = []float32{10.3, 11.3}
	m["music"] = []float32{11, 12}
	var minStart float32 = 9
	var res []string
	// 递归
	findBest := func(start float32) (string, float32) {
		var s, e float32
		var key string
		for k, v := range m {
			// 选最早上课
			if v[0] >= start {
				if s == 0 {
					// 第一个值
					s = v[0]
					e = v[1]
					key = k
				} else if v[0] < s {
					// 开始时间更早
					s = v[0]
					e = v[1]
					key = k
				} else if v[0] == s {
					// 开始时间相同
					if v[1] < e {
						// 更早结束
						s = v[0]
						e = v[1]
						key = k
					}
				}
			}
		}
		return key, e
	}
	for {
		k, e := findBest(minStart)
		if k != "" {
			res = append(res, k)
			minStart = e
		} else {
			break
		}
	}
	fmt.Println(res)
}

// 最小生成树: Kruskal&Prim

// PrimTree: 最小生成树
func Prim() {
	// N个点用N-1条边连接成一个连通块，形成的图形只可能是树
	type route struct {
		p0, p1 int
		value  float32
	}
	// prepare data
	routes := []*route{
		{1, 2, 2},
		{2, 5, 9},
		{5, 4, 7},
		{4, 1, 10},
		{1, 3, 12},
		{4, 3, 6},
		{5, 3, 3},
		{2, 3, 8},
	}
	mark := map[int]struct{}{
		1: {},
	}
	unmark := []int{2, 3, 4, 5}
	var res []*route
	for {
		var candidate []*route
		// 从未连接点与连接点间的路径中，选择值最小的
		for mk := range mark {
			// 与unmark的节点
			for _, u := range unmark {
				fmt.Println("check", mk, u)
				for _, r := range routes {
					if r.p0 == u && r.p1 == mk {
						candidate = append(candidate, r)
					} else if r.p1 == u && r.p0 == mk {
						candidate = append(candidate, r)
					}
				}
			}
		}
		fmt.Println("candidate ", len(candidate))
		if len(candidate) < 1 {
			break
		}
		min := &route{}
		for _, c := range candidate {
			if min.value == 0 || min.value > c.value {
				min = c
			}
		}
		fmt.Println("min", min.p0, min.p1, min.value)
		// 添加到mark
		mark[min.p0] = struct{}{}
		mark[min.p1] = struct{}{}
		// 从unmark中移除
		var unmark2 []int
		for _, u := range unmark {
			if u != min.p0 && u != min.p1 {
				unmark2 = append(unmark2, u)
			}
		}
		unmark = unmark2
		res = append(res, min)
	}
	for _, r := range res {
		fmt.Println(r.p0, r.p1, r.value)
	}
}

func Kruskal() {
	// 边权重排序
	// 选择未连接的
}
