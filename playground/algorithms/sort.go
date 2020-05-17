package main

import (
	"log"
	"math/rand"
	"time"
)

func randNums(n int) []int {
	rand.Seed(time.Now().UnixNano())
	nums := make([]int, n)
	for i := range nums {
		nums[i] = rand.Int()%n + 1
	}
	return nums
}

func checkSort(nums []int) {
	length := len(nums)
	for i := 0; i < length-1; i++ {
		if nums[i] > nums[i+1] {
			log.Println("error")
			log.Println(nums)
			return
		}
	}
}

// 冒泡排序
// 复杂度 O(n^2)
// 稳定
func BubbleSort(nums []int) {

}

// 选择排序
// 不稳定 5 8 5 2 9
func SelectionSort(nums []int) {
	length := len(nums)
	if length < 1 {
		return
	}
	for index := range nums {
		// 从剩下的数据找最小的
		minIndex := index
		for i := index + 1; i < length; i++ {
			// 值相等的情况。。
			if nums[i] < nums[minIndex] {
				minIndex = i
			}
		}
		// 按顺序排放
		if minIndex != index {
			nums[index], nums[minIndex] = nums[minIndex], nums[index]
		}
	}
}

// 插入排序
// 假设第一个是拍好序的
func InsertSort(nums []int) {

}

// 希尔排序
// 划分区间，排序，逐渐减少区间
// 复杂度: O(nlogn)
// 不稳定
func ShellSort(nums []int) {
	length := len(nums)
	if length < 1 {
		return
	}
	for step := length / 2; step > 0; step = step / 2 {
		// 步长step
		// 排序1 step,step-step
		// 排序2 step+2, step+2-step
		// ..
		// 排序n step+step+step,step+step,step
		// 每组的元素逐渐增多，只将最后一个元素插入到合适的位置
		for i := step; i < length; i++ {
			for s := i - step; s >= 0; s -= step {
				// 分组内，插入排序
				if nums[s+step] < nums[s] {
					nums[s+step], nums[s] = nums[s], nums[s+step]
				} else {
					break
				}
			}
		}
	}
}

// 分区：将数据根据nums[p]划分成2部分
func quick_partition(nums []int, s, e int, p int) int {
	pivot := nums[p]
	// 按大小划分2个区域
	// 挖坑法
	i, j := s, e-1
	for {
		// 从右往左找比锚点小的
		for ; j > i; j-- {
			if nums[j] < pivot {
				nums[p] = nums[j]
				p = j
				break
			}
		}
		if i == j {
			p = i
			nums[p] = pivot
			break
		}
		for ; i < j; i++ {
			if nums[i] > pivot {
				nums[p] = nums[i]
				p = i
				break
			}
		}
		if i == j {
			p = i
			nums[p] = pivot
			break
		}
	}
	return p
}

func quick_quick(nums []int, s, e int) {
	if s >= e {
		return
	}
	// 随机取锚点
	p := s
	p = quick_partition(nums, s, e, p)
	// 继续分区
	quick_quick(nums, s, p)
	quick_quick(nums, p+1, e)
}

// 快速排序
// 取锚点，分成左右分区，继续划分
// 不稳定
func QuickSort(nums []int) {
	length := len(nums)
	if length < 1 {
		return
	}
	quick_quick(nums, 0, length)
}

func merge_merge(left, right []int) []int {
	leftLen, rightLen := len(left), len(right)
	var leftIndex, rightIndex int
	res := make([]int, 0, leftLen+rightLen)
	for {
		if leftIndex == leftLen {
			if rightIndex == rightLen {
				break
			} else {
				res = append(res, right[rightIndex])
				rightIndex++
			}
		}
		if rightIndex == rightLen {
			if leftIndex == leftLen {
				break
			} else {
				res = append(res, left[leftIndex])
				leftIndex++
			}
		}
		if left[leftIndex] < right[rightIndex] {
			res = append(res, left[leftIndex])
			leftIndex++
		} else {
			res = append(res, right[rightIndex])
			rightIndex++
		}
	}
	return res
}

func merge_split(nums []int) []int {
	length := len(nums)
	// 不用再拆了
	if length == 1 {
		return nums
	}
	// 划分成左右部分
	split := length / 2
	left := nums[:split]
	right := nums[split:]
	// 分别拆左右，然后在合并

	res := merge_merge(merge_split(left), merge_split(right))
	return res
}

// 归并排序
// 分治法
func MergeSort(nums []int) {
	// 分2半，排序后再合并
	res := merge_split(nums)
	// 写入
	for i := range res {
		nums[i] = res[i]
	}
}
