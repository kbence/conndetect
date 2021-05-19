package connrt

import (
	"errors"

	"github.com/golang/mock/gomock"
	"github.com/gookit/event"
	"github.com/kbence/conndetect/internal/connlib"
	"github.com/kbence/conndetect/internal/connlib_mock"
	"github.com/kbence/conndetect/internal/ext_mock"
	. "gopkg.in/check.v1"
)

var _ = Suite(&ConnectionReaderTestSuite{})

type ConnectionReaderTestSuite struct{}

func (s *ConnectionReaderTestSuite) TestNewConnectionReaderReturnsError(c *C) {
	ctrl := gomock.NewController(c)
	defer ctrl.Finish()

	eventManagerMock := ext_mock.NewMockManagerFace(ctrl)
	connSrcMock := connlib_mock.NewMockConnectionSource(ctrl)

	expectedError := errors.New("some error")

	connSrcMock.
		EXPECT().
		ReadEstablishedTCPConnections("/path/to/tcp").
		Return(nil, expectedError)

	_, err := NewConnectionReader(eventManagerMock, "/path/to/tcp", connSrcMock)

	c.Check(err, Equals, expectedError)
}

func (s *ConnectionReaderTestSuite) TestConnectionReader(c *C) {
	ctrl := gomock.NewController(c)
	defer ctrl.Finish()

	eventManagerMock := ext_mock.NewMockManagerFace(ctrl)
	connSrcMock := connlib_mock.NewMockConnectionSource(ctrl)

	connections := connlib.ConnectionList{
		connlib.Connection{
			Local:  connlib.Endpoint{IP: connlib.IPv4Address{1, 2, 3, 4}, Port: 45678},
			Remote: connlib.Endpoint{IP: connlib.IPv4Address{5, 6, 7, 8}, Port: 443},
		},
	}
	expectedConnection := connlib.DirectionalConnection{
		Source:      connections[0].Local,
		Destination: connections[0].Remote,
	}

	connSrcMock.
		EXPECT().
		ReadEstablishedTCPConnections("/path/to/tcp").
		Return(&connlib.CategorizedConnections{}, nil)
	connSrcMock.
		EXPECT().
		ReadEstablishedTCPConnections("/path/to/tcp").
		Return(&connlib.CategorizedConnections{Established: connections}, nil)
	eventManagerMock.EXPECT().On("tick", gomock.Any())
	eventManagerMock.EXPECT().Fire("newConnection", event.M{"connection": expectedConnection})

	reader, _ := NewConnectionReader(eventManagerMock, "/path/to/tcp", connSrcMock)
	reader.connectionSource = connSrcMock

	err := reader.Handle(event.NewBasic("tick", event.M{}))

	c.Check(err, IsNil)
}
