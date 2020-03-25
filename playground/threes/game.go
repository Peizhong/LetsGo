package main

import "math/rand"

const (
	MoveUp    = 8
	MoveDown  = 2
	MoveLeft  = 4
	MoveRight = 6
)

func shiftUp() {
	for c := 1; c < 5; c++ {
		for i := 1; i < 3; {
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
				// 下面的上来
				for s := i + 1; s < 5; s++ {
					box[s][c] = box[s+1][c]
				}
				break
			}
			if box[i][c]+box[i+1][c] == 3 {
				// 合并
				box[i][c] = 3
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
				// 下面的上来
				for s := i - 1; s > 0; s-- {
					box[s][c] = box[s-1][c]
				}
				break
			}
			if box[i][c]+box[i-1][c] == 3 {
				// 合并
				box[i][c] = 3
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
		for i := 1; i < 3; {
			if box[r][i] == 0 {
				// 左移
				for s := i; s < 5; s++ {
					box[r][s] = box[r][s+1]
				}
				break
			}
			if box[r][i] == box[r][i+1] && box[r][i] >= 3 {
				box[r][i] = (box[r][i]) * 2
				// 左移
				for s := i + 1; s < 5; s++ {
					box[r][s] = box[r][s+1]
				}
				break
			}
			if box[r][i]+box[r][i+1] == 3 {
				box[r][i] = 3
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
				// 左移
				for s := i - 1; s > 0; s-- {
					box[r][s] = box[r][s-1]
				}
				break
			}
			if box[r][i]+box[r][i-1] == 3 {
				box[r][i] = 3
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
func addItem() bool {
	// 从空的地方插入
	v := rand.Intn(2) + 1
	for {
		dir := rand.Intn(4) + 1
		switch dir {
		case 0: // 从上
			col := rand.Intn(4) + 1
			if box[0][col] == 0 {
				box[0][col] = v
				return true
			}
		case 1: // 从下
		case 2: // 从左
		case 3: // 从右
		}
	}
	return true
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
