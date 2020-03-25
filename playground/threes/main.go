package main

import (
	"bufio"
	"github.com/peizhong/letsgo/internal"
	"math/rand"
	"os"
	"strconv"
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
	go func() {
		// get input
		reader := bufio.NewScanner(os.Stdin)
		for reader.Scan() {
			mv, err := strconv.Atoi(reader.Text())
			if err == nil {
				ch <- mv
			}
		}
	}()
	internal.Host(func() {
		show()
		for mv := range ch {
			move(mv)
			show()
		}
	}, nil)
}
