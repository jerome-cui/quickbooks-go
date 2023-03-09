package examples

import (
	"fmt"
	"testing"

	"github.com/jerome-cui/quickbooks-go"
	"github.com/stretchr/testify/require"
)

func TestReuseToken(t *testing.T) {
	qbClient, err := newTestClient()
	require.NoError(t, err)

	// Make a request!
	info, err := qbClient.FindCompanyInfo()
	require.NoError(t, err)
	fmt.Println(info)
}

func TestRefreshToken(t *testing.T) {
	qbClient, err := newTestClient()
	require.NoError(t, err)

	// Make a request!
	token, err := qbClient.RefreshToken("xxx")
	require.NoError(t, err)
	fmt.Println(token)
}

func TestRevokeToken(t *testing.T) {
	qbClient, err := newTestClient()
	require.NoError(t, err)

	// Make a request!
	err = qbClient.RevokeToken("xxx")
	require.NoError(t, err)
	fmt.Println(err)
}

func newTestClient() (*quickbooks.Client, error) {
	clientId := "xx"
	clientSecret := "xx"
	realmId := "xx"

	token := quickbooks.BearerToken{
		RefreshToken: "xx",
		AccessToken:  "xxx",
	}

	return quickbooks.NewClient(clientId, clientSecret, realmId, false, "", &token)
}
