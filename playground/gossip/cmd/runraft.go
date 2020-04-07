package cmd

import (
	"github.com/peizhong/letsgo/playground/gossip/app"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
)

var (
	id    uint64
	peers []uint
)

var raftCmd = &cobra.Command{
	Use:   "raft",
	Short: "another cluster communication",
	Long: `run is for printing anything back to the screen.
For many years people have printed back to the screen.`,
	Run: func(cmd *cobra.Command, args []string) {
		runraft()
	},
}

func addRunRaftCmd() {
	raftCmd.Flags().Uint64VarP(&id, "id", "i", 1, "cluster id")
	raftCmd.Flags().IntVarP(&port, "port", "p", 8001, "port to serve")
	raftCmd.Flags().UintSliceVarP(&peers, "peers", "j", nil, "peers to join")
	raftCmd.Flags().IntVarP(&broadcast, "broadcast", "b", 0, "api service")

	rootCmd.AddCommand(raftCmd)
}

func runraft() {
	app.NewRaftStorage(id, peers)
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, os.Kill)
	<-c
	println("recv close signal")
}
