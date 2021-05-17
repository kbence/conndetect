package connlib

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"strings"
	"testing"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type ProcNetTcpSuite struct{}

var _ = Suite(&ProcNetTcpSuite{})

func (s *ProcNetTcpSuite) TestParseTCPFile(c *C) {
	fileContents := `
   1: 00000000:1F90 00000000:0000 0A 00000000:00000000 00:00000000 00000000     0        0 112218 1 0000000000000000 100 0 0 10 0
   2: 0100007F:04D2 00000000:0000 0A 00000000:00000000 00:00000000 00000000     0        0 115249 1 0000000000000000 100 0 0 10 0
   3: 0100007F:8640 0100007F:1F90 01 00000000:00000000 02:00000FAE 00000000  1000        0 135463 2 0000000000000000 20 4 0 10 -1
`

	connections, err := ParseTCPFile(strings.NewReader(fileContents))

	c.Assert(err, IsNil)
	c.Assert(len(connections), Equals, 3)
	c.Assert(connections[0], Equals, Connection{Local: Endpoint{Port: 8080}})
	c.Assert(connections[1], Equals, Connection{
		Local: Endpoint{IP: IPv4Address{127, 0, 0, 1}, Port: 1234},
	})
	c.Assert(connections[2], Equals, Connection{
		Local:  Endpoint{IP: IPv4Address{127, 0, 0, 1}, Port: 34368},
		Remote: Endpoint{IP: IPv4Address{127, 0, 0, 1}, Port: 8080},
	})
}

func (s *ProcNetTcpSuite) TestPartTCPFileWithLongSequenceNumbers(c *C) {
	fileContents := `
1001: 00000000:1F90 00000000:0000 0A 00000000:00000000 00:00000000 00000000     0        0 112218 1 0000000000000000 100 0 0 10 0
1002: 0100007F:04D2 00000000:0000 0A 00000000:00000000 00:00000000 00000000     0        0 115249 1 0000000000000000 100 0 0 10 0
1003: 0100007F:8640 0100007F:1F90 01 00000000:00000000 02:00000FAE 00000000  1000        0 135463 2 0000000000000000 20 4 0 10 -1
`

	connections, err := ParseTCPFile(strings.NewReader(fileContents))

	c.Assert(err, IsNil)
	c.Assert(len(connections), Equals, 3)
	c.Assert(connections[0], Equals, Connection{Local: Endpoint{Port: 8080}})
	c.Assert(connections[1], Equals, Connection{
		Local: Endpoint{IP: IPv4Address{127, 0, 0, 1}, Port: 1234},
	})
	c.Assert(connections[2], Equals, Connection{
		Local:  Endpoint{IP: IPv4Address{127, 0, 0, 1}, Port: 34368},
		Remote: Endpoint{IP: IPv4Address{127, 0, 0, 1}, Port: 8080},
	})
}

func (s *ProcNetTcpSuite) TestReadEstablishedTCPConnections(c *C) {
	// TODO: this tests ParseTCPFile again (although it doesn't care about the results)
	// Let's refactor this if the time allows to separate the two logic into testable units
	fileContents := `
	1: 00000000:1F90 00000000:0000 0A 00000000:00000000 00:00000000 00000000     0        0 112218 1 0000000000000000 100 0 0 10 0
	2: 0100007F:04D2 00000000:0000 0A 00000000:00000000 00:00000000 00000000     0        0 115249 1 0000000000000000 100 0 0 10 0
	3: 0100007F:8640 0100007F:1F90 01 00000000:00000000 02:00000FAE 00000000  1000        0 135463 2 0000000000000000 20 4 0 10 -1
`

	testFile := fmt.Sprintf("%s/tcp", c.MkDir())
	ioutil.WriteFile(testFile, []byte(fileContents), 0644)

	connections, err := ReadEstablishedTCPConnections(testFile)

	c.Check(err, Equals, nil)
	c.Check(len(connections.Established), Equals, 1)
	c.Check(len(connections.Listening), Equals, 2)
}

func (s *ProcNetTcpSuite) TestReadEstablishedTCPConnectionsFileDoesNotExist(c *C) {
	testFile := fmt.Sprintf("%s/doesnotexist", c.MkDir())

	connections, err := ReadEstablishedTCPConnections(testFile)

	c.Check(err, FitsTypeOf, &fs.PathError{})
	c.Check(connections, Equals, (*CategorizedConnections)(nil))
}

func (s *ProcNetTcpSuite) TestCalculateDirection(c *C) {
	localServer1 := Endpoint{IP: IPv4Address{0, 0, 0, 0}, Port: 1234}
	localServer2 := Endpoint{IP: IPv4Address{5, 6, 7, 8}, Port: 1234}
	localServer3 := Endpoint{IP: IPv4Address{1, 2, 3, 4}, Port: 1234}

	localClient1 := Endpoint{IP: IPv4Address{127, 0, 0, 1}, Port: 33432}
	localClient2 := Endpoint{IP: IPv4Address{1, 2, 3, 4}, Port: 47384}
	remoteServer := Endpoint{IP: IPv4Address{54, 82, 13, 45}, Port: 1235}
	remoteClient := Endpoint{IP: IPv4Address{54, 82, 13, 46}, Port: 56238}

	listeners := ConnectionList{
		Connection{Local: localServer1},
		Connection{Local: localServer2},
		Connection{Local: localServer3},
	}

	testCases := []struct {
		Connection Connection
		Expected   DirectionalConnection
	}{
		{
			// Local connection where the client is the local
			Connection: Connection{Local: localClient1, Remote: localServer2},
			Expected:   DirectionalConnection{localClient1, localServer2},
		},
		{
			// Local connection where the client is the remote
			Connection: Connection{Local: localServer3, Remote: localClient2},
			Expected:   DirectionalConnection{localClient2, localServer3},
		},
		{
			// Outgoing connection
			Connection: Connection{Local: localClient1, Remote: remoteServer},
			Expected:   DirectionalConnection{localClient1, remoteServer},
		},
		{
			// Incoming connection
			Connection: Connection{Local: localServer2, Remote: remoteClient},
			Expected:   DirectionalConnection{remoteClient, localServer2},
		},
		{
			// Another incoming connection
			Connection: Connection{Local: remoteServer, Remote: remoteClient},
			Expected:   DirectionalConnection{remoteClient, remoteServer},
		},
	}

	for _, testCase := range testCases {
		fmt.Println(listeners, testCase.Connection)
		c.Check(CalculateDirection(listeners, testCase.Connection), Equals, testCase.Expected)
	}
}
