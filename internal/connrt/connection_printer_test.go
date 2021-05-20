package connrt

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/gookit/event"
	"github.com/kbence/conndetect/internal/connlib"
	"github.com/kbence/conndetect/internal/ext_mock"
	"github.com/kbence/conndetect/internal/utils_mock"
	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

var _ = Suite(&ConnectionPrinterTestSuite{})

type ConnectionPrinterTestSuite struct{}

func (s *ConnectionPrinterTestSuite) TestConnectionPrinter(c *C) {
	ctrl := gomock.NewController(c)
	defer ctrl.Finish()

	eventManagerMock := ext_mock.NewMockManagerFace(ctrl)
	printerMock := utils_mock.NewMockPrinter(ctrl)
	timeMock := utils_mock.NewMockTime(ctrl)

	eventManagerMock.EXPECT().On("newConnection", gomock.Any())

	printer := NewConnectionPrinter(eventManagerMock)
	printer.printer = printerMock
	printer.time = timeMock

	connection := connlib.DirectionalConnection{
		Source:      connlib.Endpoint{IP: connlib.IPv4Address{1, 2, 3, 4}, Port: 45678},
		Destination: connlib.Endpoint{IP: connlib.IPv4Address{5, 6, 7, 8}, Port: 443},
	}

	fakeTime, _ := time.Parse(TIME_FORMAT, "2021-05-17 12:34:56")
	timeMock.EXPECT().Now().AnyTimes().Return(fakeTime)
	printerMock.
		EXPECT().
		Printf("%s: New connection: %s -> %s\n", "2021-05-17 12:34:56", "1.2.3.4:45678", "5.6.7.8:443")

	err := printer.Handle(event.NewBasic(eventNewConnection, event.M{"connection": connection}))

	c.Check(err, IsNil)
}

func (s *ConnectionPrinterTestSuite) TestConnectionPrinterDoesNothingWithNoConnectionInEvent(c *C) {
	ctrl := gomock.NewController(c)
	defer ctrl.Finish()

	eventManagerMock := ext_mock.NewMockManagerFace(ctrl)
	printerMock := utils_mock.NewMockPrinter(ctrl)
	timeMock := utils_mock.NewMockTime(ctrl)

	eventManagerMock.EXPECT().On("newConnection", gomock.Any())

	printer := NewConnectionPrinter(eventManagerMock)
	printer.printer = printerMock
	printer.time = timeMock

	err := printer.Handle(event.NewBasic(eventNewConnection, event.M{}))

	c.Check(err, IsNil)
}
