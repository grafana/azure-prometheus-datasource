import { AzureCredentials } from '@grafana/azure-sdk';
import { DataSourceJsonData } from '@grafana/data';
import { Trans } from '@grafana/i18n';
import { config } from '@grafana/runtime';
import React from 'react';
import { useEffectOnce } from 'react-use';

import { getAzureCloudOptions } from './AzureCredentialsConfig';
import { AzureCredentialsForm } from './AzureCredentialsForm';

export interface HttpSettingsBaseProps<JSONData extends DataSourceJsonData = any, SecureJSONData = any> {
  credentials: AzureCredentials;
  onCredentialsChange: (updatedCredentials: AzureCredentials) => void;
  disabled?: boolean;
}

export const AzureAuthSettings = (props: HttpSettingsBaseProps) => {
  const { credentials, onCredentialsChange, disabled } = props;

  // The auth type needs to be set on the first load of the data source
  useEffectOnce(() => {
    if (!credentials.authType) {
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
        disabled={disabled}
      />
    </>
  );
};

export default AzureAuthSettings;
