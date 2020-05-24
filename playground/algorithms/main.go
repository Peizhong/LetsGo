package main

import "log"

func main() {
	p1()
	sort := func(alg func(nums []int)) {
		nums := randNums(10)
		log.Println(nums)
		alg(nums)
		checkSort(nums)
	}
	sort(SelectionSort)
}
