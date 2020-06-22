package main

func UseSlice() {
	root := []int{1, 2, 3, 4, 5, 6}
	r1 := root[1:5]
	r2 := root[2:4]

	bak := make([]int, 10)
	// min(len(bak),len(r1))
	//copy(dest,srt)，拷贝长度不超过src
	copy(bak, r1)
	println(cap(r1), cap(r2))
	r1[2] = 100
	if r1[2] == r2[2] {
		println("same data")
	}
	// append后 r1,r2 底层数组还是一样的
	r1 = append(r1, 200)
	r2 = append(r2, 300)
	r2 = append(r2, 400)
	// 达到5之后，再扩容就是不同底层数组了
	r2 = append(r2, 500)
	r2[2] = 200
}

func UseString() {
	str := "hello 你好!"
	for _, s := range []rune(str) {
		println(string(s))
	}
}
