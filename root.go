package main

import (
	"time"

	"github.com/kbence/conndetect/internal/connlib"
	"github.com/kbence/conndetect/internal/utils"
	"github.com/spf13/cobra"
)

const PROC_NET_TCP_FILE = "/proc/net/tcp"
const TIME_FORMAT = "2006-01-02 15:04:05"

type rootCmdParameters struct {
	Wait    int
	TCPFile string
}

var rootCmdArgs = rootCmdParameters{}

type rootCmdImpl struct {
	connectionSource connlib.ConnectionSource
	printer          utils.Printer
	time             utils.Time
	args             rootCmdParameters
	exit             chan interface{}
}

func newRootCmd() *rootCmdImpl {
	return &rootCmdImpl{
		connectionSource: connlib.NewConnectionSource(),
		printer:          utils.NewPrinter(),
		time:             utils.NewTime(),
		args:             rootCmdArgs,
		exit:             make(chan interface{}),
	}
}

func (c *rootCmdImpl) RunE() error {
	var oldConnections, connections *connlib.CategorizedConnections
	var err error

	ticker := time.Tick(time.Duration(c.args.Wait) * time.Second)

	if oldConnections, err = c.connectionSource.ReadEstablishedTCPConnections(c.args.TCPFile); err != nil {
		return err
	}

mainLoop:
	for {
		// Note: although this has been added to make it easy to test, it'll
		// also come handy if we want to handle signals in the future.
		select {
		case <-ticker:
			break
		case <-c.exit:
			break mainLoop
		}

		if connections, err = c.connectionSource.ReadEstablishedTCPConnections(c.args.TCPFile); err != nil {
			return err
		}

		oldConnectionsMap := oldConnections.Established.ToMap()

		for _, conn := range connections.Established {
			if _, found := oldConnectionsMap[conn]; !found {
				dirConn := connlib.CalculateDirection(connections.Listening, conn)

				c.printer.Printf(
					"%s: New connection: %s -> %s\n",
					c.time.Now().Format(TIME_FORMAT),
					dirConn.Source.String(),
					dirConn.Destination.String(),
				)
			}
		}

		oldConnections = connections
	}

	return nil
}

var rootCmd cobra.Command = cobra.Command{
	Use: "conndetect",
	RunE: func(cmd *cobra.Command, args []string) error {
		return newRootCmd().RunE()
	},
}

func init() {
	rootCmd.Flags().IntVarP(&rootCmdArgs.Wait, "wait", "w", 10, "wait time between scans")
	rootCmd.Flags().StringVarP(&rootCmdArgs.TCPFile, "tcp-file", "f", PROC_NET_TCP_FILE, "file to parse IPs out from")
}
