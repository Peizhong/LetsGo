package cmd

import (
	"log"
	"os"
	"os/signal"

	"github.com/peizhong/letsgo/playground/gossip/app"

	"github.com/spf13/cobra"
)

var (
	name      string
	port      int
	broadcast int
	join      string
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
	store := app.NewStorage(name, port, join)
	if err := store.Load(""); err != nil {
		log.Println(err.Error())
		return
	}
	if broadcast > 0 {
		if err := serve(broadcast, store); err != nil {
			log.Println(err.Error())
			return
		}
	}
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, os.Kill)
	<-c
	println("recv close signal")
}
