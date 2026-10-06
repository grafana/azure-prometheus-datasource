package schema_test

import (
	_ "embed"
	"net/http"
	"testing"

	dsconfigschema "github.com/grafana/dsconfig/schema"
	"github.com/grafana/grafana-plugin-sdk-go/experimental/pluginschema"
	"github.com/grafana/grafana-prometheus-datasource/pkg/promlib/models"
	"k8s.io/kube-openapi/pkg/spec3"

	"github.com/grafana/azure-prometheus-datasource/pkg/schema"
)

//go:embed dsconfig.json
var configSchemaJSON []byte

// secure builds the secure value map for a DataSource resource's top-level "secure"
// field (github.com/grafana/grafana/pkg/apis/datasource/v0alpha1.DataSource.Secure,
// common.InlineSecureValues): one {"create": "<value>"} entry per secret key, matching
// what a POST/PUT of a new secret looks like through the v0alpha1 datasource API.
func secure(values map[string]string) map[string]any {
	out := make(map[string]any, len(values))
	for key, value := range values {
		out[key] = map[string]any{"create": value}
	}
	return out
}

//go:generate go test -run TestPlugin -generateArtifacts
func TestPlugin(t *testing.T) {
	dsconfigschema.RunPluginTests(t, dsconfigschema.PluginUnderTest{
		ID:                "grafana-azureprometheus-datasource",
		ConfigSchemaJSON:  configSchemaJSON,
		SettingsJSONModel: schema.AzurePromOptions{},
		SecureKeys:        []string{"azureClientSecret", "clientSecret", "tlsCACert", "tlsClientCert", "tlsClientKey"},
		// Each example's Value is shaped like a v0alpha1 DataSource resource body
		// (github.com/grafana/grafana/pkg/apis/datasource/v0alpha1.DataSource): "spec"
		// holds the fields declared in dsconfig.json's root/jsonData targets (url,
		// jsonData — never secrets, per SchemaSpecHasNoSecureJSON), and "secure" holds
		// the write-only secret values, keyed by name with a "create" payload.
		SettingsExamples: &pluginschema.SettingsExamples{
			Examples: map[string]*spec3.Example{
				"": {
					ExampleProps: spec3.ExampleProps{
						Summary:     "App Registration (default)",
						Description: "The most common setup: an Azure AD App Registration with a client secret. Only spec.url and the tenant/client IDs need filling in.",
						Value: map[string]any{
							"spec": map[string]any{
								"url": "https://myworkspace-abcd.eastus.prometheus.monitor.azure.com",
								"jsonData": map[string]any{
									"azureCredentials": map[string]any{
										"authType":   "clientsecret",
										"azureCloud": "AzureCloud",
										"tenantId":   "00000000-0000-0000-0000-000000000000",
										"clientId":   "00000000-0000-0000-0000-000000000000",
									},
									"httpMethod": http.MethodPost,
								},
							},
							"secure": secure(map[string]string{
								"azureClientSecret": "REPLACE_WITH_CLIENT_SECRET",
							}),
						},
					},
				},
				"managedIdentity": {
					ExampleProps: spec3.ExampleProps{
						Summary:     "Managed Identity",
						Description: "Authenticate with the managed identity assigned to the Grafana instance. Only available when the Grafana instance has azure.managedIdentityEnabled. No secrets required.",
						Value: map[string]any{
							"spec": map[string]any{
								"url": "https://myworkspace-abcd.eastus.prometheus.monitor.azure.com",
								"jsonData": map[string]any{
									"azureCredentials": map[string]any{
										"authType": "msi",
									},
									"httpMethod": http.MethodPost,
								},
							},
						},
					},
				},
				"workloadIdentity": {
					ExampleProps: spec3.ExampleProps{
						Summary:     "Workload Identity",
						Description: "Authenticate with Azure AD Workload Identity federation. Only available when the Grafana instance has azure.workloadIdentityEnabled. No secrets required.",
						Value: map[string]any{
							"spec": map[string]any{
								"url": "https://myworkspace-abcd.eastus.prometheus.monitor.azure.com",
								"jsonData": map[string]any{
									"azureCredentials": map[string]any{
										"authType": "workloadidentity",
									},
									"httpMethod": http.MethodPost,
								},
							},
						},
					},
				},
				"currentUserWithFallback": {
					ExampleProps: spec3.ExampleProps{
						Summary:     "Current User with App Registration fallback",
						Description: "Forward the signed-in Grafana user's Azure identity to Prometheus. serviceCredentials provides a fallback App Registration used for alerting, recorded queries and reporting, which run without a signed-in user. Only available when the Grafana instance has azure.userIdentityEnabled.",
						Value: map[string]any{
							"spec": map[string]any{
								"url": "https://myworkspace-abcd.eastus.prometheus.monitor.azure.com",
								"jsonData": map[string]any{
									"azureCredentials": map[string]any{
										"authType":                  "currentuser",
										"serviceCredentialsEnabled": true,
										"serviceCredentials": map[string]any{
											"authType":   "clientsecret",
											"azureCloud": "AzureCloud",
											"tenantId":   "00000000-0000-0000-0000-000000000000",
											"clientId":   "00000000-0000-0000-0000-000000000000",
										},
									},
									"httpMethod": http.MethodPost,
								},
							},
							"secure": secure(map[string]string{
								"azureClientSecret": "REPLACE_WITH_CLIENT_SECRET",
							}),
						},
					},
				},
				"legacyClientSecret": {
					ExampleProps: spec3.ExampleProps{
						Summary:     "Legacy client secret key",
						Description: "A datasource provisioned before the credential migration to azureClientSecret. The backend still reads secureJsonData.clientSecret as a fallback when azureClientSecret is absent.",
						Value: map[string]any{
							"spec": map[string]any{
								"url": "https://myworkspace-abcd.eastus.prometheus.monitor.azure.com",
								"jsonData": map[string]any{
									"azureCredentials": map[string]any{
										"authType":   "clientsecret",
										"azureCloud": "AzureCloud",
										"tenantId":   "00000000-0000-0000-0000-000000000000",
										"clientId":   "00000000-0000-0000-0000-000000000000",
									},
									"httpMethod": http.MethodPost,
								},
							},
							"secure": secure(map[string]string{
								"clientSecret": "REPLACE_WITH_CLIENT_SECRET",
							}),
						},
					},
				},
				"migratedFromPrometheus": {
					ExampleProps: spec3.ExampleProps{
						Summary:     "Migrated from vanilla Prometheus",
						Description: "A data source migrated from the vanilla Prometheus plugin. prometheus-type-migration=true makes the editor render a 'Data source migrated' warning banner until an operator dismisses or replaces it.",
						Value: map[string]any{
							"spec": map[string]any{
								"url": "https://myworkspace-abcd.eastus.prometheus.monitor.azure.com",
								"jsonData": map[string]any{
									"azureCredentials": map[string]any{
										"authType":   "clientsecret",
										"azureCloud": "AzureCloud",
										"tenantId":   "00000000-0000-0000-0000-000000000000",
										"clientId":   "00000000-0000-0000-0000-000000000000",
									},
									"httpMethod":                http.MethodPost,
									"prometheus-type-migration": true,
									"prometheusType":            models.PromApplicationPrometheus,
								},
							},
							"secure": secure(map[string]string{
								"azureClientSecret": "REPLACE_WITH_CLIENT_SECRET",
							}),
						},
					},
				},
				"tlsMutualAuth": {
					ExampleProps: spec3.ExampleProps{
						Summary:     "App Registration with mTLS",
						Description: "Azure AD auth combined with TLS client authentication in front of the Prometheus endpoint. TLS is independent of the Azure auth method and can be layered on top of any authType.",
						Value: map[string]any{
							"spec": map[string]any{
								"url": "https://myworkspace-abcd.eastus.prometheus.monitor.azure.com",
								"jsonData": map[string]any{
									"azureCredentials": map[string]any{
										"authType":   "clientsecret",
										"azureCloud": "AzureCloud",
										"tenantId":   "00000000-0000-0000-0000-000000000000",
										"clientId":   "00000000-0000-0000-0000-000000000000",
									},
									"httpMethod": http.MethodPost,
									"tlsAuth":    true,
									"serverName": "myworkspace-abcd.eastus.prometheus.monitor.azure.com",
								},
							},
							"secure": secure(map[string]string{
								"azureClientSecret": "REPLACE_WITH_CLIENT_SECRET",
								"tlsClientCert":     "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----",
								"tlsClientKey":      "-----BEGIN RSA PRIVATE KEY-----\n...\n-----END RSA PRIVATE KEY-----",
							}),
						},
					},
				},
			},
		},
	})
}
