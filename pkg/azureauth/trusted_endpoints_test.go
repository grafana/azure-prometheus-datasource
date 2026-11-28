package azureauth

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"github.com/grafana/grafana-azure-sdk-go/v2/azsettings"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	sdkhttpclient "github.com/grafana/grafana-plugin-sdk-go/backend/httpclient"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/grafana/grafana-plugin-sdk-go/config"
	"github.com/grafana/grafana-plugin-sdk-go/experimental/featuretoggles"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetTrustedEndpoints(t *testing.T) {
	const publicResourceId = "https://prometheus.monitor.azure.com"

	t.Run("should enforce trusted endpoints by default", func(t *testing.T) {
		endpoints, err := getTrustedEndpoints(context.Background(), publicResourceId)
		require.NoError(t, err)
		assert.Equal(t, []string{"https://*.prometheus.monitor.azure.com"}, endpoints)
	})

	t.Run("should enforce trusted endpoints when the setting is empty", func(t *testing.T) {
		t.Setenv(enforceTrustedEndpointsEnv, "")

		endpoints, err := getTrustedEndpoints(context.Background(), publicResourceId)
		require.NoError(t, err)
		assert.Equal(t, []string{"https://*.prometheus.monitor.azure.com"}, endpoints)
	})

	t.Run("should enforce trusted endpoints when the setting is true", func(t *testing.T) {
		t.Setenv(enforceTrustedEndpointsEnv, "true")

		endpoints, err := getTrustedEndpoints(context.Background(), publicResourceId)
		require.NoError(t, err)
		assert.Equal(t, []string{"https://*.prometheus.monitor.azure.com"}, endpoints)
	})

	t.Run("should not enforce trusted endpoints when the setting is false", func(t *testing.T) {
		t.Setenv(enforceTrustedEndpointsEnv, "false")

		endpoints, err := getTrustedEndpoints(context.Background(), publicResourceId)
		require.NoError(t, err)
		assert.Nil(t, endpoints)
	})

	t.Run("should not enforce trusted endpoints when the feature toggle is enabled", func(t *testing.T) {
		endpoints, err := getTrustedEndpoints(withFeatureToggles("someOtherToggle,"+disableTrustedEndpointsFeatureToggle), publicResourceId)
		require.NoError(t, err)
		assert.Nil(t, endpoints)
	})

	t.Run("should enforce trusted endpoints when other feature toggles are enabled", func(t *testing.T) {
		endpoints, err := getTrustedEndpoints(withFeatureToggles("someOtherToggle"), publicResourceId)
		require.NoError(t, err)
		assert.Equal(t, []string{"https://*.prometheus.monitor.azure.com"}, endpoints)
	})

	t.Run("should let an explicit setting override the feature toggle", func(t *testing.T) {
		tests := []struct {
			name     string
			value    string
			expected []string
		}{
			{name: "true", value: "true", expected: []string{"https://*.prometheus.monitor.azure.com"}},
			{name: "empty", value: "", expected: []string{"https://*.prometheus.monitor.azure.com"}},
			{name: "false", value: "false", expected: nil},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Setenv(enforceTrustedEndpointsEnv, tt.value)

				endpoints, err := getTrustedEndpoints(withFeatureToggles(disableTrustedEndpointsFeatureToggle), publicResourceId)
				require.NoError(t, err)
				assert.Equal(t, tt.expected, endpoints)
			})
		}
	})

	t.Run("should use the Prometheus resource ID of the cloud", func(t *testing.T) {
		tests := map[string]string{
			"https://prometheus.monitor.azure.us":                 "https://*.prometheus.monitor.azure.us",
			"https://prometheus.monitor.azure.cn":                 "https://*.prometheus.monitor.azure.cn",
			"https://prometheus.monitor.azure.cloud1.contoso.com": "https://*.prometheus.monitor.azure.cloud1.contoso.com",
		}
		for resourceId, expected := range tests {
			endpoints, err := getTrustedEndpoints(context.Background(), resourceId)
			require.NoError(t, err)
			assert.Equal(t, []string{expected}, endpoints)
		}
	})

	t.Run("should add user trusted endpoints when allowed", func(t *testing.T) {
		t.Setenv(enforceTrustedEndpointsEnv, "true")
		t.Setenv(allowUserTrustedEndpointsEnv, "true")
		t.Setenv(userTrustedEndpointsEnv, "https://*.custom1.com, https://custom2.com:8443,")

		endpoints, err := getTrustedEndpoints(context.Background(), publicResourceId)
		require.NoError(t, err)
		assert.Equal(t, []string{
			"https://*.prometheus.monitor.azure.com",
			"https://*.custom1.com",
			"https://custom2.com:8443",
		}, endpoints)
	})

	t.Run("should add user trusted endpoints when allowed and enforcement is on by default", func(t *testing.T) {
		t.Setenv(allowUserTrustedEndpointsEnv, "true")
		t.Setenv(userTrustedEndpointsEnv, "https://*.custom1.com")

		endpoints, err := getTrustedEndpoints(context.Background(), publicResourceId)
		require.NoError(t, err)
		assert.Equal(t, []string{"https://*.prometheus.monitor.azure.com", "https://*.custom1.com"}, endpoints)
	})

	t.Run("should ignore user trusted endpoints when not allowed", func(t *testing.T) {
		t.Setenv(userTrustedEndpointsEnv, "https://*.custom1.com")

		endpoints, err := getTrustedEndpoints(context.Background(), publicResourceId)
		require.NoError(t, err)
		assert.Equal(t, []string{"https://*.prometheus.monitor.azure.com"}, endpoints)
	})

	t.Run("should ignore user trusted endpoints when enforcement is off", func(t *testing.T) {
		t.Setenv(enforceTrustedEndpointsEnv, "false")
		t.Setenv(allowUserTrustedEndpointsEnv, "true")
		t.Setenv(userTrustedEndpointsEnv, "https://*.custom1.com")

		endpoints, err := getTrustedEndpoints(context.Background(), publicResourceId)
		require.NoError(t, err)
		assert.Nil(t, endpoints)
	})

	t.Run("should return error for an invalid enforce setting", func(t *testing.T) {
		t.Setenv(enforceTrustedEndpointsEnv, "invalid")

		_, err := getTrustedEndpoints(context.Background(), publicResourceId)
		assert.ErrorContains(t, err, enforceTrustedEndpointsEnv)
	})

	t.Run("should return error for an invalid allow setting", func(t *testing.T) {
		t.Setenv(allowUserTrustedEndpointsEnv, "invalid")

		_, err := getTrustedEndpoints(context.Background(), publicResourceId)
		assert.ErrorContains(t, err, allowUserTrustedEndpointsEnv)
	})

	t.Run("should return error for an invalid resource ID", func(t *testing.T) {
		_, err := getTrustedEndpoints(context.Background(), "prometheus.monitor.azure.com")
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

	send := func(t *testing.T, ctx context.Context) error {
		opts := &sdkhttpclient.Options{CustomOptions: map[string]any{}}
		err := ConfigureAzureAuthentication(ctx, settings, azureSettings, opts, log.NewNullLogger())
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

	reset := func() {
		tokenRequests.Store(0)
		targetRequests.Store(0)
	}

	t.Run("should block a request to an endpoint that is not trusted", func(t *testing.T) {
		reset()

		err := send(t, context.Background())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "is not allowed by the datasource")
		assert.Zero(t, tokenRequests.Load(), "no token should be requested")
		assert.Zero(t, targetRequests.Load(), "the endpoint should not be called")
	})

	t.Run("should pass the endpoint check for a user trusted endpoint", func(t *testing.T) {
		reset()
		t.Setenv(enforceTrustedEndpointsEnv, "true")
		t.Setenv(allowUserTrustedEndpointsEnv, "true")
		t.Setenv(userTrustedEndpointsEnv, target.URL)

		// The request has no signed-in user, so getting the token fails. It must fail after the
		// endpoint check, not because of it.
		err := send(t, context.Background())
		require.Error(t, err)
		assert.NotContains(t, err.Error(), "is not allowed by the datasource")
		assert.Zero(t, targetRequests.Load(), "the endpoint should not be called without a token")
	})

	t.Run("should not check the endpoint when enforcement is off", func(t *testing.T) {
		reset()
		t.Setenv(enforceTrustedEndpointsEnv, "false")

		err := send(t, context.Background())
		require.Error(t, err)
		assert.NotContains(t, err.Error(), "is not allowed by the datasource")
	})

	t.Run("should return error for an invalid user trusted endpoint", func(t *testing.T) {
		t.Setenv(allowUserTrustedEndpointsEnv, "true")
		t.Setenv(userTrustedEndpointsEnv, "custom.example.com")

		opts := &sdkhttpclient.Options{CustomOptions: map[string]any{}}
		err := ConfigureAzureAuthentication(context.Background(), settings, azureSettings, opts, log.NewNullLogger())
		assert.ErrorContains(t, err, "invalid trusted endpoints configuration")
	})

	t.Run("should not check the endpoint when the feature toggle is enabled", func(t *testing.T) {
		reset()

		err := send(t, withFeatureToggles(disableTrustedEndpointsFeatureToggle))
		require.Error(t, err)
		assert.NotContains(t, err.Error(), "is not allowed by the datasource")
	})
}

func withFeatureToggles(toggles string) context.Context {
	cfg := config.NewGrafanaCfg(map[string]string{featuretoggles.EnabledFeatures: toggles})
	return config.WithGrafanaConfig(context.Background(), cfg)
}
