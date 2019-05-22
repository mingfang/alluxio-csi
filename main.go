package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/mingfang/alluxio-csi/alluxio"

	"github.com/spf13/cobra"
)

var (
	endpoint string
	nodeID   string
)

func init() {
	flag.Set("logtostderr", "true")
}

func main() {

	flag.CommandLine.Parse([]string{})

	cmd := &cobra.Command{
		Use:   "Alluxio",
		Short: "CSI based Alluxio driver",
		Run: func(cmd *cobra.Command, args []string) {
			handle()
		},
	}

	cmd.Flags().AddGoFlagSet(flag.CommandLine)

	cmd.PersistentFlags().StringVar(&nodeID, "nodeid", "", "node id")
	cmd.MarkPersistentFlagRequired("nodeid")

	cmd.PersistentFlags().StringVar(&endpoint, "endpoint", "", "CSI endpoint")
	cmd.MarkPersistentFlagRequired("endpoint")

	cmd.ParseFlags(os.Args[1:])
	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%s", err.Error())
		os.Exit(1)
	}

	os.Exit(0)
}

func handle() {
	d := alluxio.NewDriver(nodeID, endpoint)
	d.Run()
}