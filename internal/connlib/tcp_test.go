package connlib

import (
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
