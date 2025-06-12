import { AzureCredentials } from '@grafana/azure-sdk';
import { DataSourceJsonData, DataSourceSettings } from '@grafana/data';
import { Trans } from '@grafana/i18n';
import { config } from '@grafana/runtime';
import React, { useMemo } from 'react';
import { useEffectOnce } from 'react-use';

import { getAzureCloudOptions, getCredentials, updateCredentials } from './AzureCredentialsConfig';
import { AzureCredentialsForm } from './AzureCredentialsForm';

export interface HttpSettingsBaseProps<JSONData extends DataSourceJsonData = any, SecureJSONData = any> {
  /** The configuration object of the data source */
  dataSourceConfig: DataSourceSettings<JSONData, SecureJSONData>;
  /** Callback for handling changes to the configuration object */
  onChange: (config: DataSourceSettings<JSONData, SecureJSONData>) => void;
  /** Show the Forward OAuth identity option */
  showForwardOAuthIdentityOption?: boolean;
}

export const AzureAuthSettings = (props: HttpSettingsBaseProps) => {
  const { dataSourceConfig, onChange } = props;

  const credentials = useMemo(() => getCredentials(dataSourceConfig), [dataSourceConfig]);

  const onCredentialsChange = (credentials: AzureCredentials): void => {
    onChange(updateCredentials(dataSourceConfig, credentials));
  };

  // The auth type needs to be set on the first load of the data source
  useEffectOnce(() => {
    if (!dataSourceConfig.jsonData.authType) {
      onCredentialsChange(credentials);
    }
  });

  return (
    <>
      <h6>
        <Trans i18nKey="configuration.azure-auth-settings.azure-authentication">Azure authentication</Trans>
      </h6>
      <AzureCredentialsForm
        managedIdentityEnabled={config.azure.managedIdentityEnabled}
        workloadIdentityEnabled={config.azure.workloadIdentityEnabled}
        userIdentityEnabled={config.azure.userIdentityEnabled}
        credentials={credentials}
        azureCloudOptions={getAzureCloudOptions()}
        onCredentialsChange={onCredentialsChange}
        disabled={dataSourceConfig.readOnly}
      />
    </>
  );
};

export default AzureAuthSettings;
