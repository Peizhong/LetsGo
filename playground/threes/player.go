package main

import (
	"bufio"
	"os"
	"strconv"
)

func manplay(ch chan<- int) {
	go func() {
		reader := bufio.NewScanner(os.Stdin)
		for reader.Scan() {
			mv, err := strconv.Atoi(reader.Text())
			if err == nil {
				ch <- mv
			}
		}
	}()
}

func aiplay(ch chan<- int) {

}
