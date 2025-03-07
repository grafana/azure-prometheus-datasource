import { cx } from '@emotion/css';
import { AzureCredentials } from '@grafana/azure-sdk';
import { DataSourceJsonData, DataSourceSettings } from '@grafana/data';
import { config } from '@grafana/runtime';
import { InlineField, InlineFieldRow, InlineSwitch, Input } from '@grafana/ui';
import React, { FormEvent, useMemo, useState } from 'react';
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

  const [overrideAudienceAllowed] = useState<boolean>(!!config.featureToggles.prometheusAzureOverrideAudience);
  const [overrideAudienceChecked, setOverrideAudienceChecked] = useState<boolean>(
    !!dataSourceConfig.jsonData.azureEndpointResourceId
  );

  const credentials = useMemo(() => getCredentials(dataSourceConfig), [dataSourceConfig]);

  const onCredentialsChange = (credentials: AzureCredentials): void => {
    onChange(updateCredentials(dataSourceConfig, credentials));
  };

  const onOverrideAudienceChange = (ev: FormEvent<HTMLInputElement>): void => {
    setOverrideAudienceChecked(ev.currentTarget.checked);
    if (!ev.currentTarget.checked) {
      onChange({
        ...dataSourceConfig,
        jsonData: { ...dataSourceConfig.jsonData, azureEndpointResourceId: undefined },
      });
    }
  };

  const onResourceIdChange = (ev: FormEvent<HTMLInputElement>): void => {
    if (overrideAudienceChecked) {
      onChange({
        ...dataSourceConfig,
        jsonData: { ...dataSourceConfig.jsonData, azureEndpointResourceId: ev.currentTarget.value },
      });
    }
  };

  // The auth type needs to be set on the first load of the data source
  useEffectOnce(() => {
    if (!dataSourceConfig.jsonData.authType) {
      onCredentialsChange(credentials);
    }
  });

  return (
    <>
      <h6>Azure authentication</h6>
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
