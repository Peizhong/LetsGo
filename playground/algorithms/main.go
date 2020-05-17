package main

import "log"

func main() {
	sort := func(alg func(nums []int)) {
		nums := randNums(10)
		log.Println(nums)
		alg(nums)
		checkSort(nums)
	}
	sort(SelectionSort)
}
