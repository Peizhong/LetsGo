package cmd

import (
	"bufio"
	"fmt"
	"github.com/peizhong/letsgo/playground/gossip/app"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var (
	name string
	port int
	join string
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "run anything to the screen",
	Long: `run is for printing anything back to the screen.
For many years people have printed back to the screen.`,
	Run: func(cmd *cobra.Command, args []string) {
		run()
	},
}

func addRunCmd() {
	runCmd.Flags().StringVarP(&name, "name", "n", "", "service name")
	runCmd.Flags().IntVarP(&port, "port", "p", 8080, "port to serve")
	runCmd.Flags().StringVarP(&join, "join", "j", "", "host to join")

	rootCmd.AddCommand(runCmd)
}

func run() {
	store := app.NewStorage(name, port, join)
	if err := store.Load(); err != nil {
		log.Println(err.Error())
		return
	}
	reader := bufio.NewScanner(os.Stdin)
	for reader.Scan() {
		line := reader.Text()
		cmd := strings.Split(string(line), " ")
		l := len(cmd)
		switch cmd[0] {
		case "set":
			if l >= 3 {
				store.Set(cmd[1], cmd[2])
			}
			break
		case "del":
			if l >= 2 {
				store.Delete(cmd[1])
			}
			break
		case "get":
			if l >= 2 {
				if v, ok := store.Get(cmd[1]); ok {
					fmt.Println(v)
				} else {
					fmt.Println("not found")
				}
			}
			break
		case "gets":
			for _, v := range store.Gets() {
				fmt.Println(v)
			}
			break
		case "members":
			for _, member := range store.Members() {
				fmt.Println(member)
			}
			break
		case "exit":
			fmt.Println("bye")
			if err := store.Save(); err != nil {
				log.Println(err.Error())
			}
			return
		default:
			fmt.Println("I don't know", cmd[0])
			break
		}
	}
}
