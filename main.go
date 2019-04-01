package main

import (
	"log"
	"os"

	"github.com/peizhong/letsgo/framework"
)

// 如果导入了多个包，先初始化包的参数，然后init()，最后执行package的main()
func init() {
	// 将日志输出到标准输出
	log.SetOutput(os.Stdout)
}

// 每个package必须有个main
func main() {
	framework.RunServer(8080)
}
