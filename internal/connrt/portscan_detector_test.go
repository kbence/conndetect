package connrt

import (
	"math/rand"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/gookit/event"
	"github.com/kbence/conndetect/internal/connlib"
	"github.com/kbence/conndetect/internal/ext_mock"
	"github.com/kbence/conndetect/internal/utils_mock"
	. "gopkg.in/check.v1"
)

var _ = Suite(&PortscanDetectorTestSuite{})

type PortscanDetectorTestSuite struct{}

var testPortscanSettings = NewPortscanSettings(3, 60*time.Second)

func (s *PortscanDetectorTestSuite) newConnectionToPort(port int) connlib.DirectionalConnection {
	return connlib.DirectionalConnection{
		Source:      connlib.Endpoint{IP: connlib.IPv4Address{1, 2, 3, 4}, Port: uint16(30000 + rand.Int()%20000)},
		Destination: connlib.Endpoint{IP: connlib.IPv4Address{5, 6, 7, 8}, Port: uint16(port)},
	}
}

func (s *PortscanDetectorTestSuite) newRandomConnection() connlib.DirectionalConnection {
	sourceIP := connlib.IPv4Address{byte(rand.Int() % 255), byte(rand.Int() % 255),
		byte(rand.Int() % 255), byte(rand.Int() % 255)}
	destIP := connlib.IPv4Address{byte(rand.Int() % 255), byte(rand.Int() % 255),
		byte(rand.Int() % 255), byte(rand.Int() % 255)}

	return connlib.DirectionalConnection{
		Source:      connlib.Endpoint{IP: sourceIP, Port: uint16(30000 + rand.Int()%20000)},
		Destination: connlib.Endpoint{IP: destIP, Port: uint16(1 + rand.Int()%20000)},
	}
}

func (s *PortscanDetectorTestSuite) getTime(timeStr string) time.Time {
	fakeTime, _ := time.Parse(TIME_FORMAT, timeStr)
	return fakeTime
}

func (s *PortscanDetectorTestSuite) TestPortscanDetectorDetectsScan(c *C) {
	ctrl := gomock.NewController(c)
	defer ctrl.Finish()

	eventManagerMock := ext_mock.NewMockManagerFace(ctrl)
	printerMock := utils_mock.NewMockPrinter(ctrl)
	timeMock := utils_mock.NewTimeTravelingMock(s.getTime("2021-05-19 09:59:34"))

	eventManagerMock.EXPECT().On("newConnection", gomock.Any())

	detector := NewPortscanDetector(eventManagerMock, testPortscanSettings)
	detector.printer = printerMock
	detector.time = timeMock

	printerMock.
		EXPECT().
		Printf("%s: Port scan detected: %s -> %s on ports %s\n",
			"2021-05-19 09:59:34",
			connlib.IPv4Address{1, 2, 3, 4},
			connlib.IPv4Address{5, 6, 7, 8},
			"80,443,1234")

	detector.Handle(event.NewBasic("newConnection", event.M{"connection": s.newConnectionToPort(80)}))
	timeMock.ForwardBy(2 * time.Second)
	detector.Handle(event.NewBasic("newConnection", event.M{"connection": s.newConnectionToPort(443)}))
	timeMock.ForwardBy(57 * time.Second)

	err := detector.Handle(event.NewBasic("newConnection", event.M{"connection": s.newConnectionToPort(1234)}))

	c.Check(err, IsNil)
}

func (s *PortscanDetectorTestSuite) TestPortscanDetectorReportsOnlyOnce(c *C) {
	ctrl := gomock.NewController(c)
	defer ctrl.Finish()

	eventManagerMock := ext_mock.NewMockManagerFace(ctrl)
	printerMock := utils_mock.NewMockPrinter(ctrl)
	timeMock := utils_mock.NewTimeTravelingMock(s.getTime("2021-05-19 09:59:34"))

	eventManagerMock.EXPECT().On("newConnection", gomock.Any())

	detector := NewPortscanDetector(eventManagerMock, testPortscanSettings)
	detector.printer = printerMock
	detector.time = timeMock

	printerMock.
		EXPECT().
		Printf("%s: Port scan detected: %s -> %s on ports %s\n",
			"2021-05-19 09:59:34",
			connlib.IPv4Address{1, 2, 3, 4},
			connlib.IPv4Address{5, 6, 7, 8},
			"80,443,1234")

	detector.Handle(event.NewBasic("newConnection", event.M{"connection": s.newConnectionToPort(80)}))
	timeMock.ForwardBy(2 * time.Second)
	detector.Handle(event.NewBasic("newConnection", event.M{"connection": s.newConnectionToPort(443)}))
	timeMock.ForwardBy(57 * time.Second)
	detector.Handle(event.NewBasic("newConnection", event.M{"connection": s.newConnectionToPort(1234)}))

	err := detector.Handle(event.NewBasic("newConnection", event.M{"connection": s.newConnectionToPort(5432)}))

	c.Check(err, IsNil)
}

func (s *PortscanDetectorTestSuite) TestPortscanDetectorDoesNothingOnConnectionsFromDifferentIPs(c *C) {
	ctrl := gomock.NewController(c)
	defer ctrl.Finish()

	eventManagerMock := ext_mock.NewMockManagerFace(ctrl)
	printerMock := utils_mock.NewMockPrinter(ctrl)
	timeMock := utils_mock.NewTimeTravelingMock(s.getTime("2021-05-19 09:59:34"))

	eventManagerMock.EXPECT().On("newConnection", gomock.Any())

	detector := NewPortscanDetector(eventManagerMock, testPortscanSettings)
	detector.printer = printerMock
	detector.time = timeMock

	detector.Handle(event.NewBasic("newConnection", event.M{"connection": s.newRandomConnection()}))
	timeMock.ForwardBy(2 * time.Second)
	detector.Handle(event.NewBasic("newConnection", event.M{"connection": s.newRandomConnection()}))
	timeMock.ForwardBy(57 * time.Second)

	err := detector.Handle(event.NewBasic("newConnection", event.M{"connection": s.newRandomConnection()}))

	c.Check(err, IsNil)
}

func (s *PortscanDetectorTestSuite) TestPortscanDetectorDoesNothingOnRandomConnections(c *C) {
	ctrl := gomock.NewController(c)
	defer ctrl.Finish()

	eventManagerMock := ext_mock.NewMockManagerFace(ctrl)
	printerMock := utils_mock.NewMockPrinter(ctrl)
	timeMock := utils_mock.NewTimeTravelingMock(s.getTime("2021-05-19 09:59:34"))

	eventManagerMock.EXPECT().On("newConnection", gomock.Any())

	detector := NewPortscanDetector(eventManagerMock, testPortscanSettings)
	detector.printer = printerMock
	detector.time = timeMock

	detector.Handle(event.NewBasic("newConnection", event.M{"connection": s.newRandomConnection()}))
	timeMock.ForwardBy(2 * time.Second)
	detector.Handle(event.NewBasic("newConnection", event.M{"connection": s.newRandomConnection()}))
	timeMock.ForwardBy(57 * time.Second)

	err := detector.Handle(event.NewBasic("newConnection", event.M{"connection": s.newRandomConnection()}))

	c.Check(err, IsNil)
}
