package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var (
	name    string
	port    int
	join    string
	storage map[string]interface{}
	friends map[string]string
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "run anything to the screen",
	Long: `run is for printing anything back to the screen.
For many years people have printed back to the screen.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(fmt.Sprintf("run: at :%v, join to  %v", port, join))
		run()
	},
}

func loadRunCmd() {
	runCmd.Flags().StringVarP(&name, "name", "n", "default", "service name")
	runCmd.Flags().IntVarP(&port, "port", "p", 1, "port to serve")
	runCmd.Flags().StringVarP(&join, "join", "j", "", "host to join")

	rootCmd.AddCommand(runCmd)
}

func run() {
	storage = make(map[string]interface{})
	friends = map[string]string{
		name: fmt.Sprintf(":%v", port),
	}
	reader := bufio.NewReader(os.Stdin)
	for {
		line, _, err := reader.ReadLine()
		// 去掉换行
		if err != nil {
			fmt.Println(err)
			return
		}
		cmd := strings.Split(string(line), " ")
		len := len(cmd)

		switch cmd[0] {
		case "set":
			if len >= 3 {
				if _, ok := storage[cmd[1]]; ok {
					fmt.Println(0)
				} else {
					fmt.Println(1)
				}
				storage[cmd[1]] = cmd[2]
			}
			break
		case "get":
			if len >= 2 {
				if v, ok := storage[cmd[1]]; ok {
					fmt.Println(v)
				} else {
					fmt.Println("not found")
				}
			}
			break
		case "list":
			for k, v := range storage {
				fmt.Println(fmt.Sprintf("%v:lis %v", k, v))
			}
			break
		case "friends":
			for k, v := range friends {
				fmt.Println(fmt.Sprintf("%v:lis %v", k, v))
			}
			break
		case "exit":
			fmt.Println("bye")
			return
		default:
			fmt.Println("I don't know", cmd[0])
			break
		}
	}
}
