package argparse

// Copyright (c) 2017 by Gilbert Ramirez <gram@alumni.rice.edu>

import (
	"log"
	"testing"

	. "gopkg.in/check.v1"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) {
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Lshortfile)
	TestingT(t)
}

type MySuite struct{}

var _ = Suite(&MySuite{})
