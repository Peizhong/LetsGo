package main

import "math/rand"

const (
	MoveUp    = 8
	MoveDown  = 2
	MoveLeft  = 4
	MoveRight = 6
)

type point struct {
	x, y int
}

var (
	nextValue int
	edge      = []point{
		{1, 1}, {1, 2}, {1, 3}, {1, 4},
		{2, 1}, {2, 4},
		{3, 1}, {3, 4},
		{4, 1}, {4, 2}, {4, 3}, {4, 4},
	}
)

func shiftUp() {
	for c := 1; c < 5; c++ {
		for i := 1; i < 4; {
			if box[i][c] == 0 {
				// 0 下面的上来
				for s := i; s < 5; s++ {
					box[s][c] = box[s+1][c]
				}
				break
			}
			if box[i][c] == box[i+1][c] && box[i][c] >= 3 {
				// 合并
				box[i][c] = (box[i][c]) * 2
				score += box[i][c]
				// 下面的上来
				for s := i + 1; s < 5; s++ {
					box[s][c] = box[s+1][c]
				}
				break
			}
			if box[i][c]+box[i+1][c] == 3 {
				// 合并
				box[i][c] = 3
				score += box[i][c]
				// 下面的上来
				for s := i + 1; s < 5; s++ {
					box[s][c] = box[s+1][c]
				}
				break
			}
			i++
		}
	}
}

func shiftDown() {
	for c := 1; c < 5; c++ {
		for i := 4; i > 1; {
			if box[i][c] == 0 {
				// 0 上面的下去
				for s := i; s > 0; s-- {
					box[s][c] = box[s-1][c]
				}
				break
			}
			if box[i][c] == box[i-1][c] && box[i][c] >= 3 {
				// 合并
				box[i][c] = (box[i][c]) * 2
				score += box[i][c]
				// 下面的上来
				for s := i - 1; s > 0; s-- {
					box[s][c] = box[s-1][c]
				}
				break
			}
			if box[i][c]+box[i-1][c] == 3 {
				// 合并
				box[i][c] = 3
				score += box[i][c]
				// 下面的上来
				for s := i - 1; s > 0; s-- {
					box[s][c] = box[s-1][c]
				}
				break
			}
			i--
		}
	}
}

func shiftLeft() {
	for r := 1; r < 5; r++ {
		for i := 1; i < 4; {
			if box[r][i] == 0 {
				// 左移
				for s := i; s < 5; s++ {
					box[r][s] = box[r][s+1]
				}
				break
			}
			if box[r][i] == box[r][i+1] && box[r][i] >= 3 {
				box[r][i] = (box[r][i]) * 2
				score += box[r][i]
				// 左移
				for s := i + 1; s < 5; s++ {
					box[r][s] = box[r][s+1]
				}
				break
			}
			if box[r][i]+box[r][i+1] == 3 {
				box[r][i] = 3
				score += box[r][i]
				// 左移
				// 左移
				for s := i + 1; s < 5; s++ {
					box[r][s] = box[r][s+1]
				}
				break
			}
			i++
		}
	}
}

func shiftRight() {
	for r := 1; r < 5; r++ {
		for i := 4; i > 0; {
			if box[r][i] == 0 {
				// 左移
				for s := i; s > 0; s-- {
					box[r][s] = box[r][s-1]
				}
				break
			}
			if box[r][i] == box[r][i-1] && box[r][i] >= 3 {
				box[r][i] = (box[r][i]) * 2
				score += box[r][i]
				// 左移
				for s := i - 1; s > 0; s-- {
					box[r][s] = box[r][s-1]
				}
				break
			}
			if box[r][i]+box[r][i-1] == 3 {
				box[r][i] = 3
				score += box[r][i]
				// 左移
				for s := i - 1; s > 0; s-- {
					box[r][s] = box[r][s-1]
				}
				break
			}
			i--
		}
	}
}

// 随机方向插入
func addItem(dir int) bool {
	// 可以插入位置
	var available []int
	switch dir {
	case MoveUp:
		for i := 1; i < 5; i++ {
			if box[4][i] == 0 {
				available = append(available, i)
			}
		}
		if len(available) > 0 {
			pos := rand.Intn(len(available))
			box[4][available[pos]] = nextValue
			return true
		}
	case MoveDown:
		for i := 1; i < 5; i++ {
			if box[1][i] == 0 {
				available = append(available, i)
			}
		}
		if len(available) > 0 {
			pos := rand.Intn(len(available))
			box[1][available[pos]] = nextValue
			return true
		}
	case MoveLeft:
		for i := 1; i < 5; i++ {
			if box[i][4] == 0 {
				available = append(available, i)
			}
		}
		if len(available) > 0 {
			pos := rand.Intn(len(available))
			box[available[pos]][4] = nextValue
			return true
		}
	case MoveRight:
		for i := 1; i < 5; i++ {
			if box[i][1] == 0 {
				available = append(available, i)
			}
		}
		if len(available) > 0 {
			pos := rand.Intn(len(available))
			box[available[pos]][1] = nextValue
			return true
		}
	default:
		return true
	}
	return false
}

// 还有可以合并的项目
func canMove() bool {
	for r := 1; r < 5; r++ {
		for c := 1; c < 5; c++ {
			// 水平
			if box[r][c]+box[r][c+1] == 3 && c != 4 {
				println("hor_a", r, c)
				return true
			}
			if box[r][c] == box[r][c+1] && box[r][c] >= 3 {
				println("hor_x", r, c)
				return true
			}
			// 垂直
			if box[r][c]+box[r+1][c] == 3 && r != 4 {
				println("ver_a", r, c)
				return true
			}
			if box[r][c] == box[r+1][c] && box[r][c] >= 3 {
				println("ver_x", r, c)
				return true
			}
		}
	}
	return false
}

var v = 1

func next() {
	nextValue = v%2 + 1
	v++
}

func move(dir int) {
	switch dir {
	case MoveUp:
		println("Up")
		shiftUp()
	case MoveDown:
		println("Down")
		shiftDown()
	case MoveLeft:
		println("Left")
		shiftLeft()
	case MoveRight:
		println("Right")
		shiftRight()
	default:
		return
	}
	round++
}
