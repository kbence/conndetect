package connlib

import (
	"testing"

	. "gopkg.in/check.v1"
)

func TestConnectionSource(t *testing.T) { TestingT(t) }

type ConnectionSourceTestSuite struct{}

var _ = Suite(&ConnectionSourceTestSuite{})
