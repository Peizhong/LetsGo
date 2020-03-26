package main

import (
	"github.com/peizhong/letsgo/internal"
	"math/rand"
)

var (
	round int
	score int
	box   = [6][6]int{}
)

func init() {
	for i := 0; i < 3; i++ {
		x := rand.Intn(4) + 1
		y := rand.Intn(4) + 1
		box[x][y] = i%2 + 1
	}
}

func main() {
	ch := make(chan int)
	manplay(ch)
	internal.Host(func() {
		next()
		show()
		for mv := range ch {
			move(mv)
			if addItem(mv) {
				next()
				show()
			} else {
				println("no new item")
				if !canMove() {
					println("game over")
					return
				}
			}
		}
	}, nil)
}
