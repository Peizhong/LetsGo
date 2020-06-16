package main

var pint *int

func hi() int {
	n1, n2 := 1, 2
	pint = &n1
	pint = &n2

	return *pint
}

// go build -gcflags "-N -l"
// -N: 禁止优化
// -l: 禁止内联
// -E: 导出debug信息
// 生成过程中的汇编：go tool compile -S -N -l main.go
// 最终代码的汇编：go tool objdump -s 'main\.hi' -S ./min
func main() {
	type s struct {
		k, v             int
		v2, v3, v4       uint64
		buf1, buf2, buf3 [4096][4096]byte
	}
	var ps *s
	var psold *s
	for i := 0; i < 10000; i++ {
		for j := 0; j < 1000; j++ {
			ps = &s{k: i, v: j}
			if psold != nil {
				ps.v2 = psold.v2
			}
		}
		psold = ps
	}
	println(ps.k, ps.v)
}
