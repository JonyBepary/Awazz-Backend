package model

import (
	"testing"

	"github.com/SohelAhmedJoni/Awazz-Backend/internal/durable"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type TokenTestSuite struct {
	suite.Suite
}

func (suite *TokenTestSuite) SetupTest() {
	// Set up any test data here
}

func (suite *TokenTestSuite) TearDownTest() {
	// Clean up any test data here
}

func (suite *TokenTestSuite) TestSaveToken() {
	// Set up a new Token object with some data
	token := &Token{
		UserName:     "testuser",
		Token:        "testtoken",
		GenerateTime: 1234567890,
	}

	// Call the SaveToken method
	err := token.SaveToken()
	require.NoError(suite.T(), err)

	// Check that the token was saved correctly
	db, err := durable.CreateDatabase("Database/", "Common", "Shard_0.sqlite")
	require.NoError(suite.T(), err)
	defer db.Close()

	var savedToken Token
	err = db.QueryRow("SELECT UserName, Token, GenerateTime FROM TOKEN WHERE UserName = ?", token.UserName).Scan(&savedToken.UserName, &savedToken.Token, &savedToken.GenerateTime)
	require.NoError(suite.T(), err)

	assert.Equal(suite.T(), token.UserName, savedToken.UserName)
	assert.Equal(suite.T(), token.Token, savedToken.Token)
	assert.Equal(suite.T(), token.GenerateTime, savedToken.GenerateTime)
}

func TestTokenTestSuite(t *testing.T) {
	suite.Run(t, new(TokenTestSuite))
}
