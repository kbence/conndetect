package main

import (
	"fmt"
	"time"

	"github.com/kbence/conndetect/internal/connlib"
	"github.com/spf13/cobra"
)

const PROC_NET_TCP_FILE = "/proc/net/tcp"
const TIME_FORMAT = "2006-01-02 15:04:05"

var rootCmdParams = struct {
	Wait    int
	TCPFile string
}{}

var rootCmd cobra.Command = cobra.Command{
	Use: "conndetect",
	RunE: func(cmd *cobra.Command, args []string) error {
		var oldConnections, connections *connlib.CategorizedConnections
		var err error

		ticker := time.Tick(time.Duration(rootCmdParams.Wait) * time.Second)

		if oldConnections, err = connlib.ReadEstablishedTCPConnections(rootCmdParams.TCPFile); err != nil {
			return err
		}

		for {
			<-ticker

			if connections, err = connlib.ReadEstablishedTCPConnections(rootCmdParams.TCPFile); err != nil {
				return err
			}

			oldConnectionsMap := oldConnections.Established.ToMap()

			for _, conn := range connections.Established {
				if _, found := oldConnectionsMap[conn]; !found {
					dirConn := connlib.CalculateDirection(connections.Listening, conn)

					fmt.Printf(
						"%s: New connection: %s -> %s\n",
						time.Now().Format(TIME_FORMAT),
						dirConn.Source.String(),
						dirConn.Destination.String(),
					)
				}
			}

			oldConnections = connections
		}
	},
}

func init() {
	rootCmd.Flags().IntVarP(&rootCmdParams.Wait, "wait", "w", 10, "wait time between scans")
	rootCmd.Flags().StringVarP(&rootCmdParams.TCPFile, "tcp-file", "f", PROC_NET_TCP_FILE, "file to parse IPs out from")
}
