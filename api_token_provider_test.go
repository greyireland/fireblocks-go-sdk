package fireblocksdk

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const (
	// https://phpseclib.com/docs/rsa-keys
	privateKey = "-----BEGIN PRIVATE KEY-----\nMIIBVAIBADANBgkqhkiG9w0BAQEFAASCAT4wggE6AgEAAkEAqPfgaTEWEP3S9w0t\ngsicURfo+nLW09/0KfOPinhYZ4ouzU+3xC4pSlEp8Ut9FgL0AgqNslNaK34Kq+NZ\njO9DAQIDAQABAkAgkuLEHLaqkWhLgNKagSajeobLS3rPT0Agm0f7k55FXVt743hw\nNgkp98bMNrzy9AQ1mJGbQZGrpr4c8ZAx3aRNAiEAoxK/MgGeeLui385KJ7ZOYktj\nhLBNAB69fKwTZFsUNh0CIQEJQRpFCcydunv2bENcN/oBTRw39E8GNv2pIcNxZkcb\nNQIgbYSzn3Py6AasNj6nEtCfB+i1p3F35TK/87DlPSrmAgkCIQDJLhFoj1gbwRbH\n/bDRPrtlRUDDx44wHoEhSDRdy77eiQIgE6z/k6I+ChN1LLttwX0galITxmAYrOBh\nBVl433tgTTQ=\n-----END PRIVATE KEY-----"
)

func TestAuthTokenSuite(t *testing.T) {
	suite.Run(t, new(AuthTokenSuite))
}

type AuthTokenSuite struct {
	suite.Suite
	auth         IAuthProvider
	apiKey       string
	apiSecretKey string
	timeProvider ITimeProvider
}

type testTimeProvider struct{}

func (tp *testTimeProvider) Now() time.Time {
	return time.Unix(1000, 0)
}

func (suite *AuthTokenSuite) SetupTest() {
	suite.apiKey = "apiKey"
	suite.apiSecretKey = privateKey
	suite.timeProvider = &testTimeProvider{}
	suite.auth, _ = NewAuthProvider(
		suite.apiKey,
		suite.apiSecretKey,
		WithTimeProvider(suite.timeProvider),
	)
}

func (suite *AuthTokenSuite) TestWithCorrectPrivateKey() {
	token, err := suite.auth.SignJwt("", "")
	require.NoError(suite.T(), err)
	require.Equal(suite.T(), "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJib2R5SGFzaCI6IjMxMzI2MTY1MzMzMjYzNjIzMTY1NjMzMDMyNjQzMDMxNjU2NDYxMzMzNTM4MzE2MjMxMzIzNzYzMzE2NjY1NjUzMzYyMzA2NDYzMzUzMzM1MzczMjY1NjQzNjYyNjE2NjMyMzMzOTM3MzIzMTYxMzAzMzY0MzgzMjY1MzEzMjM2IiwiZXhwIjoxMDEwLCJub25jZSI6MTAwMCwibm93IjoxMDAwLCJzdWIiOiJhcGlLZXkiLCJ1cmkiOiIifQ.iwU8RpQMnkuqNOBRfXMSW9hEotw5tHWH1hfsNOJdFD9NAcLPtKEVRHwtGOOpO1DpTWYvy8zomgNh8CB25SmPxQ", token)
}

func (suite *AuthTokenSuite) TestMustFailWithNotRSAPrivateKey() {
	auth, err := NewAuthProvider(
		suite.apiKey,
		"fakeKey",
		WithTimeProvider(suite.timeProvider),
	)
	require.NoError(suite.T(), err)

	_, err = auth.SignJwt("", "")
	require.ErrorAs(suite.T(), err, &jwt.ErrNotRSAPrivateKey)
}

func (suite *AuthTokenSuite) TestTimeProvider() {
	unix := suite.timeProvider.Now().Unix()
	require.Equal(suite.T(), int64(1000), unix)
}

func (suite *AuthTokenSuite) TestGetApiKey() {
	require.Equal(suite.T(), suite.auth.GetApiKey(), suite.apiKey)
}