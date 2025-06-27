package dbsuite

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestDBSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}
