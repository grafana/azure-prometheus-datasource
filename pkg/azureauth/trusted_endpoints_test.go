package azureauth

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"github.com/grafana/grafana-azure-sdk-go/v2/azcredentials"
	"github.com/grafana/grafana-azure-sdk-go/v2/azsettings"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	sdkhttpclient "github.com/grafana/grafana-plugin-sdk-go/backend/httpclient"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetTrustedEndpoints(t *testing.T) {
	const publicResourceId = "https://prometheus.monitor.azure.com"

	currentUser := &azcredentials.AadCurrentUserCredentials{}
	obo := &azcredentials.AzureClientSecretOboCredentials{}
	msi := &azcredentials.AzureManagedIdentityCredentials{}
	clientSecret := &azcredentials.AzureClientSecretCredentials{}

	t.Run("should always restrict endpoints for current user credentials", func(t *testing.T) {
		endpoints, err := getTrustedEndpoints(currentUser, publicResourceId)
		require.NoError(t, err)
		assert.Equal(t, []string{"https://*.prometheus.monitor.azure.com"}, endpoints)
	})

	t.Run("should restrict endpoints for current user credentials when enforcement is off", func(t *testing.T) {
		t.Setenv(enforceTrustedEndpointsEnv, "false")

		endpoints, err := getTrustedEndpoints(currentUser, publicResourceId)
		require.NoError(t, err)
		assert.Equal(t, []string{"https://*.prometheus.monitor.azure.com"}, endpoints)
	})

	t.Run("should always restrict endpoints for on-behalf-of credentials", func(t *testing.T) {
		endpoints, err := getTrustedEndpoints(obo, publicResourceId)
		require.NoError(t, err)
		assert.Equal(t, []string{"https://*.prometheus.monitor.azure.com"}, endpoints)
	})

	t.Run("should not restrict endpoints for service credentials by default", func(t *testing.T) {
		for _, credentials := range []azcredentials.AzureCredentials{msi, clientSecret} {
			endpoints, err := getTrustedEndpoints(credentials, publicResourceId)
			require.NoError(t, err)
			assert.Nil(t, endpoints)
		}
	})

	t.Run("should restrict endpoints for service credentials when enforcement is on", func(t *testing.T) {
		t.Setenv(enforceTrustedEndpointsEnv, "true")

		for _, credentials := range []azcredentials.AzureCredentials{msi, clientSecret} {
			endpoints, err := getTrustedEndpoints(credentials, publicResourceId)
			require.NoError(t, err)
			assert.Equal(t, []string{"https://*.prometheus.monitor.azure.com"}, endpoints)
		}
	})

	t.Run("should use the Prometheus resource ID of the cloud", func(t *testing.T) {
		tests := map[string]string{
			"https://prometheus.monitor.azure.us":                 "https://*.prometheus.monitor.azure.us",
			"https://prometheus.monitor.azure.cn":                 "https://*.prometheus.monitor.azure.cn",
			"https://prometheus.monitor.azure.cloud1.contoso.com": "https://*.prometheus.monitor.azure.cloud1.contoso.com",
		}
		for resourceId, expected := range tests {
			endpoints, err := getTrustedEndpoints(currentUser, resourceId)
			require.NoError(t, err)
			assert.Equal(t, []string{expected}, endpoints)
		}
	})

	t.Run("should add user trusted endpoints when allowed", func(t *testing.T) {
		t.Setenv(allowUserTrustedEndpointsEnv, "true")
		t.Setenv(userTrustedEndpointsEnv, "https://*.custom1.com, https://custom2.com:8443,")

		endpoints, err := getTrustedEndpoints(currentUser, publicResourceId)
		require.NoError(t, err)
		assert.Equal(t, []string{
			"https://*.prometheus.monitor.azure.com",
			"https://*.custom1.com",
			"https://custom2.com:8443",
		}, endpoints)
	})

	t.Run("should ignore user trusted endpoints when not allowed", func(t *testing.T) {
		t.Setenv(userTrustedEndpointsEnv, "https://*.custom1.com")

		endpoints, err := getTrustedEndpoints(currentUser, publicResourceId)
		require.NoError(t, err)
		assert.Equal(t, []string{"https://*.prometheus.monitor.azure.com"}, endpoints)
	})

	t.Run("should not add user trusted endpoints when no restriction applies", func(t *testing.T) {
		t.Setenv(allowUserTrustedEndpointsEnv, "true")
		t.Setenv(userTrustedEndpointsEnv, "https://*.custom1.com")

		endpoints, err := getTrustedEndpoints(msi, publicResourceId)
		require.NoError(t, err)
		assert.Nil(t, endpoints)
	})

	t.Run("should return error for an invalid bool setting", func(t *testing.T) {
		t.Setenv(enforceTrustedEndpointsEnv, "invalid")

		_, err := getTrustedEndpoints(currentUser, publicResourceId)
		assert.ErrorContains(t, err, enforceTrustedEndpointsEnv)
	})

	t.Run("should return error for an invalid resource ID", func(t *testing.T) {
		_, err := getTrustedEndpoints(currentUser, "prometheus.monitor.azure.com")
		assert.Error(t, err)
	})
}

func TestConfigureAzureAuthenticationTrustedEndpoints(t *testing.T) {
	var tokenRequests atomic.Int32
	tokenServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenRequests.Add(1)
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer tokenServer.Close()

	var targetRequests atomic.Int32
	target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		targetRequests.Add(1)
		w.WriteHeader(http.StatusOK)
	}))
	defer target.Close()

	azureSettings := &azsettings.AzureSettings{
		AzureAuthEnabled:    true,
		UserIdentityEnabled: true,
		UserIdentityTokenEndpoint: &azsettings.TokenEndpointSettings{
			TokenUrl:     tokenServer.URL,
			ClientId:     "client-id",
			ClientSecret: "client-secret",
		},
	}
	settings := backend.DataSourceInstanceSettings{
		URL:      target.URL,
		JSONData: []byte(`{ "azureCredentials": { "authType": "currentuser" } }`),
	}

	send := func(t *testing.T) error {
		opts := &sdkhttpclient.Options{CustomOptions: map[string]any{}}
		err := ConfigureAzureAuthentication(settings, azureSettings, opts, log.NewNullLogger())
		require.NoError(t, err)

		client, err := sdkhttpclient.New(*opts)
		require.NoError(t, err)

		req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, target.URL+"/api/v1/query", nil)
		require.NoError(t, err)

		resp, err := client.Do(req)
		if resp != nil {
			_ = resp.Body.Close()
		}
		return err
	}

	t.Run("should block a current user request to an endpoint that is not trusted", func(t *testing.T) {
		tokenRequests.Store(0)
		targetRequests.Store(0)

		err := send(t)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "is not allowed by the datasource")
		assert.Zero(t, tokenRequests.Load(), "no token should be requested")
		assert.Zero(t, targetRequests.Load(), "the endpoint should not be called")
	})

	t.Run("should allow a current user request to a user trusted endpoint", func(t *testing.T) {
		t.Setenv(allowUserTrustedEndpointsEnv, "true")
		t.Setenv(userTrustedEndpointsEnv, target.URL)

		err := send(t)
		// The request has no signed-in user, so getting the token fails. It must fail after the
		// endpoint check, not because of it.
		if err != nil {
			assert.NotContains(t, err.Error(), "is not allowed by the datasource")
		}
	})
}
