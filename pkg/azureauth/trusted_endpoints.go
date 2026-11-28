package azureauth

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/grafana/grafana-plugin-sdk-go/config"
)

// These settings have the same names, defaults and behavior as in the Azure Data Explorer data
// source. Grafana passes the keys of the [plugin.grafana-azureprometheus-datasource] section in
// grafana.ini to this plugin only, as GF_PLUGIN_<KEY>, so the names do not clash with other plugins.
const (
	enforceTrustedEndpointsEnv   = "GF_PLUGIN_ENFORCE_TRUSTED_ENDPOINTS"
	allowUserTrustedEndpointsEnv = "GF_PLUGIN_ALLOW_USER_TRUSTED_ENDPOINTS"
	userTrustedEndpointsEnv      = "GF_PLUGIN_USER_TRUSTED_ENDPOINTS"

	// disableTrustedEndpointsFeatureToggle turns enforcement off, the same as the
	// adxDisableTrustedEndpoints feature toggle of the Azure Data Explorer data source.
	disableTrustedEndpointsFeatureToggle = "azurePrometheusDisableTrustedEndpoints"
)

// getTrustedEndpoints returns the endpoints the data source may send Azure credentials to.
// A nil result means enforcement is disabled and any endpoint is allowed.
//
// Enforcement is on by default for all authentication methods. The
// azurePrometheusDisableTrustedEndpoints feature toggle disables it. An explicitly set
// enforce_trusted_endpoints overrides the default and the feature toggle. The endpoints in
// user_trusted_endpoints are added only when both enforce_trusted_endpoints and
// allow_user_trusted_endpoints are true.
func getTrustedEndpoints(ctx context.Context, prometheusResourceId string) ([]string, error) {
	enforce := !config.GrafanaConfigFromContext(ctx).FeatureToggles().IsEnabled(disableTrustedEndpointsFeatureToggle)

	if _, ok := os.LookupEnv(enforceTrustedEndpointsEnv); ok {
		var err error
		if enforce, err = envBoolOrDefault(enforceTrustedEndpointsEnv, true); err != nil {
			return nil, err
		}
	}

	if !enforce {
		return nil, nil
	}

	endpoint, err := trustedEndpointFromResourceId(prometheusResourceId)
	if err != nil {
		return nil, err
	}
	endpoints := []string{endpoint}

	allowUserEndpoints, err := envBoolOrDefault(allowUserTrustedEndpointsEnv, false)
	if err != nil {
		return nil, err
	}

	if allowUserEndpoints {
		for _, userEndpoint := range strings.Split(os.Getenv(userTrustedEndpointsEnv), ",") {
			if userEndpoint = strings.TrimSpace(userEndpoint); userEndpoint != "" {
				endpoints = append(endpoints, userEndpoint)
			}
		}
	}

	return endpoints, nil
}

// trustedEndpointFromResourceId turns the Prometheus resource ID of a cloud into the endpoint of
// its Azure Monitor workspaces. Workspace query endpoints are subdomains of the resource ID host,
// for example https://<workspace>.<region>.prometheus.monitor.azure.com.
func trustedEndpointFromResourceId(prometheusResourceId string) (string, error) {
	resourceId, err := url.Parse(prometheusResourceId)
	if err != nil || resourceId.Scheme == "" || resourceId.Hostname() == "" {
		return "", fmt.Errorf("endpoint resource ID (audience) '%s' invalid", prometheusResourceId)
	}

	return fmt.Sprintf("%s://*.%s", resourceId.Scheme, resourceId.Host), nil
}

func envBoolOrDefault(key string, defaultValue bool) (bool, error) {
	strValue := os.Getenv(key)
	if strValue == "" {
		return defaultValue, nil
	}

	value, err := strconv.ParseBool(strValue)
	if err != nil {
		return false, fmt.Errorf("environment variable '%s' has invalid bool value '%s'", key, strValue)
	}

	return value, nil
}
