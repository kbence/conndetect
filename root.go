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
	Run: func(cmd *cobra.Command, args []string) {
		ticker := time.Tick(10 * time.Second)

		oldConnections := connlib.ReadEstabilishedTCPConnections().ToMap()

		for {
			<-ticker
			connections := connlib.ReadEstabilishedTCPConnections().ToMap()

			for conn := range connections {
				if _, found := oldConnections[conn]; !found {
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
