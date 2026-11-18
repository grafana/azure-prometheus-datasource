package azureauth

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/grafana/grafana-azure-sdk-go/v2/azcredentials"
)

// These are the same settings as the Azure Data Explorer plugin. Grafana passes the keys of the
// [plugin.grafana-azureprometheus-datasource] section in grafana.ini to this plugin only, as
// GF_PLUGIN_<KEY>, so the names do not clash with other plugins.
const (
	enforceTrustedEndpointsEnv   = "GF_PLUGIN_ENFORCE_TRUSTED_ENDPOINTS"
	allowUserTrustedEndpointsEnv = "GF_PLUGIN_ALLOW_USER_TRUSTED_ENDPOINTS"
	userTrustedEndpointsEnv      = "GF_PLUGIN_USER_TRUSTED_ENDPOINTS"
)

// getTrustedEndpoints returns the endpoints the data source may send Azure credentials to, the
// same control the Azure Data Explorer data source has. A nil result means no restriction.
//
// The list is always on for user identity (Current User and On-Behalf-Of). For other
// credentials it is on when enforce_trusted_endpoints is set.
func getTrustedEndpoints(credentials azcredentials.AzureCredentials, prometheusResourceId string) ([]string, error) {
	enforce, err := envBoolOrDefault(enforceTrustedEndpointsEnv, false)
	if err != nil {
		return nil, err
	}

	if !enforce && !isUserIdentity(credentials) {
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

func isUserIdentity(credentials azcredentials.AzureCredentials) bool {
	switch credentials.AzureAuthType() {
	case azcredentials.AzureAuthCurrentUserIdentity, azcredentials.AzureAuthClientSecretObo:
		return true
	default:
		return false
	}
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
