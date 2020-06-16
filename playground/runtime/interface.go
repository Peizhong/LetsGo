package main

type a struct{}
type b struct{}

type ia interface{}
type ib interface{}
type ic interface{ Hi() }

func isNIl(i interface{}) bool {
	return i == nil
}

func doInterface() {
	var ai ia = a{}
	var bi ib = b{}
	var ci ic
	var di ib

	var pa *a
	println("isnil", isNIl(pa))

	if func(i, j interface{}) bool {
		return i == j
	}(ai, bi) == true {
		println("same")
	} else {
		println("diff")
	}
	if func(i, j interface{}) bool {
		return i == j
	}(di, ci) == true {
		println("same")
	} else {
		println("diff")
	}
}
