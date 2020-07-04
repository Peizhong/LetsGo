package main

func main() {
	r := &GoRedis{}
	r.Init()
	r.Blpop("a")
}
