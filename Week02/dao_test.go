package Week02

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestSuite(t *testing.T) {
	suite.Run(t, &DaoSuite{})
}

type DaoSuite struct {
	suite.Suite
}

func (s *DaoSuite) TestSetUser() {
	err := SetUser()
	s.Equal(IsErrNoRows(err), true)
}
