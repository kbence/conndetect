package main

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/kbence/conndetect/internal/connlib"
	"github.com/kbence/conndetect/internal/connlib_mock"
	"github.com/kbence/conndetect/internal/utils_mock"
	. "gopkg.in/check.v1"
)

func TestRootCmd(t *testing.T) { TestingT(t) }

type RootCmdTestSuite struct{}

var _ = Suite(&RootCmdTestSuite{})

func (s *RootCmdTestSuite) TestRootCmd(c *C) {
	ctrl := gomock.NewController(c)
	defer ctrl.Finish()

	// pass all the mocks
	connSrc := connlib_mock.NewMockConnectionSource(ctrl)
	printer := utils_mock.NewMockPrinter(ctrl)
	tm := utils_mock.NewMockTime(ctrl)

	cmd := newRootCmd()
	cmd.connectionSource = connSrc
	cmd.printer = printer
	cmd.time = tm
	cmd.args = rootCmdParameters{Wait: 1, TCPFile: "/proc/net/tcp"}
	cmd.exit = make(chan interface{})

	// exit after 1.5s to give chance to run through
	go func() { time.Sleep(1500 * time.Millisecond); cmd.exit <- nil }()

	connections := connlib.ConnectionList{
		connlib.Connection{
			Local:  connlib.Endpoint{IP: connlib.IPv4Address{1, 2, 3, 4}, Port: 45678},
			Remote: connlib.Endpoint{IP: connlib.IPv4Address{5, 6, 7, 8}, Port: 443},
		},
	}

	fakeTime, _ := time.Parse(TIME_FORMAT, "2021-05-17 12:34:56")
	tm.EXPECT().Now().AnyTimes().Return(fakeTime)
	connSrc.
		EXPECT().
		ReadEstablishedTCPConnections("/proc/net/tcp").
		Return(&connlib.CategorizedConnections{}, nil)
	connSrc.
		EXPECT().
		ReadEstablishedTCPConnections("/proc/net/tcp").
		Return(&connlib.CategorizedConnections{Established: connections}, nil)
	printer.
		EXPECT().
		Printf("%s: New connection: %s -> %s\n", "2021-05-17 12:34:56", "1.2.3.4:45678", "5.6.7.8:443")

	err := cmd.RunE()

	c.Check(err, IsNil)
}
