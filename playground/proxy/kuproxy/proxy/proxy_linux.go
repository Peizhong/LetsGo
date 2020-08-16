// +build linux,amd64

package proxy

type epoll struct {
	fd int // epoll fd
}

// AddRead ...
func (p *epoll) AddRead(fd int) {
	if err := syscall.EpollCtl(p.fd, syscall.EPOLL_CTL_ADD, fd,
		&syscall.EpollEvent{Fd: int32(fd),
			Events: syscall.EPOLLIN,
		},
	); err != nil {
		panic(err)
	}
}

func (p *epoll) AddReadWrite(fd int) {
	syscall.Kqueue()
	if err := syscall.EpollCtl(p.fd, syscall.EPOLL_CTL_ADD, fd,
		&syscall.EpollEvent{Fd: int32(fd),
			Events: syscall.EPOLLIN | syscall.EPOLLOUT,
		},
	); err != nil {
		panic(err)
	}
}

func (p *epoll) Wait(fn func(fd int) error) error {
	events := make([]syscall.EpollEvent, 64)
	for {
		n, err := syscall.EpollWait(p.fd, events, 100)
		if err != nil && err != syscall.EINTR {
			log.Println("wait err", err.Error())
			return err
		}
		for i := 0; i < n; i++ {
			fd := int(events[i].Fd)
			if err := fn(fd); err != nil {
				log.Println("fn err", err.Error())
				return err
			}
		}
	}
}

// 使用epoll
func (p *proxy) serve() {
	// 获得tcp监听的文件描述符
	file, err := p.ln.File()
	if err != nil {
		panic(err)
	}
	lnFd := int(file.Fd())
	workerCnt := runtime.GOMAXPROCS(0)
	var index = 0
	acceptBalancer := func(i int) bool {
		if i == (index % workerCnt) {
			index += 1
			return true
		}
		return false
	}
	log.Println("start ", workerCnt, "workers")
	for i := 0; i < 1; i++ {
		index := i
		go func() {
			// 设置epoll
			epollFd, err := syscall.EpollCreate1(0)
			if err != nil {
				panic(err)
			}
			ep := epoll{fd: epollFd}
			// 将tcp监听描述符设为非阻塞
			syscall.SetNonblock(lnFd, true)
			// 将tcp监听描述符添加到epoll
			ep.AddRead(lnFd)

			fdconns := make(map[int]int)
			ep.Wait(func(fd int) error {
				if !acceptBalancer(index) {
					// return nil
				}
				if _, ok := fdconns[fd]; !ok {
					// accept新的连接
					nfd, _, err := syscall.Accept(fd)
					if err != nil {
						log.Println("accept err", err)
						if err == syscall.EAGAIN {
							return nil
						}
						// 惊群，accept失败
						return nil
					}
					log.Println("accept", nfd)
					header := make([]byte, 512)
					n, err := syscall.Read(nfd, header)
					if err != nil {
						log.Println(err.Error())
					}
					if n < 0 {
						log.Println("can't read")
						return nil
					}
					header = header[:n]
					if err := syscall.SetNonblock(nfd, true); err != nil {
						log.Println("set noblocking err", err)
						return err
					}
					ep.AddRead(nfd)

					// do upstream
					upFd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
					if err != nil {
						log.Println("socket err", err.Error())
						return err
					}
					err = syscall.Connect(upFd, &syscall.SockaddrInet4{
						Port: 3000,
						Addr: [...]byte{192, 168, 3, 143},
					})
					if err != nil {
						log.Println("conn err", err.Error())
						return err
					}
					log.Println("connect ", upFd)
					log.Println("write back", string(header))
					syscall.Write(upFd, header)
					if err := syscall.SetNonblock(upFd, true); err != nil {
						log.Println("set noblocking err", err)
						return err
					}
					fdconns[nfd] = upFd
					fdconns[upFd] = nfd
					ep.AddRead(upFd)
				} else {
					// 1. 直接发送  2. 看返回值, 没发送完才挂上EPOLLOUT  3. 发送完就把EPOLLOUT 取消掉
					in := make([]byte, 512)
					n, err := syscall.Read(fd, in)
					if n == 0 || err != nil {
						if err != nil {
							if err == syscall.EAGAIN {
							} else {
								log.Println("unknow read err", err.Error())
							}
						}
						//close(fd)后，epoll中注册的相关sock_fd，就会自动被清理？
						log.Println("close", fd)
						syscall.Close(fd)
						if tfd, ok := fdconns[fd]; ok {
							syscall.Close(tfd)
							delete(fdconns, tfd)
							log.Println("close2 ", tfd)
						}
						delete(fdconns, fd)
						return nil
					}
					in = in[:n]
					if othFd, ok := fdconns[fd]; ok {
						_, err := syscall.Write(othFd, in)
						if err != nil {
							log.Println("write err", err.Error())
							return err
						}
					}
					return nil
				}
				return nil
			})
		}()
	}
	select {}
}
