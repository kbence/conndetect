package main

import (
	"fmt"
	"time"

	"github.com/kbence/conndetect/internal/connlib"
	"github.com/spf13/cobra"
)

const TIME_FORMAT = "2006-01-02 15:04:05"

var rootCmd cobra.Command = cobra.Command{
	Use: "conndetect",
	RunE: func(cmd *cobra.Command, args []string) error {
		var oldConnections, connections connlib.ConnectionList
		var err error

		ticker := time.Tick(10 * time.Second)

		if oldConnections, err = connlib.ReadEstabilishedTCPConnections(); err != nil {
			return err
		}

		for {
			<-ticker

			if connections, err = connlib.ReadEstabilishedTCPConnections(); err != nil {
				return err
			}

			oldConnectionsMap := oldConnections.ToMap()

			for _, conn := range connections {
				if _, found := oldConnectionsMap[conn]; !found {
					fmt.Printf(
						"%s: New connection: %s -> %s\n",
						time.Now().Format(TIME_FORMAT),
						conn.Remote.String(),
						conn.Local.String(),
					)
				}
			}

			oldConnections = connections
		}
	},
}
