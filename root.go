package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gookit/event"
	"github.com/kbence/conndetect/internal/connlib"
	"github.com/kbence/conndetect/internal/connrt"
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

func (c *rootCmdImpl) RunE(cmd *cobra.Command, args []string) error {
	eventManager := event.NewManager("conndetect")

	ticker := connrt.NewTicker(eventManager, rootCmdArgs.Wait)
	connrt.NewConnectionReader(eventManager, rootCmdArgs.TCPFile)
	connrt.NewConnectionPrinter(eventManager)

	exitSignal := make(chan os.Signal)
	signal.Notify(exitSignal, syscall.SIGINT)
	signal.Notify(exitSignal, syscall.SIGTERM)

	go func() {
		signal := <-exitSignal
		fmt.Printf("Received signal '%s', exiting...\n", signal)
		ticker.Stop()
		event.Reset()
	}()

	ticker.Run()

	return nil
}

var rootCmd cobra.Command = cobra.Command{
	Use:  "conndetect",
	RunE: newRootCmd().RunE,
}

func init() {
	rootCmd.Flags().IntVarP(&rootCmdArgs.Wait, "wait", "w", 10, "wait time between scans")
	rootCmd.Flags().StringVarP(&rootCmdArgs.TCPFile, "tcp-file", "f", PROC_NET_TCP_FILE, "file to parse IPs out from")
}
