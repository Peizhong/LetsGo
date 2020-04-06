package main

// runtime/chan.go

func WhatChan() {
	// 在堆上分配
	// CAS: compare and swap: 内存地址，旧的预期值，目标值
	// 重复关闭channel，panic
	// 向关闭的channel发送，panic
	// 从关闭的channel读，返回默认值

	// 结构hchan
	// waitq->sudog 携程
	// qcount, dataqsiz 缓存数量
	// buf缓存

	// hchan的类型，分配内存

	// send:
	// 1. 如果有recv等待：从接收waitq取出，send，goready
	// 2. 如果buf有空间，typedmemmove数据
	// 3. buf满了，gopark挂起

	// close:
	// 唤醒recvq和sendq的阻塞goroutine

	// select:
	//

	// checkdead(): 检查运行的M，如果等于0，dead
}
