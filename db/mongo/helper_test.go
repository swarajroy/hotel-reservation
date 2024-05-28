package mongo

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type HelperSuite struct {
	suite.Suite
}

func (hs *HelperSuite) TestNewTestMongoClient() {
	var dbName = "hotel-reservation-test"
	client, err := NewTestMongoClient(dbName)

	hs.Nil(err)
	hs.NotNil(client)
	hs.NotNil(client.Container)
	hs.NotEmpty(client.Uri)
	hs.Equal(dbName, client.DbName)
}

func TestHelperSuite(t *testing.T) {
	suite.Run(t, new(HelperSuite))
}
