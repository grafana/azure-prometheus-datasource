package azureauth

import (
	"fmt"
	"net/url"
	"path"

	"github.com/grafana/grafana-azure-sdk-go/v2/azcredentials"
	"github.com/grafana/grafana-azure-sdk-go/v2/azhttpclient"
	"github.com/grafana/grafana-azure-sdk-go/v2/azsettings"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	sdkhttpclient "github.com/grafana/grafana-plugin-sdk-go/backend/httpclient"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"

	"github.com/grafana/grafana/pkg/promlib/utils"
)

func ConfigureAzureAuthentication(settings backend.DataSourceInstanceSettings, azureSettings *azsettings.AzureSettings, clientOpts *sdkhttpclient.Options, log log.Logger) error {
	jsonData, err := utils.GetJsonData(settings)
	if err != nil {
		return fmt.Errorf("failed to get jsonData: %w", err)
	}
	credentials, err := azcredentials.FromDatasourceData(jsonData, settings.DecryptedSecureJSONData)
	if err != nil {
		err = fmt.Errorf("invalid Azure credentials: %w", err)
		return err
	}

	if credentials != nil {
		var scopes []string

		if scopes, err = getPrometheusScopes(azureSettings, credentials); err != nil {
			return err
		}

		authOpts := azhttpclient.NewAuthOptions(azureSettings)
		authOpts.AllowUserIdentity()
		authOpts.Scopes(scopes)
		azhttpclient.AddAzureAuthentication(clientOpts, authOpts, credentials)
	}

	return nil
}

func getPrometheusScopes(settings *azsettings.AzureSettings, credentials azcredentials.AzureCredentials) ([]string, error) {
	// Extract cloud from credentials
	azureCloud, err := azcredentials.GetAzureCloud(settings, credentials)
	if err != nil {
		return nil, err
	}

	cloudSettings, err := settings.GetCloud(azureCloud)
	if err != nil {
		return nil, err
	}

	// Get scopes for the given cloud
	resourceIdS, ok := cloudSettings.Properties["prometheusResourceId"]
	if !ok {
		err := fmt.Errorf("the Azure cloud '%s' doesn't have configuration for Prometheus", azureCloud)
		return nil, err
	}
	return audienceToScopes(resourceIdS)
}

func audienceToScopes(audience string) ([]string, error) {
	resourceId, err := url.Parse(audience)
	if err != nil || resourceId.Scheme == "" || resourceId.Host == "" {
		err = fmt.Errorf("endpoint resource ID (audience) '%s' invalid", audience)
		return nil, err
	}

	resourceId.Path = path.Join(resourceId.Path, ".default")
	scopes := []string{resourceId.String()}
	return scopes, nil
}
