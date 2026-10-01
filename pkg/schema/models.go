// Package schema holds the dsconfig single source of truth for this plugin's
// data source configuration surface. See dsconfig.json and dsconfig_test.go.
package schema

import (
	"encoding/json"

	"github.com/grafana/grafana-prometheus-datasource/pkg/promlib/models"
)

// AzurePromOptions is promlib's PromOptions plus this plugin's Azure fields,
// checked against dsconfig.json by dsconfig_test.go. AzureCredentials is
// json.RawMessage (not parsed here) so its "any" schema type passes the
// conformance suite's type check.
type AzurePromOptions struct {
	models.PromOptions

	AzureCredentials        json.RawMessage `json:"azureCredentials"`
	AzureEndpointResourceID string          `json:"azureEndpointResourceId"`
	PrometheusTypeMigration bool            `json:"prometheus-type-migration"`
}
