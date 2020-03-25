package main

import (
	"fmt"
	"log"
)

func show() {
	log.Println(fmt.Sprintf("round: %d, score: %d", round, score))
	for r := 1; r < 5; r++ {
		for c := 1; c < 5; c++ {
			print(box[r][c], " ")
		}
		println()
	}
	println()
}
