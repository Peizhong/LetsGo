package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"

	"github.com/peizhong/letsgo/pkg/data"
	"github.com/peizhong/letsgo/playground/gossip/app"

	"github.com/spf13/cobra"
)

var (
	name      string
	port      int
	broadcast int
	join      string

	store *app.Storage
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
	runCmd.Flags().IntVarP(&port, "port", "p", 8001, "port to serve")
	runCmd.Flags().StringVarP(&join, "join", "j", "", "host to join")
	runCmd.Flags().IntVarP(&broadcast, "broadcast", "b", 0, "api service")

	rootCmd.AddCommand(runCmd)
}

func run() {
	store = app.NewStorage(name, port, join)
	if err := store.Load(); err != nil {
		log.Println(err.Error())
		return
	}
	if broadcast > 0 {
		if err := serve(broadcast); err != nil {
			log.Println(err.Error())
			return
		}
	}
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, os.Kill)
	<-c
	println("recv close signal")
	if false {
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
			case "bench":
				if l >= 2 {
					if count, ok := data.IntTryParse(cmd[1]); ok {
						store.Benchmark(count)
					}
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
}
